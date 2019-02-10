package service

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mesg-foundation/core/x/xerrors"
	"github.com/mesg-foundation/core/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
	yaml "gopkg.in/yaml.v2"
)

// DefaultServicefileName is the default config filename for mesg.
const DefaultServicefileName = "mesg.yml"

// validate is used for mesg service config.
var validate = func() *validator.Validate {
	validate := validator.New()
	validate.RegisterValidation("domain", xvalidator.IsDomainName)
	return validate
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
				e := fmt.Errorf("%s with value %q is invalid: %s %s", verr.Namespace(), verr.Value(), verr.Tag(), verr.Param())
				errs = append(errs, e)
			}
		default:
			errs = append(errs, err)
		}
	}

	for key, dep := range service.Dependencies {
		if key == MainServiceKey {
			errs = append(errs, &validateError{msg: MainServiceKey + " as key is forbidden for dependency"})
		}
		if dep == nil || dep.Image == "" {
			errs = append(errs, &validateError{key: key, msg: "image is empty"})
		}
	}

	validatePorts := func(key string, ports []string) {
		for _, port := range ports {
			parts := strings.Split(port, ":")
			if i, err := strconv.ParseUint(parts[0], 10, 64); err != nil || i == 0 || i > 65535 {
				errs = append(errs, &validateError{key: key, msg: "port " + port + " is invalid"})
				continue
			}

			if len(parts) > 1 {
				if i, err := strconv.ParseUint(parts[1], 10, 64); err != nil || i == 0 || i > 65535 {
					errs = append(errs, &validateError{key: key, msg: "port " + port + " is invalid"})
					continue
				}
			}
		}
	}

	for _, depVolumeKey := range service.Configuration.VolumesFrom {
		if depVolumeKey == MainServiceKey {
			errs = append(errs, &validateError{
				key: MainServiceKey,
				msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
			})
		}
		if _, ok := service.Dependencies[depVolumeKey]; depVolumeKey != MainServiceKey && !ok {
			errs = append(errs, &validateError{
				key: MainServiceKey,
				msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
			})
		}
	}
	validatePorts(MainServiceKey, service.Configuration.Ports)

	for key, dep := range service.Dependencies {
		if dep == nil {
			continue
		}
		for _, depVolumeKey := range dep.VolumesFrom {
			if depVolumeKey == key {
				errs = append(errs, &validateError{
					key: key,
					msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
				})
			}
			if _, ok := service.Dependencies[depVolumeKey]; depVolumeKey != MainServiceKey && !ok {
				errs = append(errs, &validateError{
					key: key,
					msg: "volumesFrom is invalid: " + depVolumeKey + " dependency dose not exist",
				})
			}
		}

		validatePorts(key, dep.Ports)
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
		return fmt.Sprintf("configuration %s", e.msg)
	case "":
		return e.msg
	default:
		return fmt.Sprintf("dependency[%s] %s", e.key, e.msg)
	}
}
