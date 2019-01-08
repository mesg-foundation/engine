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

package deployer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/systemservices"
)

func newDeployer(t *testing.T) (*Deployer, func()) {
	conf, err := config.Global()
	require.NoError(t, err)

	var (
		systemServicesPath  = filepath.Join(conf.Core.Path, conf.SystemServices.RelativePath)
		serviceDatabasePath = "service.db.test"
		execDatabasePath    = "exec.db.test"
	)

	serviceDB, err := database.NewServiceDB(serviceDatabasePath)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execDatabasePath)
	require.NoError(t, err)

	ss := systemservices.New()

	a, err := api.New(serviceDB, execDB, ss)
	require.NoError(t, err)

	d := New(a, systemServicesPath, ss)
	require.NotZero(t, d)

	closer := func() {
		require.NoError(t, serviceDB.Close())
		require.NoError(t, execDB.Close())
		require.NoError(t, os.RemoveAll(serviceDatabasePath))
		require.NoError(t, os.RemoveAll(execDatabasePath))
	}

	return d, closer
}

func TestNotExisting(t *testing.T) {
	d, closer := newDeployer(t)
	defer closer()
	err := d.Deploy([]string{"noExisting"})
	expectedErr := &systemservices.SystemServiceNotFoundError{Name: "noExisting"}
	require.EqualError(t, err, expectedErr.Error())
}

// TODO: Should have a mock on API in order to convert this integration test to unit test
func TestDeploy(t *testing.T) {
	d, closer := newDeployer(t)
	defer closer()
	err := d.Deploy([]string{systemservices.ResolverService})
	// TODO: should stop the deployed ss
	require.NoError(t, err)
}
