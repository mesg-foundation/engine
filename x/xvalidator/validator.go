package xvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/x/xnet"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

const (
	minPort = 1
	maxPort = 65535

	portSeparator = ":"
	envSeparator  = "="
)

// Default global value to used from package level.
var Validate, Translator = New()

// NewWithPrefix returns a new instance of 'validate' with more validation fields prefixed with 'prefix'.
func NewWithPrefix(prefix string) (*validator.Validate, ut.Translator) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()

	validate.RegisterValidation("env", IsEnv)
	validate.RegisterTranslation("env", trans, func(ut ut.Translator) error {
		return ut.Add("env", "{0} must be a valid env variable name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("portmap", IsPortMapping)
	validate.RegisterTranslation("portmap", trans, func(ut ut.Translator) error {
		return ut.Add("portmap", "{0} must be a valid port mapping. eg: 80 or 80:80", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("domain", IsDomainName)
	validate.RegisterTranslation("domain", trans, func(ut ut.Translator) error {
		return ut.Add("domain", "{0} must respect domain-style notation. eg: author.name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field())
		return t
	})

	validate.RegisterValidation("hash", IsHash)
	validate.RegisterTranslation("hash", trans, func(ut ut.Translator) error {
		return ut.Add("hash", "{0} must be a valid hash in hex format", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("accaddress", IsAccAddress)
	validate.RegisterTranslation("accaddress", trans, func(ut ut.Translator) error {
		return ut.Add("hash", "{0} must be a valid address", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	en_translations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}

// New returns a new instance of 'validate' with more validation fields.
func New() (*validator.Validate, ut.Translator) {
	return NewWithPrefix("")
}

// IsHash validates if given field is valid hash.
// It checks both string and slice of bytes.
func IsHash(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		_, err := hash.Decode(v.String())
		return err == nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			// if it's not slice of bytes then break
			break
		}
		_, err := hash.DecodeFromBytes(v.Bytes())
		return err == nil
	}
	return false
}

// IsAccAddress validates if given field is valid cosmos account address.
func IsAccAddress(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		_, err := cosmostypes.AccAddressFromBech32(fl.Field().String())
		return err == nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			// if it's not slice of bytes then break
			break
		}
		return cosmostypes.VerifyAddressFormat(v.Bytes()) == nil
	}
	return false
}

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

var envNameRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

// IsEnv validates if given field is valid env variable declaration.
// The valid formats are:
// - ENV
// - ENV=
// - ENV=value
func IsEnv(fl validator.FieldLevel) bool {
	e := strings.Split(fl.Field().String(), envSeparator)
	return (len(e) == 1 || len(e) == 2) && envNameRegexp.MatchString(e[0])
}
