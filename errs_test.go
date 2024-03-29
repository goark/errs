package errs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"syscall"
	"testing"
)

type testError struct {
	Msg string
	Err error
}

func (t *testError) Error() string {
	return strings.Join([]string{t.Msg, t.Err.Error()}, ": ")
}
func (t *testError) Unwrap() error {
	return t.Err
}
func (t *testError) MarshalJSON() ([]byte, error) {
	if t == nil {
		return []byte("null"), nil
	}
	elms := []string{}
	elms = append(elms, fmt.Sprintf(`"Type":%q`, fmt.Sprintf("%T", t)))
	msgBuf := &bytes.Buffer{}
	json.HTMLEscape(msgBuf, []byte(fmt.Sprintf(`"Msg":%q`, t.Error())))
	elms = append(elms, msgBuf.String())
	if t.Err != nil && !reflect.ValueOf(t.Err).IsZero() {
		elms = append(elms, fmt.Sprintf(`"Err":%s`, EncodeJSON(t.Err)))
	}

	return []byte("{" + strings.Join(elms, ",") + "}"), nil
}

var (
	nilErr         = New("") // nil object
	nilValueErr    = (*Error)(nil)
	errTest        = New("\"Error\" for test")
	wrapedErrTest  = Wrap(errTest)
	wrapedErrTest2 = &testError{Msg: "test for testError", Err: wrapedErrTest}
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
			err:     nilErr,
			typeStr: "<nil>",
			ptr:     "%!p(<nil>)",
			msg:     "<nil>",
			detail:  `<nil>`,
			json:    `<nil>`,
			badStr:  `%!d(<nil>)`,
		},
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

