package errs

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

type testError string

func (e testError) Error() string {
	return string(e)
}

var errTest = testError("\"Error\" for test")
var wrapedErrTest = Wrap(errTest, "")

func TestWrap(t *testing.T) {
	testCases := []struct {
		err  error
		msg  string
		json string
	}{
		{err: nil, msg: "<nil>", json: "<nil>"},
		{err: errTest, msg: "wrapped message: \"Error\" for test", json: `{"Msg":"wrapped message","Cause":{"Msg":"\"Error\" for test"},"Params":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap"}}`},
		{err: wrapedErrTest, msg: "wrapped message: \"Error\" for test", json: `{"Msg":"wrapped message","Cause":{"Msg":"","Cause":{"Msg":"\"Error\" for test"},"Params":{"function":"github.com/spiegel-im-spiegel/errs.init"}},"Params":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap"}}`},
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

	}
}

func TestIs(t *testing.T) {
	testCases := []struct {
		err   error
		res   bool
		cause error
	}{
		{err: nil, res: true, cause: nil},
		{err: Wrap(nil, ""), res: true, cause: nil},
		{err: nil, res: false, cause: errTest},
		{err: errTest, res: false, cause: nil},
		{err: errTest, res: true, cause: errTest},
		{err: errTest, res: false, cause: os.ErrInvalid},
		{err: Wrap(errTest, "wrapped error"), res: true, cause: errTest},
		{err: Wrap(errTest, "wrapped error"), res: true, cause: wrapedErrTest},
		{err: Wrap(errTest, "wrapped error"), res: false, cause: os.ErrInvalid},
	}

	for _, tc := range testCases {
		if ok := errors.Is(tc.err, tc.cause); ok != tc.res {
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
		{err: errTest, res: true, cause: errTest},
		{err: Wrap(errTest, "wrapping error"), res: true, cause: errTest},
	}

	for _, tc := range testCases {
		var cs testError
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
