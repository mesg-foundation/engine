package instancesdk

import (
	"strconv"
	"strings"

	"github.com/mesg-foundation/engine/container"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

// instanceNamespace returns the namespace of the service.
func instanceNamespace(hash hash.Hash) string {
	return hash.String()
}

// dependencyNamespace builds the namespace of a dependency.
func dependencyNamespace(instanceNamespace string, dependencyKey string) string {
	return hash.Dump(instanceNamespace + dependencyKey).String()
}

func extractPorts(d *service.Dependency) []container.Port {
	ports := make([]container.Port, len(d.Ports))
	for i, p := range d.Ports {
		split := strings.Split(p, ":")
		from, _ := strconv.ParseUint(split[0], 10, 64)
		to := from
		if len(split) > 1 {
			to, _ = strconv.ParseUint(split[1], 10, 64)
		}
		ports[i] = container.Port{
			Target:    uint32(to),
			Published: uint32(from),
		}
	}
	return ports
}

// TODO: add test and hack for MkDir in CircleCI
func extractVolumes(s *service.Service, d *service.Dependency) []container.Mount {
	volumes := make([]container.Mount, 0)
	for _, volume := range d.Volumes {
		volumes = append(volumes, container.Mount{
			Source: volumeKey(s, d.Key, volume),
			Target: volume,
		})
	}
	return volumes
}

func extractVolumesFrom(s *service.Service, d *service.Dependency) ([]container.Mount, error) {
	volumesFrom := make([]container.Mount, 0)
	for _, depName := range d.VolumesFrom {
		dep, err := s.GetDependency(depName)
		if err != nil {
			if depName == service.MainServiceKey {
				dep = s.Configuration
			} else {
				return nil, err
			}
		}
		for _, volume := range dep.Volumes {
			volumesFrom = append(volumesFrom, container.Mount{
				Source: volumeKey(s, depName, volume),
				Target: volume,
			})
		}
	}
	return volumesFrom, nil
}

// volumeKey creates a key for service's volume based on the sid to make sure that the volume
// will stay the same for different versions of the service.
func volumeKey(s *service.Service, dependency, volume string) string {
	return hash.Dump([]string{
		s.Sid,
		dependency,
		volume,
	}).String()
}
