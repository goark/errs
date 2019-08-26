package errs

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
)

var nilErr = New("")
var errTest = New("\"Error\" for test")
var wrapedErrTest = Wrap(errTest, "")

func TestWrap(t *testing.T) {
	testCases := []struct {
		err  error
		msg  string
		json string
	}{
		{err: nilErr, msg: "<nil>", json: "<nil>"},
		{err: os.ErrInvalid, msg: "wrapped message: invalid argument", json: `{"Type":"*errs.Error","Msg":"wrapped message: invalid argument","Params":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}`},
		{err: errTest, msg: "wrapped message: \"Error\" for test", json: `{"Type":"*errs.Error","Msg":"wrapped message: \"Error\" for test","Params":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap"},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Params":{"function":"github.com/spiegel-im-spiegel/errs.init"}}}`},
		{err: wrapedErrTest, msg: "wrapped message: \"Error\" for test", json: `{"Type":"*errs.Error","Msg":"wrapped message: \"Error\" for test","Params":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap"},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Params":{"function":"github.com/spiegel-im-spiegel/errs.init"},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Params":{"function":"github.com/spiegel-im-spiegel/errs.init"}}}}`},
	}

	for _, tc := range testCases {
		err := Wrap(tc.err, "wrapped message", WithParam("foo", "bar"))
		str := fmt.Sprintf("%v", err)
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		str = fmt.Sprintf("%#v", err)
		if str != tc.json {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
		}
		str = fmt.Sprintf("%+v", err)
		if str != tc.json {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
		}
		if err != nil {
			str = EncodeJSON(err)
			if str != tc.json {
				t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
			}
		}
	}
}

func TestIs(t *testing.T) {
	testCases := []struct {
		err    error
		res    bool
		target error
	}{
		{err: nil, res: true, target: nil},
		{err: New("error"), res: false, target: nil},
		{err: Wrap(nil, ""), res: true, target: nil},
		{err: nil, res: false, target: errTest},
		{err: errTest, res: false, target: nil},
		{err: errTest, res: true, target: errTest},
		{err: errTest, res: false, target: os.ErrInvalid},
		{err: Wrap(os.ErrInvalid, "wrapped error"), res: true, target: os.ErrInvalid},
		{err: Wrap(os.ErrInvalid, "wrapped error"), res: false, target: errTest},
		{err: Wrap(errTest, "wrapped error"), res: true, target: errTest},
		{err: Wrap(errTest, "wrapped error"), res: true, target: wrapedErrTest},
		{err: Wrap(errTest, "wrapped error"), res: false, target: os.ErrInvalid},
	}

	for _, tc := range testCases {
		if ok := errors.Is(tc.err, tc.target); ok != tc.res {
			t.Errorf("result if Is(\"%v\") is %v, want %v", tc.err, ok, tc.res)
		}
	}
}

func TestAs(t *testing.T) {
	testCases := []struct {
		err   error
		res   bool
		cause error
	}{
		{err: nil, res: false, cause: nil},
		{err: Wrap(syscall.ENOENT, "wrapping error"), res: true, cause: syscall.ENOENT},
	}

	for _, tc := range testCases {
		var cs syscall.Errno
		if ok := errors.As(tc.err, &cs); ok != tc.res {
			t.Errorf("result if As(\"%v\") is %v, want %v", tc.err, ok, tc.res)
			if ok && cs != tc.cause {
				t.Errorf("As(\"%v\") = \"%v\", want \"%v\"", tc.err, cs, tc.cause)
			}
		}
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
