package xjson

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

const testfilenmae = "data.test"

func TestReadFile(t *testing.T) {
	if err := ioutil.WriteFile(testfilenmae, []byte("{}"), os.ModePerm); err != nil {
		t.Fatalf("write error: %s", err)
	}
	defer os.Remove(testfilenmae)

	content, err := ReadFile(testfilenmae)
	if err != nil {
		t.Fatalf("read error: %s", err)
	}

	if string(content) != "{}" {
		t.Fatalf("invalid content - got: %s, want: {}", string(content))
	}
}

func TestReadFileNoFile(t *testing.T) {
	_, err := ReadFile("nodata.test")
	if _, ok := err.(*os.PathError); !ok {
		t.Fatalf("read expected os.PathError - got: %s", err)
	}
}

func TestReadFileInvalidJSON(t *testing.T) {
	if err := ioutil.WriteFile(testfilenmae, []byte("{"), os.ModePerm); err != nil {
		t.Fatalf("write error: %s", err)
	}
	defer os.Remove(testfilenmae)

	_, err := ReadFile(testfilenmae)
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("read expected json.SyntaxError - got: %s", err)
	}
}
