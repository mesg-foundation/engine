package applicationtest

import "encoding/json"

type Execute struct {
	serviceID string
	task      string
	data      string
}

func (e *Execute) ServiceID() string {
	return e.serviceID
}

func (e *Execute) Task() string {
	return e.task
}

func (e *Execute) Decode(out interface{}) error {
	return json.Unmarshal([]byte(e.data), out)
}
