package manager

import (
	"github.com/mesg-foundation/core/service"
)

// Manager is responsible for managing Docker Containers of MESG services.
// it can be implemented for any container orchestration tool.
type Manager interface {
	// Logs gives log streams of service.
	Logs(s *service.Service, dependencies ...string) ([]*service.Log, error)
}
