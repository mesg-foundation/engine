package structhash

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/crypto/sha3"
)

// Sha1 takes a data structure and returns its sha1 hash.
func Sha1(v interface{}) [sha1.Size]byte {
	return sha1.Sum(serialize(v))
}

// Sha3 takes a data structure and returns its sha3 hash.
func Sha3(v interface{}) [64]byte {
	return sha3.Sum512(serialize(v))
}

// Md5 takes a data structure and returns its md5 hash.
func Md5(v interface{}) [md5.Size]byte {
	return md5.Sum(serialize(v))
}

// Dump takes a data structure and returns its string representation.
func Dump(v interface{}) []byte {
	return serialize(v)
}

func serialize(v interface{}) []byte {
	return []byte(valueToString(reflect.ValueOf(v)))
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

func valueToString(v reflect.Value) string {
	return write(new(bytes.Buffer), v).String()
}

//nolint:gocyclo
func write(buf *bytes.Buffer, v reflect.Value) *bytes.Buffer {
	switch v.Kind() {
	case reflect.Bool:
		buf.WriteString(strconv.FormatBool(v.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		buf.WriteString(strconv.FormatInt(v.Int(), 16))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		buf.WriteString(strconv.FormatUint(v.Uint(), 16))
	case reflect.Float32, reflect.Float64:
		buf.WriteString(strconv.FormatFloat(v.Float(), 'e', -1, 64))
	case reflect.String:
		buf.WriteString("\"" + v.String() + "\"")
	case reflect.Interface:
		if v.IsNil() {
			buf.WriteString("nil")
			return buf
		}
		if v.CanInterface() {
			write(buf, v.Elem())
		}
	case reflect.Struct:
		vt := v.Type()
		items := make([]string, 0)
		for i := 0; i < v.NumField(); i++ {
			sf := vt.Field(i)
			to := parseTag(sf)
			if to.skip || to.omitempty && isEmptyValue(v.Field(i)) {
				continue
			}

			str := valueToString(v.Field(i))
			// if field string == "" then it is chan,func or invalid type
			// and skip it
			if str == "" {
				continue
			}
			items = append(items, to.name+":"+str)
		}
		sort.Strings(items)

		buf.WriteByte('{')
		for i := range items {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i])
		}
		buf.WriteByte('}')
	case reflect.Map:
		if v.IsNil() {
			buf.WriteString("()nil")
			return buf
		}
		buf.WriteByte('(')

		keys := v.MapKeys()
		items := make([]string, len(keys))

		// Extract and sort the keys.
		for i, key := range keys {
			items[i] = valueToString(key) + ":" + valueToString(v.MapIndex(key))
		}
		sort.Strings(items)

		for i := range items {
			if i > 0 {
				buf.WriteByte(',')
			}
			buf.WriteString(items[i])
		}
		buf.WriteByte(')')
	case reflect.Slice:
		if v.IsNil() {
			buf.WriteString("[]nil")
			return buf
		}
		fallthrough
	case reflect.Array:
		buf.WriteByte('[')
		for i := 0; i < v.Len(); i++ {
			if i != 0 {
				buf.WriteByte(',')
			}
			write(buf, v.Index(i))
		}
		buf.WriteByte(']')
	case reflect.Ptr:
		if v.IsNil() {
			buf.WriteString("nil")
			return buf
		}
		write(buf, v.Elem())
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		buf.WriteString(strconv.FormatFloat(real(c), 'e', -1, 64))
		buf.WriteString(strconv.FormatFloat(imag(c), 'e', -1, 64))
		buf.WriteString("i")
	}
	return buf
}

// tagOptions is the string struct field's "hash"
type tagOptions struct {
	name      string
	omitempty bool
	skip      bool
}

// parseTag splits a struct field's hash tag into its name and
// comma-separated options.
func parseTag(f reflect.StructField) tagOptions {
	tag := f.Tag.Get("hash")
	if tag == "-" {
		return tagOptions{skip: true}
	}

	to := tagOptions{name: f.Name}
	if tag == "" {
		return to
	}

	options := strings.Split(tag, ",")
	for _, option := range options {
		switch {
		case option == "omitempty":
			to.omitempty = true
		case strings.HasPrefix(option, "name:"):
			to.name = option[len("name:"):]
		default:
			panic(fmt.Sprintf("structhash: field %s with tag hash:%q has invalid option %q", f.Name, tag, option))
		}
	}
	return to
}
