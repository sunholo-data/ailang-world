//go:build !darwin && !linux

package agentrun

func ProcessIdentity(int) (string, bool) { return "", false }
