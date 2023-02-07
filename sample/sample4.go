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
	return errs.New("error with multiple causes", errs.WithCause(errors.Join(os.ErrInvalid, io.EOF)))
}

func main() {
	err := generateMultiError()
	fmt.Printf("%+v\n", err)            // {"Type":"*errs.Error","Err":{"Type":"*errors.errorString","Msg":"error with multiple causes"},"Context":{"function":"main.generateMultiError"},"Cause":{"Type":"*errors.joinError","Msg":"invalid argument\nEOF","Cause":[{"Type":"*errors.errorString","Msg":"invalid argument"},{"Type":"*errors.errorString","Msg":"EOF"}]}}
	fmt.Println(errors.Is(err, io.EOF)) // true
}
