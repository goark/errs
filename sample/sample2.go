//go:build run
// +build run

package main

import (
	"fmt"
	"os"

	"github.com/goark/errs"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errs.Wrap(
			err,
			errs.WithContext("path", path),
		)
	}
	defer file.Close()

	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		fmt.Printf("%v\n", err)             // open not-exist.txt: no such file or directory
		fmt.Printf("%#v\n", err)            // *errs.Error{Err:&fs.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}, Cause:<nil>, Context:map[string]interface {}{"function":"main.checkFileOpen", "path":"not-exist.txt"}}
		fmt.Printf("%+v\n", err)            // {"Type":"*errs.Error","Err":{"Type":"*fs.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}},"Context":{"function":"main.checkFileOpen","path":"not-exist.txt"}}
		fmt.Printf("%v\n", errs.Cause(err)) // no such file or directory
	}
}
