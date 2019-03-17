package xvalidator

import (
	"strconv"
	"strings"

	"github.com/mesg-foundation/core/x/xnet"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	minPort = 1
	maxPort = 65535

	portSeparator = ":"
	envSeparator  = "="
)

// IsDomainName validates if given field is valid domain name.
func IsDomainName(fl validator.FieldLevel) bool {
	return xnet.IsDomainName(fl.Field().String())
}

// IsPortMapping validates if given field is valid port mapping.
func IsPortMapping(fl validator.FieldLevel) bool {
	ports := strings.Split(fl.Field().String(), portSeparator)
	if len(ports) != 1 && len(ports) != 2 {
		return false
	}

	for _, port := range ports {
		i, err := strconv.Atoi(port)
		if err != nil || !(minPort <= i && i <= maxPort) {
			return false
		}
	}
	return true
}

// IsEnv validates if given field is valid env variable declaration.
func IsEnv(fl validator.FieldLevel) bool {
	return len(strings.Split(fl.Field().String(), envSeparator)) == 2
}
