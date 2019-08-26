package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

//Error is wrapper error class
type Error struct {
	Msg    string
	Cause  error
	Params map[string]string
}

//ErrorParamsFunc is self-referential function for functional options pattern
type ErrorParamsFunc func(*Error)

//New returns new Error instance
func New(msg string, opts ...ErrorParamsFunc) error {
	if len(msg) == 0 {
		return nil
	}
	return newError(nil, msg, 2, opts...)
}

//Wrap returns wrapping error instance
func Wrap(err error, msg string, opts ...ErrorParamsFunc) error {
	if err == nil {
		return nil
	}
	return newError(err, msg, 2, opts...)
}

//newError returns error instance (internal)
func newError(err error, msg string, depth int, opts ...ErrorParamsFunc) error {
	we := &Error{Msg: msg, Cause: err, Params: map[string]string{}}
	//caller function name
	if fname, _, _ := caller(depth); len(fname) > 0 {
		we.SetParam("function", fname)
	}
	//other params
	for _, opt := range opts {
		opt(we)
	}
	return we
}

//WithScheme returns function for setting scheme
func WithParam(name, value string) ErrorParamsFunc {
	return func(e *Error) {
		e.SetParam(name, value)
	}
}

//SetParam sets parameter to Error instance
func (e *Error) SetParam(name, value string) {
	if e == nil {
		return
	}
	if e.Params == nil {
		e.Params = map[string]string{}
	}
	e.Params[name] = value
	return
}

//Unwrap method for errors.Unwrap function
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

//Is method for errors.Is function
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

//Error returns message string of Error
func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Msg
	}
	if len(e.Msg) == 0 {
		return e.Cause.Error()
	}
	return fmt.Sprintf("%v: %v", e.Msg, e.Cause)
}

//String returns message string of Error
func (e *Error) String() string {
	return e.Error()
}

//JSON returns string with JSON format
func (e *Error) JSON() string {
	elms := []string{}
	elms = append(elms, fmt.Sprintf(`"Type":%s`, strconv.Quote(fmt.Sprintf("%T", e))))
	elms = append(elms, fmt.Sprintf(`"Msg":%s`, strconv.Quote(e.Error())))
	if len(e.Params) > 0 {
		if b, err := json.Marshal(e.Params); err == nil {
			elms = append(elms, fmt.Sprintf(`"Params":%s`, string(b)))
		}
	}
	if e.Cause != nil {
		elms = append(elms, fmt.Sprintf(`"Cause":%s`, EncodeJSON(e.Cause)))
	}
	return "{" + strings.Join(elms, ",") + "}"
}

//Format formats Error instance
func (e *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('#'), s.Flag('+'):
			s.Write([]byte(e.JSON()))
		default:
			s.Write([]byte(e.Error()))
		}
	case 's':
		s.Write([]byte(e.String()))
	}
}

//Cause returns cause error instance
func Cause(err error) error {
	for {
		unwraped := errors.Unwrap(err)
		if unwraped == nil {
			return err
		}
		err = unwraped
	}
}

//EncodeJSON dumps out error instance with JSON format
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
	elms = append(elms, fmt.Sprintf(`"Type":%s`, strconv.Quote(fmt.Sprintf("%T", err))))
	elms = append(elms, fmt.Sprintf(`"Msg":%s`, strconv.Quote(err.Error())))
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
