package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xstrings"
	"github.com/mesg-foundation/core/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	yaml "gopkg.in/yaml.v2"
)

// DefaultServicefileName is the default config filename for mesg.
const DefaultServicefileName = "mesg.yml"

// service namespace prefix.
const namespacePrefix = "Service."

// validate is used for mesg service definition.
var validate, translate = func() (*validator.Validate, ut.Translator) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	validate.RegisterValidation("portmap", xvalidator.IsPortMapping)
	validate.RegisterTranslation("portmap", trans, func(ut ut.Translator) error {
		return ut.Add("portmap", "{0} must be a valid port mapping. Eg: 80 or 80:80", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("portmap", fe.Field(), namespacePrefix)
		return t
	})
	validate.RegisterValidation("domain", xvalidator.IsDomainName)
	validate.RegisterTranslation("domain", trans, func(ut ut.Translator) error {
		return ut.Add("domain", "{0} must respect domain-style notation. Eg: author.name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("domain", fe.Field())
		return t
	})
	en_translations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}()

// ReadDefinition validates a service at a source.
func ReadDefinition(contextDir string) (*Service, error) {
	mesgfile := filepath.Join(contextDir, DefaultServicefileName)
	f, err := os.Open(mesgfile)
	if err != nil {
		return nil, fmt.Errorf("open mesg.yml error: %s", err)
	}

	service, err := DecodeDefinition(f)
	if err != nil {
		return nil, err
	}

	if err := ValidateDefinition(service); err != nil {
		return nil, err
	}

	return service, nil
}

// DecodeDefinition uses yaml parser to decode
// service definition from given reader.
func DecodeDefinition(r io.Reader) (*Service, error) {
	var (
		service Service
		dec     = yaml.NewDecoder(r)
	)
	dec.SetStrict(true)
	if err := dec.Decode(&service); err != nil {
		return nil, fmt.Errorf("parse mesg.yml error: %s", err)
	}
	return &service, nil
}

// ValidateDefinition validates service definition.
//nolint:gocyclo
func ValidateDefinition(service *Service) error {
	var errs xerrors.Errors
	if err := validate.Struct(service); err != nil {
		switch verrs := err.(type) {
		case validator.ValidationErrors:
			for _, verr := range verrs {
				// Remove the name of the struct and the field from namespace
				trimmedNamespace := strings.TrimPrefix(verr.Namespace(), namespacePrefix)
				trimmedNamespace = strings.TrimSuffix(trimmedNamespace, verr.Field())
				// Only use it when in-cascade field
				namespace := ""
				if verr.Field() != trimmedNamespace {
					namespace = strings.ToLower(trimmedNamespace)
				}
				errs = append(errs, fmt.Errorf("%s%s", namespace, xstrings.FirstToLower(verr.Translate(translate))))
			}
		default:
			errs = append(errs, err)
		}
	}

	for key, dep := range service.Dependencies {
		if dep == nil {
			continue
		}
		if dep.Image == "" {
			errs = append(errs, &validateError{key: key, msg: "image is empty"})
		}
	}

	for _, depVolumeKey := range service.Configuration.VolumesFrom {
		if depVolumeKey == MainServiceKey {
			errs = append(errs, &validateError{
				key: MainServiceKey,
				msg: "volumesFrom is invalid: cyclic volume import",
			})
		}
		if _, ok := service.Dependencies[depVolumeKey]; depVolumeKey != MainServiceKey && !ok {
			errs = append(errs, &validateError{
				key: MainServiceKey,
				msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
			})
		}
	}

	for key, dep := range service.Dependencies {
		if dep == nil {
			continue
		}
		for _, depVolumeKey := range dep.VolumesFrom {
			if depVolumeKey == key {
				errs = append(errs, &validateError{
					key: key,
					msg: "volumesFrom is invalid: cyclic volume import",
				})
			}
			if _, ok := service.Dependencies[depVolumeKey]; depVolumeKey != MainServiceKey && !ok {
				errs = append(errs, &validateError{
					key: key,
					msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
				})
			}
		}
	}
	return errs.ErrorOrNil()
}

// validateError is an error used just for definition validation.
// It returns error message based on dependency key.
type validateError struct {
	key string
	msg string
}

func (e *validateError) Error() string {
	switch e.key {
	case MainServiceKey:
		return fmt.Sprintf("configuration.%s", e.msg)
	case "":
		return e.msg
	default:
		return fmt.Sprintf("dependency[%s].%s", e.key, e.msg)
	}
}
