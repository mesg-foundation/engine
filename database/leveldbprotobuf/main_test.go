package leveldbprotobuf

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stvp/assert"
)

type testRecord struct {
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Foo string `protobuf:"bytes,2,opt,name=foo" json:"foo,omitempty"`
	Bar string `protobuf:"bytes,3,opt,name=bar" json:"bar,omitempty"`
}

func (m *testRecord) Reset()         { *m = testRecord{} }
func (m *testRecord) String() string { return proto.CompactTextString(m) }
func (*testRecord) ProtoMessage()    {}

const testCollection = "test"
const testCollection2 = "test2"

var Db = &Database{}

var keyCounter = 0

func insert(collection string) (key string, err error) {
	keyCounter++
	keyString := strconv.Itoa(keyCounter)
	key = "id_" + keyString
	data := testRecord{
		Key: key,
		Foo: "hello " + keyString,
		Bar: "bye " + keyString,
	}
	err = Db.Insert(collection, key, &data)
	return
}

func TestOpen(t *testing.T) {
	err := Db.Open()
	defer Db.Close()
	assert.Nil(t, err)
}

func TestClose(t *testing.T) {
	Db.Open()
	err := Db.Close()
	assert.Nil(t, err)
}

func TestInsert(t *testing.T) {
	Db.Open()
	defer Db.Close()
	key, err := insert(testCollection)
	assert.Nil(t, err)
	Db.Delete(testCollection, key)
}

func TestDelete(t *testing.T) {
	Db.Open()
	defer Db.Close()
	key, _ := insert(testCollection)
	err := Db.Delete(testCollection, key)
	assert.Nil(t, err)
}

func TestFind(t *testing.T) {
	Db.Open()
	defer Db.Close()
	key, _ := insert(testCollection)
	data := testRecord{}
	err := Db.Find(testCollection, key, &data)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.NotEqual(t, data.Key, "")
	assert.NotEqual(t, data.Bar, "")
	assert.NotEqual(t, data.Foo, "")
	Db.Delete(testCollection, key)
}

func TestKeysFiltering(t *testing.T) {
	Db.Open()
	defer Db.Close()

	key11, _ := insert(testCollection)
	key12, _ := insert(testCollection)
	key21, _ := insert(testCollection2)
	key22, _ := insert(testCollection2)

	keys1, err := Db.Keys(testCollection)
	fmt.Println("keys1", keys1)
	assert.Nil(t, err)
	assert.NotNil(t, keys1)
	assert.Equal(t, len(keys1), 2)
	assert.NotEqual(t, keys1[0], "")
	assert.NotEqual(t, keys1[1], "")

	keys2, err := Db.Keys(testCollection2)
	fmt.Println("keys2", keys2)
	assert.Nil(t, err)
	assert.NotNil(t, keys2)
	assert.Equal(t, len(keys2), 2)
	assert.NotEqual(t, keys2[0], "")
	assert.NotEqual(t, keys2[1], "")

	Db.Delete(testCollection, key11)
	Db.Delete(testCollection, key12)
	Db.Delete(testCollection, key21)
	Db.Delete(testCollection, key22)
}

func TestKeys(t *testing.T) {
	Db.Open()
	defer Db.Close()

	key1, _ := insert(testCollection)
	key2, _ := insert(testCollection)

	keys, err := Db.Keys(testCollection)
	fmt.Println("keys", keys)
	assert.Nil(t, err)
	assert.NotNil(t, keys)
	assert.Equal(t, len(keys), 2)
	assert.NotEqual(t, keys[0], "")
	assert.NotEqual(t, keys[1], "")

	Db.Delete(testCollection, key1)
	Db.Delete(testCollection, key2)
}

func TestAll(t *testing.T) {
	Db.Open()
	defer Db.Close()

	key1, _ := insert(testCollection)
	key2, _ := insert(testCollection)

	new := func() proto.Message {
		return new(testRecord)
	}
	data, err := Db.All(testCollection, new)
	fmt.Println("data", data)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 2)

	typedData := make([]testRecord, len(data))
	for i, record := range data {
		typedData[i] = *record.(*testRecord)
		assert.NotNil(t, typedData[i])
		assert.NotEqual(t, typedData[i].Key, "")
		assert.NotEqual(t, typedData[i].Foo, "")
		assert.NotEqual(t, typedData[i].Bar, "")
	}
	fmt.Println("typedData", typedData)
	assert.NotNil(t, typedData)
	assert.Equal(t, len(typedData), 2)
	assert.Equal(t, len(typedData), len(data))
	assert.NotEqual(t, typedData[0].Key, typedData[1].Key)
	assert.NotEqual(t, typedData[0].Bar, typedData[1].Bar)
	assert.NotEqual(t, typedData[0].Foo, typedData[1].Foo)

	Db.Delete(testCollection, key1)
	Db.Delete(testCollection, key2)
}
