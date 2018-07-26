package logger

import (
	"encoding/json"

	"github.com/ilgooz/mesg-go/service"
)

func (l *Logger) handler(req *service.Request) service.Response {
	var data logRequest
	if err := req.Decode(&data); err != nil {
		return service.Response{
			"error": errorResponse{err.Error()},
		}
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		return service.Response{
			"error": errorResponse{err.Error()},
		}
	}

	l.log.Printf("%s: %s", data.ServiceID, string(bytes))

	return service.Response{
		"success": successResponse{"ok"},
	}
}

type logRequest struct {
	ServiceID string      `json:"serviceID"`
	Data      interface{} `json:"data"`
}

type successResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}
