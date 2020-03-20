package process

import (
	"fmt"

	"github.com/mesg-foundation/engine/protobuf/types"
)

// Resolve returns the value specified in the reference.
func (path *Process_Node_Reference_Path) Resolve(data *types.Struct) (*types.Value, error) {
	if path == nil {
		return &types.Value{Kind: &types.Value_StructValue{StructValue: data}}, nil
	}

	var v *types.Value
	key, ok := path.Selector.(*Process_Node_Reference_Path_Key)
	if !ok {
		return nil, fmt.Errorf("first selector in the path must be a key")
	}

	v, ok = data.Fields[key.Key]
	if !ok {
		return nil, fmt.Errorf("key %s not found", key.Key)
	}

	for p := path.Path; p != nil; p = p.Path {
		switch s := p.Selector.(type) {
		case *Process_Node_Reference_Path_Key:
			str, ok := v.GetKind().(*types.Value_StructValue)
			if !ok {
				return nil, fmt.Errorf("can't get key from non-struct value")
			}
			if str.StructValue.GetFields() == nil {
				return nil, fmt.Errorf("can't get key from nil-struct")
			}
			v, ok = str.StructValue.Fields[s.Key]
			if !ok {
				return nil, fmt.Errorf("key %s not found", s.Key)
			}
		case *Process_Node_Reference_Path_Index:
			list, ok := v.GetKind().(*types.Value_ListValue)
			if !ok {
				return nil, fmt.Errorf("can't get index from non-list value")
			}

			if len(list.ListValue.GetValues()) <= int(s.Index) {
				return nil, fmt.Errorf("index %d out of range", s.Index)
			}
			v = list.ListValue.Values[s.Index]
		default:
			return nil, fmt.Errorf("unknown selector type %T", v)
		}
	}

	return v, nil
}
