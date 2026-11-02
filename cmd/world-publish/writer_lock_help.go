package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// writerLockHelp explains how to release writer authority before publishing.
func writerLockHelp(verb, storePath string) []string {
	lines := []string{
		fmt.Sprintf("world-publish %s: the world writer lock is held by another process", verb),
		"  why: publishing writes the store, and only one writer may hold it",
		fmt.Sprintf("  next: stop the daemon serving %s, publish, then start it again", storePath),
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return lines
	}
	data, err := os.ReadFile(filepath.Join(home, ".ailang-world", "service.json"))
	if err != nil {
		return lines
	}
	var record struct {
		Schema      string `json:"schema"`
		Store       string `json:"store"`
		Address     string `json:"address"`
		Uninstalled bool   `json:"uninstalled"`
	}
	if json.Unmarshal(data, &record) != nil || record.Schema != "world/service/v1" || record.Uninstalled || record.Store == "" {
		return lines
	}
	resolve := func(path string) string {
		if resolved, err := filepath.EvalSymlinks(path); err == nil {
			return resolved
		}
		return filepath.Clean(path)
	}
	if resolve(record.Store) == resolve(storePath) {
		lines[0] = fmt.Sprintf("world-publish %s: the world writer lock is held by the World service (%s)", verb, record.Address)
		lines[2] = fmt.Sprintf("  next: ailang-worldd service stop; world-publish %s …; ailang-worldd service start", verb)
	}
	return lines
}
