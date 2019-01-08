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

package service

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceLogs(t *testing.T) {
	testDependencyLogs(t, func(s *Service, dependencyKey string) (rstd, rerr io.ReadCloser,
		err error) {
		l, err := s.Logs(dependencyKey)
		require.NoError(t, err)
		require.Len(t, l, 1)
		return l[0].Standard, l[0].Error, nil
	})
}
