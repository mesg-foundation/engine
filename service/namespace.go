package service

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/mesg-foundation/core/utils/hash"
)

// NAMESPACE is the namespace used for the docker services.
const eventChannel string = "Event"
const taskChannel string = "Task"
const resultChannel string = "Result"

// namespace returns the namespace of the service.
func (service *Service) namespace() []string {
	sum := sha1.Sum([]byte(service.Hash))
	return []string{hex.EncodeToString(sum[:])}
}

// namespace builds the namespace of a dependency.
func (d *Dependency) namespace(serviceNamespace []string) []string {
	return append(serviceNamespace, d.Key)
}

// EventSubscriptionChannel returns the channel to listen for events from this service.
func (service *Service) EventSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		eventChannel,
	))
}

// TaskSubscriptionChannel returns the channel to listen for tasks from this service.
func (service *Service) TaskSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		taskChannel,
	))
}

// ResultSubscriptionChannel returns the channel to listen for tasks from this service.
func (service *Service) ResultSubscriptionChannel() string {
	return hash.Calculate(append(
		service.namespace(),
		resultChannel,
	))
}
