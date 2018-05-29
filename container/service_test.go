package container

import (
	"bytes"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/stvp/assert"
)

func startTestService(name []string) (dockerService *swarm.Service, err error) {
	namespace := Namespace(name)
	return StartService(docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: namespace,
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace,
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: "nginx",
					Labels: map[string]string{
						"com.docker.stack.namespace": namespace,
					},
				},
			},
		},
	})
}

func TestStartService(t *testing.T) {
	namespace := []string{"TestStartService"}
	service, err := startTestService(namespace)
	defer StopService(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestStopService(t *testing.T) {
	namespace := []string{"TestStopService"}
	startTestService(namespace)
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestStopServiceAlreadyStopped(t *testing.T) {
	namespace := []string{"TestStopServiceAlreadyStopped"}
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestServiceStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestServiceStatusNeverStarted"}
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
	running, err := IsServiceRunning(namespace)
	assert.Nil(t, err)
	assert.Equal(t, running, false)
	stopped, err := IsServiceStopped(namespace)
	assert.Nil(t, err)
	assert.Equal(t, stopped, true)
}

func TestServiceStatusRunning(t *testing.T) {
	namespace := []string{"TestServiceStatusRunning"}
	startTestService(namespace)
	defer StopService(namespace)
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, RUNNING)
	running, err := IsServiceRunning(namespace)
	assert.Nil(t, err)
	assert.Equal(t, running, true)
	stopped, err := IsServiceStopped(namespace)
	assert.Nil(t, err)
	assert.Equal(t, stopped, false)
}

func TestServiceStatusStopped(t *testing.T) {
	namespace := []string{"TestServiceStatusStopped"}
	startTestService(namespace)
	StopService(namespace)
	status, err := ServiceStatus(namespace)
	assert.Nil(t, err)
	assert.Equal(t, status, STOPPED)
	running, err := IsServiceRunning(namespace)
	assert.Nil(t, err)
	assert.Equal(t, running, false)
	stopped, err := IsServiceStopped(namespace)
	assert.Nil(t, err)
	assert.Equal(t, stopped, true)
}

func TestFindServiceNotExisting(t *testing.T) {
	namespace := []string{"TestFindServiceNotExisting"}
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.Nil(t, service)
}

func TestFindService(t *testing.T) {
	namespace := []string{"TestFindService"}
	startTestService(namespace)
	defer StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestFindServiceCloseName(t *testing.T) {
	namespace := []string{"TestFindServiceCloseName", "name"}
	namespace1 := []string{"TestFindServiceCloseName", "name2"}
	startTestService(namespace)
	defer StopService(namespace)
	startTestService(namespace1)
	defer StopService(namespace1)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, service)
}

func TestFindServiceStopped(t *testing.T) {
	namespace := []string{"TestFindServiceStopped"}
	startTestService(namespace)
	StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.Nil(t, service)
}

func TestListServices(t *testing.T) {
	namespace := Namespace([]string{"TestListServices"})
	namespace2 := Namespace([]string{"TestListServiceswithValue2"})
	StartService(docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: namespace,
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace,
					"label_name":                 "value_1",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: "nginx",
					Labels: map[string]string{
						"com.docker.stack.namespace": namespace,
					},
				},
			},
		},
	})
	StartService(docker.CreateServiceOptions{
		ServiceSpec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{
				Name: namespace2,
				Labels: map[string]string{
					"com.docker.stack.namespace": namespace2,
					"label_name_2":               "value_2",
				},
			},
			TaskTemplate: swarm.TaskSpec{
				ContainerSpec: &swarm.ContainerSpec{
					Image: "nginx",
					Labels: map[string]string{
						"com.docker.stack.namespace": namespace2,
					},
				},
			},
		},
	})
	defer StopService([]string{"TestListServices"})
	defer StopService([]string{"TestListServiceswithValue2"})
	services, err := ListServices("label_name")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace([]string{"TestListServices"}), services[0].Spec.Name)
}

func TestServiceLogs(t *testing.T) {
	namespace := []string{"TestServiceLogs"}
	startTestService(namespace)
	defer StopService(namespace)
	var stream bytes.Buffer
	err := ServiceLogs(namespace, &stream)
	assert.Nil(t, err)
	assert.NotNil(t, stream)
}
