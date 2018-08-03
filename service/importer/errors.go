package importer

// ValidationError is the error type for the Validation of service.
type ValidationError struct{}

func (v *ValidationError) Error() string {
	return "Service is not valid"
}
