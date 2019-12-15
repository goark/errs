package errs_test

import (
	"fmt"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

// func checkFileOpen(path string) error {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return Wrap(
// 			err,
// 			"file open error",
// 			WithContext("path", path),
// 		)
// 	}
// 	defer file.Close()
//
// 	return nil
// }
//
// func ExampleErrs() {
// 	var lastErr error
// 	if err := checkFileOpen("not-exist.txt"); err != nil {
// 		lastErr = fmt.Errorf("detect error!: %w", err)
// 	}
// 	fmt.Printf("%+v\n", Wrap(lastErr, ""))
// 	// Output:
// 	// {"Type":"*errs.Error","Msg":"detect error!: file open error: open not-exist.txt: no such file or directory","Context":{"function":"github.com/spiegel-im-spiegel/errs_test.ExampleErrs"},"Cause":{"Type":"*fmt.wrapError","Msg":"detect error!: file open error: open not-exist.txt: no such file or directory","Cause":{"Type":"*errs.Error","Msg":"file open error: open not-exist.txt: no such file or directory","Context":{"function":"github.com/spiegel-im-spiegel/errs_test.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*os.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}}}
// }

func ExampleNew() {
	fmt.Printf("%+v", errs.New("custom error"))
	// Output:
	// {"Type":"*errs.Error","Msg":"custom error","Context":{"function":"github.com/spiegel-im-spiegel/errs_test.ExampleNew"}}
}

func ExampleWrap() {
	fmt.Printf("%+v", errs.Wrap(os.ErrInvalid, "wrapper error"))
	// Output:
	// {"Type":"*errs.Error","Msg":"wrapper error: invalid argument","Context":{"function":"github.com/spiegel-im-spiegel/errs_test.ExampleWrap"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}
}

func ExampleWithContext() {
	err := errs.Wrap(
		os.ErrInvalid,
		"wrapper error",
		errs.WithContext("foo", "bar"),
	)
	fmt.Printf("%+v", err)
	// Output:
	// {"Type":"*errs.Error","Msg":"wrapper error: invalid argument","Context":{"foo":"bar","function":"github.com/spiegel-im-spiegel/errs_test.ExampleWithContext"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}
}

func ExampleError() {
	err := errs.Wrap(
		os.ErrInvalid,
		"wrapper error",
		errs.WithContext("foo1", "bar1"),
	)
	_ = err.(*errs.Error).SetContext("foo2", "bar2")
	fmt.Printf("%+v", err)
	// Output:
	// {"Type":"*errs.Error","Msg":"wrapper error: invalid argument","Context":{"foo1":"bar1","foo2":"bar2","function":"github.com/spiegel-im-spiegel/errs_test.ExampleError"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}
}

func ExampleCause() {
	_, err := os.Open("not-exist.txt")
	fmt.Printf("%v", errs.Cause(err))
	// Output:
	// no such file or directory
}

func ExampleEncodeJSON() {
	_, err := os.Open("not-exist.txt")
	fmt.Printf("%v", errs.EncodeJSON(err))
	// Output:
	// {"Type":"*os.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}
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
