package xos

import (
	"sort"
	"strings"
)

// EnvMergeSlices merges multiple slices into single one.
// If the same key exist multiple time, it will be added in occurrence order.
func EnvMergeSlices(values ...[]string) []string {
	envs := make(map[string]string)
	for _, value := range values {
		for _, v := range value {
			if e := strings.SplitN(v, "=", 2); len(e) == 1 {
				envs[e[0]] = ""
			} else {
				envs[e[0]] = e[1]
			}
		}
	}

	env := make([]string, 0, len(values))
	for k, v := range envs {
		env = append(env, k+"="+v)
	}

	// Make sure envs are sorted to give repeatable output
	// It is important for hash instance calculation
	sort.Stable(sort.StringSlice(env))
	return env
}
