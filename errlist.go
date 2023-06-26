package errs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

// Errors is multiple error instance.
type Errors struct {
	mu   sync.RWMutex
	errs []error
}

// Join function returns Errors instance.
func Join(errlist ...error) error {
	if len(errlist) == 0 {
		return nil
	}
	ct := 0
	for _, err := range errlist {
		if err != nil {
			ct++
		}
	}
	if ct == 0 {
		return nil
	}
	es := &Errors{errs: make([]error, 0, ct)}
	for _, err := range errlist {
		if err != nil {
			es.errs = append(es.errs, err)
		}
	}
	return es
}

// Add method adds errors to Errors.
func (es *Errors) Add(errlist ...error) {
	if es == nil {
		return
	}
	es.mu.Lock()
	defer es.mu.Unlock()
	for _, err := range errlist {
		if err != nil {
			es.errs = append(es.errs, err)
		}
	}
}

// ErrorOrNil method returns this as a error type.
func (es *Errors) ErrorOrNil() error {
	if es == nil {
		return nil
	}
	es.mu.RLock()
	defer es.mu.RUnlock()
	if len(es.errs) == 0 {
		return nil
	}
	return es
}

// Error method returns error message.
// This method is a implementation of error interface.
func (es *Errors) Error() string {
	if es == nil {
		return nilAngleString
	}
	es.mu.RLock()
	defer es.mu.RUnlock()
	if len(es.errs) == 0 {
		return nilAngleString
	}
	var b []byte
	for i, err := range es.errs {
		if i > 0 {
			b = append(b, '\n')
		}
		b = append(b, err.Error()...)
	}
	return string(b)
}

// String method returns error message.
// This method is a implementation of fmt.Stringer interface.
func (es *Errors) String() string {
	return es.Error()
}

// GoString method returns serialize string of Errors.
// This method is a implementation of fmt.GoStringer interface.
func (es *Errors) GoString() string {
	if es == nil {
		return nilAngleString
	}
	es.mu.RLock()
	defer es.mu.RUnlock()
	if len(es.errs) == 0 {
		return nilAngleString
	}
	return fmt.Sprintf("%T{Errs:%#v}", es, es.errs)
}

// MarshalJSON method returns serialize string of Errors with JSON format.
// This method is implementation of json.Marshaler interface.
func (es *Errors) MarshalJSON() ([]byte, error) {
	return []byte(es.EncodeJSON()), nil
}

// Format method returns formatted string of Errors instance.
// This method is a implementation of fmt.Formatter interface.
func (es *Errors) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('#'):
			_, _ = strings.NewReader(es.GoString()).WriteTo(s)
		case s.Flag('+'):
			_, _ = strings.NewReader(es.EncodeJSON()).WriteTo(s)
		default:
			_, _ = strings.NewReader(es.Error()).WriteTo(s)
		}
	case 's':
		_, _ = strings.NewReader(es.String()).WriteTo(s)
	default:
		fmt.Fprintf(s, `%%!%c(%s)`, verb, es.GoString())
	}
}

// EncodeJSON method returns serialize string of Errors with JSON format.
func (es *Errors) EncodeJSON() string {
	if es == nil {
		return "null"
	}
	es.mu.RLock()
	defer es.mu.RUnlock()
	elms := []string{}
	elms = append(elms, strings.Join([]string{`"Type":`, strconv.Quote(reflect.TypeOf(es).String())}, ""))
	if len(es.errs) > 0 {
		elms2 := []string{}
		for _, err := range es.errs {
			msgBuf := &bytes.Buffer{}
			json.HTMLEscape(msgBuf, []byte(EncodeJSON(err)))
			elms2 = append(elms2, msgBuf.String())
		}
		elms = append(elms, strings.Join([]string{`"Errs":[`, strings.Join(elms2, ","), "]"}, ""))
	}
	return strings.Join([]string{"{", strings.Join(elms, ","), "}"}, "")
}

// Unwrap method returns error list in Errors instance.
// This method is used in errors.Unwrap function.
func (es *Errors) Unwrap() []error {
	if es == nil {
		return nil
	}
	es.mu.RLock()
	defer es.mu.RUnlock()
	if len(es.errs) == 0 {
		return nil
	}
	cpy := make([]error, len(es.errs), cap(es.errs))
	copy(cpy, es.errs)
	return cpy
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
