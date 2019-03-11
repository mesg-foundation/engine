package xvalidator

import (
	"strconv"
	"strings"

	"github.com/mesg-foundation/core/x/xnet"
	validator "gopkg.in/go-playground/validator.v9"
)

// IsDomainName validates if given field is valid domain name.
func IsDomainName(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}
	return xnet.IsDomainName(fl.Field().String())
}

// IsPortMapping validates if given field is valid port mapping.
func IsPortMapping(fl validator.FieldLevel) bool {
	ports := strings.Split(fl.Field().String(), ":")
	if len(ports) != 1 && len(ports) != 2 {
		return false
	}
	for _, port := range ports {
		i, err := strconv.Atoi(port)
		if !(err == nil && 0 < i && i <= 65535) {
			return false
		}
	}
	return true
}
