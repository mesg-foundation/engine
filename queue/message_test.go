package queue

import (
	"testing"

	"github.com/stvp/assert"
)

func TestMessageWithDefaultJSON(t *testing.T) {
	type TestMessageWithDefaultJSONType struct {
		Foo string
		Bar int
	}
	msg, err := message(TestMessageWithDefaultJSONType{
		Foo: "hello",
		Bar: 2,
	})
	assert.Nil(t, err)
	assert.Equal(t, msg.ContentType, "application/json")
	assert.Equal(t, string(msg.Body), "{\"Foo\":\"hello\",\"Bar\":2}")
}

func TestMessageWithJsonType(t *testing.T) {
	type TestMessageWithJSONSpecificType struct {
		Foo string `json:"aa"`
		Bar int    `json:"bb"`
	}
	msg, err := message(TestMessageWithJSONSpecificType{
		Foo: "hello",
		Bar: 2,
	})
	assert.Nil(t, err)
	assert.Equal(t, msg.ContentType, "application/json")
	assert.Equal(t, string(msg.Body), "{\"aa\":\"hello\",\"bb\":2}")
}

func TestMessageWithNonExportedJSON(t *testing.T) {
	type TestMessageWithNonExportedJSONType struct {
		Foo         string
		Bar         int
		nonExported bool
	}
	msg, err := message(TestMessageWithNonExportedJSONType{
		Foo:         "hello",
		Bar:         2,
		nonExported: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, msg.ContentType, "application/json")
	assert.Equal(t, string(msg.Body), "{\"Foo\":\"hello\",\"Bar\":2}")
}
