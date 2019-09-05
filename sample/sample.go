package main

import (
	"fmt"
	"os"

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
		fmt.Printf("%#v\n", err)            //&errs.Error{Msg:"file open error", Params:map[string]string{"function":"main.checkFileOpen", "path":"not-exist.txt"}, Cause:&os.PathError{Op:"open", Path:"not-exist.txt", Err:0x2}}
		fmt.Printf("%+v\n", err)            //{"Type":"*errs.Error","Msg":"file open error: open not-exist.txt: no such file or directory","Params":{"function":"main.checkFileOpen","path":"not-exist.txt"},"Cause":{"Type":"*os.PathError","Msg":"open not-exist.txt: no such file or directory","Cause":{"Type":"syscall.Errno","Msg":"no such file or directory"}}}
		fmt.Printf("%v\n", errs.Cause(err)) //no such file or directory
	}
}
