package service

import (
	"errors"
	"strconv"
	"testing"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/container/mocks"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/stretchr/testify/require"
)

func TestExtractPortEmpty(t *testing.T) {
	dep := Dependency{}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 0)
}

func TestExtractPorts(t *testing.T) {
	dep := &Dependency{
		Ports: []string{
			"80",
			"3000:8080",
		},
	}
	ports := dep.extractPorts()
	require.Equal(t, len(ports), 2)
	require.Equal(t, ports[0].Target, uint32(80))
	require.Equal(t, ports[0].Published, uint32(80))
	require.Equal(t, ports[1].Target, uint32(8080))
	require.Equal(t, ports[1].Published, uint32(3000))
}

func TestExtractVolumes(t *testing.T) {
	s := &Service{
		Dependencies: []*Dependency{{
			Key:     "test",
			Volumes: []string{"foo", "bar"},
		}},
	}
	volumes := s.Dependencies[0].extractVolumes(s)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(s, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(s, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)

	s = &Service{
		Dependencies: []*Dependency{{
			VolumesFrom: []string{"test"},
		}},
	}
	_, err := s.Dependencies[0].extractVolumesFrom(s)
	require.Error(t, err)

	s = &Service{
		Dependencies: []*Dependency{
			{
				Key:     "test",
				Volumes: []string{"foo", "bar"},
			},
			{
				VolumesFrom: []string{"test"},
			}},
	}
	volumes, err = s.Dependencies[1].extractVolumesFrom(s)
	require.NoError(t, err)
	require.Len(t, volumes, 2)
	require.Equal(t, volumeKey(s, "test", "foo"), volumes[0].Source)
	require.Equal(t, "foo", volumes[0].Target)
	require.Equal(t, false, volumes[0].Bind)
	require.Equal(t, volumeKey(s, "test", "bar"), volumes[1].Source)
	require.Equal(t, "bar", volumes[1].Target)
	require.Equal(t, false, volumes[1].Bind)
}

func TestStartService(t *testing.T) {
	var (
		containerServiceID = "1"
		dependencyKey      = "2"
		serviceName        = "TestStartService"
		Sid                = "SidStartService"
		networkID          = "3"
		sharedNetworkID    = "4"
		s                  = &Service{
			Name: serviceName,
			Sid:  Sid,
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.STOPPED, nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Once().Return(sharedNetworkID, nil)
	mockStartService(s, d, mc, networkID, sharedNetworkID, containerServiceID, nil)

	serviceIDs, err := s.Start(mc)
	require.NoError(t, err)
	require.Len(t, serviceIDs, 1)
	require.Equal(t, containerServiceID, serviceIDs[0])

	mc.AssertExpectations(t)
}

func TestStartWith2Dependencies(t *testing.T) {
	var (
		containerServiceIDs = []string{"1", "2"}
		dependencyKey       = "3"
		dependencyKey2      = "4"
		dependencyImage     = "5"
		dependencyImage2    = "6"
		networkID           = "7"
		sharedNetworkID     = "8"
		serviceName         = "TestStartWith2Dependencies"
		s                   = &Service{
			Name: serviceName,
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: dependencyImage,
				},
				{
					Key:   dependencyKey2,
					Image: dependencyImage2,
				},
			},
		}
		mc = &mocks.Container{}
	)

	var (
		d, _  = s.getDependency(dependencyKey)
		d2, _ = s.getDependency(dependencyKey2)
		ds    = []*Dependency{d, d2}
	)

	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.STOPPED, nil)
	mc.On("Status", d2.namespace(s.namespace())).Once().Return(container.STOPPED, nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Twice().Return(sharedNetworkID, nil)

	for i, d := range ds {
		mockStartService(s, d, mc, networkID, sharedNetworkID, containerServiceIDs[i], nil)
	}

	serviceIDs, err := s.Start(mc)
	require.NoError(t, err)
	require.Len(t, serviceIDs, len(s.Dependencies))

	for i := range ds {
		require.Contains(t, serviceIDs, containerServiceIDs[i])
	}

	mc.AssertExpectations(t)
}

func TestStartServiceRunning(t *testing.T) {
	var (
		dependencyKey = "1"
		s             = &Service{
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "2",
				},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)
	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.RUNNING, nil)

	dockerServices, err := s.Start(mc)
	require.NoError(t, err)
	require.Len(t, dockerServices, 0)

	mc.AssertExpectations(t)
}

