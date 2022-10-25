import React from 'react';
import {closeIcon} from '@jupyterlab/ui-components';
import {useDatum} from './hooks/useDatum';
import {caretLeftIcon, caretRightIcon} from '@jupyterlab/ui-components';
import {
  CrossInputSpec,
  CurrentDatumResponse,
  PfsInput,
} from 'plugins/mount/types';

type DatumProps = {
  showDatum: boolean;
  setShowDatum: (shouldShow: boolean) => Promise<void>;
  keepMounted: boolean;
  setKeepMounted: (keep: boolean) => void;
  refresh: (path: string) => void;
  pollRefresh: () => Promise<void>;
  currentDatumInfo?: CurrentDatumResponse;
  repoViewInputSpec: CrossInputSpec | PfsInput;
};

const placeholderText = `pfs:
  repo: images
  glob: /*
`;

const Datum: React.FC<DatumProps> = ({
  showDatum,
  setShowDatum,
  keepMounted,
  setKeepMounted,
  refresh,
  pollRefresh,
  currentDatumInfo,
  repoViewInputSpec,
}) => {
  const {
    loading,
    shouldShowCycler,
    currDatum,
    currIdx,
    setCurrIdx,
    inputSpec,
    setInputSpec,
    callMountDatums,
    callUnmountAll,
    errorMessage,
    saveInputSpec,
    initialInputSpec,
  } = useDatum(
    showDatum,
    keepMounted,
    refresh,
    pollRefresh,
    repoViewInputSpec,
    currentDatumInfo,
  );

  return (
    <div className="pachyderm-mount-datum-base">
      <div className="pachyderm-mount-datum-back">
        <button
          data-testid="Datum__back"
          className="pachyderm-button-link"
          onClick={async () => {
            await callUnmountAll();
            saveInputSpec();
            setKeepMounted(false);
            await setShowDatum(false);
          }}
        >
          Back{' '}
          <closeIcon.react
            tag="span"
            className="pachyderm-mount-icon-padding"
          />
        </button>
      </div>

      <span className="pachyderm-mount-datum-subheading">Test Datums</span>

      <div className="pachyderm-mount-datum-input-wrapper">
        <label className="pachyderm-mount-datum-label" htmlFor="inputSpec">
          Input spec
        </label>
        <textarea
          className="pachyderm-input"
          data-testid="Datum__inputSpecInput"
          style={{minHeight: '200px'}}
          name="inputSpec"
          value={inputSpec}
          onChange={(e: any) => {
            setInputSpec(e.target.value);
          }}
          disabled={loading}
          placeholder={
            Object.keys(initialInputSpec).length === 0
              ? placeholderText
              : undefined
          }
        ></textarea>
        <span
          className="pachyderm-mount-datum-error"
          data-testid="Datum__errorMessage"
        >
          {errorMessage}
        </span>
        <button
          data-testid="Datum__mountDatums"
          className="pachyderm-button-link"
          onClick={callMountDatums}
          style={{padding: '0.5rem'}}
        >
          Mount Datums
        </button>
        {shouldShowCycler && (
          <div
            className="pachyderm-mount-datum-cycler"
            data-testid="Datum__cycler"
          >
            Datum
            <div style={{display: 'flex'}}>
              <button
                className="pachyderm-button-link"
                data-testid="Datum__cyclerLeft"
                disabled={currIdx <= 0}
                onClick={() => {
                  if (currIdx >= 1) {
                    setCurrIdx(currIdx - 1);
                  }
                }}
              >
                <caretLeftIcon.react
                  tag="span"
                  className="pachyderm-mount-datum-left"
                />
              </button>
              {'(' + (currIdx + 1) + '/' + currDatum.num_datums + ')'}
              <button
                className="pachyderm-button-link"
                data-testid="Datum__cyclerRight"
                disabled={currIdx >= currDatum.num_datums - 1}
                onClick={() => {
                  if (currIdx < currDatum.num_datums - 1) {
                    setCurrIdx(currIdx + 1);
                  }
                }}
              >
                <caretRightIcon.react
                  tag="span"
                  className="pachyderm-mount-datum-right"
                />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Datum;
