// Package errs implements functions to manipulate error instances.
package errs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	nilAngleString = "<nil>"
)

// Error type is a implementation of error interface.
// This type is for wrapping cause error instance.
type Error struct {
	wrapFlag bool
	Err      error
	Cause    error
	Context  map[string]interface{}
}

var _ error = (*Error)(nil)          //Error type is compatible with error interface
var _ fmt.Stringer = (*Error)(nil)   //Error type is compatible with fmt.Stringer interface
var _ fmt.GoStringer = (*Error)(nil) //Error type is compatible with fmt.GoStringer interface
var _ fmt.Formatter = (*Error)(nil)  //Error type is compatible with fmt.Formatter interface
var _ json.Marshaler = (*Error)(nil) //Error type is compatible with json.Marshaler interface

// ErrorContextFunc type is self-referential function type for New and Wrap functions. (functional options pattern)
type ErrorContextFunc func(*Error)

// New function returns an error instance with message and context informations.
func New(msg string, opts ...ErrorContextFunc) error {
	if len(msg) == 0 {
		return nil
	}
	return newError(errors.New(msg), false, 2, opts...)
}

// Wrap function returns a wrapping error instance with context informations.
func Wrap(err error, opts ...ErrorContextFunc) error {
	if err == nil {
		return nil
	}
	return newError(err, true, 2, opts...)
}

// newError returns error instance. (internal)
func newError(err error, wrapFlag bool, depth int, opts ...ErrorContextFunc) error {
	we := &Error{Err: err, wrapFlag: wrapFlag}
	//caller function name
	if fname, _, _ := caller(depth); len(fname) > 0 {
		we = we.SetContext("function", fname)
	}
	//other params
	for _, opt := range opts {
		opt(we)
	}
	return we
}

// WithContext function returns ErrorContextFunc function value.
// This function is used in New and Wrap functions that represents context (key/value) data.
func WithContext(name string, value interface{}) ErrorContextFunc {
	return func(e *Error) {
		_ = e.SetContext(name, value)
	}
}

// WithCause function returns ErrorContextFunc function value.
// This function is used in New and Wrap functions that represents context (key/value) data.
func WithCause(err error) ErrorContextFunc {
	return func(e *Error) {
		_ = e.SetCause(err)
	}
}

// SetContext method sets context information
func (e *Error) SetContext(name string, value interface{}) *Error {
	if e == nil {
		return e
	}
	if e.Context == nil {
		e.Context = map[string]interface{}{}
	}
	if len(name) > 0 {
		e.Context[name] = value
	}
	return e
}

// SetCause method sets cause error instance
func (e *Error) SetCause(err error) *Error {
	if e == nil {
		return e
	}
	e.Cause = err
	return e
}

// Unwrap method returns cause error in Error instance.
// This method is used in errors.Unwrap function.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	if e.Cause == nil {
		if e.wrapFlag {
			return e.Err
		}
		return errors.Unwrap(e.Err)
	}
	return e.Cause
}

// Is method reports whether any error in error's chain matches cause of target error.
// This method is used in errors.Is function.
func (e *Error) Is(target error) bool {
	if e == target {
		return true
	}
	if e != nil {
		if errors.Is(e.Err, target) {
			return true
		}
		if errors.Is(e.Cause, target) {
			return true
		}
	}
	return false
}

// Error method returns error message.
// This method is a implementation of error interface.
func (e *Error) Error() string {
	if e == nil {
		return nilAngleString
	}
	errMsg := e.Err.Error()
	var causeMsg string
	if e.Cause != nil {
		causeMsg = e.Cause.Error()
	}
	if len(causeMsg) == 0 {
		return errMsg
	}
	if len(errMsg) == 0 {
		return causeMsg
	}
	return strings.Join([]string{errMsg, causeMsg}, ": ")
}

// String method returns error message.
// This method is a implementation of fmt.Stringer interface.
func (e *Error) String() string {
	return e.Error()
}

// GoString method returns serialize string of Error.
// This method is a implementation of fmt.GoStringer interface.
func (e *Error) GoString() string {
	if e == nil {
		return nilAngleString
	}
	return fmt.Sprintf("%T{Err:%#v, Cause:%#v, Context:%#v}", e, e.Err, e.Cause, e.Context)
}

// MarshalJSON method returns serialize string of Error with JSON format.
// This method is implementation of json.Marshaler interface.
func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(e.EncodeJSON()), nil
}

