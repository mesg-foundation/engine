package hashserializer

import (
	"fmt"
)

type testSerializable struct {
	a int
	b string
	c []string
}

func (data testSerializable) HashSerialize() string {
	ser := New()
	ser.AddInt("1", data.a)
	ser.AddString("2", data.b)
	ser.AddStringSlice("3", data.c)
	return ser.HashSerialize()
}

func ExampleHashSerializer() {
	ser := New()
	ser.AddBool("1", true)
	ser.AddBool("2", false)
	ser.AddFloat("3", 3.14159265359)
	ser.AddFloat("4", 0.0)
	ser.AddFloat("5", -42.14159265359)
	ser.AddInt("6", 42)
	ser.AddInt("7", 0.0)
	ser.AddInt("8", -42)
	ser.AddString("9", "")
	ser.AddString("10", "hello")
	ser.AddStringSlice("11", []string{"c", "b", "a"})
	ser.AddStringSlice("12", []string{})
	ser.AddStringSlice("13", []string{"c", "", "a"})
	ser.Add("14", testSerializable{42, "hello", []string{"c", "b", "a"}})
	ser.Add("15", testSerializable{})
	fmt.Println(ser.HashSerialize())
	// Output: 1:true;3:3.14159265359;5:-42.14159265359;6:42;8:-42;10:hello;11:0:c;1:b;2:a;;13:0:c;2:a;;14:1:42;2:hello;3:0:c;1:b;2:a;;;
}

func ExampleHashSerializer_AddBool_true() {
	ser := New()
	ser.AddBool("1", true)
	fmt.Println(ser.HashSerialize())
	// Output: 1:true;
}

func ExampleHashSerializer_AddBool_false() {
	ser := New()
	ser.AddBool("2", false)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddFloat() {
	ser := New()
	ser.AddFloat("3", 3.14159265359)
	fmt.Println(ser.HashSerialize())
	// Output: 3:3.14159265359;
}

func ExampleHashSerializer_AddFloat_zero() {
	ser := New()
	ser.AddFloat("4", 0.0)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddFloat_negative() {
	ser := New()
	ser.AddFloat("5", -42.14159265359)
	fmt.Println(ser.HashSerialize())
	// Output: 5:-42.14159265359;
}

func ExampleHashSerializer_AddInt() {
	ser := New()
	ser.AddInt("6", 42)
	fmt.Println(ser.HashSerialize())
	// Output: 6:42;
}

func ExampleHashSerializer_AddInt_zero() {
	ser := New()
	ser.AddInt("7", 0.0)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddInt_negative() {
	ser := New()
	ser.AddInt("8", -42)
	fmt.Println(ser.HashSerialize())
	// Output: 8:-42;
}

func ExampleHashSerializer_AddString_empty() {
	ser := New()
	ser.AddString("9", "")
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddString() {
	ser := New()
	ser.AddString("10", "hello")
	fmt.Println(ser.HashSerialize())
	// Output: 10:hello;
}

func ExampleHashSerializer_AddStringSlice() {
	ser := New()
	ser.AddStringSlice("11", []string{"c", "b", "a"})
	fmt.Println(ser.HashSerialize())
	// Output: 11:0:c;1:b;2:a;;
}

func ExampleHashSerializer_AddStringSlice_empty() {
	ser := New()
	ser.AddStringSlice("12", []string{})
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddStringSlice_empty_value() {
	ser := New()
	ser.AddStringSlice("13", []string{"c", "", "a"})
	fmt.Println(ser.HashSerialize())
	// Output: 13:0:c;2:a;;
}

func ExampleHashSerializer_Add() {
	ser := New()
	ser.Add("14", testSerializable{42, "hello", []string{"c", "b", "a"}})
	fmt.Println(ser.HashSerialize())
	// Output: 14:1:42;2:hello;3:0:c;1:b;2:a;;;
}

func ExampleHashSerializer_Add_empty() {
	ser := New()
	ser.Add("15", testSerializable{})
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleStringSlice() {
	ser := New()
	ser.Add("1", StringSlice([]string{"a", "b", "c"}))
	fmt.Println(ser.HashSerialize())
	// Output: 1:0:a;1:b;2:c;;
}

func ExampleHashSerializable_HashSerialize() {
	s := testSerializable{
		a: 42,
		b: "hello",
		c: []string{"c", "b", "a"},
	}
	fmt.Println(s.HashSerialize())
	//Output: 1:42;2:hello;3:0:c;1:b;2:a;;
}
