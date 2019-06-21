package core

import (
	"github.com/mesg-foundation/core/protobuf/coreapi"
	service "github.com/mesg-foundation/core/service"
)

func toProtoServiceStatusType(s service.StatusType) coreapi.Service_Status {
	switch s {
	default:
		return coreapi.Service_UNKNOWN
	case service.STOPPED:
		return coreapi.Service_STOPPED
	case service.STARTING:
		return coreapi.Service_STARTING
	case service.PARTIAL:
		return coreapi.Service_PARTIAL
	case service.RUNNING:
		return coreapi.Service_RUNNING
	}
}
