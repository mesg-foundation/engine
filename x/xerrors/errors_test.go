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
	"errors"
	"testing"
)

func TestAppend(t *testing.T) {
	var errs Errors
	errs = append(errs, errors.New("a"))
	errs = append(errs, errors.New("b"))
	if len(errs) != 2 {
		t.Fatalf("invalid errors count - got: %d, want: 2", len(errs))
	}
}

func TestError(t *testing.T) {
	var errs Errors
	if errs.Error() != "" {
		t.Fatalf("invalid error message - got: %q, want: %q", errs.Error(), "")
	}

	errs = append(errs, errors.New("a"))
	errs = append(errs, errors.New("b"))
	if errs.Error() != "a\nb" {
		t.Fatalf("invalid error message - got: %q, want: %q", errs.Error(), "a\nb")
	}
}
