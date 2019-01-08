// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xerrors

import (
	"strings"
	"sync"
)

// Errors is an error for tracing multiple errors.
type Errors []error

// ErrorOrNil returns an error if there is more then 0 error, nil otherwise.
func (e Errors) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

func (e Errors) Error() string {
	var s []string
	for _, err := range e {
		if err != nil {
			s = append(s, err.Error())
		}
	}

	return strings.Join(s, "\n")
}

// SyncErrors is an error for tracing multiple errors safe to use in multiple goroutines.
type SyncErrors struct {
	mx sync.Mutex
	Errors
}

// Append appends given err.
func (e *SyncErrors) Append(err error) {
	e.mx.Lock()
	e.Errors = append(e.Errors, err)
	e.mx.Unlock()
}
