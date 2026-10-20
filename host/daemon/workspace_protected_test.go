package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
)

func protectedProfile(t *testing.T, project string, protected ...string) *broker.ExecProfile {
	t.Helper()
	m := map[string]any{"profile": broker.ExecProfileVersion, "project": project,
		"path": []string{"/usr/bin", "/bin"}, "timeout_ms": 5000, "protected": protected,
		"commands": map[string]any{"where": map[string]any{"argv": []string{"/bin/pwd"}}}}
	if protected == nil {
		delete(m, "protected")
	}
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	p, err := broker.ParseExecProfile(data)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

// The fake records argv's policy bytes. It reports actual summary fields and
// digest to exercise the production run-variant verifier; it enforces nothing.
func protectedTool(t *testing.T, f wsFixture) string {
	t.Helper()
	bin := fakeToolBin(t, f.logDir, ToolBinaryRelease)
	data, err := os.ReadFile(bin)
	if err != nil {
		t.Fatal(err)
	}
	script := strings.Replace(string(data), `  req=$(cat)`, `  req=$(cat)
  cp "$3" "$log/policy-tool.toml"`, 1)
	start := strings.Index(script, "      printf '")
	end := strings.Index(script[start:], "      exit 0") + start
	script = script[:start] + `      mode=$(awk -F ' = ' '$1 == "security_mode" {print $2}' "$3")
      caps=$(awk -F ' = ' '$1 == "allowed_caps" {print $2}' "$3")
      root=$(awk -F ' = ' '$1 == "fs_sandbox" {print $2}' "$3")
      digest=$(shasum -a 256 "$3" | cut -d ' ' -f 1)
      printf '{"ok":true,"summary":{"security_mode":%s,"caps":%s,"fs_sandbox":%s,"policy_digest":"%s","cli":["check"]}}' "$mode" "$caps" "$root" "$digest"
` + script[end:]
	script = strings.Replace(script, `echo x >> "$log/dispatches"
exit 0`, `cp "$3" "$log/run.toml"
echo x >> "$log/dispatches"
printf 'policy: {"ok":true,"decision":{"ok":true}}\n' >&2
exit 0`, 1)
	writeFile(t, bin, script)
	return bin
}

func protectedRead(t *testing.T, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func requireProtectedPolicy(t *testing.T, path, sandbox string, want []string) {
	t.Helper()
	expected, err := broker.RenderEpisodePolicy(sandbox, want)
	if err != nil {
		t.Fatal(err)
	}
	if got := protectedRead(t, path); !bytes.Equal(got, expected) {
		t.Fatalf("policy bytes differ: got\n%s\nwant\n%s", got, expected)
	}
}

func protectedCall(t *testing.T, ctx context.Context, h broker.Handler, effect string, payload any) []byte {
	t.Helper()
	if h == nil {
		t.Fatal("registry has no AILANG handler")
	}
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	out, err := h.Execute(ctx, broker.EffectRequest{Effect: effect, Scope: broker.WorkspaceScope, Cost: 1}, body)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func TestEpisodePolicyProtectedRebased(t *testing.T) {
	for _, tc := range []struct {
		name, root     string
		patterns, want []string
	}{
		{"under-sub", "sub", []string{"sub/tests/**", "*_test.go", "other/**"}, []string{"*_test.go", "tests/**"}},
		{"ancestor", "sub", []string{"sub/**"}, []string{"*"}},
		{"ancestor-above", "sub/deep", []string{"sub/**"}, []string{"*"}},
		{"disjoint", "sub", []string{"other/**"}, nil},
		{"basename", "sub", []string{"*_test.go"}, []string{"*_test.go"}},
		{"root-dot", ".", []string{"sub/tests/**", "*_test.go", "other/**"}, []string{"*_test.go", "other/**", "sub/tests/**"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := boundedTestContext(t)
			f := newWSFixture(t)
			sandbox := filepath.Join(f.root, "ep1", tc.root)
			mkdirs(t, sandbox)
			d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: tc.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease)})
			d.workspace.exec = &workspaceExec{profile: protectedProfile(t, "alpha", tc.patterns...)}
			protectedCall(t, ctx, d.workspace.registry("ep1")[broker.EffectWorkspaceRead], broker.EffectWorkspaceRead, map[string]string{"op": "read", "path": "x"})
			requireProtectedPolicy(t, filepath.Join(f.stateDir, "policies/ep1.toml"), sandbox, tc.want)
		})
	}
	t.Run("no-cross-episode-leak", func(t *testing.T) {
		ctx := boundedTestContext(t)
		f := newWSFixture(t)
		mkdirs(t, filepath.Join(f.root, "ep1/sub"), filepath.Join(f.root, "ep2/lib"))
		d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: fakeToolBin(t, f.logDir, ToolBinaryRelease), WorkspaceEpisodeModuleRoots: []string{"ep1=sub", "ep2=lib"}})
		d.workspace.exec = &workspaceExec{byProject: map[string]*broker.ExecProfile{"alpha": protectedProfile(t, "alpha", "sub/alpha/**"), "beta": protectedProfile(t, "beta", "lib/beta/**")}, episodeProject: map[string]string{"ep1": "alpha", "ep2": "beta"}}
		for i, ep := range []string{"ep1", "ep2"} {
			protectedCall(t, ctx, d.workspace.registry(ep)[broker.EffectWorkspaceRead], broker.EffectWorkspaceRead, map[string]string{"op": "read", "path": "x"})
			requireProtectedPolicy(t, filepath.Join(f.stateDir, "policies", ep+".toml"), filepath.Join(f.root, ep, []string{"sub", "lib"}[i]), []string{[]string{"alpha/**", "beta/**"}[i]})
		}
	})
}

