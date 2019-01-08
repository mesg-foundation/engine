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

// Package resolver is a wrapper for system resolver service to map tasks data.
package resolver

import (
	"fmt"

	"github.com/mesg-foundation/core/execution"
)

// Resolvers tasks.
const (
	AddPeersTask = "addPeers"
	ResolveTask  = "resolve"
)

// AddPeersInputs map add peer task inputs.
func AddPeersInputs(addresses []string) map[string]interface{} {
	// TODO: this hack is not something that we should do but
	// it's needed because *parameterValidator is not able to identify
	// string slices for now.
	var addressesInterface []interface{}
	for _, address := range addresses {
		addressesInterface = append(addressesInterface, address)
	}
	return map[string]interface{}{"addresses": addressesInterface}
}

// AddPeersOutputs map add peer task outputs.
func AddPeersOutputs(e *execution.Execution) error {
	switch e.OutputKey {
	case "success":
		return nil
	case "error":
		return fmt.Errorf("resolver: %s", e.OutputData["message"])
	}
	return fmt.Errorf("resolver: task add peers has unknown output %s", e.OutputKey)
}

// ResolveInputs map resolve task inputs.
func ResolveInputs(serviceID string) map[string]interface{} {
	return map[string]interface{}{"serviceID": serviceID}
}

// ResolveOutputs map resolve task outputs.
func ResolveOutputs(e *execution.Execution) (peerAddress string, err error) {
	switch e.OutputKey {
	case "found":
		return e.OutputData["address"].(string), nil
	case "notFound":
		return "", fmt.Errorf("resolver: peer address could not be found for %s service", e.OutputData["serviceID"])
	case "error":
		return "", fmt.Errorf("resolver: %s", e.OutputData["message"])
	}
	return "", fmt.Errorf("resolver: task resolve has unknown output %s", e.OutputKey)
}