func TestNewWithCause(t *testing.T) {
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
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
		{
			err:     nilValueErr,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: <nil>",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
		{
			err:     os.ErrInvalid,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: invalid argument",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errors.errorString{s:"invalid argument"}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errors.errorString{s:"invalid argument"}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
		{
			err:     errTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1},"Cause":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
		{
			err:     wrapedErrTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1},"Cause":{"Type":"*errs.Error","Err":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}},"Context":{"function":"github.com/goark/errs.init"}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
		{
			err:     wrapedErrTest2,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: test for testError: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errs.testError{Msg:"test for testError", Err:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestNewWithCause","num":1},"Cause":{"Type":"*errs.testError","Msg":"test for testError: \"Error\" for test","Err":{"Type":"*errs.Error","Err":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}},"Context":{"function":"github.com/goark/errs.init"}}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errs.testError{Msg:"test for testError", Err:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestNewWithCause", "num":1}})`,
		},
	}

	for _, tc := range testCases {
		err := New("wrapped message", WithCause(tc.err), WithContext("foo", "bar"), WithContext("num", 1))
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
		str = fmt.Sprintf("%s", err)
		if str != tc.msg {
			t.Errorf("Wrap(\"%v\") is %v, want %v", tc.err, str, tc.msg)
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

func TestWrapWithCause(t *testing.T) {
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
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
		{
			err:     nilValueErr,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: <nil>",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:<nil>, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
		{
			err:     os.ErrInvalid,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: invalid argument",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errors.errorString{s:"invalid argument"}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errors.errorString{s:"invalid argument"}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
		{
			err:     errTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1},"Cause":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
		{
			err:     wrapedErrTest,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1},"Cause":{"Type":"*errs.Error","Err":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}},"Context":{"function":"github.com/goark/errs.init"}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
		{
			err:     wrapedErrTest2,
			typeStr: "*errs.Error",
			ptr:     "0x0",
			msg:     "wrapped message: test for testError: \"Error\" for test",
			detail:  `*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errs.testError{Msg:"test for testError", Err:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}}`,
			json:    `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapped message"},"Context":{"foo":"bar","function":"github.com/goark/errs.TestWrapWithCause","num":1},"Cause":{"Type":"*errs.testError","Msg":"test for testError: \"Error\" for test","Err":{"Type":"*errs.Error","Err":{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"\"Error\" for test"},"Context":{"function":"github.com/goark/errs.init"}},"Context":{"function":"github.com/goark/errs.init"}}}}`,
			badStr:  `%!d(*errs.Error{Err:&errors.errorString{s:"wrapped message"}, Cause:&errs.testError{Msg:"test for testError", Err:*errs.Error{Err:*errs.Error{Err:&errors.errorString{s:"\"Error\" for test"}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}, Cause:<nil>, Context:map[string]interface {}{"function":"github.com/goark/errs.init"}}}, Context:map[string]interface {}{"foo":"bar", "function":"github.com/goark/errs.TestWrapWithCause", "num":1}})`,
		},
	}

	for _, tc := range testCases {
		err := Wrap(errors.New("wrapped message"), WithCause(tc.err), WithContext("foo", "bar"), WithContext("num", 1))
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

func TestMultiError(t *testing.T) {
	testCases := []struct {
		err  error
		want string
	}{
		{err: Wrap(os.ErrInvalid, WithCause(errors.Join(io.EOF, io.ErrUnexpectedEOF))), want: `{"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"invalid argument"},"Context":{"function":"github.com/goark/errs.TestMultiError"},"Cause":{"Type":"*errors.joinError","Msg":"EOF\nunexpected EOF","Cause":[{"Type":"*errors.errorString","Msg":"EOF"},{"Type":"*errors.errorString","Msg":"unexpected EOF"}]}}`},
		{err: Wrap(errors.Join(io.EOF, io.ErrUnexpectedEOF), WithCause(os.ErrInvalid)), want: `{"Type":"*errs.Error","Err":{"Type":"*errors.joinError","Msg":"EOF\nunexpected EOF","Cause":[{"Type":"*errors.errorString","Msg":"EOF"},{"Type":"*errors.errorString","Msg":"unexpected EOF"}]},"Context":{"function":"github.com/goark/errs.TestMultiError"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}`},
	}

	for _, tc := range testCases {
		if got := EncodeJSON(tc.err); got != tc.want {
			t.Errorf("result EncodeJSON(\"%v\") is %v, want %v", tc.err, got, tc.want)
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
		{err: New(""), res: true, target: nil},
		{err: Wrap(nil), res: true, target: nil},
		{err: nil, res: false, target: errTest},
		{err: errTest, res: false, target: nil},
		{err: errTest, res: true, target: errTest},
		{err: errTest, res: false, target: os.ErrInvalid},
		{err: New("wrapped error", WithCause(os.ErrInvalid)), res: true, target: os.ErrInvalid},
		{err: New("wrapped error", WithCause(os.ErrInvalid)), res: false, target: errTest},
		{err: New("wrapped error", WithCause(errTest)), res: true, target: errTest},
		{err: New("wrapped error", WithCause(errTest)), res: false, target: os.ErrInvalid},
		{err: New("wrapped error", WithCause(errors.Join(io.EOF, os.ErrInvalid))), res: true, target: os.ErrInvalid},
		{err: Wrap(errors.Join(io.EOF, os.ErrInvalid)), res: true, target: os.ErrInvalid},
	}

	for _, tc := range testCases {
		if ok := Is(tc.err, tc.target); ok != tc.res {
			t.Errorf("result Is(\"%v\", \"%v\") is %v, want %v", tc.err, tc.target, ok, tc.res)
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
		{err: New("wrapped error", WithCause(syscall.ENOENT)), res: true, cause: syscall.ENOENT},
		{err: New("wrapped error", WithCause(errors.Join(syscall.ENOENT, os.ErrInvalid))), res: true, cause: syscall.ENOENT},
		{err: Wrap(errors.Join(syscall.ENOENT, os.ErrInvalid)), res: true, cause: syscall.ENOENT},
	}

	for _, tc := range testCases {
		var cs syscall.Errno
		if ok := As(tc.err, &cs); ok != tc.res {
			t.Errorf("result if As(\"%v\") is %v, want %v", tc.err, ok, tc.res)
			if ok && cs != tc.cause {
				t.Errorf("As(\"%v\") = \"%v\", want \"%v\"", tc.err, cs, tc.cause)
			}
		}
	}
}

func TestUnwrap(t *testing.T) {
	testCases := []struct {
		err    error
		unwrap string
		cause  error
	}{
		{err: nil, unwrap: "", cause: nil},
		{err: syscall.ENOENT, unwrap: "", cause: nil},
		{err: New("wrapped error"), unwrap: "", cause: nil},
		{err: New("wrapped error", WithCause(syscall.ENOENT)), unwrap: "no such file or directory", cause: syscall.ENOENT},
		{err: Wrap(syscall.ENOENT), unwrap: "no such file or directory", cause: syscall.ENOENT},
	}

	for _, tc := range testCases {
		cs := Unwrap(tc.err)
		if cs != nil {
			str := cs.Error()
			if str != tc.unwrap {
				t.Errorf("Unwrap(\"%v\") = \"%v\", want \"%v\"", tc.err, str, tc.unwrap)
			}
		}
		if cs != tc.cause {
			t.Errorf("As(\"%v\") = \"%v\", want \"%v\"", tc.err, cs, tc.cause)
		}
	}
}

func TestUnwraps(t *testing.T) {
	testCases := []struct {
		err      error
		unwraped []error
	}{
		{err: nil, unwraped: nil},
		{err: os.ErrInvalid, unwraped: nil},
		{err: Wrap(os.ErrInvalid), unwraped: []error{os.ErrInvalid}},
		{err: errors.Join(syscall.ENOENT, os.ErrInvalid), unwraped: []error{syscall.ENOENT, os.ErrInvalid}},
	}

	for _, tc := range testCases {
		got := Unwraps(tc.err)
		if len(got) != len(tc.unwraped) {
			t.Errorf("Unwrap(\"%v\") = \"%v\", want \"%v\"", tc.err, got, tc.unwraped)
		}
		for i, e := range tc.unwraped {
			if e != got[i] {
				t.Errorf("Unwrap(\"%v\") = \"%v\", want \"%v\"", tc.err, got, tc.unwraped)
			}
		}
	}
}

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
