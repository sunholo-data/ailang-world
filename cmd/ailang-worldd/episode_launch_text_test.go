package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func TestLaunchSharedVectors(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("..", "..", "scripts", "trial", "testdata", "launch_vectors.json"))
	if err != nil {
		t.Fatal(err)
	}
	var vectors struct {
		Contexts []struct {
			Name        string          `json:"name"`
			SessionSelf json.RawMessage `json:"session_self"`
			Addr        string          `json:"addr"`
			Expected    string          `json:"expected"`
		} `json:"contexts"`
		Prompts []struct {
			Task     string   `json:"task"`
			Commands []string `json:"commands"`
			Expected string   `json:"expected"`
		} `json:"prompts"`
		Developer struct {
			Context string `json:"context"`
			Value   string `json:"value"`
		} `json:"developer_instructions"`
	}
	if err := json.Unmarshal(data, &vectors); err != nil {
		t.Fatal(err)
	}
	if len(vectors.Contexts) < 5 || len(vectors.Prompts) < 3 {
		t.Fatal("empty vectors")
	}
	for _, c := range vectors.Contexts {
		got, err := launchContext(c.SessionSelf, c.Addr, strings.Repeat("a", 64))
		if err != nil || got != c.Expected {
			t.Errorf("%s: %v\ngot %q\nwant %q", c.Name, err, got, c.Expected)
		}
	}
	for _, p := range vectors.Prompts {
		if got := launchPrompt(p.Task, p.Commands); got != p.Expected {
			t.Errorf("prompt: %q != %q", got, p.Expected)
		}
	}
	if got := pythonJSONString(vectors.Developer.Context); got != vectors.Developer.Value {
		t.Errorf("developer instructions: %s != %s", got, vectors.Developer.Value)
	}
	if _, err := launchContext([]byte(`{"episode":"secret"}`), "127.0.0.1:7644", "secret"); err == nil {
		t.Fatal("context token accepted")
	}
	if got := pythonJSONString("<>& · 🚀\x7f"); got != `"<>& \u00b7 \ud83d\ude80\u007f"` {
		t.Fatal(got)
	}
}

func TestLaunchArgvFloorGoldens(t *testing.T) {
	grants, err := sessionNewGrants("se-tools", nil, 50)
	if err != nil {
		t.Fatal(err)
	}
	var floor []broker.Capability
	for _, g := range grants {
		if g.Effect != broker.EffectWorkspaceDelete && g.Effect != broker.EffectWorkspaceExec {
			floor = append(floor, g)
		}
	}
	for _, agent := range []string{"claude", "codex"} {
		got := launchArgv(agent, "$M", "$P", "<void/ep>", "<private>/<ep>.mcp.json", "127.0.0.1:7644", "", "5.00", floor)
		data, err := os.ReadFile(filepath.Join("..", "..", "scripts", "floor", "testdata", "argv_"+agent+"_world.json"))
		if err != nil {
			t.Fatal(err)
		}
		// The goldens use one-space JSON indentation.
		var encoded bytes.Buffer
		encoder := json.NewEncoder(&encoded)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", " ")
		_ = encoder.Encode(got)
		if encoded.String() != string(data) {
			t.Errorf("%s golden differs:\n%s", agent, encoded.String())
		}
		emptyModel := launchArgv(agent, "", "$P", "/void", "/private/mcp", "127.0.0.1:7644", "context", "", floor)
		for _, a := range emptyModel {
			if a == "-m" || a == "--model" || a == "--max-budget-usd" {
				t.Fatal("unexpected model/budget", emptyModel)
			}
		}
		if agent == "codex" && !reflect.DeepEqual(emptyModel[len(emptyModel)-3:], []string{"-c", "developer_instructions=\"context\"", "$P"}) {
			t.Fatal(emptyModel)
		}
	}
}
