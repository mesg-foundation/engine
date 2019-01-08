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

package dockertest

// NotFoundErr satisfies docker client's notFound interface.
// docker.IsErrNotFound(err) will return true with NotFoundErr.
type NotFoundErr struct{}

// NotFound indicates that this error is a not found error.
func (e NotFoundErr) NotFound() bool {
	return true
}

// Error returns the string representation of error.
func (e NotFoundErr) Error() string {
	return "not found"
}
