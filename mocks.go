package main

import "github.com/docker/docker/client"

// CommonAPIClient is the common methods between stable and experimental versions of APIClient.
type CommonAPIClient interface {
	client.CommonAPIClient
}

// disble no-used linter
var _ = CommonAPIClient(nil)
