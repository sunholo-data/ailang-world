package store

import (
	"context"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/hashref"
)

var storeSurfaceAllow = map[string]bool{"BusyTimeout": true, "Close": true}
var intendedNonIO = []string{"BusyTimeout", "Close"}

func guardedMethods(t *testing.T) []reflect.Method {
	t.Helper()
	typ := reflect.TypeOf(&Store{})
	var names []string
	var methods []reflect.Method
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		if storeSurfaceAllow[method.Name] {
			continue
		}
		names = append(names, method.Name)
		methods = append(methods, method)
		if method.Type.NumIn() < 2 || method.Type.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			t.Errorf("%s: first argument is not context.Context", method.Name)
		}
	}
	slices.Sort(names)
	_ = names
	return methods
}

func TestStoreSurfaceTakesContext(t *testing.T) {
	typ := reflect.TypeOf(&Store{})
	var excluded []string
	for i := 0; i < typ.NumMethod(); i++ {
		if storeSurfaceAllow[typ.Method(i).Name] {
			excluded = append(excluded, typ.Method(i).Name)
		}
	}
	slices.Sort(excluded)
	if !slices.Equal(excluded, intendedNonIO) {
		t.Fatalf("non-I/O surface = %v, want %v", excluded, intendedNonIO)
	}
	if len(storeSurfaceAllow) != len(intendedNonIO) {
		t.Fatalf("allow-list grew: %v", storeSurfaceAllow)
	}
	for _, m := range guardedMethods(t) {
		if m.Type.NumIn() < 2 || m.Type.In(1) != reflect.TypeOf((*context.Context)(nil)).Elem() {
			t.Fatalf("%s lacks context", m.Name)
		}
	}
}

type guardFixture func(*testing.T, *Store, context.Context) []any

func guardFixtures() map[string]guardFixture {
	ref := func(s string) hashref.HashRef { return hashref.SumSHA256([]byte(s)) }
	simple := func(args ...any) guardFixture { return func(*testing.T, *Store, context.Context) []any { return args } }
	obj1 := obj("guard object", "guard")
	world := World{Ref: ref("world"), Revision: 1, StateRoot: ref("state"), LogHead: ref("head")}
	intent := JournalIntent{InvocationID: "guard", WorldRef: ref("world"), EntryHash: ref("entry"), PrevEntryHash: ref("prev"), TransitionFn: ref("transition"), TransitionRef: ref("body"), Interpreter: ref("interpreter")}
	effect := EffectIntent{EpisodeID: "episode", Effect: "publish", Scope: "scope", RequestRef: ref("request")}
	commit := Commit{NextWorld: world, Entry: LogEntry{Header: LogHeader{EntryIndex: 1, SemanticsEpoch: 1, TransitionFn: ref("transition"), Interpreter: ref("interpreter"), PrevEntryHash: ref("prev")}, EntryHash: ref("entry"), TransitionRef: ref("body")}}
	return map[string]guardFixture{
		"AppendIntent":              simple("guard", intent),
		"AppendNextEffectIntent":    simple("episode", effect),
		"AppendClaimedEffectIntent": simple("episode", effect, ref("approval"), ref("request")),
		"AppendOutcome": func(t *testing.T, s *Store, ctx context.Context) []any {
			if _, _, err := s.AppendIntent(ctx, "guard", intent); err != nil {
				t.Fatal(err)
			}
			return []any{"guard", JournalOutcome{InvocationID: "guard", Status: "succeeded", ResultRef: ref("result")}}
		},
		"AppendEffectOutcome": func(t *testing.T, s *Store, ctx context.Context) []any {
			id, _, err := s.AppendNextEffectIntent(ctx, "episode", effect)
			if err != nil {
				t.Fatal(err)
			}
			return []any{id, EffectOutcome{InvocationID: id, Status: "succeeded", RecordRef: ref("record")}}
		},
		"Commit": simple(commit), "CommitLanded": simple(commit),
		"CompareAndSetRegistryHead": func(t *testing.T, s *Store, ctx context.Context) []any {
			next := obj("next", "guard")
			if err := s.PutObject(ctx, next); err != nil {
				t.Fatal(err)
			}
			return []any{"name", hashref.HashRef{}, next.Hash}
		},
		"GetEffectReceipt": simple("effect:episode:0"), "GetLogEntry": simple(int64(1)),
		"GetObject": simple(ref("missing")), "GetReceipt": simple("guard"),
		"GetRegistryHead": simple("name"), "GetVerifyResult": simple(ref("transition"), ref("interpreter")),
		"GetWorld": simple(ref("missing")), "MintSession": simple(SessionRow{CredentialID: "credential", EpisodeID: "episode", GrantsJSON: "[]"}),
		"ObjectCommits":        simple(ref("missing"), int64(0), 1),
		"ObjectReferences":     simple(ref("missing"), (*ObjectReferenceCursor)(nil), 1),
		"ObjectsBySemanticID":  simple("missing", "", 1),
		"PendingEffectIntents": simple(1), "PendingIntents": simple(1),
		"PutObject": simple(obj1), "PutVerifyResult": simple(VerifyResult{TransitionFn: ref("transition"), Interpreter: ref("interpreter"), Verified: true}),
		"PutWorld": simple(world), "ReadObject": simple(ref("missing"), int64(1024)),
		"ResolveSession": simple("missing"), "RevokeSession": simple("missing"),
		"ScanUnreadableLog": simple(int64(0), 1), "ScanUnreadableWorlds": simple("", 1),
		"SelectHead": simple(ref("world")), "SelectedHead": simple(),
		"SetRegistryHead": simple("name", ref("object")),
	}
}

