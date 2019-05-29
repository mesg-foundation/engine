package ethwallet

import "errors"

var (
	errAccountNotFound     = errors.New("Account not found")
	errCannotParseValue    = errors.New("Cannot parse value")
	errCannotParseGasPrice = errors.New("Cannot parse gasPrice")
)
