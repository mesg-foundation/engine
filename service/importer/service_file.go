package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/mesg-foundation/core/x/xvalidator"
	validator "gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
	yaml "gopkg.in/yaml.v2"
)

const (
	namespacePrefix = "ServiceDefinition."
)

func readServiceFile(path string) ([]byte, error) {
	file := filepath.Join(path, "mesg.yml")
	return ioutil.ReadFile(file)
}

// validateServiceFile returns a list of warnings.
func validateServiceFile(data []byte) []string {
	var service ServiceDefinition
	if err := yaml.UnmarshalStrict(data, &service); err != nil {
		errs, ok := err.(*yaml.TypeError)
		if !ok {
			return []string{err.Error()}
		}
		return errs.Errors
	}
	return validateServiceStruct(&service)
}

func validateServiceStruct(service *ServiceDefinition) []string {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()
	validate.RegisterValidation("port", xvalidator.IsPortMapping)
	validate.RegisterTranslation("port", trans, func(ut ut.Translator) error {
		return ut.Add("port", "{0} must be a valid port. Eg: 80 or 80:80", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("port", fe.Field(), namespacePrefix)
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

	err := validate.Struct(service)
	warnings := []string{}
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// Remove the name of the struct and the field from namespace
			trimmedNamespace := strings.TrimPrefix(e.Namespace(), namespacePrefix)
			trimmedNamespace = strings.TrimSuffix(trimmedNamespace, e.Field())
			// Only use it when in-cascade field
			namespace := ""
			if e.Field() != trimmedNamespace {
				namespace = trimmedNamespace
			}
			warnings = append(warnings, fmt.Sprintf("%s%s", namespace, e.Translate(trans)))
		}
	}

	for k, d := range service.Dependencies {
		if d.Image == "" {
			warnings = append(
				warnings,
				fmt.Sprintf("Dependencies[%s].Image is a required field", k),
			)
		}
	}

	for _, depVolumeKey := range service.Configuration.VolumesFrom {
		if _, ok := service.Dependencies[depVolumeKey]; !ok {
			warnings = append(
				warnings,
				fmt.Sprintf("configuration.volumesfrom is invalid: dependency %q does not exist", depVolumeKey),
			)
		}
	}
	for key, dep := range service.Dependencies {
		for _, depVolumeKey := range dep.VolumesFrom {
			if _, ok := service.Dependencies[depVolumeKey]; !ok && depVolumeKey != "service" {
				warnings = append(
					warnings,
					fmt.Sprintf("dependencies[%s].volumesfrom is invalid: dependency %q does not exist", key, depVolumeKey),
				)
			}
		}
	}

	return warnings
}
