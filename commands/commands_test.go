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

package commands

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/mesg-foundation/core/commands/mocks"
	"github.com/stretchr/testify/require"
)

// captureStd is helper function that captures Stdout and Stderr and returns function
// that returns standard output and standard error as string.
func captureStd(t *testing.T) func() (stdout string, stderr string) {
	var (
		bufout strings.Builder
		buferr strings.Builder
		wg     sync.WaitGroup

		stdout = os.Stdout
		stderr = os.Stderr
	)

	or, ow, err := os.Pipe()
	require.NoError(t, err)

	er, ew, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = ow
	os.Stderr = ew

	wg.Add(1)
	// copy out and err to buffers
	go func() {
		_, err := io.Copy(&bufout, or)
		require.NoError(t, err)
		or.Close()

		_, err = io.Copy(&buferr, er)
		require.NoError(t, err)
		er.Close()

		wg.Done()
	}()

	return func() (string, string) {
		// close writers and wait for copy to finish
		ow.Close()
		ew.Close()
		wg.Wait()

		// set back oginal stdout and stderr
		os.Stdout = stdout
		os.Stderr = stderr

		// return stdout and stderr
		return bufout.String(), buferr.String()
	}
}

// newMockExecutor returns an Executor mock for testing.
func newMockExecutor() *mocks.Executor {
	return &mocks.Executor{}
}

func TestBaseCommandCmd(t *testing.T) {
	// NOTE: this test is only to satisfy structcheck
	// as it doesn't handle embedded structs yet.
	// It is still very usesful linter so
	// DO NOT REMOVE this test
	_ = baseCmd{}.cmd
}

func readLine(t *testing.T, r *bufio.Reader) []byte {
	line, _, err := r.ReadLine()
	require.NoError(t, err)
	return line
}
