package container

import (
	"testing"

	"github.com/stvp/assert"
)

func TestServiceOptionMergeNamespace(t *testing.T) {
	namespace := []string{"name1", "name2"}
	options := &ServiceOptions{
		Namespace: namespace,
	}
	expectedNamespace := Namespace(namespace)
	service := options.toSwarmServiceSpec()
	assert.Equal(t, expectedNamespace, service.Annotations.Name)
	assert.Equal(t, expectedNamespace, service.Annotations.Labels["com.docker.stack.namespace"])
	assert.Equal(t, expectedNamespace, service.TaskTemplate.ContainerSpec.Labels["com.docker.stack.namespace"])
}

func TestServiceOptionMergeImage(t *testing.T) {
	image := "nginx"
	options := &ServiceOptions{
		Image: image,
	}
	service := options.toSwarmServiceSpec()
	assert.Equal(t, image, service.Annotations.Labels["com.docker.stack.image"])
	assert.Equal(t, image, service.TaskTemplate.ContainerSpec.Image)
}

func TestServiceOptionMergeLabels(t *testing.T) {
	options := &ServiceOptions{
		Labels: map[string]string{
			"label1": "foo",
			"label2": "bar",
		},
	}
	service := options.toSwarmServiceSpec()
	assert.Equal(t, "foo", service.Annotations.Labels["label1"])
	assert.Equal(t, "bar", service.Annotations.Labels["label2"])
}

func TestServiceOptionMergePorts(t *testing.T) {
	options := &ServiceOptions{
		Ports: []Port{
			Port{
				Published: 50503,
				Target:    50501,
			},
			Port{
				Published: 30503,
				Target:    30501,
			},
		},
	}
	service := options.toSwarmServiceSpec()
	ports := service.EndpointSpec.Ports
	assert.Equal(t, 2, len(ports))
	assert.Equal(t, uint32(50503), ports[0].PublishedPort)
	assert.Equal(t, uint32(50501), ports[0].TargetPort)
	assert.Equal(t, uint32(30503), ports[1].PublishedPort)
	assert.Equal(t, uint32(30501), ports[1].TargetPort)
}

func TestServiceOptionMergeMounts(t *testing.T) {
	options := &ServiceOptions{
		Mounts: []Mount{
			Mount{
				Source: "source/file",
				Target: "target/file",
			},
		},
	}
	service := options.toSwarmServiceSpec()
	mounts := service.TaskTemplate.ContainerSpec.Mounts
	assert.Equal(t, 1, len(mounts))
	assert.Equal(t, "source/file", mounts[0].Source)
	assert.Equal(t, "target/file", mounts[0].Target)
}

func TestServiceOptionMergeEnv(t *testing.T) {
	options := &ServiceOptions{
		Env: []string{"env1", "env2"},
	}
	service := options.toSwarmServiceSpec()
	env := service.TaskTemplate.ContainerSpec.Env
	assert.Equal(t, 2, len(env))
	assert.Equal(t, "env1", env[0])
	assert.Equal(t, "env2", env[1])
}

func TestServiceOptionMergeNetworks(t *testing.T) {
	options := &ServiceOptions{
		NetworksID: []string{"network1", "network2"},
	}
	service := options.toSwarmServiceSpec()
	networks := service.TaskTemplate.Networks
	assert.Equal(t, 2, len(networks))
	assert.Equal(t, "network1", networks[0].Target)
	assert.Equal(t, "network2", networks[1].Target)
}
