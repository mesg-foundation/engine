// +build integration

package service

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/stretchr/testify/require"
)

func TestServiceNotExists(t *testing.T) {
	server, closer := newServer(t)
	defer closer()

	_, err := server.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     "TestServiceNotExists",
		EventKey:  "test",
		EventData: "{}",
	})
	require.Error(t, err)
}

func TestSubmitWithInvalidID(t *testing.T) {
	var (
		outputData     = "{}"
		executionHash  = "1"
		server, closer = newServer(t)
	)
	defer closer()

	_, err := server.SubmitResult(context.Background(), &serviceapi.SubmitResultRequest{
		ExecutionHash: executionHash,
		Result: &serviceapi.SubmitResultRequest_OutputData{
			OutputData: outputData,
		},
	})
	require.Error(t, err)
}
