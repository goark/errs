//go:build run
// +build run

package main

import (
	"fmt"
	"sync"

	"github.com/goark/errs"
)

func generateMultiError() error {
	errlist := &errs.Errors{}
	var wg sync.WaitGroup
	for i := 1; i <= 2; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			errlist.Add(fmt.Errorf("error %d", i))
		}()
	}
	wg.Wait()
	return errlist.ErrorOrNil()
}

func main() {
	err := generateMultiError()
	fmt.Printf("%+v\n", err) // {"Type":"*errs.Errors","Errs":[{"Type":"*errors.errorString","Msg":"error 2"},{"Type":"*errors.errorString","Msg":"error 1"}]}
}
