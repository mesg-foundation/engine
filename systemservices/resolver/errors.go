package resolver

// errorOutputKey is the common error output key.
const errorOutputKey = "error"

// ErrorOutput is the common error output.
type ErrorOutput struct {
	// Message is the error message.
	Message string `json:"message"`
}
