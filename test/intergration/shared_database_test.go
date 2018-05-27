package intergration_test

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/mesg-foundation/core/daemon"
// 	"github.com/mesg-foundation/core/database/services"
// 	"github.com/mesg-foundation/core/docker"
// 	"github.com/stvp/assert"

// 	"github.com/mesg-foundation/core/api/core"
// 	"github.com/mesg-foundation/core/config"
// 	"github.com/mesg-foundation/core/service"
// 	"github.com/spf13/viper"
// 	"google.golang.org/grpc"
// )

// func testForceAndWaitForDaemonToStart() (wait chan error) {
// 	start := time.Now()
// 	timeout := time.Minute
// 	wait = make(chan error, 1)
// 	go func() {
// 		for {
// 			taskErrors, err := docker.TasksError([]string{"daemon"})
// 			if err != nil {
// 				wait <- err
// 				return
// 			}
// 			if taskErrors != nil {
// 				fmt.Println("taskErrors", taskErrors)
// 			}
// 			isRunning, err := daemon.IsRunning()
// 			if err != nil {
// 				wait <- err
// 				return
// 			}
// 			fmt.Println("IsRunning", isRunning)
// 			_, err = daemon.Start()
// 			if err != nil {
// 				wait <- err
// 				return
// 			}
// 			status, err := daemon.ContainerStatus()
// 			if err != nil {
// 				wait <- err
// 				return
// 			}
// 			fmt.Println("status", status)
// 			if status == docker.RUNNING {
// 				close(wait)
// 				return
// 			}
// 			diff := time.Now().Sub(start)
// 			if diff.Nanoseconds() >= int64(timeout) {
// 				wait <- errors.New("Wait too long for the daemon to runs, timeout reached")
// 				return
// 			}
// 			time.Sleep(500 * time.Millisecond)
// 		}
// 	}()
// 	return
// }

// func TestSharedDatabase(t *testing.T) {
// 	err := <-testForceAndWaitForDaemonToStart()
// 	assert.Nil(t, err)
// 	time.Sleep(2 * time.Second)
// 	connection, err := grpc.Dial(viper.GetString(config.APIClientTarget), grpc.WithInsecure())
// 	assert.Nil(t, err)
// 	cli := core.NewCoreClient(connection)
// 	reply, err := cli.DeployService(context.Background(), &core.DeployServiceRequest{
// 		Service: &service.Service{
// 			Name: "TestSharedDatabase",
// 			Dependencies: map[string]*service.Dependency{
// 				"test": &service.Dependency{
// 					Image: "nginx",
// 				},
// 			},
// 		},
// 	})
// 	assert.Nil(t, err)
// 	defer services.Delete(reply.ServiceID)
// 	service, err := services.Get(reply.ServiceID)
// 	assert.Nil(t, err)
// 	assert.Equal(t, service.Name, "TestSharedDatabase")
// }
