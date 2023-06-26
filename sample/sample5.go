//go:build run
// +build run

package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/goark/errs"
)

func generateMultiError() error {
	return errs.Join(os.ErrInvalid, io.EOF)
}

func main() {
	err := generateMultiError()
	fmt.Printf("%+v\n", err)            // {"Type":"*errs.Errors","Errs":[{"Type":"*errors.errorString","Msg":"invalid argument"},{"Type":"*errors.errorString","Msg":"EOF"}]}
	fmt.Println(errors.Is(err, io.EOF)) // true
}
