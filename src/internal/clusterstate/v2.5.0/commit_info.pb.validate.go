// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: internal/clusterstate/v2.5.0/commit_info.proto

package v2_5_0

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CommitInfo with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CommitInfo) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CommitInfo with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CommitInfoMultiError, or
// nil if none found.
func (m *CommitInfo) ValidateAll() error {
	return m.validate(true)
}

func (m *CommitInfo) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCommit()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Commit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Commit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCommit()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Commit",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetOrigin()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Origin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Origin",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetOrigin()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Origin",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Description

	if all {
		switch v := interface{}(m.GetParentCommit()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "ParentCommit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "ParentCommit",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetParentCommit()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "ParentCommit",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetChildCommits() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CommitInfoValidationError{
						field:  fmt.Sprintf("ChildCommits[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommitInfoValidationError{
						field:  fmt.Sprintf("ChildCommits[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommitInfoValidationError{
					field:  fmt.Sprintf("ChildCommits[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetStarted()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Started",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Started",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetStarted()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Started",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetFinishing()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Finishing",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Finishing",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFinishing()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Finishing",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetFinished()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Finished",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Finished",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetFinished()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Finished",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	for idx, item := range m.GetDirectProvenance() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CommitInfoValidationError{
						field:  fmt.Sprintf("DirectProvenance[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CommitInfoValidationError{
						field:  fmt.Sprintf("DirectProvenance[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CommitInfoValidationError{
					field:  fmt.Sprintf("DirectProvenance[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Error

	// no validation rules for SizeBytesUpperBound

	if all {
		switch v := interface{}(m.GetDetails()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Details",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommitInfoValidationError{
					field:  "Details",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDetails()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommitInfoValidationError{
				field:  "Details",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CommitInfoMultiError(errors)
	}

	return nil
}

// CommitInfoMultiError is an error wrapping multiple validation errors
// returned by CommitInfo.ValidateAll() if the designated constraints aren't met.
type CommitInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CommitInfoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CommitInfoMultiError) AllErrors() []error { return m }

// CommitInfoValidationError is the validation error returned by
// CommitInfo.Validate if the designated constraints aren't met.
type CommitInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CommitInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CommitInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CommitInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CommitInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CommitInfoValidationError) ErrorName() string { return "CommitInfoValidationError" }

// Error satisfies the builtin error interface
func (e CommitInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCommitInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CommitInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CommitInfoValidationError{}