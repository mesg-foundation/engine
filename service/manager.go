package service

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"sync"

	"github.com/mesg-foundation/core/config"
	"github.com/mesg-foundation/core/container"
	"github.com/mesg-foundation/core/utils/dirhash"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xnet"
	"github.com/mesg-foundation/core/x/xos"
	"github.com/mesg-foundation/core/x/xstrings"
)

// Manager is an interface that defines every service manager.
type Manager interface {
	// Deploy deploys the service according to specific manager.
	Deploy(service *Service, contextDir string, env map[string]string) error
	// Start runs the service and all its dependencies.
	Start(service *Service) error
	// Stop stops the proviously started service and all its Dependencies.
	Stop(service *Service) error
	// Delete deletes all persisent data used by service.
	Delete(service *Service) error
	// Status returns the service running status.
	Status(service *Service) error
	// Logs returns logs for service and its dependencies.
	Logs(service *Service, dependencies []string) ([]*LogReader, error)
}

var containerToServiceStatus = map[container.StatusType]Status{
	container.UNKNOWN:  StatusUnknown,
	container.STOPPED:  StatusStopped,
	container.STARTING: StatusStarting,
	container.RUNNING:  StatusRunning,
}

// ContainerManager is service manager that uses docker container
// as platfor for deploy, start, ... services.
type ContainerManager struct {
	c   container.Container
	cfg *config.Config
}

// NewContainerManager returns new container manager.
func NewContainerManager(c container.Container, cfg *config.Config) *ContainerManager {
	return &ContainerManager{c: c, cfg: cfg}
}

// Deploy see Manager.Deploy.
func (m *ContainerManager) Deploy(service *Service, contextDir string, env map[string]string) error {
	if err := service.validateConfigurationEnv(env); err != nil {
		return err
	}

	hash, err := serviceHash(contextDir, env)
	if err != nil {
		return err
	}
	service.Hash = hex.EncodeToString(hash)

	// replace default env with new one.
	defenv := xos.EnvSliceToMap(service.Configuration.Env)
	service.Configuration.Env = xos.EnvMapToSlice(xos.EnvMergeMaps(defenv, env))

	image, err := m.c.Build(contextDir)
	if err != nil {
		return err
	}

	service.Configuration.Image = image
	return nil
}

// Start see Manager.Start.
func (m *ContainerManager) Start(service *Service) error {
	if err := m.Status(service); err != nil {
		return err
	}
	if service.Status == StatusRunning {
		return nil
	}
	if service.Status == StatusPartial {
		if err := m.Stop(service); err != nil {
			return err
		}
	}

	networkID, err := m.c.CreateNetwork(service.namespace())
	if err != nil {
		return err
	}
	sharedNetworkID, err := m.c.SharedNetworkID()
	if err != nil {
		return err
	}
	_, port, _ := xnet.SplitHostPort(m.cfg.Server.Address)
	endpoint := m.cfg.Core.Name + ":" + strconv.Itoa(port)
	start := func(dep *Dependency, namespace []string) error {
		_, err := m.c.StartService(container.ServiceOptions{
			Namespace: namespace,
			Labels: map[string]string{
				"mesg.service": service.Name,
				"mesg.hash":    service.Hash,
				"mesg.sid":     service.Sid,
				"mesg.core":    m.cfg.Core.Name,
			},
			Image:   dep.Image,
			Args:    dep.Args,
			Command: dep.Command,
			Env: xos.EnvMergeSlices(dep.Env, []string{
				"MESG_TOKEN=" + service.Hash,
				"MESG_ENDPOINT=" + endpoint,
				"MESG_ENDPOINT_TCP=" + endpoint,
			}),
			Mounts: service.volumes(dep.Key),
			Ports:  dep.ports(),
			Networks: []container.Network{
				{ID: networkID, Alias: dep.Key},
				{ID: sharedNetworkID},
			},
		})
		return err
	}

	if err := start(service.Configuration, service.namespace()); err != nil {
		return err
	}

	// BUG: https://github.com/mesg-foundation/core/issues/382
	// After solving this by docker, switch back to deploy in parallel
	for _, dep := range service.Dependencies {
		if err := start(dep, depNamespace(service.Hash, dep.Key)); err != nil {
			return err
		}

		if err != nil {
			m.Stop(service)
			return err
		}
	}
	service.Status = StatusRunning
	return nil
}

