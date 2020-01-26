package errs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
)

var (
	nilErr        = New("")
	nilValueErr   = (*Error)(nil)
	errTest       = New("\"Error\" for test")
	wrapedErrTest = Wrap(errTest, "")
)

func TestNil(t *testing.T) {
	testCases := []struct {
		err     error
		typeStr string
		ptr     string
		msg     string
		detail  string
		json    string
		badStr  string
	}{
		{
			err:     nilValueErr,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "<nil>",
			detail:  `<nil>`,
			json:    `null`,
			badStr:  `%!d(<nil>)`,
		},
	}

	for _, tc := range testCases {
		str := fmt.Sprintf("%T", tc.err)
		if str != tc.typeStr {
			t.Errorf("Type of Wrap(\"%v\") is %v, want %v", tc.err, str, tc.typeStr)
		}
		str = fmt.Sprintf("%p", tc.err)
		if str != tc.ptr {
			t.Errorf("Pointer of Wrap(\"%v\") is %v, want %v", tc.err, str, tc.ptr)
		}
		str = fmt.Sprintf("%s", tc.err)
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		str = fmt.Sprintf("%v", tc.err)
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		str = fmt.Sprintf("%#v", tc.err)
		if str != tc.detail {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.detail)
		}
		str = fmt.Sprintf("%+v", tc.err)
		if str != tc.json {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
		}
		str = fmt.Sprintf("%d", tc.err)
		if str != tc.badStr {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.badStr)
		}
	}
}

func TestWrap(t *testing.T) {
	testCases := []struct {
		err     error
		typeStr string
		ptr     string
		msg     string
		detail  string
		json    string
		badStr  string
	}{
		{
			err:     nilErr,
			typeStr: "<nil>",
			ptr:     "0x0",
			msg:     "<nil>",
			detail:  "<nil>",
			json:    "<nil>",
			badStr:  `%!d(<nil>)`,
		},
		{
			err:     nilValueErr,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: <nil>",
			detail:  `*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:<nil>}`,
			json:    `{"Type":"*errs.Error","Msg":"wrapped message: \u003cnil\u003e","Context":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap","num":1}}`,
			badStr:  `%!d(*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:<nil>})`,
		},
		{
			err:     os.ErrInvalid,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: invalid argument",
			detail:  `*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:&errors.errorString{s:"invalid argument"}}`,
			json:    `{"Type":"*errs.Error","Msg":"wrapped message: invalid argument","Context":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap","num":1},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}`,
			badStr:  `%!d(*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:&errors.errorString{s:"invalid argument"}})`,
		},
		{
			err:     errTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:*errs.Error{Msg:"\"Error\" for test", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:<nil>}}`,
			json:    `{"Type":"*errs.Error","Msg":"wrapped message: \"Error\" for test","Context":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap","num":1},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Context":{"function":"github.com/spiegel-im-spiegel/errs.init"}}}`,
			badStr:  `%!d(*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:*errs.Error{Msg:"\"Error\" for test", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:<nil>}})`,
		},
		{
			err:     wrapedErrTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:*errs.Error{Msg:"", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:*errs.Error{Msg:"\"Error\" for test", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:<nil>}}}`,
			json:    `{"Type":"*errs.Error","Msg":"wrapped message: \"Error\" for test","Context":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs.TestWrap","num":1},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Context":{"function":"github.com/spiegel-im-spiegel/errs.init"},"Cause":{"Type":"*errs.Error","Msg":"\"Error\" for test","Context":{"function":"github.com/spiegel-im-spiegel/errs.init"}}}}`,
			badStr:  `%!d(*errs.Error{Msg:"wrapped message", Context:map[string]interface {}{"foo":"bar", "function":"github.com/spiegel-im-spiegel/errs.TestWrap", "num":1}, Cause:*errs.Error{Msg:"", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:*errs.Error{Msg:"\"Error\" for test", Context:map[string]interface {}{"function":"github.com/spiegel-im-spiegel/errs.init"}, Cause:<nil>}}})`,
		},
	}

	for _, tc := range testCases {
		err := Wrap(tc.err, "wrapped message", WithContext("foo", "bar"), WithContext("num", 1))
		str := fmt.Sprintf("%T", err)
		if str != tc.typeStr {
			t.Errorf("Type of Wrap(\"%v\") is %v, want %v", tc.err, str, tc.typeStr)
		}
		str = fmt.Sprintf("%p", err)
		if str == tc.ptr {
			t.Errorf("Pointer of Wrap(\"%v\") is %v, not want %v", tc.err, str, tc.ptr)
		} else {
			fmt.Println("Info:", str)
		}
		str = fmt.Sprintf("%v", err)
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
		}
		str = fmt.Sprintf("%#v", err)
		if str != tc.detail {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.detail)
		}
		str = fmt.Sprintf("%+v", err)
		if str != tc.json {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
		}
		str = fmt.Sprintf("%d", err)
		if str != tc.badStr {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.badStr)
		}
		if err != nil {
			b, e := json.Marshal(err)
			if e != nil {
				t.Errorf("json.Marshal(\"%v\") is %v, want <nil>", tc.err, e)
			} else if string(b) != tc.json {
				t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, string(b), tc.json)
			}
			str = EncodeJSON(err)
			if str != tc.json {
				t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.json)
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
		{err: os.ErrInvalid, cause: os.ErrInvalid},
		{err: Wrap(os.ErrInvalid, "wrapped error"), cause: os.ErrInvalid},
	}

	for _, tc := range testCases {
		c := Cause(tc.err)
		if c != tc.cause {
			t.Errorf("result Cause(\"%v\") is \"%v\", want %v", tc.err, c, tc.cause)
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
			t.Errorf("result if Is(\"%v\", \"%v\") is %v, want %v", tc.err, tc.target, ok, tc.res)
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

/* Copyright 2019,2020 Spiegel
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
