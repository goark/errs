package zapobject

import (
	"errors"
	"testing"

	"github.com/goark/errs"
	"go.uber.org/zap/zapcore"
)

type addReflectedErrEncoder struct {
	*zapcore.MapObjectEncoder
}

func (e *addReflectedErrEncoder) AddReflected(_ string, _ interface{}) error {
	return errors.New("add reflected failed")
}

func TestMarshalLogObject_AddReflectedError(t *testing.T) {
	err := errs.New("wrapped", errs.WithContext("bad", func() {}))
	enc := &addReflectedErrEncoder{MapObjectEncoder: zapcore.NewMapObjectEncoder()}

	if got := New(err).MarshalLogObject(enc); got == nil {
		t.Fatal("MarshalLogObject() = nil, want non-nil error")
	}
}
