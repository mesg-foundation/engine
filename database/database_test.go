package database

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stvp/assert"
)

type testData struct {
	Foo string
	Bar string
}

const collection = "test"

var keyCounter = 0

func insertOne() (key string, err error) {
	keyCounter++
	keyString := strconv.Itoa(keyCounter)
	key = "id_" + keyString
	data := &testData{
		Foo: "hello " + keyString,
		Bar: "bye " + keyString,
	}
	err = Db.Insert(collection, key, data)
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
	data := testData{}
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

	data, err := Db.All(collection)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, len(data), 2)

	typedData := make([]testData, len(data))
	for i, record := range data {
		var recordTyped testData
		err = json.Unmarshal([]byte(record), &recordTyped)
		assert.Nil(t, err)
		typedData[i] = recordTyped
		assert.NotNil(t, recordTyped)
		assert.NotEqual(t, recordTyped.Foo, "")
		assert.NotEqual(t, recordTyped.Bar, "")
	}
	assert.NotNil(t, typedData)
	assert.Equal(t, len(typedData), 2)
	assert.Equal(t, len(typedData), len(data))
}
