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
