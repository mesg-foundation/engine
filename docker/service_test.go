package docker

import (
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stvp/assert"
)

func TestServiceMatch(t *testing.T) {
	namespace := "TestDockerServiceMatch"
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
	res1 := serviceMatch(dockerServices, []string{namespace, "test"})
	assert.Nil(t, res1)
	res2 := serviceMatch(dockerServices, []string{namespace, "test1"})
	assert.Equal(t, res2, &dockerServices[0])
	res3 := serviceMatch(dockerServices, []string{namespace, "test2"})
	assert.Equal(t, res3, &dockerServices[1])
}
