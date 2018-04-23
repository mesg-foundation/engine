package service

import (
	"strings"

	"github.com/mesg-foundation/application/utils/hash"
)

// NAMESPACE is the namespace used for the docker services
const NAMESPACE string = "MESG"
const eventChannel string = "Event"
const taskChannel string = "Task"

func (service *Service) namespace() string {
	return strings.Join([]string{
		NAMESPACE,
		strings.Replace(service.Name, " ", "-", -1),
	}, "-")
}

// EventSubscriptionChannel returns the channel to listen for events from this service
func (service *Service) EventSubscriptionChannel() string {
	return hash.Calculate([]string{
		service.namespace(),
		eventChannel,
	})
}

// TaskSubscriptionChannel returns the channel to listen for tasks from this service
func (service *Service) TaskSubscriptionChannel() string {
	return hash.Calculate([]string{
		service.namespace(),
		taskChannel,
	})
}
