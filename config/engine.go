package config

// engine is a generic interface that any kind of engine should implement
type engine interface {
	setDefaultValue(key string, defaultValue string) error
	setValue(key string, value string) error
	getValue(key string) (value string, err error)
	getEnvKey(key string) (envKey string)
}
