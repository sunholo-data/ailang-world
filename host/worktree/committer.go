package worktree

import (
	"fmt"
	"strings"
)

// ParseCommitter validates an optional preset identity without consulting git.
func ParseCommitter(value string) (name, email string, err error) {
	if value == "" {
		return "", "", nil
	}
	bad := func() (string, string, error) {
		return "", "", fmt.Errorf("worktree: invalid committer")
	}
	if strings.ContainsAny(value, "\r\n") || strings.Count(value, "<") != 1 || strings.Count(value, ">") != 1 {
		return bad()
	}
	left := strings.IndexByte(value, '<')
	right := strings.IndexByte(value, '>')
	if right <= left || strings.TrimSpace(value[right+1:]) != "" {
		return bad()
	}
	name, email = strings.TrimSpace(value[:left]), strings.TrimSpace(value[left+1:right])
	if name == "" || email == "" {
		return bad()
	}
	return name, email, nil
}
