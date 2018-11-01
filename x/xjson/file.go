package xjson

import (
	"encoding/json"
	"io/ioutil"
)

// ReadFile reads the file named by filename and returns the contents.
// It returns error if file is not well formated json.
func ReadFile(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := json.Unmarshal(content, &raw); err != nil {
		return nil, err
	}
	return raw, nil
}
