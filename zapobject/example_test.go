package zapobject_test

import (
	"os"

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

func Example() {
	logger := zap.NewExample()
	defer logger.Sync()

	if err := checkFileOpen("not-exist.txt"); err != nil {
		logger.Error("err", zap.Object("error", zapobject.New(err)))
	}
	// Output:
	// {"level":"error","msg":"err","error":{"type":"*errors.errorString","msg":"file open error","cause":{"type":"*fs.PathError","msg":"open not-exist.txt: no such file or directory","cause":{"type":"syscall.Errno","msg":"no such file or directory"}},"context":{"function":"github.com/goark/errs/zapobject_test.checkFileOpen","path":"not-exist.txt"}}}
}