// EncodeJSON method returns serialize string of Error with JSON format.
func (e *Error) EncodeJSON() string {
	if e == nil {
		return "null"
	}
	elms := []string{}
	elms = append(elms, strings.Join([]string{`"Type":`, strconv.Quote(reflect.TypeOf(e).String())}, ""))
	msgBuf := &bytes.Buffer{}
	json.HTMLEscape(msgBuf, bytes.Join([][]byte{[]byte(`"Err":`), []byte(EncodeJSON(e.Err))}, []byte{}))
	elms = append(elms, msgBuf.String())
	if len(e.Context) > 0 {
		if b, err := json.Marshal(e.Context); err == nil {
			elms = append(elms, string(bytes.Join([][]byte{[]byte(`"Context":`), b}, []byte{})))
		}
	}
	if e.Cause != nil && !reflect.ValueOf(e.Cause).IsZero() {
		elms = append(elms, strings.Join([]string{`"Cause":`, EncodeJSON(e.Cause)}, ""))
	}
	return strings.Join([]string{"{", strings.Join(elms, ","), "}"}, "")
}

// Format method returns formatted string of Error instance.
// This method is a implementation of fmt.Formatter interface.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('#'):
			_, _ = strings.NewReader(e.GoString()).WriteTo(s)
		case s.Flag('+'):
			_, _ = strings.NewReader(e.EncodeJSON()).WriteTo(s)
		default:
			_, _ = strings.NewReader(e.Error()).WriteTo(s)
		}
	case 's':
		_, _ = strings.NewReader(e.String()).WriteTo(s)
	default:
		fmt.Fprintf(s, `%%!%c(%s)`, verb, e.GoString())
	}
}

// Cause function finds cause error in target error instance.
//
// Deprecated: should not be used
func Cause(err error) error {
	for err != nil {
		unwraped := errors.Unwrap(err)
		if unwraped == nil {
			return err
		}
		err = unwraped
	}
	return err
}

// Unwraps function finds cause errors ([]error slice) in target error instance.
func Unwraps(err error) []error {
	if err == nil {
		return nil
	}
	if es, ok := err.(interface {
		Unwrap() []error
	}); ok {
		return es.Unwrap()
	}
	if e := errors.Unwrap(err); e != nil {
		return []error{e}
	}
	return nil
}

// caller returns caller info.
func caller(depth int) (string, string, int) {
	pc, src, line, ok := runtime.Caller(depth + 1)
	if !ok {
		return "", "", 0
	}
	return runtime.FuncForPC(pc).Name(), src, line
}

// EncodeJSON function dumps out error instance with JSON format.
func EncodeJSON(err error) string {
	if e, ok := err.(*Error); ok {
		return e.EncodeJSON()
	}
	if e, ok := err.(json.Marshaler); ok {
		b, ee := json.Marshal(e)
		if ee != nil {
			return encodeJSON(err)
		}
		return strings.TrimSpace(string(b))
	}
	return encodeJSON(err)
}

func encodeJSON(err error) string {
	if err == nil {
		return "null"
	}
	elms := []string{}
	elms = append(elms, strings.Join([]string{`"Type":`, strconv.Quote(reflect.TypeOf(err).String())}, ""))
	msgBuf := &bytes.Buffer{}
	json.HTMLEscape(msgBuf, bytes.Join([][]byte{[]byte(`"Msg":`), []byte(strconv.Quote(err.Error()))}, []byte{}))
	elms = append(elms, msgBuf.String())
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		unwraped := x.Unwrap()
		if err != nil {
			elms = append(elms, strings.Join([]string{`"Cause":`, EncodeJSON(unwraped)}, ""))
		}
	case interface{ Unwrap() []error }:
		unwraped := x.Unwrap()
		if len(unwraped) > 0 {
			causes := []string{}
			for _, c := range unwraped {
				causes = append(causes, EncodeJSON(c))
			}
			elms = append(elms, strings.Join([]string{`"Cause":[`, strings.Join(causes, ","), "]"}, ""))
		}
	}
	return strings.Join([]string{"{", strings.Join(elms, ","), "}"}, "")
}

// Is is conpatible with errors.Is.
func Is(err, target error) bool { return errors.Is(err, target) }

// As is conpatible with errors.As.
func As(err error, target interface{}) bool { return errors.As(err, target) }

// Unwrap is conpatible with errors.Unwrap.
func Unwrap(err error) error { return errors.Unwrap(err) }

/* Copyright 2019-2023 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
