package main

import (
	"fmt"
	"sort"
	"strings"
)

func finishProfileEnvLine(profile string, env map[string]string, path []string) string {
	if profile == "" {
		return ""
	}
	keys := make([]string, 0, len(env))
	for key := range env {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	names := strings.Join(keys, ", ")
	if names == "" {
		names = "(none)"
	}
	return fmt.Sprintf("  with profile %s env: %s and PATH +%d\n", profile, names, len(path))
}
