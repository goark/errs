package errs_test

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/goark/errs"
)

func ExampleNew() {
	err := errs.New(
		"wrapper error",
		errs.WithCause(os.ErrInvalid),
		errs.WithContext("foo", "bar"),
	)
	fmt.Printf("%+v", err)
	// Output:
	// {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"wrapper error"},"Context":{"foo":"bar","function":"github.com/goark/errs_test.ExampleNew"},"Cause":{"Type":"*errors.errorString","Msg":"invalid argument"}}
}

func ExampleError() {
	err := errs.Wrap(
		os.ErrInvalid,
		errs.WithContext("foo1", "bar1"),
	)
	_ = err.(*errs.Error).SetContext("foo2", "bar2")
	fmt.Printf("%+v", err)
	// Output:
	// {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"invalid argument"},"Context":{"foo1":"bar1","foo2":"bar2","function":"github.com/goark/errs_test.ExampleError"}}
}

func ExampleEncodeJSON() {
	_, err := os.Open("not-exist.txt")
	fmt.Printf("%v", errs.EncodeJSON(err))
	// Output:
	// {"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}
}

func ExampleJoin() {
	err := errs.Join(errors.New("error 1"), errors.New("error 2"))
	fmt.Println(err)
	errlist, ok := err.(*errs.Errors)
	if !ok {
		return
	}
	errlist.Add(errors.New("error 3"))
	fmt.Printf("%+v\n", errlist.ErrorOrNil())
	// Output:
	// error 1
	// error 2
	// {"Type":"*errs.Errors","Errs":[{"Type":"*errors.errorString","Msg":"error 1"},{"Type":"*errors.errorString","Msg":"error 2"},{"Type":"*errors.errorString","Msg":"error 3"}]}
}

func ExampleErrors() {
	errlist := &errs.Errors{}
	var wg sync.WaitGroup
	for i := 1; i <= 100000; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			errlist.Add(fmt.Errorf("error %d", i))
		}()
	}
	wg.Wait()
	fmt.Println("error ount =", len(errlist.Unwrap()))
	// Output:
	// error ount = 100000
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
