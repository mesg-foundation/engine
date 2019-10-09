package container

import (
	"context"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/require"
)

const nstestprefix = "enginetest"

func TestBuild(t *testing.T) {
	c, err := New(nstestprefix)
	require.NoError(t, err)
	tag, err := c.Build("testdata/test-image")
	require.NoError(t, err)
	require.NotEmpty(t, tag)
}

func TestNetwork(t *testing.T) {
	const netname = "test-net"

	c, err := New(nstestprefix)
	require.NoError(t, err)
	defer c.Cleanup()

	t.Run("shared-network", func(t *testing.T) {
		// check if is set after New
		require.NotEmpty(t, c.SharedNetworkID())

		// create 2nd time - 1st time is created with New
		require.NoError(t, c.createSharedNetwork())
	})

	t.Run("create", func(t *testing.T) {
		nid, err := c.CreateNetwork(netname)
		require.NoError(t, err)
		require.NotEmpty(t, nid)
	})

	t.Run("delete", func(t *testing.T) {
		require.NoError(t, c.DeleteNetwork(netname))
		_, err := c.client.NetworkInspect(context.Background(), c.namespace(netname), types.NetworkInspectOptions{})
		require.Error(t, err)
	})
}

func TestService(t *testing.T) {
	const servicename = "test-service"

	c, err := New(nstestprefix)
	require.NoError(t, err)
	defer c.Cleanup()

	t.Run("start", func(t *testing.T) {
		id, err := c.StartService(ServiceOptions{
			Image:     "busybox",
			Namespace: servicename,
			Ports: []Port{
				{
					Target:    50053,
					Published: 50053,
				},
				{
					Target:    50054,
					Published: 50054,
				},
			},
			Mounts: []Mount{
				{
					Source: "testdata",
					Target: "/testdata",
				},
			},
			Env:     []string{"foo=bar"},
			Args:    []string{"hello"},
			Command: "echo",
			Networks: []Network{
				{
					ID:    c.SharedNetworkID(),
					Alias: "test-net",
				},
			},
			Labels: map[string]string{"label": "test"},
		})
		require.NoError(t, err)

		resp, _, err := c.client.ServiceInspectWithRaw(context.Background(), id, types.ServiceInspectOptions{})
		require.NoError(t, err)

		require.Equal(t, "test", resp.Spec.Labels["label"])
		require.Equal(t, c.namespace(servicename), resp.Spec.Labels["com.docker.stack.namespace"])
		require.Equal(t, "busybox", resp.Spec.Labels["com.docker.stack.image"])
		require.Len(t, resp.Spec.EndpointSpec.Ports, 2)
		require.Equal(t, resp.Spec.EndpointSpec.Ports[0].TargetPort, uint32(50053))
		require.Equal(t, resp.Spec.EndpointSpec.Ports[0].PublishedPort, uint32(50053))
		require.Equal(t, resp.Spec.EndpointSpec.Ports[1].TargetPort, uint32(50054))
		require.Equal(t, resp.Spec.EndpointSpec.Ports[1].PublishedPort, uint32(50054))
		require.Len(t, resp.Spec.TaskTemplate.Networks, 1)
		require.Equal(t, resp.Spec.TaskTemplate.Networks[0].Aliases, []string{"test-net"})
		require.Len(t, resp.Spec.TaskTemplate.ContainerSpec.Mounts, 1)
		require.Equal(t, resp.Spec.TaskTemplate.ContainerSpec.Mounts[0].Target, "/testdata")
		require.Equal(t, resp.Spec.TaskTemplate.ContainerSpec.Mounts[0].Source, "testdata")
		require.Equal(t, resp.Spec.TaskTemplate.ContainerSpec.Command, []string{"echo"})
		require.Equal(t, resp.Spec.TaskTemplate.ContainerSpec.Args, []string{"hello"})
		require.Equal(t, resp.Spec.TaskTemplate.ContainerSpec.Env, []string{"foo=bar"})
	})

	t.Run("stop", func(t *testing.T) {
		require.NoError(t, c.StopService(servicename))
	})
}

func TestNamespace(t *testing.T) {
	c, _ := New("engine")
	require.Equal(t, c.namespace("foo"), "engine-foo")
}

func TestReadDockerIgnoreFile(t *testing.T) {
	// non existing .dockerignore file
	files, err := readDockerIgnoreFile("testdata/test-image")
	require.NoError(t, err)
	require.Empty(t, files)

	// existing .dockerignore file
	files, err = readDockerIgnoreFile("testdata/test-image-with-ignore")
	require.NoError(t, err)
	require.Equal(t, []string{"file.txt"}, files)
}

func TestTagFromResponse(t *testing.T) {
	var tests = []struct {
		name string
		resp string
		tag  string
		err  bool
	}{
		{
			"ok",
			`{"stream":"sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`,
			"sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			false,
		},
		{
			"empty response",
			"",
			"",
			true,
		},
		{
			"invalid json",
			`-`,
			"",
			true,
		},
		{
			"no sha256 prefix",
			`{"stream":"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"}`,
			"",
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tag, err := tagFromResponse(strings.NewReader(tt.resp))
			if tt.err && err == nil {
				t.Errorf("want error")
			}
			if tt.tag != tag {
				t.Errorf("invalid tag: want %s, got %s", tt.tag, tag)
			}
		})
	}
}
