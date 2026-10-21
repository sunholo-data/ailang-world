package surface

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var humanEmail = regexp.MustCompile(`[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[a-z]{2,}`)
var boldSpan = regexp.MustCompile(`\*\*([^*]+)\*\*`)
var humanWord = regexp.MustCompile(`\b[A-Za-z]+\b`)

func capHuman(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	r = r[:n-1]
	for i := len(r) - 1; i >= 0; i-- {
		if unicode.IsSpace(r[i]) {
			return strings.TrimSpace(string(r[:i])) + "…"
		}
	}
	return "…"
}

func humanName(slug string) string {
	return upperFirst(strings.ReplaceAll(strings.TrimPrefix(slug, "w-"), "-", " "))
}

func upperFirst(s string) string {
	r := []rune(s)
	if len(r) > 0 {
		r[0] = unicode.ToUpper(r[0])
	}
	return string(r)
}

// ChartRowGoals reads only measured problem statements on charter row lines.
func ChartRowGoals(charter string) map[int]string {
	out := map[int]string{}
	codeSpan := regexp.MustCompile("`[^`]*`")
	cases := map[string]string{}
	for _, w := range strings.Fields("World AILANG MCP CLI TTY PR CI API JSON Go Codex Claude") {
		cases[strings.ToLower(w)] = w
	}
	for _, line := range strings.Split(charter, "\n") {
		m := queueRow.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		row, _ := strconv.Atoi(m[1])
		// Keep byte offsets so each matched span can restore its original code text.
		masked := codeSpan.ReplaceAllStringFunc(line, func(code string) string {
			return strings.Repeat("\x00", len(code))
		})
		for _, span := range boldSpan.FindAllStringSubmatchIndex(masked, -1) {
			s := line[span[2]:span[3]]
			letters, upper := 0, 0
			for _, r := range s {
				if unicode.IsLetter(r) {
					letters++
					if unicode.IsUpper(r) {
						upper++
					}
				}
			}
			if len([]rune(s)) < 25 || letters == 0 || upper*100 < letters*80 {
				continue
			}
			s = humanWord.ReplaceAllStringFunc(strings.ToLower(s), func(w string) string {
				if c := cases[w]; c != "" {
					return c
				}
				return w
			})
			s = upperFirst(s)
			s = strings.TrimRight(strings.TrimSpace(s), ".!?") + "."
			if _, ok := out[row]; !ok {
				out[row] = s
				if len([]rune(s)) > 300 {
					out[row] = capHuman(s, 299) + "."
				}
			}
			break
		}
	}
	return out
}

// summaries redacts before either bounded text is stored.
func summaries(s string) (string, string) {
	s = strings.Map(func(r rune) rune {
		if strings.ContainsRune("*_`#>", r) {
			return -1
		}
		return r
	}, s)
	s = humanEmail.ReplaceAllString(s, "[email]")
	tokens := strings.Fields(s)
	for i, t := range tokens {
		if strings.HasPrefix(t, "/") || strings.HasPrefix(t, "~/") || strings.Contains(t, "/Users/") || strings.Contains(t, "/home/") {
			tokens[i] = "[path]"
		}
	}
	s = strings.Join(tokens, " ")
	first := s
	r := []rune(s)
	for i, c := range r {
		if i+1 >= 20 && strings.ContainsRune(".!?", c) && (i+1 == len(r) || r[i+1] == ' ') {
			first = string(r[:i+1])
			break
		}
	}
	return capHuman(first, 300), capHuman(s, 4000)
}

// RegradePassed decodes no explanatory or output fields.
func RegradePassed(data []byte, episode, grader string) bool {
	if len(data) > 64<<10 {
		return false
	}
	var r struct {
		Schema  string `json:"schema"`
		Episode string `json:"episode"`
		Grader  string `json:"grader"`
		Passed  bool   `json:"passed"`
		RC      *int   `json:"rc"`
	}
	return json.Unmarshal(data, &r) == nil && r.Schema == "trial-regrade/1" && r.Episode == episode && r.Grader == grader && r.Passed && r.RC != nil && *r.RC == 0
}

func validateHuman(ix Index) error {
	check := func(field, s string, cap int) error {
		if len([]rune(s)) > cap || strings.Contains(s, "/Users/") || strings.Contains(s, "/home/") || humanEmail.MatchString(s) {
			return fmt.Errorf("surface index: %s exceeds cap or contains private text", field)
		}
		return nil
	}
	for _, s := range ix.States {
		if err := check("title", s.Title, 200); err != nil {
			return err
		}
	}
	for _, b := range ix.Branes {
		if err := check("goal", b.Goal, 300); err != nil {
			return err
		}
		if err := check("name", b.Name, 300); err != nil {
			return err
		}
	}
	for _, s := range ix.Steps {
		if err := check("summary", s.Summary, 300); err != nil {
			return err
		}
		if err := check("summary_full", s.SummaryFull, 4000); err != nil {
			return err
		}
		if d := s.Delivered; d != nil {
			if d.Passed < 0 || d.Passed > d.Total || d.Total != s.Graders.Total {
				return fmt.Errorf("surface index: delivered counts invalid")
			}
			seen := map[string]bool{}
			for _, name := range d.Regraded {
				found := false
				for _, failed := range s.Failed {
					found = found || failed == name
				}
				if !found || seen[name] {
					return fmt.Errorf("surface index: regraded name was not a failing grader")
				}
				seen[name] = true
			}
		}
	}
	return nil
}
