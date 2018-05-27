package docker

import (
	"fmt"
	"testing"

	"github.com/stvp/assert"
)

func TestListTasks(t *testing.T) {
	namespace := []string{"TestListTasks"}
	startTestService(namespace)
	defer StopService(namespace)
	tasks, err := ListTasks(namespace)
	assert.Nil(t, err)
	assert.NotNil(t, tasks)
	fmt.Println("number of", len(tasks))
	for _, task := range tasks {
		fmt.Println("task", task.ID, task.Status, task.ServiceID)
	}
}
