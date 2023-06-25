package zapobject_test

import (
	"fmt"
	"os"
	"sync"

	"github.com/goark/errs"
	"github.com/goark/errs/zapobject"
	"go.uber.org/zap"
)

func checkFileOpen(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return errs.New(
			"file open error",
			errs.WithCause(err),
			errs.WithContext("path", path),
		)
	}
	defer file.Close()

	return nil
}

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

func Example() {
	logger := zap.NewExample()
	defer logger.Sync()

	if err := checkFileOpen("not-exist.txt"); err != nil {
		logger.Error("err", zap.Object("error", zapobject.New(err)))
	}
	if err := generateMultiError(); err != nil {
		logger.Error("err", zap.Object("error", zapobject.New(err)))
	}
	// Output:
	// {"level":"error","msg":"err","error":{"type":"*errs.Error","msg":"file open error: open not-exist.txt: no such file or directory","error":{"type":"*errors.errorString","msg":"file open error"},"cause":{"type":"*fs.PathError","msg":"open not-exist.txt: no such file or directory","cause":{"type":"syscall.Errno","msg":"no such file or directory"}},"context":{"function":"github.com/goark/errs/zapobject_test.checkFileOpen","path":"not-exist.txt"}}}
	// {"level":"error","msg":"err","error":{"type":"*errs.Errors","msg":"error 2\nerror 1","causes":[{"type":"*errors.errorString","msg":"error 2"},{"type":"*errors.errorString","msg":"error 1"}]}}
}
