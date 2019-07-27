package errs

import (
	errors "golang.org/x/xerrors"
)

//Is function is compatible with errors.Is function.
func Is(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}
	return errors.Is(err, target)
}

//As function is compatible with errors.As function.
func As(err error, target interface{}) bool {
	if err == nil {
		return false
	}
	return errors.As(err, target)
}

//Cause returns cause error in target error
func Cause(err error) error {
	for {
		unwrap := errors.Unwrap(err)
		if unwrap == nil {
			return err
		}
		err = unwrap
	}
}

/* Copyright 2019 Spiegel
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
