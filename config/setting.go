package config

// Setting is the generic interface that any kind of setting should implement
type setting interface {
	setValue(string) error
	getValue() (string, error)
	getEnvKey() string
}
