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

package main

import (
	mesg "github.com/mesg-foundation/go-service"
)

const (
	addPeersSuccessOutputsKey = "success"
)

type addPeersInputs struct {
	Addresses []string `json:"addresses"`
}

type addPeersSuccessOutputs struct {
	Addresses []string `json:"addresses"`
}

func addPeersHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs addPeersInputs
	if err := execution.Data(&inputs); err != nil {
		return newOutputsError(err)
	}
	addPeers(inputs.Addresses...)
	return addPeersSuccessOutputsKey, &addPeersSuccessOutputs{
		Addresses: inputs.Addresses,
	}
}
