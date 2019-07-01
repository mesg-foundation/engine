package client

import "encoding/json"

// Marshal returns the MESG format encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal parses the MESG-encoded data input and stores the result in the value pointed to by v.
func Unmarshal(input []byte, v interface{}) error {
	return json.Unmarshal(input, v)
}
