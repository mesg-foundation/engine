package xerrors

import "strings"

type Errors []error

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
