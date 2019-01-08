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

package chunker

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogReader(t *testing.T) {
	var (
		chunk1 = []byte{1}
		chunk2 = []byte{2}
	)

	s := NewStream()

	go func() {
		s.Provide(chunk1)
		s.Provide(chunk2)
		s.Close()
	}()

	data, err := ioutil.ReadAll(s)
	require.NoError(t, err)
	require.Len(t, data, 2)
	require.Equal(t, chunk1, []byte{data[0]})
	require.Equal(t, chunk2, []byte{data[1]})
}
