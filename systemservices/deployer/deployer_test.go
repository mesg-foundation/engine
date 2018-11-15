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
	expectedErr := &systemServiceNotFoundError{name: "noExisting"}
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
