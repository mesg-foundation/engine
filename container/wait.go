package container

import (
	"errors"
	"strings"
	"time"
)

// WaitForStatus wait for the container to have the provided status. Returns error as soon as possible
func WaitForStatus(namespace []string, status StatusType) (err error) {
	var tasksErrors []string
	var currentStatus StatusType
	for {
		tasksErrors, err = TasksError(namespace)
		if err != nil {
			break
		}
		if len(tasksErrors) > 0 {
			err = errors.New(strings.Join(tasksErrors, ", "))
			break
		}
		currentStatus, err = Status(namespace)
		if err != nil || currentStatus == status {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	return
}
