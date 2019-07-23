package client

import (
	"encoding/json"

	"github.com/mitchellh/mapstructure"
)

// Marshal returns the MESG format encoding of v.
func Marshal(input interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var output map[string]interface{}
	if err := json.Unmarshal(b, &output); err != nil {
		return nil, err
	}

	return output, nil
}

// Unmarshal parses the MESG-encoded data input and stores the result in the value pointed to by v.
func Unmarshal(input map[string]interface{}, output interface{}) error {
	return mapstructure.Decode(input, output)
}
