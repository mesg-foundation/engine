package dockertest

// NotFoundErr satisfies docker client's NotFoundErr interface.
// docker.IsErrNotFound(err) will return true with NotFoundErr.
type NotFoundErr struct{}

// NotFound indicates that this error is a not found error.
func (e NotFoundErr) NotFound() bool {
	return true
}

// Error returns the string representation of error.
func (e NotFoundErr) Error() string {
	return "not found"
}
