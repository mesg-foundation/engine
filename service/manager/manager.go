package manager

import (
	"github.com/mesg-foundation/core/service"
)

// Manager is responsible for managing Docker Containers of MESG services.
// it can be implemented for any container orchestration tool.
type Manager interface {
	// Start starts service and returns service ids related to service.
	Start(s *service.Service) (serviceIDs []string, err error)

	// Stop stops service.
	Stop(s *service.Service) error

	// Status gives status of service.
	Status(s *service.Service) (service.StatusType, error)

	// Logs gives log streams of service.
	Logs(s *service.Service, dependencies ...string) ([]*service.Log, error)

	// Delete deletes anything related to service and its persistent data.
	Delete(s *service.Service) error
}
