package service

import (
	"bytes"
	"encoding/hex"
	"io/ioutil"

	"github.com/mesg-foundation/core/x/xos"
	"github.com/mesg-foundation/core/x/xstrings"
)

// FakeManager is dummy manager for tests purpose.
type FakeManager struct {
	services []*Service
}

// NewFakeManager returns fake manager that only tracks status of the services
// without deploys, starts or deletes them.
func NewFakeManager() *FakeManager {
	return &FakeManager{}
}

// Deploy see Manager.Deploy.
func (m *FakeManager) Deploy(service *Service, contextDir string, env map[string]string) error {
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
	m.update(service, StatusStopped)
	return nil
}

// Start see Manager.Start.
func (m *FakeManager) Start(service *Service) error {
	m.update(service, StatusRunning)
	return nil
}

// Stop see Manager.Stop.
func (m *FakeManager) Stop(service *Service) error {
	m.update(service, StatusStopped)
	return nil
}

// Delete see Manager.Delete.
func (m *FakeManager) Delete(service *Service) error {
	m.update(service, StatusDeleted)
	return nil
}

// Status see Manager.Status.
func (m *FakeManager) Status(service *Service) error {
	for i := range m.services {
		if m.services[i].Hash == service.Hash {
			service.Status = m.services[i].Status
		}
	}
	return nil
}

// Logs see Manager.Logs.
func (m *FakeManager) Logs(service *Service, dependencies []string) ([]*LogReader, error) {
	var (
		lrs []*LogReader
		all = len(dependencies) == 0
	)

	if all || xstrings.SliceContains(dependencies, MainServiceKey) {
		lrs = append(lrs, &LogReader{
			key: MainServiceKey,
			r:   ioutil.NopCloser(&bytes.Buffer{}),
		})
	}

	for _, dep := range service.Dependencies {
		if !all || !xstrings.SliceContains(dependencies, dep.Key) {
			continue
		}
		lrs = append(lrs, &LogReader{
			key: dep.Key,
			r:   ioutil.NopCloser(&bytes.Buffer{}),
		})
	}
	return lrs, nil
}

// update set service status and add service to tracking if wasn't add before.
func (m *FakeManager) update(service *Service, status Status) {
	service.Status = status

	for i := range m.services {
		if m.services[i].Hash == service.Hash {
			m.services[i].Status = status
			return
		}
	}
	m.services = append(m.services, service)
}
