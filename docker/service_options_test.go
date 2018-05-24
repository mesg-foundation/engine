package docker

import (
	"testing"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	godocker "github.com/fsouza/go-dockerclient"
	"github.com/stvp/assert"
)

func TestServiceOptionMergeNamespace(t *testing.T) {
	namespace := []string{"name1", "name2"}
	options := &ServiceOptions{
		Namespace: namespace,
	}
	expectedNamespace := Namespace(namespace)
	options.merge()

	service := options.CreateServiceOptions
	assert.Equal(t, expectedNamespace, service.Annotations.Name)
	assert.Equal(t, expectedNamespace, service.Annotations.Labels["com.docker.stack.namespace"])
	assert.Equal(t, expectedNamespace, service.TaskTemplate.ContainerSpec.Labels["com.docker.stack.namespace"])
}

func TestServiceOptionMergeImage(t *testing.T) {
	image := "nginx"
	options := &ServiceOptions{
		Image: image,
	}
	options.merge()

	service := options.CreateServiceOptions
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
	options.merge()

	service := options.CreateServiceOptions
	assert.Equal(t, "foo", service.Annotations.Labels["label1"])
	assert.Equal(t, "bar", service.Annotations.Labels["label2"])
}

func TestServiceOptionMergeLabelsWithExisting(t *testing.T) {
	options := &ServiceOptions{
		Labels: map[string]string{
			"label1": "foo",
			"label2": "bar",
		},
		CreateServiceOptions: &godocker.CreateServiceOptions{
			ServiceSpec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Labels: map[string]string{
						"label1": "should be replaced",
						"label3": "bar",
					},
				},
			},
		},
	}
	options.merge()

	service := options.CreateServiceOptions
	assert.Equal(t, "foo", service.Annotations.Labels["label1"])
	assert.Equal(t, "bar", service.Annotations.Labels["label2"])
	assert.Equal(t, "bar", service.Annotations.Labels["label3"])
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
	options.merge()

	service := options.CreateServiceOptions
	ports := service.EndpointSpec.Ports
	assert.Equal(t, 2, len(ports))
	assert.Equal(t, uint32(50503), ports[0].PublishedPort)
	assert.Equal(t, uint32(50501), ports[0].TargetPort)
	assert.Equal(t, uint32(30503), ports[1].PublishedPort)
	assert.Equal(t, uint32(30501), ports[1].TargetPort)
}

func TestServiceOptionMergePortsWithExisting(t *testing.T) {
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
		CreateServiceOptions: &godocker.CreateServiceOptions{
			ServiceSpec: swarm.ServiceSpec{
				EndpointSpec: &swarm.EndpointSpec{
					Ports: []swarm.PortConfig{
						swarm.PortConfig{
							Protocol:      swarm.PortConfigProtocolTCP,
							PublishMode:   swarm.PortConfigPublishModeIngress,
							PublishedPort: uint32(231),
							TargetPort:    uint32(232),
						},
						swarm.PortConfig{
							Protocol:      swarm.PortConfigProtocolTCP,
							PublishMode:   swarm.PortConfigPublishModeIngress,
							PublishedPort: uint32(131),
							TargetPort:    uint32(132),
						},
					},
				},
			},
		},
	}
	options.merge()

	service := options.CreateServiceOptions
	ports := service.EndpointSpec.Ports
	assert.Equal(t, 4, len(ports))
	assert.Equal(t, uint32(231), ports[0].PublishedPort)
	assert.Equal(t, uint32(232), ports[0].TargetPort)
	assert.Equal(t, uint32(131), ports[1].PublishedPort)
	assert.Equal(t, uint32(132), ports[1].TargetPort)
	assert.Equal(t, uint32(50503), ports[2].PublishedPort)
	assert.Equal(t, uint32(50501), ports[2].TargetPort)
	assert.Equal(t, uint32(30503), ports[3].PublishedPort)
	assert.Equal(t, uint32(30501), ports[3].TargetPort)
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
	options.merge()

	service := options.CreateServiceOptions
	mounts := service.TaskTemplate.ContainerSpec.Mounts
	assert.Equal(t, 1, len(mounts))
	assert.Equal(t, "source/file", mounts[0].Source)
	assert.Equal(t, "target/file", mounts[0].Target)
}