func TestWorkspaceProtectedPolicyWiring(t *testing.T) {
	for _, name := range []string{"episode-selection", "no-exec", "no-protected", "no-profile", "rebind", "binding-error"} {
		t.Run(name, func(t *testing.T) {
			ctx := boundedTestContext(t)
			f := newWSFixture(t)
			var log bytes.Buffer
			d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: protectedTool(t, f), ErrorLog: &log})
			a, b := protectedProfile(t, "alpha", "alpha/**"), protectedProfile(t, "beta", "beta/**")
			d.workspace.exec = &workspaceExec{profile: a, byProject: map[string]*broker.ExecProfile{"alpha": a, "beta": b}, episodeProject: map[string]string{"ep1": "alpha", "ep2": "beta"}}
			var current execBinding
			d.workspace.bindingFor = func(string) (execBinding, error) { return current, nil }
			switch name {
			case "no-exec":
				d.workspace.exec = nil
			case "no-protected":
				d.workspace.exec = &workspaceExec{profile: protectedProfile(t, "alpha")}
			case "no-profile":
				d.workspace.exec = &workspaceExec{}
			case "binding-error":
				// Both store errors and conflicting live credentials must fail closed.
				for _, cause := range []error{errors.New("binding read failed"), &bindingConflictError{episode: "ep1", bindings: []string{"a", "b"}}} {
					d.workspace.bindingFor = func(string) (execBinding, error) { return execBinding{}, cause }
					if h := d.workspace.registry("ep1")[broker.EffectWorkspaceRead]; h != nil {
						t.Fatal("binding error left an AILANG handler bound")
					}
					if !strings.Contains(log.String(), cause.Error()) {
						t.Fatalf("binding error not reported: %s", &log)
					}
				}
				if f.summaries(t) != 0 || f.dispatches(t) != 0 {
					t.Fatal("binding error spawned tool")
				}
				return
			}
			read := func(ep string) broker.Handler {
				h := d.workspace.registry(ep)[broker.EffectWorkspaceRead]
				protectedCall(t, ctx, h, broker.EffectWorkspaceRead, map[string]string{"op": "read", "path": "x"})
				return h
			}
			read("ep1")
			want := []string{"alpha/**"}
			if name == "no-exec" || name == "no-protected" || name == "no-profile" {
				want = nil
			}
			requireProtectedPolicy(t, filepath.Join(f.stateDir, "policies/ep1.toml"), filepath.Join(f.root, "ep1"), want)
			if name == "episode-selection" {
				read("ep2")
				requireProtectedPolicy(t, filepath.Join(f.stateDir, "policies/ep2.toml"), filepath.Join(f.root, "ep2"), []string{"beta/**"})
			}
			if name == "rebind" {
				current = execBinding{Project: "beta", Digest: b.Digest}
				read("ep1")
				requireProtectedPolicy(t, filepath.Join(f.logDir, "policy-tool.toml"), filepath.Join(f.root, "ep1"), []string{"beta/**"})
				summaries := f.summaries(t)
				read("ep1")
				if summaries != 2 || f.summaries(t) != summaries {
					t.Fatalf("rebind/cache summaries = %d then %d, want 2 then 2", summaries, f.summaries(t))
				}
			}
		})
	}
}

// A caller whose protected set differs from the in-flight build's must not
// join that build: it waits for it, then builds its own handler under its own
// set (MUT-JOIN-IGNORES-SET, which would hand it the held build's policy).
func TestWorkspaceProtectedPolicyJoinRequiresSameSet(t *testing.T) {
	withHandlerBudget(t, 5*time.Second)
	r := newFamilyRig(t, Config{})
	var pending []*pendingRegistry
	r.joinOnCleanup([]string{"ep1"}, &pending)
	a, b := protectedProfile(t, "alpha", "alpha/**"), protectedProfile(t, "beta", "beta/**")
	ex := r.d.workspace.exec
	ex.profile, ex.byProject = a, map[string]*broker.ExecProfile{"alpha": a, "beta": b}
	var mu sync.Mutex
	current := execBinding{}
	r.d.workspace.bindingFor = func(string) (execBinding, error) {
		mu.Lock()
		defer mu.Unlock()
		return current, nil
	}

	r.set("hold", "ep1")
	first := r.startRegistry("ep1")
	pending = append(pending, first)
	waitFor(t, "ep1's first summary", 30*time.Second, func() bool { return r.summaries("ep1") == 1 })
	mu.Lock()
	current = execBinding{Project: "beta", Digest: b.Digest}
	mu.Unlock()
	second := r.startRegistry("ep1")
	pending = append(pending, second)
	// The second caller must wait without joining; a joiner shows up at once.
	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) && r.buildWaiters("ep1") == 0 {
		time.Sleep(10 * time.Millisecond)
	}
	if n := r.buildWaiters("ep1"); n != 0 {
		t.Fatalf("a caller under another protected set joined the in-flight build (waiters=%d)", n)
	}
	r.clear("hold", "ep1")
	<-first.done
	<-second.done
	if n := r.summaries("ep1"); n != 2 {
		t.Fatalf("summaries = %d, want 2: one build per protected set", n)
	}
	h1, h2 := first.reg[broker.EffectWorkspaceRead], second.reg[broker.EffectWorkspaceRead]
	if h1 == nil || h2 == nil || h1 == h2 {
		t.Fatalf("handlers = %v, %v: want two distinct AILANG handlers", h1, h2)
	}
	requireProtectedPolicy(t, filepath.Join(r.f.stateDir, "policies/ep1.toml"), filepath.Join(r.f.root, "ep1"), []string{"beta/**"})
}

