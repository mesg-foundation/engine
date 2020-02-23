package structhash

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/mr-tron/base58"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

// Dump takes a data structure and returns its string representation.
func Dump(v interface{}) [sha256.Size]byte {
	return sha256.Sum256([]byte(serialize(v)))
}

func serialize(v interface{}) string {
	return valueToString(reflect.ValueOf(v))
}

func valueToString(v reflect.Value) string {
	buf := bufPool.Get().(*bytes.Buffer)
	write(buf, v)
	s := buf.String()
	buf.Reset()
	bufPool.Put(buf)
	return s
}

//nolint:gocyclo
func write(buf *bytes.Buffer, v reflect.Value) {
	if isEmptyValue(v) {
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(v.Int(), 16))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteString(strconv.FormatUint(v.Uint(), 16))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(v.Float(), 'f', -1, 64))
	case reflect.String:
		buf.WriteString(v.String())
	case reflect.Ptr, reflect.Interface:
		write(buf, v.Elem())
	case reflect.Struct:
		vt := v.Type()
		items := make(map[int]string)
		isOneof := make(map[int]bool)
		keys := make([]int, 0, v.NumField())
		for i := 0; i < v.NumField(); i++ {
			vf := v.Field(i)
			if isEmptyValue(vf) {
				continue
			}

			to := parseTag(vt.Field(i))
			if to.skip {
				continue
			}

			if val := valueToString(vf); val != "" {
				// if field string == "" then it is chan,
				// func or invalid type and it can be skipped
				if to.index == 0 {
					// oneof interface - get the index of the field.
					str := strings.Split(val, ":")[0]
					index, err := strconv.ParseInt(str, 10, 64)
					if err != nil {
						panic(fmt.Sprintf("structhash: value in field %s is not a struct: %v", vt.Field(i).Name, err))
					}

					to.index = int(index)
					isOneof[to.index] = true
				}

				switch vf.Kind() {
				case reflect.Array, reflect.Slice, reflect.Map, reflect.Interface:
					h := sha256.Sum256([]byte(val))
					items[to.index] = base58.Encode(h[:])
				case reflect.Ptr:
					if vf.Elem().Kind() == reflect.Struct {
						h := sha256.Sum256([]byte(val))
						items[to.index] = base58.Encode(h[:])
						break
					}
					fallthrough
				default:
					items[to.index] = val
				}
				keys = append(keys, to.index)
			}
		}

		sort.Ints(keys)
		for _, index := range keys {
			if isOneof[index] {
				buf.WriteString(items[index])
			} else {
				buf.WriteString(strconv.Itoa(index))
				buf.WriteByte(':')
				buf.WriteString(items[index])
				buf.WriteByte(';')
			}
		}
	case reflect.Map:
		keys := v.MapKeys()
		if len(keys) == 0 {
			return
		}
		sort.Sort(byValue(keys))

		// Extract and sort the keys.
		for _, key := range keys {
			if val := valueToString(v.MapIndex(key)); val != "" {
				buf.WriteString(valueToString(key))
				buf.WriteByte(':')
				buf.WriteString(val)
				buf.WriteByte(';')
			}
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if val := valueToString(v.Index(i)); val != "" {
				buf.WriteString(strconv.FormatInt(int64(i), 16))
				buf.WriteByte(':')
				buf.WriteString(val)
				buf.WriteByte(';')
			}
		}
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return real(c) == 0 && imag(c) == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

// tagOptions is the string struct field's "hash"
type tagOptions struct {
	index int
	oneof bool
	skip  bool
}

// parseTag splits a struct field's protobuf tag into its name and index.
func parseTag(f reflect.StructField) tagOptions {
	tag := f.Tag.Get("protobuf_oneof")
	if tag != "" {
		return tagOptions{oneof: true}
	}

	if tag = f.Tag.Get("hash"); tag == "-" {
		return tagOptions{skip: true}
	}

	tag = f.Tag.Get("protobuf")
	if tag == "" {
		return tagOptions{skip: true}
	}

	options := strings.Split(tag, ",")
	if len(options) < 2 {
		return tagOptions{skip: true}
	}

	index, err := strconv.Atoi(options[1])
	if err != nil {
		return tagOptions{skip: true}
	}

	return tagOptions{index: index}
}

type byValue []reflect.Value

func (a byValue) Len() int      { return len(a) }
func (a byValue) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a byValue) Less(i, j int) bool {
	switch a[i].Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return a[i].Int() < a[j].Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return a[i].Uint() < a[j].Uint()
	case reflect.String:
		return a[i].String() < a[j].String()
	default:
		panic("value sort: type not supported")
	}
}
