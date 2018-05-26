package intergration_test

import (
	"context"
	"fmt"
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

func TestSharedDatabse(t *testing.T) {
	daemon.Start()
	defer daemon.Stop()
	<-daemon.WaitForContainerToRun()
	time.Sleep(1 * time.Second)
	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
	cli := core.NewCoreClient(connection)
	reply, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
		Service: &service.Service{
			Name: "TestSharedDatabse",
			Dependencies: map[string]*service.Dependency{
				"test": &service.Dependency{
					Image: "nginx",
				},
			},
		},
	})
	fmt.Println(err)
	service, err := services.Get(reply.ServiceID)
	defer services.Delete(reply.ServiceID)
	assert.Nil(t, err)
	assert.Equal(t, service.Name, "TestSharedDatabse")
}
