package dockertest

import "errors"

var (
	ErrContainerDoesNotExists = errors.New("containers does not exists")
)

type NotFoundErr struct {
}

func NewNotFoundErr() NotFoundErr {
	return NotFoundErr{}
}

func (e NotFoundErr) NotFound() bool {
	return true
}

func (e NotFoundErr) Error() string {
	return "not found"
}
