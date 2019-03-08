package xvalidator

import (
	"strconv"

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

// IsPort validates if given field is valid port.
func IsPort(fl validator.FieldLevel) bool {
	i, err := strconv.Atoi(fl.Field().String())
	return err == nil && 0 < i && i <= 65535
}
