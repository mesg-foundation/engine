package service

import (
	"io"

	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"
)

// Dependency represents a Docker container and it holds instructions about
// how it should run.
type Dependency struct {
	// Key is the key of dependency.
	Key string `hash:"1"`

	// Image is the Docker image.
	Image string `hash:"name:2"`

	// Volumes are the Docker volumes.
	Volumes []string `hash:"name:3"`

	// VolumesFrom are the docker volumes-from from.
	VolumesFrom []string `hash:"name:4"`

	// Ports holds ports configuration for container.
	Ports []string `hash:"name:5"`

	// Command is the Docker command which will be executed when container started.
	Command string `hash:"name:6"`

	// service is the dependency's service.
	service *Service `hash:"-"`
}

// Logs gives the dependency logs. rstd stands for standard logs and rerr stands for
// error logs.
func (d *Dependency) Logs() (rstd, rerr io.ReadCloser, err error) {
	var reader io.ReadCloser
	reader, err = d.service.docker.ServiceLogs(d.namespace())
	if err != nil {
		return nil, nil, err
	}
	sr, sw := io.Pipe()
	er, ew := io.Pipe()
	go func() {
		if _, err := stdcopy.StdCopy(sw, ew, reader); err != nil {
			reader.Close()
			logrus.Errorln(err)
		}
	}()
	return sr, er, nil
}
