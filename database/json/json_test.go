package json

import (
	"encoding/json"
	"testing"

	"github.com/stvp/assert"
)

type TestData struct {
	Foo string
	Bar string
}

const collection = "test"

var Db = &Database{}

func TestOpenClose(t *testing.T) {
	err := Db.Open()
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

	data := &TestData{
		Foo: "hello",
		Bar: "bye",
	}
	err := Db.Insert(collection, "id_2", data)
	assert.Nil(t, err)
}

func TestFind(t *testing.T) {
	Db.Open()
	defer Db.Close()

	data := new(TestData)
	err := Db.Find(collection, "id_1", data)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.NotEqual(t, data.Bar, "")
	assert.NotEqual(t, data.Foo, "")
}

func TestAll(t *testing.T) {
	Db.Open()
	defer Db.Close()

	data, err := Db.All(collection)
	assert.Nil(t, err)

	typedData := make([]TestData, len(data))
	for i, record := range data {
		var recordTyped TestData
		err = json.Unmarshal([]byte(record), &recordTyped)
		assert.Nil(t, err)
		typedData[i] = recordTyped
	}

	assert.NotNil(t, typedData)
	assert.NotEqual(t, len(typedData), 0)
	assert.Equal(t, len(typedData), len(data))
}
