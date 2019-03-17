package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xgit"
)

// StatusType indicates the type of status message.
type StatusType int

// Deploy Statuses.
const (
	Running StatusType = iota
	Success
	Failed
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    StatusType
}

// DeployService deploys a service from a gzipped tarball.
func (api *API) DeployService(r io.Reader, env map[string]string, statusC chan DeployStatus) (*service.Service, error) {
	contextDir, err := ioutil.TempDir("", "mesg-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(contextDir)

	if err := archive.Untar(r, contextDir, nil); err != nil {
		return nil, err
	}
	return api.deploy(contextDir, env)
}

// DeployServiceFromURL deploys a service living at a Git host.
// Supported URL types:
// - https://github.com/mesg-foundation/service-ethereum
// - https://github.com/mesg-foundation/service-ethereum#branchName
func (api *API) DeployServiceFromURL(url string, env map[string]string, statusC chan DeployStatus) (*service.Service, error) {
	contextDir, err := ioutil.TempDir("", "mesg-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(contextDir)

	sendStatus(statusC, "Downloading service...", Running)
	if xgit.IsGitURL(url) {
		if err := xgit.Clone(url, contextDir); err != nil {
			return nil, err
		}
	} else {
		// if not git repo then it must be tarball
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if err := archive.Untar(resp.Body, contextDir, nil); err != nil {
			return nil, err
		}
	}
	sendStatus(statusC, "Service downloaded with success", Success)

	return api.deploy(contextDir, env)
}

// deploy deploys a service in path.
func (api *API) deploy(contextDir string, env map[string]string) (*service.Service, error) {
	var err error
	contextDir, err = formalizeContextDir(contextDir)
	if err != nil {
		return nil, err
	}

	s, err := service.ReadDefinition(contextDir)
	if err != nil {
		return nil, err
	}
	if err := api.sm.Deploy(s, contextDir, env); err != nil {
		return nil, err
	}
	if err := api.db.Save(s); err != nil {
		return nil, err
	}
	return s, nil
}

func formalizeContextDir(contextDir string) (string, error) {
	// NOTE: remove .git folder from repo.
	// It makes docker build image id same between repo clones.
	if err := os.RemoveAll(filepath.Join(contextDir, ".git")); err != nil {
		return "", err
	}

	// NOTE: if there is only one directory inside service context enter it.
	dirs, err := ioutil.ReadDir(contextDir)
	if err != nil {
		return "", err
	}
	if len(dirs) == 1 && dirs[0].IsDir() {
		contextDir = filepath.Join(contextDir, dirs[0].Name())
	}
	return contextDir, nil
}

func sendStatus(statusC chan DeployStatus, message string, typ StatusType) {
	if statusC != nil {
		statusC <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}
