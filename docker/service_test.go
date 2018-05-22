package docker

import (
	"strings"
	"testing"

	"github.com/docker/docker/api/types/swarm"
	"github.com/stvp/assert"
)

func TestDockerServiceMatch(t *testing.T) {
	namespace := strings.Join([]string{"MESG", "TestDockerServiceMatch"}, "_")
	dockerServices := []swarm.Service{
		swarm.Service{
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: strings.Join([]string{namespace, "test1"}, "_"),
				},
			},
		},
		swarm.Service{
			Spec: swarm.ServiceSpec{
				Annotations: swarm.Annotations{
					Name: strings.Join([]string{namespace, "test2"}, "_"),
				},
			},
		},
	}
	res1 := serviceMatch(dockerServices, namespace, "test")
	assert.Equal(t, res1, swarm.Service{})
	res2 := serviceMatch(dockerServices, namespace, "test1")
	assert.Equal(t, res2, dockerServices[0])
	res3 := serviceMatch(dockerServices, namespace, "test2")
	assert.Equal(t, res3, dockerServices[1])
}
