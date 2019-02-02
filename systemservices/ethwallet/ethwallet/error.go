package ethwallet

import "errors"

var (
	errAccountNotFound     = errors.New("Account not found")
	errCannotParseValue    = errors.New("Cannot parse value")
	errCannotParseGasPrice = errors.New("Cannot parse gasPrice")
)

type outputError struct {
	Message string `json:"message"`
}

// OutputError is a helper to create an error output from err.
func OutputError(err error) (string, interface{}) {
	return "error", outputError{
		Message: err.Error(),
	}
}
