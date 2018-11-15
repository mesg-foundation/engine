package deployer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/database"
	"github.com/mesg-foundation/core/systemservices"
)

func TestNew(t *testing.T) {
	conf, err := config.Global()
	require.NoError(t, err)

	var (
		systemServicesPath = filepath.Join(conf.Core.Path, conf.SystemServices.RelativePath)
		databasePath       = filepath.Join(conf.Core.Path, conf.Core.Database.ServiceRelativePath)
		execDatabasePath   = filepath.Join(conf.Core.Path, conf.Core.Database.ExecutionRelativePath)
	)

	serviceDB, err := database.NewServiceDB(databasePath)
	require.NoError(t, err)

	execDB, err := database.NewExecutionDB(execDatabasePath)
	require.NoError(t, err)

	ss := systemservices.New()

	a, err := api.New(serviceDB, execDB, ss)
	require.NoError(t, err)

	d := New(a, systemServicesPath, ss)
	require.NotZero(t, d)
}
