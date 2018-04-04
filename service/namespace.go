package service

import "strings"

// NAMESPACE is the namespace used for the docker services
const NAMESPACE string = "MESG"

// Namespace of a given service
func (service *Service) Namespace() (res string) {
	res = strings.Join([]string{
		NAMESPACE,
		strings.Replace(service.Name, " ", "-", -1),
	}, "-")
	return
}