func TestStoreDeadlineGuardDynamic(t *testing.T) {
	methods := guardedMethods(t)
	assertGuardAtEntry(t, methods)
	fixtures := guardFixtures()
	reflected := map[string]bool{}
	for _, m := range methods {
		reflected[m.Name] = true
	}
	if len(reflected) != len(fixtures) {
		t.Fatalf("fixture key-set mismatch: reflected=%v fixtures=%v", reflected, reflect.ValueOf(fixtures).MapKeys())
	}
	for name := range reflected {
		if _, ok := fixtures[name]; !ok {
			t.Fatalf("missing fixture for %s", name)
		}
	}
	for name := range fixtures {
		if !reflected[name] {
			t.Fatalf("extra fixture for %s", name)
		}
	}
	for _, m := range methods {
		for _, kind := range []string{"nil", "background", "deadline"} {
			t.Run(m.Name+"/"+kind, func(t *testing.T) {
				defer func() {
					if caught := recover(); caught != nil {
						t.Errorf("%s/%s panicked instead of returning an error: %v", m.Name, kind, caught)
					}
				}()
				s, openErr := Open(":memory:")
				if openErr != nil {
					t.Fatal(openErr)
				}
				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()
				args := fixtures[m.Name](t, s, ctx)
				var callctx context.Context
				switch kind {
				case "background":
					callctx = context.Background()
				case "deadline":
					callctx = ctx
				}
				values := []reflect.Value{reflect.ValueOf(s), reflect.ValueOf(&callctx).Elem()}
				for _, arg := range args {
					values = append(values, reflect.ValueOf(arg))
				}
				out := m.Func.Call(values)
				defer func() { _ = s.Close() }()
				var err error
				if e := out[len(out)-1].Interface(); e != nil {
					err = e.(error)
				}
				if kind == "deadline" {
					if err != nil {
						t.Fatalf("deadline call: %v", err)
					}
				} else if !errors.Is(err, ErrNoDeadline) {
					t.Fatalf("%s call = %v, want ErrNoDeadline", kind, err)
				}
			})
		}
	}
}

// A delegated read can return ErrNoDeadline even if its own entry guard is
// removed (CommitLanded calls GetLogEntry), so pin the entry order as well.
func assertGuardAtEntry(t *testing.T, methods []reflect.Method) {
	t.Helper()
	want := map[string]bool{}
	for _, method := range methods {
		want[method.Name] = true
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), entry.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || !want[fn.Name.Name] || fn.Recv == nil {
				continue
			}
			if len(fn.Body.List) < 2 || !entryGuardCall(fn.Body.List[0], "checkQuarantine") || !entryGuardCall(fn.Body.List[1], "requireDeadline") {
				t.Errorf("%s: first two statements must check quarantine then deadline", fn.Name.Name)
			}
			delete(want, fn.Name.Name)
		}
	}
	for name := range want {
		t.Errorf("%s: no source declaration", name)
	}
}

func entryGuardCall(stmt ast.Stmt, name string) bool {
	ifStmt, ok := stmt.(*ast.IfStmt)
	if !ok {
		return false
	}
	init, ok := ifStmt.Init.(*ast.AssignStmt)
	if !ok || len(init.Rhs) != 1 {
		return false
	}
	call, ok := init.Rhs[0].(*ast.CallExpr)
	if !ok {
		return false
	}
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		return fun.Name == name
	case *ast.SelectorExpr:
		return fun.Sel.Name == name
	}
	return false
}
