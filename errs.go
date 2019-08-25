package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

//Error is wrapper error class
type Error struct {
	Msg    string
	Err    error
	Params map[string]string
}

//ErrorParamsOptFunc is self-referential function for functional options pattern
type ErrorParamsOptFunc func(*Error)

//Wrap returns wrapping error instance
func Wrap(err error, msg string, opts ...ErrorParamsOptFunc) error {
	if err == nil {
		return nil
	}
	we := &Error{Msg: msg, Err: err, Params: map[string]string{}}
	//caller function name
	if fname, _, _ := caller(); len(fname) > 0 {
		we.SetParam("function", fname)
	}
	//other params
	for _, opt := range opts {
		opt(we)
	}
	return we
}

//WithScheme returns function for setting scheme
func WithParam(name, value string) ErrorParamsOptFunc {
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
	return e.Err
}

//Is method for errors.Is function
func (e *Error) Is(target error) bool {
	if e == nil || target == nil {
		return e == target
	}
	if e == target {
		return true
	}
	cause := Cause(target)
	if cause != target {
		if errors.Is(e, cause) {
			return true
		}
		if e.Err != nil && errors.Is(Cause(e), cause) {
			return true
		}
	}
	return false
}

//Error returns message string of Error
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if len(e.Msg) == 0 {
		return e.Err.Error()
	}
	return fmt.Sprintf("%v: %v", e.Msg, e.Err)
}

//String returns message string of Error
func (e *Error) String() string {
	return e.Error()
}

//JSON returns string with JSON format
func (e *Error) JSON() string {
	msg := strconv.Quote(e.Msg)
	parms := ""
	if len(e.Params) > 0 {
		if b, err := json.Marshal(e.Params); err == nil {
			parms = string(b)
		}
	}
	cause := fmt.Sprintf(`{"Msg":%s}`, strconv.Quote(e.Err.Error()))
	var ee *Error
	if errors.As(e.Err, &ee) {
		cause = ee.JSON()
	}
	if len(parms) == 0 {
		return fmt.Sprintf(`{"Msg":%s,"Cause":%s}`, msg, cause)
	}
	return fmt.Sprintf(`{"Msg":%s,"Cause":%s,"Params":%s}`, msg, cause, parms)
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

//caller returns caller info.
func caller() (string, string, int) {
	pc, src, line, ok := runtime.Caller(2)
	if !ok {
		return "", "", 0
	}
	return runtime.FuncForPC(pc).Name(), src, line
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
