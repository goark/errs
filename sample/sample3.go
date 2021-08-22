//go:build run
// +build run

package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spiegel-im-spiegel/errs"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errs.Wrap(
			errors.New("file open error"),
			errs.WithCause(err),
			errs.WithContext("path", path),
		)
	}
	defer file.Close()

	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		fmt.Printf("%v\n", err)             // file open error: open not-exist.txt: no such file or directory
		fmt.Printf("%#v\n", err)            // *errs.Error{Err:&errors.errorString{s:"file open error"}, Cause:&fs.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}, Context:map[string]interface {}{"function":"main.checkFileOpen", "path":"not-exist.txt"}}
		fmt.Printf("%+v\n", err)            // {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"file open error"},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}
		fmt.Printf("%v\n", errs.Cause(err)) // no such file or directory
	}
}
