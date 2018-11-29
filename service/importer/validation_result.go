package importer

// ValidationResult contains the result of the validation of a service.
type ValidationResult struct {
	ServiceFileWarnings []string
	ServiceFileExist    bool
	DockerfileExist     bool
}

// IsValid returns true if the validation result is valid.
func (v *ValidationResult) IsValid() bool {
	return v.ServiceFileExist &&
		len(v.ServiceFileWarnings) == 0 &&
		v.DockerfileExist
}
