package main

// Row 243 V3 sprint 1 M3 (design §2.2, §6 VT1): the bytes finishPrompt writes at the terminal, for
// three fixtures (draft and ready, create-label, protected paths, an empty committer), pinned to the
// golden that the base tree's finishPrompt generated before the text moved to
// broker.FinishPromptText (WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1 on the base tree, plan §1.2). The
// broker's TestFinishPromptText reads the same golden; displayCommitter (episode start's preset
// prompt) is pinned to the same "" rule.

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

const finishPromptGolden = "../../host/broker/testdata/finish-prompt.golden.json"

// finishPromptFixture is one golden row: finishPrompt's inputs and the bytes it wrote.
type finishPromptFixture struct {
	Name      string
	Scope     broker.FinishApprovalScope
	Title     string
	Body      string
	Protected []string
	Committer string
	Prompt    string
}

func finishPromptFixtures() []finishPromptFixture {
	scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1",
		SHA: strings.Repeat("a1", 20), Base: "dev", TitleSHA256: strings.Repeat("b2", 32), BodySHA256: strings.Repeat("c3", 32),
		Draft: true, ExpiresAt: 1700000900}
	ready := scope
	ready.Draft, ready.CreateLabel, ready.Head, ready.Base = false, true, "world/w2", "main"
	label := scope
	label.CreateLabel = true
	return []finishPromptFixture{
		{Name: "draft, existing label, no protected path, repository committer", Scope: scope, Title: "Row 243: the law",
			Body: "## Episode w1\n\nline 3\nline 4\nline 5\nline 6\nline 7 (not shown)\nline 8\n"},
		{Name: "ready, label created, two protected paths, pinned committer", Scope: ready, Title: "Tidy <b> & \"quote\"",
			Body: "one\ntwo\nthree", Protected: []string{"cmd/ailang-worldd/main.go", "go.mod"}, Committer: "World Bot <bot@example.invalid>"},
		{Name: "draft, label created, one protected path, empty body", Scope: label, Title: "Café — naïve",
			Body: "", Protected: []string{"docs/QUICKSTART.md"}, Committer: ""},
	}
}

func TestFinishPromptIsTheGolden(t *testing.T) {
	got := finishPromptFixtures()
	for i := range got {
		var tty bytes.Buffer
		finishPrompt(&tty, got[i].Scope, got[i].Title, got[i].Body, got[i].Protected, got[i].Committer)
		got[i].Prompt = tty.String()
	}
	if os.Getenv("WORLD_UPDATE_FINISH_PROMPT_GOLDEN") == "1" {
		b, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(finishPromptGolden, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	raw, err := os.ReadFile(finishPromptGolden)
	if err != nil {
		t.Fatal(err)
	}
	var want []finishPromptFixture
	if err := json.Unmarshal(raw, &want); err != nil {
		t.Fatal(err)
	}
	if len(want) != 3 || !reflect.DeepEqual(got, want) {
		t.Fatalf("finishPrompt's bytes moved from the base's:\n got %+v\nwant %+v", got, want)
	}
	var none bytes.Buffer // the variadic committer omitted is the empty committer
	finishPrompt(&none, want[0].Scope, want[0].Title, want[0].Body, want[0].Protected)
	if none.String() != want[0].Prompt {
		t.Fatalf("no committer argument: %q", none.String())
	}
	if got := displayCommitter(""); got != "(repository git config)" { // episode start's preset prompt shares the rule
		t.Fatalf("displayCommitter(\"\") = %q", got)
	}
}
