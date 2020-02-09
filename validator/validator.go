package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/mesg-foundation/engine/hash"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

const (
	minPort = 1
	maxPort = 65535

	portSeparator = ":"
	envSeparator  = "="
)

// Validate let's you validate data.
var Validate, _ = New()

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
// Source: https://golang.org/src/net/dnsclient.go
// IsDomainName checks if a string is a presentation-format domain name
// (currently restricted to hostname-compatible "preferred name" LDH labels and
// SRV-like "underscore labels"; see golang.org/issue/12421).
//nolint:gocyclo
func IsDomainName(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
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
