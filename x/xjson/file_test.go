// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
