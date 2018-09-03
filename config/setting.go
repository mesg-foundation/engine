package config

// Setting is the generic interface that any kind of setting should implement
type Setting interface {
	GetValue() string
	GetEnvKey() string
}
