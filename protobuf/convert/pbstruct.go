package convert

import (
	"bytes"
	"encoding/json"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var marshaler = jsonpb.Marshaler{EmitDefaults: true}

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// Marshal converts protobuf struct to map[string]interface{}.
func Marshal(in proto.Message, out interface{}) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	dec := json.NewDecoder(buf)
	if err := marshaler.Marshal(buf, in); err != nil {
		return err
	}
	return dec.Decode(out)
}

// Unmarshal converts map[string]interface{} to protobuf struct.
func Unmarshal(in interface{}, out proto.Message) error {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)

	enc := json.NewEncoder(buf)
	if err := enc.Encode(in); err != nil {
		return err
	}
	return jsonpb.Unmarshal(buf, out)
}
