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

package xos

import "os"

// Remove removes all given named file or directory.
func Remove(names ...string) error {
	var err error
	for _, name := range names {
		if err1 := os.Remove(name); err == nil {
			err = err1
		}
	}
	return err
}

// RemoveAll removes all given path and any children it contains.
func RemoveAll(paths ...string) error {
	var err error
	for _, path := range paths {
		if err1 := os.RemoveAll(path); err == nil {
			err = err1
		}
	}
	return err
}
