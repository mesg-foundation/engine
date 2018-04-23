package service

import (
	"strings"

	"github.com/mesg-foundation/application/utils/hash"
)

// NAMESPACE is the namespace used for the docker services
const NAMESPACE string = "MESG"
const eventKey string = "Event"
const taskKey string = "Task"

func (service *Service) namespace() string {
	return strings.Join([]string{
		NAMESPACE,
		strings.Replace(service.Name, " ", "-", -1),
	}, "-")
}

// EventSubscriptionKey returns the key to listen for some events from this service
func (service *Service) EventSubscriptionKey() string {
	return hash.Calculate([]string{
		service.namespace(),
		eventKey,
	})
}

// TaskSubscriptionKey returns the key to listen for some tasks from this service
func (service *Service) TaskSubscriptionKey() string {
	return hash.Calculate([]string{
		service.namespace(),
		taskKey,
	})
}
