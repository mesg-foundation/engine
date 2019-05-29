package logger

import (
	"encoding/json"

	"github.com/mesg-foundation/core/client/service"
)

func (l *Logger) handler(execution *service.Execution) (interface{}, error) {
	var data logRequest
	if err := execution.Data(&data); err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(data.Data)
	if err != nil {
		return nil, err
	}

	l.log.Printf("%s: %s", data.ServiceID, string(bytes))
	return successResponse{"ok"}, nil
}

type logRequest struct {
	ServiceID string      `json:"serviceID"`
	Data      interface{} `json:"data"`
}

type successResponse struct {
	Message string `json:"message"`
}