// Stop see Manager.Stop.
func (m *ContainerManager) Stop(service *Service) error {
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors

		stop = func(namespace []string) {
			defer wg.Done()
			if err := m.c.StopService(namespace); err != nil {
				errs.Append(err)
			}
		}
	)

	wg.Add(len(service.Dependencies) + 1)
	for _, dep := range service.Dependencies {
		go stop(depNamespace(service.Hash, dep.Key))
	}
	go stop(service.namespace())
	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}
	if err := m.c.DeleteNetwork(service.namespace(), container.EventDestroy); err != nil {
		return err
	}
	service.Status = StatusStopped
	return nil
}

// Delete see Manager.Delete.
func (m *ContainerManager) Delete(service *Service) error {
	var (
		wg      sync.WaitGroup
		errs    xerrors.SyncErrors
		sources = make(map[string]bool)
	)

	// make map of all volumes, because it may occur more then
	// once if volumesFrom were passed.
	for _, dep := range service.Dependencies {
		for _, volume := range service.volumes(dep.Key) {
			sources[volume.Source] = true
		}
	}
	for _, volume := range service.volumes(mainServiceKey) {
		sources[volume.Source] = true
	}

	wg.Add(len(sources))
	for source := range sources {
		go func(source string) {
			defer wg.Done()
			if err := m.c.DeleteVolume(source); err != nil {
				errs.Append(err)
			}
		}(source)
	}
	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}
	service.Status = StatusDeleted
	return nil
}

// Status see Manager.Status.
func (m *ContainerManager) Status(service *Service) error {
	var (
		wg   sync.WaitGroup
		errs xerrors.SyncErrors

		statuses = make(map[container.StatusType]bool)
		status   = func(namespace []string) {
			defer wg.Done()
			status, err := m.c.Status(namespace)
			if err != nil {
				errs.Append(err)
			}
			statuses[status] = true
		}
	)

	wg.Add(len(service.Dependencies) + 1)
	go status(service.namespace())
	for _, dep := range service.Dependencies {
		go status(depNamespace(service.Hash, dep.Key))
	}
	wg.Wait()

	if err := errs.ErrorOrNil(); err != nil {
		return err
	}

	service.Status = pickStatus(statuses)
	return nil
}

// Logs see Manager.Logs.
func (m *ContainerManager) Logs(service *Service, dependencies []string) ([]*LogReader, error) {
	var (
		lrs []*LogReader
		all = len(dependencies) == 0
	)

	if all || xstrings.SliceContains(dependencies, mainServiceKey) {
		r, err := m.c.ServiceLogs(service.namespace())
		if err != nil {
			return nil, err
		}
		lrs = append(lrs, &LogReader{key: mainServiceKey, r: r})
	}

	for _, dep := range service.Dependencies {
		if !all && !xstrings.SliceContains(dependencies, dep.Key) {
			continue
		}
		r, err := m.c.ServiceLogs(depNamespace(service.Hash, dep.Key))
		if err != nil {
			return nil, err
		}
		lrs = append(lrs, &LogReader{
			key: dep.Key,
			r:   r,
		})
	}
	return lrs, nil
}

func pickStatus(statuses map[container.StatusType]bool) Status {
	if len(statuses) == 1 {
		for status := range statuses {
			return containerToServiceStatus[status]
		}
	}
	if statuses[container.UNKNOWN] {
		return StatusUnknown
	}
	return StatusPartial
}

func serviceHash(contextDir string, env map[string]string) ([]byte, error) {
	dh := dirhash.New(contextDir)
	envbytes := []byte(xos.EnvMapToString(env))
	return dh.Sum(envbytes)
}

func depNamespace(hash, key string) []string {
	sum := sha1.Sum([]byte(hash + key))
	return []string{hex.EncodeToString(sum[:])}
}
