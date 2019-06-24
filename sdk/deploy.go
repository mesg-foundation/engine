package sdk

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/hash"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
	"github.com/mesg-foundation/core/utils/dirhash"
	"github.com/mesg-foundation/core/x/xgit"
	"github.com/mesg-foundation/core/x/xos"
	uuid "github.com/satori/go.uuid"
)

// serviceDeployer provides functionalities to deploy a MESG service.
type serviceDeployer struct {
	// statuses receives status messages produced during deployment.
	statuses chan DeployStatus

	env map[string]string

	sdk *SDK
}

// StatusType indicates the type of status message.
type StatusType int

const (
	_ StatusType = iota // skip zero value.

	// Running indicates that status message belongs to a continuous state.
	Running

	// DonePositive indicates that status message belongs to a positive noncontinuous state.
	DonePositive

	// DoneNegative indicates that status message belongs to a negative noncontinuous state.
	DoneNegative
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    StatusType
}

// DeployServiceOption is a configuration func for deploying services.
type DeployServiceOption func(*serviceDeployer)

// DeployServiceStatusOption receives chan statuses to send deploy statuses.
func DeployServiceStatusOption(statuses chan DeployStatus) DeployServiceOption {
	return func(deployer *serviceDeployer) {
		deployer.statuses = statuses
	}
}

// DeployService deploys a service from a gzipped tarball.
func (sdk *SDK) DeployService(r io.Reader, env map[string]string, options ...DeployServiceOption) (*service.Service,
	*importer.ValidationError, error) {
	d := newServiceDeployer(sdk, env, options...)
	return d.importWithProcess(func(path string) error {
		return d.preprocessArchive(r, path)
	})
}

// DeployServiceFromURL deploys a service from a Git or tarball url.
// individual branch deployments are supported through Git, see:
// - https://github.com/mesg-foundation/service-ethereum
// - https://github.com/mesg-foundation/service-ethereum#branchName
func (sdk *SDK) DeployServiceFromURL(url string, env map[string]string, options ...DeployServiceOption) (*service.Service,
	*importer.ValidationError, error) {
	d := newServiceDeployer(sdk, env, options...)
	return d.importWithProcess(func(path string) error {
		if xgit.IsGitURL(url) {
			return d.preprocessGit(url, path)
		}
		return d.preprocessURL(url, path)
	})
}

// newServiceDeployer creates a new serviceDeployer with given sdk and options.
func newServiceDeployer(sdk *SDK, env map[string]string, options ...DeployServiceOption) *serviceDeployer {
	d := &serviceDeployer{
		sdk: sdk,
		env: env,
	}
	for _, option := range options {
		option(d)
	}
	return d
}

func (d *serviceDeployer) importWithProcess(processing func(path string) error) (*service.Service, *importer.ValidationError, error) {
	path, err := d.createTempDir()
	if err != nil {
		return nil, nil, err
	}
	defer os.RemoveAll(path)

	d.sendStatus("Receiving service context...", Running)
	defer d.closeStatus()

	if err := processing(path); err != nil {
		return nil, nil, err
	}

	contextDir, err := d.processPath(path)
	if err != nil {
		return nil, nil, err
	}
	return d.deploy(contextDir)
}

func (d *serviceDeployer) preprocessGit(url string, path string) error {
	return xgit.Clone(url, path)
}

func (d *serviceDeployer) preprocessURL(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return d.preprocessArchive(resp.Body, path)
}

func (d *serviceDeployer) preprocessArchive(r io.Reader, path string) error {
	return archive.Untar(r, path, nil)
}

func (d *serviceDeployer) processPath(path string) (string, error) {
	// XXX: remove .git folder from repo.
	// It makes docker build iamge id same between repo clones.
	if err := os.RemoveAll(filepath.Join(path, ".git")); err != nil {
		return "", err
	}

	// NOTE: this is check for tar repos, if there is only one
	// directory inside untar archive set temp path to it.
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	if len(dirs) == 1 && dirs[0].IsDir() {
		path = filepath.Join(path, dirs[0].Name())
	}

	return path, nil
}

// deploy deploys a service in contextDir.
func (d *serviceDeployer) deploy(contextDir string) (*service.Service, *importer.ValidationError, error) {
	s := &service.Service{}

	def, err := importer.From(contextDir)
	validationErr, err := d.assertValidationError(err)
	if err != nil {
		return nil, nil, err
	}
	if validationErr != nil {
		return nil, validationErr, nil
	}

	dh := dirhash.New(contextDir)
	envbytes := []byte(xos.EnvMapToString(d.env))
	h, err := dh.Sum(envbytes)
	if err != nil {
		return nil, nil, err
	}
	s.Hash = hash.Hash(h)

	injectDefinition(s, def)

	if err := s.ValidateConfigurationEnv(d.env); err != nil {
		return nil, nil, err
	}

	// replace default env with new one.
	defenv := xos.EnvSliceToMap(s.Configuration.Env)
	s.Configuration.Env = xos.EnvMapToSlice(xos.EnvMergeMaps(defenv, d.env))

	d.sendStatus("Building Docker image...", Running)

	imageHash, err := d.sdk.container.Build(contextDir)
	if err != nil {
		return nil, nil, err
	}

	d.sendStatus("Image built with success", DonePositive)

	s.Configuration.Image = imageHash
	// TODO: the following test should be moved in New function
	if s.Sid == "" {
		// make sure that sid doesn't have the same length with id.
		s.Sid = "_" + s.Hash.String()
	}

	return s, nil, d.sdk.db.Save(s)
}

func (d *serviceDeployer) createTempDir() (path string, err error) {
	return ioutil.TempDir("", "mesg-"+uuid.NewV4().String())
}

// sendStatus sends a status message.
func (d *serviceDeployer) sendStatus(message string, typ StatusType) {
	if d.statuses != nil {
		d.statuses <- DeployStatus{
			Message: message,
			Type:    typ,
		}
	}
}

// closeStatus closes statuses chan.
func (d *serviceDeployer) closeStatus() {
	if d.statuses != nil {
		close(d.statuses)
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