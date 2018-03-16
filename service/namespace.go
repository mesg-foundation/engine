package service

import "strings"

func (service *Service) namespace() string {
	return strings.Join([]string{"MESG", service.Name}, "-")
}
