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
	return errs.Wrap(errors.Join(os.ErrInvalid, io.EOF))
}

func main() {
	err := generateMultiError()
	fmt.Printf("%+v\n", err)            // {"Type":"*errs.Error","Err":{"Type":"*errors.joinError","Msg":"invalid argument\nEOF","Cause":[{"Type":"*errors.errorString","Msg":"invalid argument"},{"Type":"*errors.errorString","Msg":"EOF"}]},"Context":{"function":"main.generateMultiError"}}
	fmt.Println(errors.Is(err, io.EOF)) // true
}
