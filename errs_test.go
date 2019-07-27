package errs

import (
	"fmt"
	"os"
	"testing"
)

type testError string

func (e testError) Error() string {
	return string(e)
}

var errTest = testError("Error for test")
var errTestAnother = testError("Error for test (another instance)")

func TestWrap(t *testing.T) {
	testCases := []struct {
		err error
		msg string
	}{
		{err: errTest, msg: "wrapped message: " + string(errTest)},
	}

	for _, tc := range testCases {
		err := Wrap(tc.err, "wrapped message")
		str := err.Error()
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		fmt.Printf("Info: %+v\n", err)

		err = Wrapf(tc.err, "%s", "wrapped message")
		str = err.Error()
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		fmt.Printf("Info: %+v\n", err)
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
		{err: Wrapf(nil, "%s", "no error"), res: true, cause: nil},
		{err: nil, res: false, cause: errTest},
		{err: errTest, res: false, cause: nil},
		{err: errTest, res: true, cause: errTest},
		{err: errTest, res: false, cause: errTestAnother},
		{err: Wrap(errTest, "wrapped error"), res: true, cause: errTest},
		{err: Wrapf(errTest, "%s", "wrapped error"), res: true, cause: errTest},
		{err: Wrap(errTest, "wrapped error"), res: false, cause: errTestAnother},
	}

	for _, tc := range testCases {
		if ok := Is(tc.err, tc.cause); ok != tc.res {
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
		if ok := As(tc.err, &cs); ok != tc.res {
			t.Errorf("result if As(\"%v\") is %v, want %v", tc.err, ok, tc.res)
			if ok && cs != tc.cause {
				t.Errorf("As(\"%v\") = \"%v\", want \"%v\"", tc.err, cs, tc.cause)
			}
		}
	}
}

func TestCause(t *testing.T) {
	testCases := []struct {
		err   error
		cause error
	}{
		{err: nil, cause: nil},
		{err: errTest, cause: errTest},
		{err: Wrap(errTest, "wrapping error"), cause: errTest},
	}

	for _, tc := range testCases {
		res := Cause(tc.err)
		if res != tc.cause {
			t.Errorf("Cause in \"%v\" == \"%v\", want \"%v\"", tc.err, res, tc.cause)
		}
	}
}

func ExampleWrap() {
	err := Wrap(os.ErrInvalid, "wrapped message")
	fmt.Println(err)
	fmt.Printf("errs.Cause(err): %v\n", Cause(err))
	// Output:
	// wrapped message: invalid argument
	// errs.Cause(err): invalid argument
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
