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

package service

// Parameter describes task input parameters, output parameters of a task
// output and input parameters of an event.
type Parameter struct {
	// Key is the key of parameter.
	Key string `hash:"name:1"`

	// Name is the name of parameter.
	Name string `hash:"name:2"`

	// Description is the description of parameter.
	Description string `hash:"name:3"`

	// Type is the data type of parameter.
	Type string `hash:"name:4"`

	// Optional indicates if parameter is optional.
	Optional bool `hash:"name:5"`

	// Repeated is to have an array of this parameter
	Repeated bool `hash:"name:6"`
}
