package docker

import (
	godocker "github.com/fsouza/go-dockerclient"
)

type ServiceOptions struct {
	Image                string
	Namespace            []string
	Ports                []Port
	Mounts               []Mount
	Env                  []string
	Args                 []string
	NetworksID           []string
	Labels               map[string]string
	CreateServiceOptions *godocker.CreateServiceOptions
}

type Port struct {
	Target    uint32
	Published uint32
}

type Mount struct {
	Source string
	Target string
}
