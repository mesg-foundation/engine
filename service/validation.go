package service

import (
	"errors"
	"os"

	"github.com/xeipuuv/gojsonschema"
)

// IsValid returns true if the service is valid, false otherwise
// The all validation can be found in https://github.com/mesg-foundation/application/tree/dev/service/schema.json
func (service *Service) IsValid() (valid bool, errs []error) {
	pwd, err := os.Getwd()
	if err != nil {
		errs = append(errs, err)
		return
	}
	schema := gojsonschema.NewReferenceLoader("file://" + pwd + "/schema.json")
	data := gojsonschema.NewGoLoader(service)
	result, err := gojsonschema.Validate(schema, data)
	for _, e := range result.Errors() {
		errs = append(errs, errors.New(e.String()))
	}
	if err != nil {
		errs = append(errs, err)
	}
	valid = result.Valid()
	return
}
