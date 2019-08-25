package main

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/spiegel-im-spiegel/errs"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errs.Wrap(
			err,
			"file open error",
			errs.WithParam("path", path),
		)
	}
	defer file.Close()

	return nil
}

func main() {
	if err := checkFileOpen("not-exist.txt"); err != nil {
		fmt.Printf("%v\n", err)             //file open error: open not-exist.txt: no such file or directory
		fmt.Printf("%#v\n", err)            //{"Msg":"file open error","Cause":{"Msg":"open not-exist.txt: no such file or directory"},"Params":{"function":"main.checkFileOpen","path":"not-exist.txt"}}
		fmt.Printf("%v\n", errs.Cause(err)) //no such file or directory
		var errno syscall.Errno
		if errors.As(err, &errno) {
			fmt.Printf("%[1]v (%[1]d)\n", errno) //no such file or directory (2)
		}
	}
}
