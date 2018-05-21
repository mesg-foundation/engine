package dependency

import (
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types/mount"
	"github.com/mesg-foundation/core/config"
	"github.com/spf13/viper"
)

// Volumes extract volumes from a Dependency and transform them to a Docker Mount
func Volumes(service service, dependency dependency, name string) (volumes []mount.Mount) {
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
	volumesForm := dependency.GetVolumesfrom()
	for _, dep := range volumesForm {
		depVolumes := service.GetDependencies()[dep].GetVolumes()
		for _, volume := range depVolumes {
			path := filepath.Join(service.Namespace(), dep, volume)
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