func TestServiceOptionMergeMountsWithExisting(t *testing.T) {
	options := &ServiceOptions{
		Mounts: []Mount{
			Mount{
				Source: "source/file2",
				Target: "target/file2",
			},
		},
		CreateServiceOptions: &godocker.CreateServiceOptions{
			ServiceSpec: swarm.ServiceSpec{
				TaskTemplate: swarm.TaskSpec{
					ContainerSpec: &swarm.ContainerSpec{
						Mounts: []mount.Mount{
							mount.Mount{
								Source: "source/file1",
								Target: "target/file1",
							},
						},
					},
				},
			},
		},
	}
	options.merge()

	service := options.CreateServiceOptions
	mounts := service.TaskTemplate.ContainerSpec.Mounts
	assert.Equal(t, 2, len(mounts))
	assert.Equal(t, "source/file1", mounts[0].Source)
	assert.Equal(t, "target/file1", mounts[0].Target)
	assert.Equal(t, "source/file2", mounts[1].Source)
	assert.Equal(t, "target/file2", mounts[1].Target)
}

func TestServiceOptionMergeEnv(t *testing.T) {
	options := &ServiceOptions{
		Env: []string{"env1", "env2"},
	}
	options.merge()

	env := options.CreateServiceOptions.TaskTemplate.ContainerSpec.Env
	assert.Equal(t, 2, len(env))
	assert.Equal(t, "env1", env[0])
	assert.Equal(t, "env2", env[1])
}

func TestServiceOptionMergeEnvWithExisting(t *testing.T) {
	options := &ServiceOptions{
		Env: []string{"env1", "env2"},
		CreateServiceOptions: &godocker.CreateServiceOptions{
			ServiceSpec: swarm.ServiceSpec{
				TaskTemplate: swarm.TaskSpec{
					ContainerSpec: &swarm.ContainerSpec{
						Env: []string{"env3", "env4"},
					},
				},
			},
		},
	}
	options.merge()

	env := options.CreateServiceOptions.TaskTemplate.ContainerSpec.Env
	assert.Equal(t, 4, len(env))
	assert.Equal(t, "env3", env[0])
	assert.Equal(t, "env4", env[1])
	assert.Equal(t, "env1", env[2])
	assert.Equal(t, "env2", env[3])
}

func TestServiceOptionMergeNetworks(t *testing.T) {
	options := &ServiceOptions{
		NetworksID: []string{"network1", "network2"},
	}
	options.merge()

	networks := options.CreateServiceOptions.Networks
	assert.Equal(t, 2, len(networks))
	assert.Equal(t, "network1", networks[0].Target)
	assert.Equal(t, "network2", networks[1].Target)
}

func TestServiceOptionMergeNetworksWithExisting(t *testing.T) {
	options := &ServiceOptions{
		NetworksID: []string{"network1", "network2"},
		CreateServiceOptions: &godocker.CreateServiceOptions{
			ServiceSpec: swarm.ServiceSpec{
				Networks: []swarm.NetworkAttachmentConfig{
					swarm.NetworkAttachmentConfig{
						Target: "network3",
					},
					swarm.NetworkAttachmentConfig{
						Target: "network4",
					},
				},
			},
		},
	}
	options.merge()

	networks := options.CreateServiceOptions.Networks
	assert.Equal(t, 4, len(networks))
	assert.Equal(t, "network3", networks[0].Target)
	assert.Equal(t, "network4", networks[1].Target)
	assert.Equal(t, "network1", networks[2].Target)
	assert.Equal(t, "network2", networks[3].Target)
}
