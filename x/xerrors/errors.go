package xerrors

import "strings"

// Errors is an error for tracing multiple errors.
type Errors []error

// ErrorOrNil returns an error if there is more then 0 error, nil otherwise.
func (e Errors) ErrorOrNil() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

func (e Errors) Error() string {
	var s []string
	for _, err := range e {
		if err != nil {
			s = append(s, err.Error())
		}
	}

	return strings.Join(s, "\n")
}
