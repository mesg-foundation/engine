package service

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"
)

// DockerVolumes extract volumes from a Dependency and transform them to a Docker Mount
func DockerVolumes(service *Service, dependency *Dependency, name string) (volumes []mount.Mount) {
	volumes = make([]mount.Mount, 0)
	for _, volume := range dependency.GetVolumes() {
		path := filepath.Join(service.Namespace(), name, volume)
		source := filepath.Join(viper.GetString(config.ServicePathHost), path)
		volumes = append(volumes, mount.Mount{
			Source: source,
			Target: volume,
		})
		os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
	}
	for _, depString := range dependency.Volumesfrom {
		for _, volume := range service.Dependencies[depString].Volumes {
			path := filepath.Join(service.Namespace(), depString, volume)
			source := filepath.Join(viper.GetString(config.ServicePathHost), path)
			volumes = append(volumes, mount.Mount{
				Source: source,
				Target: volume,
			})
			os.MkdirAll(filepath.Join(viper.GetString(config.ServicePathDocker), path), os.ModePerm)
		}
	}
	return
}
