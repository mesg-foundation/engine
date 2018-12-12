package xos

import (
	"fmt"
	"os"
	"sort"
)

// GetenvDefault retrieves the value of the environment variable named by the key.
// It returns the value, which will be set to fallback if the variable is empty.
func GetenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// MapToEnv transform a map of key value to a slice of env in the form "key=value".
// Env vars are sorted by names to get an accurate order while testing.
func MapToEnv(data map[string]string) []string {
	env := make([]string, 0, len(data))
	for key, value := range data {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}
	sort.Strings(env)
	return env
}

// MergeMapEnvs merges multiple maps storing environment varialbes into single one.
// If the same key exist multiple time, it will be overwritten by the latest occurrence.
func MergeMapEnvs(envs ...map[string]string) map[string]string {
	env := make(map[string]string)
	for _, e := range envs {
		for k, v := range e {
			env[k] = v
		}
	}
	return env
}
