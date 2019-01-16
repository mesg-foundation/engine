package ethwallet

type outputError struct {
	Message string `json:"message"`
}

func OutputError(message string) (string, interface{}) {
	return "error", outputError{
		Message: message,
	}
}
