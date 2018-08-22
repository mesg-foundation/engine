package xos

import "os"

// GetenvDefault retrieves the value of the environment variable named by the key.
// It returns the value, which will be set to fallback if the variable is empty.
func GetenvDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
