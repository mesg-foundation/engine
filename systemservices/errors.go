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

package systemservices

import (
	"fmt"
)

// SystemServiceNotFoundError is returned when an expected
// system service is not found.
type SystemServiceNotFoundError struct {
	Name string
}

func (e *SystemServiceNotFoundError) Error() string {
	return fmt.Sprintf("System service %q not found", e.Name)
}
