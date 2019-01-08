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

package api

import (
	"os"
	"testing"

	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/dockertest"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/systemservices"
	"github.com/stretchr/testify/require"
)

const (
	servicedbname = "service.db.test"
	execdbname    = "exec.db.test"
)

func newAPIAndDockerTest(t *testing.T) (*API, *dockertest.Testing, func()) {
	dt := dockertest.New()

	container, err := container.New(container.ClientOption(dt.Client()))
	require.NoError(t, err)

	db, err := database.NewServiceDB(servicedbname)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execdbname)
	require.NoError(t, err)

	ss := systemservices.New()

	a, err := New(db, execDB, ss, ContainerOption(container))
	require.NoError(t, err)

	closer := func() {
		require.NoError(t, db.Close())
		require.NoError(t, execDB.Close())
		require.NoError(t, os.RemoveAll(servicedbname))
		require.NoError(t, os.RemoveAll(execdbname))
	}
	return a, dt, closer
}
