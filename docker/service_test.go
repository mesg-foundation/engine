package docker

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stvp/assert"
)

func startTestService(namespace []string) (dockerService *swarm.Service, err error) {
	return StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "nginx",
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
	namespace := []string{"TestStopService"}
	err := StopService(namespace)
	assert.Nil(t, err)
}

func TestServiceStatusNeverStarted(t *testing.T) {
	namespace := []string{"TestServiceStatusNeverStarted"}
	status := ServiceStatus(namespace)
	assert.Equal(t, status, STOPPED)
	assert.Equal(t, IsServiceRunning(namespace), false)
	assert.Equal(t, IsServiceStopped(namespace), true)
}

func TestServiceStatusRunning(t *testing.T) {
	namespace := []string{"TestServiceStatus"}
	startTestService(namespace)
	defer StopService(namespace)
	status := ServiceStatus(namespace)
	assert.Equal(t, status, RUNNING)
	assert.Equal(t, IsServiceRunning(namespace), true)
	assert.Equal(t, IsServiceStopped(namespace), false)
}

func TestServiceStatus(t *testing.T) {
	namespace := []string{"TestServiceStatus"}
	startTestService(namespace)
	StopService(namespace)
	status := ServiceStatus(namespace)
	assert.Equal(t, status, STOPPED)
	assert.Equal(t, IsServiceRunning(namespace), false)
	assert.Equal(t, IsServiceStopped(namespace), true)
}

func TestFindServiceNeverStarted(t *testing.T) {
	namespace := []string{"TestFindServiceNeverStarted"}
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

func TestFindServiceStopped(t *testing.T) {
	namespace := []string{"TestFindService"}
	startTestService(namespace)
	StopService(namespace)
	service, err := FindService(namespace)
	assert.Nil(t, err)
	assert.Nil(t, service)
}

func TestServiceMatch(t *testing.T) {
	namespace := "TestServiceMatch"
	dockerServices := []swarm.Service{
		swarm.Service{
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: Namespace([]string{namespace, "test1"}),
				},
			},
		},
		swarm.Service{
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: Namespace([]string{namespace, "test2"}),
				},
			},
		},
	}
	res1 := serviceMatch(dockerServices, []string{namespace, "wrong-name"})
	assert.Nil(t, res1)
	res2 := serviceMatch(dockerServices, []string{namespace, "test1"})
	assert.Equal(t, res2, &dockerServices[0])
	res3 := serviceMatch(dockerServices, []string{namespace, "test2"})
	assert.Equal(t, res3, &dockerServices[1])
	res4 := serviceMatch(dockerServices, []string{"not-existing-namespace", "test2"})
	assert.Nil(t, res4)
}

func TestListServices(t *testing.T) {
	namespace := []string{"TestListServices"}
	namespace2 := []string{"TestListServiceswithValue2"}
	StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "nginx",
		Labels: map[string]string{
			"label_name": "value_1",
		},
	})
	StartService(&ServiceOptions{
		Namespace: namespace2,
		Image:     "nginx",
	})
	defer StopService(namespace)
	defer StopService(namespace2)
	services, err := ListServices("label_name", "")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace(namespace), services[0].Spec.Name)
}

func TestListServiceswithValue(t *testing.T) {
	namespace := []string{"TestListServiceswithValue"}
	namespace2 := []string{"TestListServiceswithValue2"}
	StartService(&ServiceOptions{
		Namespace: namespace,
		Image:     "nginx",
		Labels: map[string]string{
			"label_name": "value_1",
		},
	})
	StartService(&ServiceOptions{
		Namespace: namespace2,
		Image:     "nginx",
		Labels: map[string]string{
			"label_name": "value_2",
		},
	})
	defer StopService(namespace)
	defer StopService(namespace2)
	services, err := ListServices("label_name", "value_2")
	assert.Nil(t, err)
	assert.Equal(t, 1, len(services))
	assert.Equal(t, Namespace(namespace2), services[0].Spec.Name)
}
