package main

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestPresetCommitterValidationAndPrompt(t *testing.T) {
	for _, value := range []string{"Test Operator <operator@example.invalid>", "", "Name", "<email>", "Name <>", "Name <email> extra", "Name <<email>>", "Name\n <email>", "Name\r <email>"} {
		t.Run(value, func(t *testing.T) {
			f := newPresetFixture(t)
			var p launchPreset
			if err := json.Unmarshal(f.data, &p); err != nil {
				t.Fatal(err)
			}
			p.Committer = value
			data, _ := json.Marshal(p)
			if err := os.WriteFile(f.file, data, 0600); err != nil {
				t.Fatal(err)
			}
			term := &promptTerm{in: strings.NewReader("n\n")}
			var out, errs bytes.Buffer
			code := runEpisodeStart(f.args(), &out, &errs, presetEnv(term))
			valid := value == "" || value == "Test Operator <operator@example.invalid>"
			if !valid {
				if code != 1 || !strings.Contains(errs.String(), "invalid committer") || term.out.Len() != 0 {
					t.Fatalf("%d %s %s", code, errs.String(), term.out.String())
				}
			} else if !strings.Contains(term.out.String(), "  committer "+displayCommitter(value)+"\n") {
				t.Fatal(term.out.String(), errs.String())
			}
		})
	}
}
