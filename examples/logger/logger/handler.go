package logger

import (
	"encoding/json"

	mesg "github.com/mesg-foundation/go-service"
)

func (l *Logger) handler(execution *mesg.Execution) (string, mesg.Data) {
	var data logRequest
	if err := execution.Data(&data); err != nil {
		return "error", errorResponse{err.Error()}
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		return "error", errorResponse{err.Error()}
	}

	l.log.Printf("%s: %s", data.ServiceID, string(bytes))

	return "success", successResponse{"ok"}
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