func TestPartiallyRunningService(t *testing.T) {
	var (
		dependencyKey       = "1"
		dependencyKey2      = "2"
		networkID           = "3"
		sharedNetworkID     = "4"
		containerServiceIDs = []string{"5", "6"}
		s                   = &Service{
			Name: "TestPartiallyRunningService",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
				{
					Key:   dependencyKey2,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
	)

	var (
		d, _  = s.getDependency(dependencyKey)
		d2, _ = s.getDependency(dependencyKey2)
		ds    = []*Dependency{d, d2}
	)

	mc.On("Status", d.namespace(s.namespace())).Return(container.STOPPED, nil)
	mc.On("StopService", d.namespace(s.namespace())).Once().Return(nil)
	mc.On("Status", d2.namespace(s.namespace())).Return(container.RUNNING, nil)
	mc.On("StopService", d2.namespace(s.namespace())).Once().Return(nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Twice().Return(sharedNetworkID, nil)

	for i, d := range ds {
		mockStartService(s, d, mc, networkID, sharedNetworkID, containerServiceIDs[i], nil)
	}

	serviceIDs, err := s.Start(mc)
	require.NoError(t, err)
	require.Len(t, serviceIDs, len(s.Dependencies))

	for i := range ds {
		require.Contains(t, serviceIDs, containerServiceIDs[i])
	}

	mc.AssertExpectations(t)
}

func TestStartDependency(t *testing.T) {
	var (
		dependencyKey      = "1"
		networkID          = "3"
		sharedNetworkID    = "4"
		containerServiceID = "5"
		s                  = &Service{
			Name: "TestStartDependency",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("SharedNetworkID").Once().Return(sharedNetworkID, nil)
	mockStartService(s, d, mc, networkID, sharedNetworkID, containerServiceID, nil)

	serviceID, err := d.Start(mc, s, networkID)
	require.NoError(t, err)
	require.Equal(t, containerServiceID, serviceID)

	mc.AssertExpectations(t)
}

func TestServiceStartError(t *testing.T) {
	var (
		dependencyKey   = "1"
		networkID       = "3"
		sharedNetworkID = "4"
		startErr        = errors.New("ops")
		s               = &Service{
			Name: "TestNetworkCreated",
			Dependencies: []*Dependency{
				{
					Key:   dependencyKey,
					Image: "http-server",
				},
			},
		}
		mc = &mocks.Container{}
	)

	d, _ := s.getDependency(dependencyKey)

	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.STOPPED, nil)
	mc.On("CreateNetwork", s.namespace()).Once().Return(networkID, nil)
	mc.On("SharedNetworkID").Once().Return(sharedNetworkID, nil)
	mockStartService(s, d, mc, networkID, sharedNetworkID, "", startErr)
	mc.On("Status", d.namespace(s.namespace())).Once().Return(container.STOPPED, nil)

	serviceIDs, err := s.Start(mc)
	require.Equal(t, startErr, err)
	require.Len(t, serviceIDs, 0)

	mc.AssertExpectations(t)
}

func mockStartService(s *Service, d *Dependency, mc *mocks.Container,
	networkID, sharedNetworkID, containerServiceID string, err error) {
	var (
		c, _           = config.Global()
		_, port, _     = xnet.SplitHostPort(c.Server.Address)
		endpoint       = c.Core.Name + ":" + strconv.Itoa(port)
		volumes        = d.extractVolumes(s)
		volumesFrom, _ = d.extractVolumesFrom(s)
	)
	mc.On("StartService", container.ServiceOptions{
		Namespace: d.namespace(s.namespace()),
		Labels: map[string]string{
			"mesg.core":    c.Core.Name,
			"mesg.sid":     s.Sid,
			"mesg.service": s.Name,
			"mesg.hash":    s.Hash,
		},
		Image:   d.Image,
		Command: d.Command,
		Args:    d.Args,
		Env: []string{
			"MESG_TOKEN=" + s.Hash,
			"MESG_ENDPOINT=" + endpoint,
			"MESG_ENDPOINT_TCP=" + endpoint,
		},
		Mounts: append(volumes, volumesFrom...),
		Ports:  d.extractPorts(),
		Networks: []container.Network{
			{ID: networkID, Alias: d.Key},
			{ID: sharedNetworkID},
		},
	}).Once().Return(containerServiceID, err)
}
