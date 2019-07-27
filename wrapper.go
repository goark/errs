package errs

import (
	"fmt"

	errors "golang.org/x/xerrors"
)

type wrapError struct {
	msg   string
	cause error
	frame errors.Frame
}

//Wrap returns wraped error instance
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &wrapError{msg: msg, cause: err, frame: errors.Caller(1)}
}

//Wrapf returns wraped error instance
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrapError{msg: fmt.Sprintf(format, args...), cause: err, frame: errors.Caller(1)}
}

//Error method for error interface
func (we *wrapError) Error() string {
	return fmt.Sprintf("%v: %v", we.msg, we.cause)
}

//Unwrap method for errors.Wrapper interface
func (e *wrapError) Unwrap() error {
	return e.cause
}

//Format method for fmt.Formatter interface
func (we *wrapError) Format(s fmt.State, v rune) {
	errors.FormatError(we, s, v)
}

//FormatError method for errors.Formatter interface
func (we *wrapError) FormatError(p errors.Printer) error {
	p.Print(we.msg)
	we.frame.Format(p)
	return we.cause
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
