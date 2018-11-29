package dockerapi

import (
	docker "github.com/docker/docker/client"
)

// Docker embeds docker.CommonAPIClient interface.
type Docker interface {
	docker.CommonAPIClient
}
