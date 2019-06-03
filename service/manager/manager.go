package manager

import (
	"github.com/mesg-foundation/core/service"
)

// Manager is responsible for managing Docker Containers of MESG services.
// it can be implemented for any container orchestration tool.
// TODO(ilgooz): discuss if these should accept service.Instance instead of service.Service.
type Manager interface {
	// Start starts service and returns related info provided by the underlying container
	// orchestration tool.
	Start(s *service.Service) (serviceIDs []string, networkID string, err error)

	// Stop stops service.
	Stop(s *service.Service) error

	// Status gives status of service.
	Status(s *service.Service) (service.StatusType, error)

	// Logs gives log streams of service.
	Logs(s *service.Service, dependencies ...string) ([]*service.Log, error)

	// Delete deletes anything related to service and its persistent data.
	Delete(s *service.Service) error
}
