package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/require"
)

func TestDefaultValue(t *testing.T) {
	home, _ := homedir.Dir()
	c, err := New()
	require.NoError(t, err)
	require.Equal(t, ":50052", c.Server.Address)
	require.Equal(t, "localhost:50052", c.Client.Address)
	require.Equal(t, "text", c.Log.Format)
	require.Equal(t, "info", c.Log.Level)
	require.Equal(t, filepath.Join(home, ".mesg"), c.Core.RootPath)
	require.Equal(t, "core", c.Core.Name)
	require.Equal(t, "/mesg", c.Docker.Core.Path)
	require.Equal(t, "/var/run/docker.sock", c.Docker.Socket)
	require.True(t, strings.HasPrefix(c.Core.Image, "mesg/core:"))
	require.Equal(t, filepath.Join(home, ".mesg", c.Core.Name), c.CorePath())
}

func TestGlobal(t *testing.T) {
	c, err := Global()
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestLoad(t *testing.T) {
	snapsnot := map[string]string{
		"MESG_SERVER_ADDRESS": "",
		"MESG_CLIENT_ADDRESS": "",
		"MESG_LOG_FORMAT":     "",
		"MESG_LOG_LEVEL":      "",
		"MESG_CORE_IMAGE":     "",
	}
	for key := range snapsnot {
		snapsnot[key] = os.Getenv(key)
	}
	defer func() {
		for key, value := range snapsnot {
			os.Setenv(key, value)
		}
	}()

	os.Setenv("MESG_SERVER_ADDRESS", "test_server_address")
	os.Setenv("MESG_CLIENT_ADDRESS", "test_client_address")
	os.Setenv("MESG_LOG_FORMAT", "test_log_format")
	os.Setenv("MESG_LOG_LEVEL", "test_log_level")
	os.Setenv("MESG_CORE_IMAGE", "test_core_image")
	c, _ := New()
	c.Load()
	require.Equal(t, "test_server_address", c.Server.Address)
	require.Equal(t, "test_client_address", c.Client.Address)
	require.Equal(t, "test_log_format", c.Log.Format)
	require.Equal(t, "test_log_level", c.Log.Level)
	require.Equal(t, "test_core_image", c.Core.Image)
}

func TestValidate(t *testing.T) {
	c, _ := New()
	require.NoError(t, c.Validate())

	c, _ = New()
	c.Log.Format = "wrongValue"
	require.Error(t, c.Validate())

	c, _ = New()
	c.Log.Level = "wrongValue"
	require.Error(t, c.Validate())
}

func TestDaemonEnv(t *testing.T) {
	c, _ := New()
	env := c.DaemonEnv()
	require.Equal(t, c.Server.Address, env["MESG_SERVER_ADDRESS"])
	require.Equal(t, c.Log.Level, env["MESG_LOG_LEVEL"])
	require.Equal(t, c.Log.Format, env["MESG_LOG_FORMAT"])
	require.Equal(t, c.Core.Name, env["MESG_CORE_NAME"])
	require.Equal(t, c.Docker.Core.Path, env["MESG_CORE_ROOTPATH"])
}
