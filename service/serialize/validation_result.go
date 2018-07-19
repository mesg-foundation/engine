package serialize

// ValidationResult contains the result of the validation of a service
type ValidationResult struct {
	ServiceFileExist    bool
	ServiceFileWarnings []string
	DockerfileExist     bool
}

// IsValid returns true if all the validation result is valid
func (v *ValidationResult) IsValid() bool {
	return v.ServiceFileExist &&
		len(v.ServiceFileWarnings) == 0 &&
		v.DockerfileExist
}
