package container

import (
	"errors"
	"strings"
	"time"
)

// waitForStatus waits for the container to have the provided status. Returns error as soon as possible.
func waitForStatus(namespace []string, status StatusType) error {
	for {
		tasksErrors, err := TasksError(namespace)
		if err != nil {
			return err
		}
		if len(tasksErrors) > 0 {
			return errors.New(strings.Join(tasksErrors, ", "))
		}
		currentStatus, err := Status(namespace)
		if err != nil || currentStatus == status {
			return err
		}
		time.Sleep(500 * time.Millisecond)
	}
}
