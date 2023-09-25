// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: task/task.proto

package task

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

// Validate checks the field values on Group with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Group) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Group with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in GroupMultiError, or nil if none found.
func (m *Group) ValidateAll() error {
	return m.validate(true)
}

func (m *Group) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Namespace

	// no validation rules for Group

	if len(errors) > 0 {
		return GroupMultiError(errors)
	}

	return nil
}

// GroupMultiError is an error wrapping multiple validation errors returned by
// Group.ValidateAll() if the designated constraints aren't met.
type GroupMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GroupMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GroupMultiError) AllErrors() []error { return m }

// GroupValidationError is the validation error returned by Group.Validate if
// the designated constraints aren't met.
type GroupValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GroupValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GroupValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GroupValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GroupValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GroupValidationError) ErrorName() string { return "GroupValidationError" }

// Error satisfies the builtin error interface
func (e GroupValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGroup.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GroupValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GroupValidationError{}

// Validate checks the field values on TaskInfo with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TaskInfo) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TaskInfo with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TaskInfoMultiError, or nil
// if none found.
func (m *TaskInfo) ValidateAll() error {
	return m.validate(true)
}

func (m *TaskInfo) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	if all {
		switch v := interface{}(m.GetGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TaskInfoValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TaskInfoValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TaskInfoValidationError{
				field:  "Group",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for State

	// no validation rules for Reason

	// no validation rules for InputType

	// no validation rules for InputData

	if len(errors) > 0 {
		return TaskInfoMultiError(errors)
	}

	return nil
}

// TaskInfoMultiError is an error wrapping multiple validation errors returned
// by TaskInfo.ValidateAll() if the designated constraints aren't met.
type TaskInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TaskInfoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TaskInfoMultiError) AllErrors() []error { return m }

// TaskInfoValidationError is the validation error returned by
// TaskInfo.Validate if the designated constraints aren't met.
type TaskInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TaskInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TaskInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TaskInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TaskInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TaskInfoValidationError) ErrorName() string { return "TaskInfoValidationError" }

// Error satisfies the builtin error interface
func (e TaskInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTaskInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TaskInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TaskInfoValidationError{}

// Validate checks the field values on ListTaskRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListTaskRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListTaskRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListTaskRequestMultiError, or nil if none found.
func (m *ListTaskRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListTaskRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetGroup() == nil {
		err := ListTaskRequestValidationError{
			field:  "Group",
			reason: "value is required",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ListTaskRequestValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ListTaskRequestValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ListTaskRequestValidationError{
				field:  "Group",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ListTaskRequestMultiError(errors)
	}

	return nil
}

// ListTaskRequestMultiError is an error wrapping multiple validation errors
// returned by ListTaskRequest.ValidateAll() if the designated constraints
// aren't met.
type ListTaskRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListTaskRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListTaskRequestMultiError) AllErrors() []error { return m }

// ListTaskRequestValidationError is the validation error returned by
// ListTaskRequest.Validate if the designated constraints aren't met.
type ListTaskRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListTaskRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListTaskRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListTaskRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListTaskRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListTaskRequestValidationError) ErrorName() string { return "ListTaskRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListTaskRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListTaskRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListTaskRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListTaskRequestValidationError{}