// Package errs implements functions to manipulate error instances.
package errs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"runtime"
	"strings"
)

const (
	nilAngleString = "<nil>"
)

//Error type is a implementation of error interface.
//This type is for wrapping cause error instance.
type Error struct {
	//Msg element is error message.
	Msg     string
	Cause   error
	Context map[string]interface{}
}

var _ error = (*Error)(nil)          //Error type is compatible with error interface
var _ fmt.Stringer = (*Error)(nil)   //Error type is compatible with fmt.Stringer interface
var _ fmt.GoStringer = (*Error)(nil) //Error type is compatible with fmt.GoStringer interface
var _ fmt.Formatter = (*Error)(nil)  //Error type is compatible with fmt.Formatter interface
var _ json.Marshaler = (*Error)(nil) //Error type is compatible with json.Marshaler interface

//ErrorContextFunc type is self-referential function type for New and Wrap functions. (functional options pattern)
type ErrorContextFunc func(*Error)

//New function returns an error instance with message and context informations.
func New(msg string, opts ...ErrorContextFunc) error {
	if len(msg) == 0 {
		return nil
	}
	return newError(nil, msg, 2, opts...)
}

//Wrap function returns a wrapping error instance with message and context informations.
func Wrap(err error, msg string, opts ...ErrorContextFunc) error {
	if err == nil {
		return nil
	}
	return newError(err, msg, 2, opts...)
}

//newError returns error instance. (internal)
func newError(err error, msg string, depth int, opts ...ErrorContextFunc) error {
	we := &Error{Msg: msg, Cause: err}
	//caller function name
	if fname, _, _ := caller(depth); len(fname) > 0 {
		we.SetContext("function", fname)
	}
	//other params
	for _, opt := range opts {
		opt(we)
	}
	return we
}

//WithContext function returns ErrorContextFunc function value.
//This function is used in New and Wrap functions that represents context (key/value) data.
func WithContext(name string, value interface{}) ErrorContextFunc {
	return func(e *Error) {
		e.SetContext(name, value)
	}
}

//SetContext method sets context information in Error instance
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

//Unwrap method returns cause error in Error instance.
//This method is used in errors.Unwrap function.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

//Is method reports whether any error in error's chain matches cause of target error.
//This method is used in errors.Is function.
func (e *Error) Is(target error) bool {
	if e == target {
		return true
	}
	cause := Cause(target)
	if cause != target && errors.Is(e, cause) {
		return true
	}
	return false
}

//Error method returns error message.
//This method is a implementation of error interface.
func (e *Error) Error() string {
	if e == nil {
		return nilAngleString
	}
	if e.Cause == nil {
		return e.Msg
	}
	if len(e.Msg) == 0 {
		return e.Cause.Error()
	}
	return fmt.Sprintf("%v: %v", e.Msg, e.Cause)
}

//String method returns error message.
//This method is a implementation of fmt.Stringer interface.
func (e *Error) String() string {
	return e.Error()
}

//GoString method returns serialize string of Error.
//This method is a implementation of fmt.GoStringer interface.
func (e *Error) GoString() string {
	if e == nil {
		return nilAngleString
	}
	return fmt.Sprintf("&errs.Error{Msg:%q, Context:%#v, Cause:%#v}", e.Msg, e.Context, e.Cause)
}

//MarshalJSON method returns serialize string of Error with JSON format.
//This method is implementation of json.Marshaler interface.
func (e *Error) MarshalJSON() ([]byte, error) {
	return []byte(e.JSON()), nil
}

//JSON method returns serialize string of Error with JSON format.
func (e *Error) JSON() string {
	if e == nil {
		return "null"
	}
	elms := []string{}
	elms = append(elms, fmt.Sprintf(`"Type":%q`, fmt.Sprintf("%T", e)))
	msgBuf := &bytes.Buffer{}
	json.HTMLEscape(msgBuf, []byte(fmt.Sprintf(`"Msg":%q`, e.Error())))
	elms = append(elms, msgBuf.String())
	if len(e.Context) > 0 {
		if b, err := json.Marshal(e.Context); err == nil {
			elms = append(elms, fmt.Sprintf(`"Context":%s`, string(b)))
		}
	}
	if e.Cause != nil && !reflect.ValueOf(e.Cause).IsZero() {
		elms = append(elms, fmt.Sprintf(`"Cause":%s`, EncodeJSON(e.Cause)))
	}
	return "{" + strings.Join(elms, ",") + "}"
}

//Format method returns formatted string of Error instance.
//This method is a implementation of fmt.Formatter interface.
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('#'):
			io.Copy(s, strings.NewReader(e.GoString()))
		case s.Flag('+'):
			io.Copy(s, strings.NewReader(e.JSON()))
		default:
			io.Copy(s, strings.NewReader(e.Error()))
		}
	case 's':
		io.Copy(s, strings.NewReader(e.String()))
	default:
		fmt.Fprintf(s, `%%!%c(%s)`, verb, e.GoString())
	}
}

//Cause function finds cause error in target error instance.
func Cause(err error) error {
	for {
		unwraped := errors.Unwrap(err)
		if unwraped == nil {
			return err
		}
		err = unwraped
	}
}

//EncodeJSON function dumps out error instance with JSON format.
func EncodeJSON(err error) string {
	switch e := err.(type) {
	case *Error:
		return e.JSON()
	default:
		return encodeJSON(err)
	}
}

//caller returns caller info.
func caller(depth int) (string, string, int) {
	pc, src, line, ok := runtime.Caller(depth + 1)
	if !ok {
		return "", "", 0
	}
	return runtime.FuncForPC(pc).Name(), src, line
}

func encodeJSON(err error) string {
	if err == nil {
		return "null"
	}
	elms := []string{}
	elms = append(elms, fmt.Sprintf(`"Type":%q`, fmt.Sprintf("%T", err)))
	msgBuf := &bytes.Buffer{}
	json.HTMLEscape(msgBuf, []byte(fmt.Sprintf(`"Msg":%q`, err.Error())))
	elms = append(elms, msgBuf.String())
	unwraped := errors.Unwrap(err)
	if unwraped != nil {
		cause := `{}`
		switch e := unwraped.(type) {
		case *Error:
			cause = e.JSON()
		default:
			cause = encodeJSON(unwraped)
		}
		elms = append(elms, fmt.Sprintf(`"Cause":%s`, cause))
	}
	return "{" + strings.Join(elms, ",") + "}"
}

/* Copyright 2019 Spiegel
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
