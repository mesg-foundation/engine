package docker

// import (
// 	"strings"

// 	"github.com/docker/docker/api/types/mount"
// 	"github.com/docker/docker/api/types/swarm"
// 	godocker "github.com/fsouza/go-dockerclient"
// 	"github.com/mesg-foundation/core/config"
// 	"github.com/spf13/viper"
// )

// type StartConfig struct {
// 	Namespace      string
// 	DependencyName string
// 	ServiceName    string
// 	Image          string
// 	Command        string
// 	DaemonIP       string
// 	Volumes        []mount.Mount
// 	Ports          []swarm.PortConfig
// 	SharedNetwork  string
// }

// func StartDaemon() (dockerService *swarm.Service, err error) {

// }

// // Start a docker service
// func StartService(c StartConfig) (dockerService *swarm.Service, err error) {
// 	client, err := Client()
// 	if err != nil {
// 		return
// 	}
// 	serviceTemplate := godocker.CreateServiceOptions{
// 		ServiceSpec: swarm.ServiceSpec{
// 			Annotations: swarm.Annotations{
// 				Name: strings.Join([]string{c.Namespace, c.DependencyName}, "_"),
// 				Labels: map[string]string{
// 					"com.docker.stack.image":     c.Image,
// 					"com.docker.stack.namespace": c.Namespace,
// 					"mesg.service":               c.ServiceName,
// 				},
// 			},
// 			TaskTemplate: swarm.TaskSpec{
// 				ContainerSpec: &swarm.ContainerSpec{
// 					Image: c.Image,
// 					Args:  strings.Fields(c.Command),
// 					Env: []string{
// 						"MESG_ENDPOINT=" + viper.GetString(config.APIServiceTargetSocket),
// 						"MESG_ENDPOINT_TCP=" + c.DaemonIP + "" + viper.GetString(config.APIClientTarget),
// 					},
// 					Mounts: append(c.Volumes, mount.Mount{
// 						Source: viper.GetString(config.APIServiceSocketPath),
// 						Target: viper.GetString(config.APIServiceTargetPath),
// 					}),
// 					Labels: map[string]string{
// 						"com.docker.stack.namespace": c.Namespace,
// 					},
// 				},
// 			},
// 			EndpointSpec: &swarm.EndpointSpec{
// 				Ports: c.Ports,
// 			},
// 			Networks: []swarm.NetworkAttachmentConfig{
// 				swarm.NetworkAttachmentConfig{
// 					Target: c.SharedNetwork,
// 				},
// 			},
// 		},
// 	}
// 	return client.CreateService(serviceTemplate)
// }
