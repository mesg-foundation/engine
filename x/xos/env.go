package xos

import (
	"os"
	"sort"
	"strings"
)

// GetenvDefault retrieves the value of the environment variable named by the key.
// It returns the value, which will be set to fallback if the variable is empty.
func GetenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// EnvMapToSlice transform a map of key value to a slice of env in the form "key=value".
// Env vars are sorted by names to get an accurate order while testing.
func EnvMapToSlice(values map[string]string) []string {
	env := make([]string, 0, len(values))
	for k, v := range values {
		env = append(env, k+"="+v)
	}
	sort.Stable(sort.StringSlice(env))
	return env
}

// EnvMapToString transform a map of key value to a string in the form "key=value;key1=value1".
// Env vars are sorted by names to get an accurate order while testing.
func EnvMapToString(values map[string]string) string {
	env := EnvMapToSlice(values)
	return strings.Join(env, ";")
}

// EnvSliceToMap transform a slice of key=value to a map.
func EnvSliceToMap(values []string) map[string]string {
	env := make(map[string]string, len(values))
	for _, v := range values {
		if e := strings.SplitN(v, "=", 2); len(e) == 1 {
			env[e[0]] = ""
		} else {
			env[e[0]] = e[1]
		}
	}
	return env
}

// EnvMergeMaps merges multiple maps into single one.
// If the same key exist multiple time, it will be overwritten by the latest occurrence.
func EnvMergeMaps(values ...map[string]string) map[string]string {
	env := make(map[string]string)
	for _, e := range values {
		for k, v := range e {
			env[k] = v
		}
	}
	return env
}

// EnvMergeSlices merges multiple slices into single one.
// If the same key exist multiple time, it will be added in occurrence order.
func EnvMergeSlices(values ...[]string) []string {
	env := make([]string, 0)
	for _, v := range values {
		env = append(env, v...)
	}
	return env
}
