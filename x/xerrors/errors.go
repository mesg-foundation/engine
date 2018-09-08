package xerrors

type Errors []error

func (e Errors) Error() string {
	if len(e) == 0 {
		return ""
	}

	var s string
	for _, err := range e {
		s += err.Error() + "\n"
	}
	// remove last new line
	return s[:len(s)-1]
}
