package main

import (
	"context"
	"testing"

	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/hash/serializer"
	"github.com/mesg-foundation/engine/instance"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/stretchr/testify/require"
)

var testInstanceHash hash.Hash

func testInstance(t *testing.T) {
	t.Run("get", func(t *testing.T) {
		t.Run("grpc", func(t *testing.T) {
			resp, err := client.InstanceClient.Get(context.Background(), &pb.GetInstanceRequest{Hash: testInstanceHash})
			require.NoError(t, err)
			require.Equal(t, testInstanceHash, resp.Hash)
			require.Equal(t, testServiceHash, resp.ServiceHash)
			require.Equal(t, hash.Dump(serializer.StringSlice([]string{"BAR=3", "FOO=1", "REQUIRED=4"})), resp.EnvHash)
		})
		t.Run("lcd", func(t *testing.T) {
			var inst *instance.Instance
			lcdGet(t, "instance/get/"+testInstanceHash.String(), &inst)
			require.Equal(t, testInstanceHash, inst.Hash)
			require.Equal(t, testServiceHash, inst.ServiceHash)
			require.Equal(t, hash.Dump(serializer.StringSlice([]string{"BAR=3", "FOO=1", "REQUIRED=4"})), inst.EnvHash)
		})
	})

	t.Run("list", func(t *testing.T) {
		t.Run("with nil filter", func(t *testing.T) {
			t.Run("grpc", func(t *testing.T) {
				resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{})
				require.NoError(t, err)
				require.Len(t, resp.Instances, 1)
			})
			t.Run("lcd", func(t *testing.T) {
				insts := make([]*instance.Instance, 0)
				lcdGet(t, "instance/list", &insts)
				require.Len(t, insts, 1)
			})
		})
		t.Run("do not match service", func(t *testing.T) {
			resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{
				Filter: &pb.ListInstanceRequest_Filter{
					ServiceHash: hash.Int(1),
				},
			})
			require.NoError(t, err)
			require.Len(t, resp.Instances, 0)
		})
		t.Run("match service", func(t *testing.T) {
			resp, err := client.InstanceClient.List(context.Background(), &pb.ListInstanceRequest{
				Filter: &pb.ListInstanceRequest_Filter{
					ServiceHash: testServiceHash,
				},
			})
			require.NoError(t, err)
			require.Len(t, resp.Instances, 1)
			require.Equal(t, testServiceHash, resp.Instances[0].ServiceHash)
			require.Equal(t, testInstanceHash, resp.Instances[0].Hash)
		})
	})
}
