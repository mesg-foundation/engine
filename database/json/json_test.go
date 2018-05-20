package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	"github.com/mesg-foundation/core/database"
	"github.com/stvp/assert"
)

type testRecord struct {
	key string
	Foo string
	Bar string
}

func (record *testRecord) Key() string {
	return record.key
}

func (record *testRecord) Encode() (bytes []byte, err error) {
	bytes, err = json.Marshal(record)
	return
}

func (record *testRecord) Decode(bytes []byte) (err error) {
	err = json.Unmarshal(bytes, &record)
	return
}

const collection = "test"

var Db = &Database{}

var keyCounter = 0

func insertOne() (key string, err error) {
	keyCounter++
	keyString := strconv.Itoa(keyCounter)
	key = "id_" + keyString
	data := testRecord{
		key: key,
		Foo: "hello " + keyString,
		Bar: "bye " + keyString,
	}
	err = Db.Insert(collection, key, &data)
	return
}

func delete(key string) (err error) {
	err = Db.Delete(collection, key)
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
	key, err := insertOne()
	defer delete(key)
	assert.Nil(t, err)
}

func TestDelete(t *testing.T) {
	Db.Open()
	defer Db.Close()
	key, _ := insertOne()
	err := delete(key)
	assert.Nil(t, err)
}

func TestFind(t *testing.T) {
	Db.Open()
	defer Db.Close()
	key, _ := insertOne()
	defer delete(key)
	data := testRecord{}
	err := Db.Find(collection, key, &data)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.NotEqual(t, data.Bar, "")
	assert.NotEqual(t, data.Foo, "")
}

func TestAll(t *testing.T) {
	Db.Open()
	defer Db.Close()

	key1, _ := insertOne()
	defer delete(key1)
	key2, _ := insertOne()
	defer delete(key2)

	new := func() database.Record {
		return new(testRecord)
	}
	data, err := Db.All(collection, new)
	fmt.Println("data", data)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 2)

	typedData := make([]testRecord, len(data))
	for i, record := range data {
		typedData[i] = *record.(*testRecord)
		assert.NotNil(t, typedData[i])
		assert.NotEqual(t, typedData[i].Foo, "")
		assert.NotEqual(t, typedData[i].Bar, "")
	}
	fmt.Println("typedData", typedData)
	assert.NotNil(t, typedData)
	assert.Equal(t, len(typedData), 2)
	assert.Equal(t, len(typedData), len(data))
	assert.NotEqual(t, typedData[0].Bar, typedData[1].Bar)
	assert.NotEqual(t, typedData[0].Foo, typedData[1].Foo)
}
