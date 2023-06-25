package zapobject

import (
	"fmt"
	"sort"

	"github.com/goark/errs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ErrObject struct {
	Err error
}

// New wrapped error instance in ErrObject.
func New(err error) ErrObject {
	return ErrObject{Err: err}
}

// MarshalLogObject method is object marshaler for go.uber.org/zap.
func (e ErrObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if e.Err == nil {
		return nil
	}
	var ee *errs.Error
	if errs.As(e.Err, &ee) {
		enc.AddString("type", fmt.Sprintf("%T", ee))
		enc.AddString("msg", ee.Error())
		if ee.Err != nil {
			if err := enc.AddObject("error", New(ee.Err)); err != nil {
				return err
			}
		}
		if ee.Cause != nil {
			if err := enc.AddObject("cause", New(ee.Cause)); err != nil {
				return err
			}
		}
		if len(ee.Context) > 0 {
			keys := make([]string, 0, len(ee.Context))
			for k := range ee.Context {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			enc.OpenNamespace("context")
			for _, k := range keys {
				_ = enc.AddReflected(k, ee.Context[k])
			}
		}
	} else {
		enc.AddString("type", fmt.Sprintf("%T", e.Err))
		enc.AddString("msg", e.Err.Error())
		if errList := errs.Unwraps(e.Err); len(errList) > 0 {
			if len(errList) == 1 {
				return enc.AddObject("cause", New(errList[0]))
			}
			elist := make([]ErrObject, 0, len(errList))
			for _, e := range errList {
				elist = append(elist, New(e))
			}
			zap.Objects("causes", elist).AddTo(enc)
		}
	}
	return nil
}

/* Copyright 2023 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
