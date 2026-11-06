package main

import (
	"bytes"
	"strconv"
	"strings"
)

// codexStateSuffix accepts only Codex-owned state after World's exact bytes.
// A nonzero badLine is the first offending physical line.
func codexStateSuffix(generated, actual []byte) (lines, badLine int) {
	if !bytes.HasPrefix(actual, generated) {
		i := 0
		for i < len(generated) && i < len(actual) && generated[i] == actual[i] {
			i++
		}
		return 0, bytes.Count(actual[:min(i, len(actual))], []byte("\n")) + 1
	}
	suffix := string(actual[len(generated):])
	if suffix == "" {
		return 0, 0
	}
	rows := strings.Split(strings.TrimSuffix(suffix, "\n"), "\n")
	base := bytes.Count(generated, []byte("\n"))
	var table []string
	for i, row := range rows {
		line := strings.TrimSpace(tomlStateComment(row))
		if line == "" {
			continue
		}
		bad := func() (int, int) { return 0, base + i + 1 }
		if strings.HasPrefix(line, "[") {
			if strings.HasPrefix(line, "[[") || !strings.HasSuffix(line, "]") {
				return bad()
			}
			var ok bool
			table, ok = tomlStateName(strings.TrimSpace(line[1 : len(line)-1]))
			if !ok {
				return bad()
			}
			switch table[0] {
			case "tui", "projects", "notice":
			case "hooks":
				if len(table) < 2 || table[1] != "state" {
					return bad()
				}
			default:
				return bad()
			}
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		keys, valid := tomlStateName(strings.TrimSpace(key))
		if len(table) == 0 || !ok || !valid || len(keys) != 1 || strings.TrimSpace(value) == "" {
			return bad()
		}
		if table[0] == "projects" && keys[0] != "trust_level" {
			return bad()
		}
	}
	return len(rows), 0
}

// Comments and dots inside quoted keys are data.
func tomlStateComment(s string) string {
	var quote byte
	escape := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if escape {
				escape = false
			} else if quote == '"' && c == '\\' {
				escape = true
			} else if c == quote {
				quote = 0
			}
		} else if c == '"' || c == '\'' {
			quote = c
		} else if c == '#' {
			return s[:i]
		}
	}
	return s
}

func tomlStateName(s string) ([]string, bool) {
	var parts []string
	for {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil, false
		}
		var part string
		if s[0] == '"' || s[0] == '\'' {
			quote := s[0]
			end, escape := 1, false
			for ; end < len(s); end++ {
				if escape {
					escape = false
				} else if quote == '"' && s[end] == '\\' {
					escape = true
				} else if s[end] == quote {
					break
				}
			}
			if end == len(s) {
				return nil, false
			}
			part = s[1:end]
			if quote == '"' {
				var err error
				part, err = strconv.Unquote(s[:end+1])
				if err != nil {
					return nil, false
				}
			}
			s = s[end+1:]
		} else {
			end := 0
			for end < len(s) && (s[end] >= 'a' && s[end] <= 'z' || s[end] >= 'A' && s[end] <= 'Z' || s[end] >= '0' && s[end] <= '9' || s[end] == '_' || s[end] == '-') {
				end++
			}
			if end == 0 {
				return nil, false
			}
			part, s = s[:end], s[end:]
		}
		parts = append(parts, part)
		s = strings.TrimSpace(s)
		if s == "" {
			return parts, true
		}
		if s[0] != '.' {
			return nil, false
		}
		s = s[1:]
	}
}
