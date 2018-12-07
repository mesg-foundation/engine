package main

import (
	mesg "github.com/mesg-foundation/go-service"
)

const errorOutputKey = "error"

type errorOutputs struct {
	Message string `json:"message"`
}

func newOutputsError(err error) (string, mesg.Data) {
	return errorOutputKey, &errorOutputs{
		Message: err.Error(),
	}
}
