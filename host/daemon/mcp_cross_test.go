package daemon

import (
	"encoding/json"
	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/projection"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func mcpCrossSets(t *testing.T, exact bool) ([][]string, [][]string) {
	t.Helper()
	d, db, _ := invocationDaemon(t)
	snap, err := transitionreg.NewReader(d.store).ReadSnapshot(boundedTestContext(t))
	if err != nil {
		t.Fatal(err)
	}
	base := snap.List()[0]
	ids := []string{"tools.echo", "world/recovery-transition/v1", "tools.hidden"}
	effects := []string{"alpha", "beta", "gamma"}
	changes := make([]transitionreg.Change, 0, 3)
	for i, id := range ids {
		desc := base
		desc.ID = id
		desc.Access.Effect = effects[i]
		changes = append(changes, transitionreg.Change{ID: id, Descriptor: &desc})
	}
	if _, err := transitionreg.NewPublisher(d.store, archive.New(db)).PublishSet(boundedTestContext(t), changes); err != nil {
		t.Fatal(err)
	}
	handler := d.Handler()
	cards, tools := [][]string{}, [][]string{}
	for i, effect := range effects[:2] {
		token := mintSessionGrants(t, d, "ep-cross-"+effect, effect)
		req := httptest.NewRequest("GET", "/.well-known/agent.json", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		var card struct{ Skills []struct{ ID string } }
		if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &card) != nil {
			t.Fatalf("card=%d %s", w.Code, w.Body)
		}
		actual := []string{}
		for _, s := range card.Skills {
			actual = append(actual, s.ID)
		}
		sort.Strings(actual)
		if exact && !reflect.DeepEqual(actual, []string{ids[i]}) {
			t.Fatalf("session%s card=%v", effect, actual)
		}
		cards = append(cards, actual)
		req = httptest.NewRequest("POST", "/mcp/", strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"tools/list"}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		w = httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		var wire struct {
			Result struct{ Tools []struct{ Name string } }
		}
		payload := strings.TrimSpace(strings.TrimPrefix(w.Body.String(), "event: message\ndata: "))
		if w.Code != 200 || json.Unmarshal([]byte(payload), &wire) != nil {
			t.Fatalf("MCP=%d %s", w.Code, w.Body)
		}
		names := []string{}
		for _, tool := range wire.Result.Tools {
			names = append(names, tool.Name)
		}
		sort.Strings(names)
		encoded, err := projection.EncodeMCPName(ids[i])
		if err != nil {
			t.Fatal(err)
		}
		if exact && !reflect.DeepEqual(names, []string{encoded}) {
			t.Fatalf("session%s MCP=%v expected=%s", effect, names, encoded)
		}
		if exact {
			decoded := []string{}
			for _, name := range names {
				decoded = append(decoded, strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(name, "_s", "/"), "_d", "."), "_u", "_"))
			}
			if !reflect.DeepEqual(decoded, actual) || len(names) != len(actual) {
				t.Fatalf("decoded MCP/card=%v/%v", decoded, actual)
			}
		}
		tools = append(tools, names)
	}
	if reflect.DeepEqual(cards[0], cards[1]) || reflect.DeepEqual(tools[0], tools[1]) {
		t.Fatal("unequal session control failed")
	}
	return cards, tools
}
func TestMCPCrossSurfaceExactSet(t *testing.T) { mcpCrossSets(t, true) }
func TestMCPAmbientExportsAbsent(t *testing.T) {
	cards, tools := mcpCrossSets(t, false)
	banned := []string{"eprintln", "exit", "flush", "print", "printErr", "println", "readLine", "writeBytes", "submit_feedback", "readBytes", "writeLine", "std/io.writeBytes", "std/io.readBytes"}
	for i := range cards {
		if len(cards[i]) == 0 || len(tools[i]) == 0 {
			t.Fatal("ambient control vacuous")
		}
		for _, bad := range banned {
			encoded, err := projection.EncodeMCPName(bad)
			if err != nil {
				t.Fatal(err)
			}
			for _, id := range cards[i] {
				if id == bad || strings.HasPrefix(id, "std/io.") || strings.HasPrefix(id, "std/io/") {
					t.Fatalf("ambient card=%s", id)
				}
			}
			for _, name := range tools[i] {
				if name == encoded {
					t.Fatalf("ambient MCP=%s", name)
				}
			}
		}
	}
}

// Run the payloads printed in the operator guide verbatim, against one actual
// mounted daemon backed by an ephemeral store and the pinned interpreter.
func TestMCPQuickstartPayloadsVerbatim(t *testing.T) {
	d, _, _ := invocationDaemon(t)
	token := mintSessionGrants(t, d, "quickstart", "world.apply")
	guide, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	section := strings.Split(string(guide), "### 8. Use MCP tools")
	if len(section) != 2 {
		t.Fatal("MCP section absent")
	}
	matches := regexp.MustCompile(`-d '([^']+)'`).FindAllStringSubmatch(section[1], -1)
	if len(matches) != 3 {
		t.Fatalf("payload count=%d", len(matches))
	}
	server := httptest.NewServer(d.Handler())
	defer server.Close()
	client := &http.Client{Timeout: DefaultClientTimeout}
	for i, match := range matches {
		r, err := http.NewRequest("POST", server.URL+"/mcp/", strings.NewReader(match[1]))
		if err != nil {
			t.Fatal(err)
		}
		r.Header.Set("Authorization", "Bearer "+token)
		r.Header.Set("Content-Type", "application/json")
		r.Header.Set("Accept", "application/json, text/event-stream")
		response, err := client.Do(r)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := io.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != 200 || response.Header.Get("Content-Type") != "text/event-stream" || !strings.HasPrefix(string(raw), "event: message\ndata: ") || strings.Contains(string(raw), `"error"`) {
			t.Fatalf("payload%d response=%d %s", i, response.StatusCode, raw)
		}
		t.Logf("VERBATIM payload=%s\nHTTP %d Content-Type:%s\n%s", match[1], response.StatusCode, response.Header.Get("Content-Type"), raw)
	}
	for index := int64(1); index <= 3; index++ {
		entry, ok, err := d.store.GetLogEntry(boundedTestContext(t), index)
		if err != nil || !ok {
			t.Fatalf("journal%d=%t/%v", index, ok, err)
		}
		t.Logf("journal entry=%d hash=%s transition=%s", index, entry.EntryHash, entry.TransitionRef)
	}
}

func TestMCPToolsListMatchesAgentCard(t *testing.T) { mcpCrossSets(t, true) }
