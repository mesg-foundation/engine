package intergration_test

import (
	"context"
	"testing"
	"time"

	"github.com/mesg-foundation/core/daemon"
	"github.com/mesg-foundation/core/database/services"
	"github.com/stvp/assert"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func TestSharedDatabase(t *testing.T) {
	daemon.Start()
	err := <-daemon.WaitForContainerToRun()
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	assert.Nil(t, err)
	cli := core.NewCoreClient(connection)
	reply, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
		Service: &service.Service{
			Name: "TestSharedDatabase",
			Dependencies: map[string]*service.Dependency{
				"test": &service.Dependency{
					Image: "nginx",
				},
			},
		},
	})
	assert.Nil(t, err)
	defer services.Delete(reply.ServiceID)
	service, err := services.Get(reply.ServiceID)
	assert.Nil(t, err)
	assert.Equal(t, service.Name, "TestSharedDatabase")
}
