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

package api

import (
	"time"

	"github.com/mesg-foundation/core/execution"
	uuid "github.com/satori/go.uuid"
)

// ExecuteAndListen executes given task and listen for result.
func (a *API) ExecuteAndListen(serviceID, task string, inputs map[string]interface{}) (*execution.Execution, error) {
	tag := uuid.NewV4().String()
	result, err := a.ListenResult(serviceID, ListenResultTagFilters([]string{tag}))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	// XXX: sleep because listen stream may not be ready to stream the data
	// and execution will done before stream is ready. In that case the response
	// wlll never come TODO: investigate
	time.Sleep(1 * time.Second)

	if _, err := a.ExecuteTask(serviceID, task, inputs, []string{tag}); err != nil {
		return nil, err
	}
	return <-result.Executions, nil

}
