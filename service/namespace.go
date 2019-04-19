package service

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/mesg-foundation/core/utils/hash"
)

// NAMESPACE is the namespace used for the docker services.
const eventTopic string = "Event"
const executionTopic string = "Execution"

// namespace returns the namespace of the service.
func (service *Service) namespace() []string {
	sum := sha1.Sum([]byte(service.Hash))
	return []string{hex.EncodeToString(sum[:])}
}

// namespace builds the namespace of a dependency.
func (d *Dependency) namespace(serviceNamespace []string) []string {
	return append(serviceNamespace, d.Key)
}

// EventSubTopic returns the topic to listen for events from this service.
func (service *Service) EventSubTopic() string {
	return hash.Calculate(append(service.namespace(), eventTopic))
}

// ExecutionSubTopic returns the topic to listen for executions from this service.
func (service *Service) ExecutionSubTopic() string {
	return hash.Calculate(append(service.namespace(), executionTopic))
}
