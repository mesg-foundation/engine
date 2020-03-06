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
	return New().
		AddInt("1", data.a).
		AddString("2", data.b).
		AddStringSlice("3", data.c).
		HashSerialize()
}

func ExampleHashSerializer() {
	ser := New().
		AddBool("1", true).                                                // 1:true;
		AddBool("2", false).                                               // discarded
		AddFloat("3", 3.14159265359).                                      // 3:3.14159265359;
		AddFloat("4", 0.0).                                                // discarded
		AddFloat("5", -42.14159265359).                                    // 5:-42.14159265359;
		AddInt("6", 42).                                                   // 6:42;
		AddInt("7", 0.0).                                                  // discarded
		AddInt("8", -42).                                                  // 8:-42;
		AddString("9", "").                                                // discarded
		AddString("10", "hello").                                          // 10:hello;
		AddStringSlice("11", []string{"c", "b", "a"}).                     // 11:0:c;1:b;2:a;;
		AddStringSlice("12", []string{}).                                  // discarded
		AddStringSlice("13", []string{"c", "", "a"}).                      // 13:0:c;2:a;;
		Add("14", testSerializable{42, "hello", []string{"c", "b", "a"}}). // 14:1:42;2:hello;3:0:c;1:b;2:a;;;
		Add("15", testSerializable{})                                      // discarded
	fmt.Println(ser.HashSerialize())
	// Output: 1:true;3:3.14159265359;5:-42.14159265359;6:42;8:-42;10:hello;11:0:c;1:b;2:a;;13:0:c;2:a;;14:1:42;2:hello;3:0:c;1:b;2:a;;;
}

func ExampleHashSerializer_AddBool_true() {
	ser := New().
		AddBool("1", true)
	fmt.Println(ser.HashSerialize())
	// Output: 1:true;
}

func ExampleHashSerializer_AddBool_false() {
	ser := New().
		AddBool("2", false)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddFloat() {
	ser := New().
		AddFloat("3", 3.14159265359)
	fmt.Println(ser.HashSerialize())
	// Output: 3:3.14159265359;
}

func ExampleHashSerializer_AddFloat_zero() {
	ser := New().
		AddFloat("4", 0.0)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddFloat_negative() {
	ser := New().
		AddFloat("5", -42.14159265359)
	fmt.Println(ser.HashSerialize())
	// Output: 5:-42.14159265359;
}

func ExampleHashSerializer_AddInt() {
	ser := New().
		AddInt("6", 42)
	fmt.Println(ser.HashSerialize())
	// Output: 6:42;
}

func ExampleHashSerializer_AddInt_zero() {
	ser := New().
		AddInt("7", 0.0)
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddInt_negative() {
	ser := New().
		AddInt("8", -42)
	fmt.Println(ser.HashSerialize())
	// Output: 8:-42;
}

func ExampleHashSerializer_AddString_empty() {
	ser := New().
		AddString("9", "")
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddString() {
	ser := New().
		AddString("10", "hello")
	fmt.Println(ser.HashSerialize())
	// Output: 10:hello;
}

func ExampleHashSerializer_AddStringSlice() {
	ser := New().
		AddStringSlice("11", []string{"c", "b", "a"})
	fmt.Println(ser.HashSerialize())
	// Output: 11:0:c;1:b;2:a;;
}

func ExampleHashSerializer_AddStringSlice_empty() {
	ser := New().
		AddStringSlice("12", []string{})
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleHashSerializer_AddStringSlice_empty_value() {
	ser := New().
		AddStringSlice("13", []string{"c", "", "a"})
	fmt.Println(ser.HashSerialize())
	// Output: 13:0:c;2:a;;
}

func ExampleHashSerializer_Add() {
	ser := New().
		Add("14", testSerializable{42, "hello", []string{"c", "b", "a"}})
	fmt.Println(ser.HashSerialize())
	// Output: 14:1:42;2:hello;3:0:c;1:b;2:a;;;
}

func ExampleHashSerializer_Add_empty() {
	ser := New().
		Add("15", testSerializable{})
	fmt.Println(ser.HashSerialize())
	// Output:
}

func ExampleStringSlice() {
	ser := New().
		Add("1", StringSlice([]string{"a", "b", "c"}))
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