func TestWorkspaceProtectedPolicyRequests(t *testing.T) {
	ctx := boundedTestContext(t)
	f := newWSFixture(t)
	mkdirs(t, filepath.Join(f.root, "ep1/sub"))
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, WorkspaceModuleRoot: "sub", ToolAilangBin: protectedTool(t, f), RunAllowCaps: []string{"Declassify"}})
	d.workspace.exec = &workspaceExec{profile: protectedProfile(t, "alpha", "sub/tests/**", "*_test.go", "other/**")}
	reg := d.workspace.registry("ep1")
	for _, op := range []string{"write", "read"} {
		effect := broker.EffectWorkspaceWrite
		payload := map[string]string{"op": op, "path": "ok.txt"}
		if op == "write" {
			payload["content"] = "hello"
		} else {
			effect = broker.EffectWorkspaceRead
		}
		before := f.dispatches(t)
		protectedCall(t, ctx, reg[effect], effect, payload)
		if f.dispatches(t) != before+1 {
			t.Fatal("request never reached policy-tool")
		}
		requireProtectedPolicy(t, filepath.Join(f.logDir, "policy-tool.toml"), filepath.Join(f.root, "ep1/sub"), []string{"*_test.go", "tests/**"})
	}
	caps := []string{"Declassify", "FS", "IO"}
	protectedCall(t, ctx, reg[broker.EffectAilangRun], broker.EffectAilangRun, map[string]any{"path": "main.ail", "caps": caps})
	want, err := broker.RenderRunVariant(filepath.Join(f.root, "ep1/sub"), caps, broker.RunCapsConfig{Allow: []string{"Declassify"}}, []string{"*_test.go", "tests/**"})
	if err != nil {
		t.Fatal(err)
	}
	if got := protectedRead(t, filepath.Join(f.logDir, "run.toml")); !bytes.Equal(got, want) {
		t.Fatalf("run variant lost protected entries: got\n%s\nwant\n%s", got, want)
	}
}

func TestWorkspaceWriteProtectedPolicyTool(t *testing.T) {
	bin := realToolBin(t) // Skip before any fixture when the optional pin is unset.
	ctx := boundedTestContext(t)
	f := newWSFixture(t)
	mkdirs(t, filepath.Join(f.root, "ep1/tests"))
	for _, path := range []string{"tests/x.ail", "a_test.go"} {
		writeFile(t, filepath.Join(f.root, "ep1", path), "seed\n")
	}
	d := mustWSDaemon(t, Config{DBPath: f.db, WorkspaceRoot: f.root, ToolAilangBin: bin})
	d.workspace.exec = &workspaceExec{profile: protectedProfile(t, "alpha", "tests/**", "*_test.go")}
	h := d.workspace.registry("ep1")[broker.EffectWorkspaceWrite]
	for _, path := range []string{"tests/x.ail", "a_test.go"} {
		for _, op := range []string{"write", "edit"} {
			payload := map[string]string{"op": op, "path": path}
			if op == "write" {
				payload["content"] = "changed\n"
			} else {
				payload["old_text"] = "seed"
				payload["new_text"] = "changed"
			}
			out := protectedCall(t, ctx, h, broker.EffectWorkspaceWrite, payload)
			var response map[string]any
			if json.Unmarshal(out, &response) != nil || response["ok"] != false || !strings.Contains(string(out), "fs_deny_write") {
				t.Fatalf("%s %s: want fs_deny_write refusal, got %s", op, path, out)
			}
			if got := protectedRead(t, filepath.Join(f.root, "ep1", path)); string(got) != "seed\n" {
				t.Fatalf("protected host bytes changed: %s", path)
			}
		}
	}
	out := protectedCall(t, ctx, h, broker.EffectWorkspaceWrite, map[string]string{"op": "write", "path": "ok.txt", "content": "allowed\n"})
	if !strings.Contains(string(out), `"ok":true`) || string(protectedRead(t, filepath.Join(f.root, "ep1/ok.txt"))) != "allowed\n" {
		t.Fatalf("unprotected write failed: %s", out)
	}
}
