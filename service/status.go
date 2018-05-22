package service

// import (
// 	"context"
// 	"strings"

// 	"github.com/docker/docker/api/types/swarm"
// 	godocker "github.com/fsouza/go-dockerclient"
// 	"github.com/mesg-foundation/core/docker"
// )

// func serviceStatus(service *Service) (status StatusType) {
// 	status = STOPPED
// 	allRunning := true
// 	for name, dependency := range service.GetDependencies() {
// 		if dependency.IsRunning(service.Namespace(), name) {
// 			status = RUNNING
// 		} else {
// 			allRunning = false
// 		}
// 	}
// 	if status == RUNNING && !allRunning {
// 		status = PARTIAL
// 	}
// 	return
// }

// // IsRunning returns true if the service is running, false otherwise
// func (service *Service) IsRunning() (running bool) {
// 	running = serviceStatus(service) == RUNNING
// 	return
// }

// // IsPartiallyRunning returns true if the service is running, false otherwise
// func (service *Service) IsPartiallyRunning() (running bool) {
// 	running = serviceStatus(service) == PARTIAL
// 	return
// }

// // IsStopped returns true if the service is stopped, false otherwise
// func (service *Service) IsStopped() (running bool) {
// 	running = serviceStatus(service) == STOPPED
// 	return
// }

// // IsRunning returns true if the dependency is running, false otherwise
// func (dependency *Dependency) IsRunning(namespace string, name string) (running bool) {
// 	running = dependencyStatus(namespace, name) == RUNNING
// 	return
// }

// // IsStopped returns true if the dependency is stopped, false otherwise
// func (dependency *Dependency) IsStopped(namespace string, name string) (running bool) {
// 	running = dependencyStatus(namespace, name) == STOPPED
// 	return
// }

// // List all the running services
// func List() (res []string, err error) {
// 	client, err := docker.Client()
// 	services, err := client.ListServices(godocker.ListServicesOptions{
// 		Context: context.Background(),
// 	})
// 	mapRes := make(map[string]uint)
// 	for _, service := range services {
// 		serviceName := service.Spec.Annotations.Labels["mesg.service"]
// 		mapRes[serviceName]++
// 	}
// 	res = make([]string, 0, len(mapRes))
// 	for k := range mapRes {
// 		res = append(res, k)
// 	}
// 	return
// }
