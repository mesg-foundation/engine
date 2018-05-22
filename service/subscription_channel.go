package service

import (
	"github.com/mesg-foundation/core/utils/hash"
)

const eventChannel string = "Event"
const taskChannel string = "Task"
const resultChannel string = "Result"

// EventSubscriptionChannel returns the channel to listen for events from this service
func (service *Service) EventSubscriptionChannel() string {
	return hash.Calculate([]string{
		service.Name,
		eventChannel,
	})
}

// TaskSubscriptionChannel returns the channel to listen for tasks from this service
func (service *Service) TaskSubscriptionChannel() string {
	return hash.Calculate([]string{
		service.Name,
		taskChannel,
	})
}

// ResultSubscriptionChannel returns the channel to listen for tasks from this service
func (service *Service) ResultSubscriptionChannel() string {
	return hash.Calculate([]string{
		service.Name,
		resultChannel,
	})
}
