package mesg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	uuid "github.com/satori/go.uuid"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type serviceDeployer struct {
	Statuses chan DeployStatus

	mesg *MESG
}

// StatusType indicates the type of status message.
type StatusType int

const (
	// RUNNING indicates that status message belongs to an active state.
	RUNNING StatusType = iota + 1

	// DONE indicates that status message belongs to completed state.
	DONE
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    StatusType
}

func newServiceDeployer(mesg *MESG) *serviceDeployer {
	return &serviceDeployer{
		mesg: mesg,
	}
}

// FromGitURL deploys a service hosted at a Git url.
func (d *serviceDeployer) FromGitURL(url string) (*service.Service, *importer.ValidationError, error) {
	d.sendStatus("Downloading service...", RUNNING)
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	if err := d.gitClone(url, path); err != nil {
		return nil, nil, err
	}
	d.sendStatus(fmt.Sprintf("%s Service downloaded with success.", aurora.Green("✔")), DONE)
	return d.deploy(path)
}

// FromGzippedTar deploys a service from a gzipped tarball.
func (d *serviceDeployer) FromGzippedTar(r io.Reader) (*service.Service, *importer.ValidationError, error) {
	d.sendStatus("Sending service context to core daemon...", RUNNING)
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	if err := archive.Untar(r, path, &archive.TarOptions{
		Compression: archive.Gzip,
	}); err != nil {
		return nil, nil, err
	}
	d.sendStatus(fmt.Sprintf("%s Service context sent to core daemon with success.", aurora.Green("✔")), DONE)
	return d.deploy(path)
}

// gitClone clones a repo hosted at repoURL to path.
func (d *serviceDeployer) gitClone(repoURL string, path string) error {
	u, err := url.Parse(repoURL)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}
	options := &git.CloneOptions{}
	if u.Fragment != "" {
		options.ReferenceName = plumbing.ReferenceName("refs/heads/" + u.Fragment)
		u.Fragment = ""
	}
	options.URL = u.String()
	_, err = git.PlainClone(path, false, options)
	return err
}

// deploy deploys a service in path.
func (d *serviceDeployer) deploy(path string) (*service.Service, *importer.ValidationError, error) {
	defer os.RemoveAll(path)

	service, err := importer.From(path)
	validationErr, err := d.assertValidationError(err)
	if err != nil {
		return nil, nil, err
	}
	if validationErr != nil {
		return nil, validationErr, nil
	}

	d.sendStatus("Building Docker image...", RUNNING)
	imageHash, err := d.mesg.container.Build(path)
	if err != nil {
		return nil, nil, err
	}
	if _, err := os.Stat(filepath.Join(path, ".mesgignore")); err == nil {
		// TODO: remove for a future release
		d.sendStatus(fmt.Sprintf("%s [DEPRECATED] Please use .dockerignore instead of .mesgignore", aurora.Red("⨯")), DONE)
	}
	d.sendStatus(fmt.Sprintf("%s Image built with success.", aurora.Green("✔")), DONE)
	d.injectConfigurationInDependencies(service, imageHash)

	d.sendStatus(fmt.Sprintf("%s Completed.", aurora.Green("✔")), DONE)
	return service, nil, services.Save(service)
}

func (d *serviceDeployer) injectConfigurationInDependencies(s *service.Service, imageHash string) {
	config := s.Configuration
	if config == nil {
		config = &service.Dependency{}
	}
	dependency := &service.Dependency{
		Command:     config.Command,
		Ports:       config.Ports,
		Volumes:     config.Volumes,
		Volumesfrom: config.Volumesfrom,
		Image:       imageHash,
	}
	if s.Dependencies == nil {
		s.Dependencies = make(map[string]*service.Dependency)
	}
	s.Dependencies["service"] = dependency
}

func (d *serviceDeployer) createTempDir() (path string, err error) {
	return ioutil.TempDir("", "mesg-"+uuid.NewV4().String())
}

// sendStatus sends a status message.
func (d *serviceDeployer) sendStatus(message string, typ StatusType) {
	if d.Statuses != nil {
		d.Statuses <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}

func (d *serviceDeployer) assertValidationError(err error) (*importer.ValidationError, error) {
	if err == nil {
		return nil, nil
	}
	if validationError, ok := err.(*importer.ValidationError); ok {
		return validationError, nil
	}
	return nil, err
}
