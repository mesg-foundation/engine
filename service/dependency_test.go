package service

// import (
// 	"strings"
// 	"testing"

// 	"github.com/mesg-foundation/core/docker"
// 	"github.com/stvp/assert"
// )

// func TestGetDockerService(t *testing.T) {
// 	namespace := strings.Join([]string{NAMESPACE, "TestGetDockerService"}, "_")
// 	name := "test"
// 	dependency := Dependency{Image: "nginx"}
// 	dependency.Start(&Service{}, dependencyDetails{
// 		namespace:      namespace,
// 		dependencyName: name,
// 		serviceName:    "TestGetDockerService",
// 	}, testDaemonIP, testSharedNetwork)
// 	res, err := docker.Service(namespace, name)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, res.ID, "")
// 	res, err = docker.Service(namespace, "textx")
// 	assert.Nil(t, err)
// 	assert.Equal(t, res.ID, "")
// 	docker.Stop(namespace, name)
// }
