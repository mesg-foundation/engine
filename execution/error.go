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

package execution

import "fmt"

// StatusError is an error when the processing is done on en execution with the wrong status
type StatusError struct {
	ExpectedStatus Status
	ActualStatus   Status
}

// Error returns the string representation of error.
func (e StatusError) Error() string {
	return fmt.Sprintf("Execution status error: %q instead of %q", e.ActualStatus, e.ExpectedStatus)
}
