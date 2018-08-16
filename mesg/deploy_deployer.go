package mesg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"

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
	Statuses chan string

	mesg *MESG
}

func newServiceDeployer(mesg *MESG) *serviceDeployer {
	return &serviceDeployer{
		mesg: mesg,
	}
}

func (d *serviceDeployer) FromGitURL(url string) (*service.Service, *importer.ValidationError, error) {
	d.sendStatus("Downloading service...")
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	if err := d.gitClone(url, path); err != nil {
		return nil, nil, err
	}
	return d.deploy(path)
}

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

func (d *serviceDeployer) FromGzippedTar(r io.Reader) (*service.Service, *importer.ValidationError, error) {
	d.sendStatus("Sending service context to core daemon...")
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	if err := archive.Untar(r, path, &archive.TarOptions{
		Compression: archive.Gzip,
	}); err != nil {
		return nil, nil, err
	}
	return d.deploy(path)
}

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

	d.sendStatus("Building Docker image...")
	imageHash, err := d.mesg.container.Build(path)
	if err != nil {
		return nil, nil, err
	}
	d.injectConfigurationInDependencies(service, imageHash)

	d.sendStatus(fmt.Sprintf("%s Completed.", aurora.Green("✔")))
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

func (d *serviceDeployer) sendStatus(message string) {
	if d.Statuses != nil {
		d.Statuses <- message
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
