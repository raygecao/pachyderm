package v2_10_0

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pachyderm/pachyderm/v2/src/internal/clusterstate/migrationutils"
	"github.com/pachyderm/pachyderm/v2/src/internal/errors"
	"github.com/pachyderm/pachyderm/v2/src/internal/log"
	"github.com/pachyderm/pachyderm/v2/src/internal/meters"
	"github.com/pachyderm/pachyderm/v2/src/internal/migrations"
	"github.com/pachyderm/pachyderm/v2/src/internal/pachsql"
	"github.com/pachyderm/pachyderm/v2/src/internal/pctx"
	"go.uber.org/zap"
)

func normalizeCommitTotals(ctx context.Context, env migrations.Env) error {
	ctx = pctx.Child(ctx, "normalizeCommitTotals", pctx.WithCounter("migration_errors", 0))
	tx := env.Tx
	if _, err := tx.ExecContext(ctx, `ALTER TABLE pfs.commits ADD COLUMN total_fileset_id UUID REFERENCES storage.filesets(id)`); err != nil {
		return errors.Wrap(err, "add total_fileset column to pfs.commits")
	}

	count := uint64(0)
	if err := tx.GetContext(ctx, &count, `SELECT count(fileset_id) FROM pfs.commit_totals;`); err != nil {
		return errors.Wrap(err, "counting rows in pfs.commit_totals")
	}
	if count == 0 {
		return nil
	}
	pageSize := uint64(10)
	totalPages := count / pageSize
	if count%pageSize > 0 {
		totalPages++
	}
	errs := make([]error, 0)
	for i := uint64(0); i < totalPages; i++ {
		log.Info(ctx, "normalizing pfs.commit_totals", zap.Uint64("current", i*pageSize), zap.Uint64("total", count))
		commitsToFilesets, err := getCommitsToTotalFilesets(ctx, tx, pageSize, i*pageSize)
		if err != nil {
			errs = append(errs, err)
			continue // attempt to migrate the other pages
		}
		if err := migrateTotals(ctx, tx, commitsToFilesets); err != nil {
			errs = append(errs, err)
		}
		if err := errIfUnmigratedTotals(ctx, tx, commitsToFilesets, pageSize, i*pageSize); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func getCommitsToTotalFilesets(ctx context.Context, tx *pachsql.Tx, pageSize, offset uint64) (map[string]string, error) {
	type commitTotalsRow struct {
		Commit  string `db:"commit_key"`
		Fileset string `db:"fileset_id"`
	}
	commitTotalsPage := make([]*commitTotalsRow, 0)
	err := sqlx.SelectContext(ctx, tx, &commitTotalsPage, `SELECT commit_id AS "commit_key", fileset_id AS "fileset_id" FROM pfs.commit_totals ORDER BY commit_id LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		meters.Inc(ctx, "migration_errors", 1)
		return nil, errors.Wrap(err, "getting commits and filesets from pfs.commit_totals")
	}
	commitsToFilesets := make(map[string]string)
	for _, row := range commitTotalsPage {
		commitsToFilesets[row.Commit] = row.Fileset
	}
	return commitsToFilesets, nil
}

func migrateTotals(ctx context.Context, tx *pachsql.Tx, commitsToFilesets map[string]string) error {
	batcher := migrationutils.NewSimplePostgresBatcher(tx)
	var errs []error
	for commit, fileset := range commitsToFilesets {
		fmt.Println("UPDATE pfs.commits SET total_fileset_id =", fileset, "WHERE commit_id = ", commit)
		if err := batcher.Add(ctx, `UPDATE pfs.commits SET total_fileset_id = $1 WHERE commit_id = $2`, fileset, commit); err != nil {
			meters.Inc(ctx, "migration_errors", 1)
			errs = append(errs, errors.Wrap(err, "adding to batcher: updating pfs.commits table, commit: "+commit))
		}
	}
	if err := batcher.Flush(ctx); err != nil {
		meters.Inc(ctx, "migration_errors", 1)
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func errIfUnmigratedTotals(ctx context.Context, tx *pachsql.Tx, commitsToFilesets map[string]string, pageSize, offset uint64) error {
	var errs []error
	count := struct {
		Migrated uint64 `db:"count_migrated"`
		Totals   uint64 `db:"count_totals"`
	}{}
	rows, err := tx.QueryxContext(ctx, `SELECT count(total_fileset_id) AS "count_migrated" FROM pfs.commits GROUP BY commit_id ORDER BY commit_id LIMIT $1 OFFSET $2`, pageSize, offset)
	if err != nil {
		meters.Inc(ctx, "migration_errors", 1)
		errs = append(errs, errors.Wrap(err, "getting number of commits migrated"))
	}
	for rows != nil && rows.Next() {
		if err := rows.Scan(&count.Migrated); err != nil {
			return errors.Wrap(err, "scanning rows")
		}
		fmt.Println(count)
	}
	defer rows.Close()
	if len(errs) > 0 {
		return errors.Join(errs...) // There's no point in validating if the count queries failed.
	}
	if count.Totals == count.Migrated {
		return nil
	}
	meters.Inc(ctx, "migration_errors", 1)
	errs = append(errs, fmt.Errorf("unmigrated total filesets: %d", count.Totals-count.Migrated))
	for commit, fileset := range commitsToFilesets {
		// Linear scanning the commits table is not ideal, but it should help the customer identify what hasn't been migrated.
		dest := "" // Sadly, a valid dest is required by GetContext, even if it is not used.
		err := tx.GetContext(ctx, &dest, `SELECT total_fileset_id FROM pfs.commits WHERE commit_id = $1 ORDER BY commit_id LIMIT $2 OFFSET $3`, commit, pageSize, offset)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				errs = append(errs, errors.Wrap(err, "getting total_fileset_id for commit "+commit))
			}
			errs = append(errs, fmt.Errorf("not migrated: commit: %s, fileset: %s", commit, fileset))
			meters.Inc(ctx, "migration_errors", 1)
		}
	}
	return errors.Join(errs...)
}
