package broker

// Row 201 A4 sprint 1 (design §2.4, §2.6 L2, §5.1, §6 BR1, BR2, R4b, R4c, XG1): the first
// decision of a request is final; a key decision never stamps a publish; who may call the decide
// entry point.

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// openDecideStore is a world store under t.TempDir().
func openDecideStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

// pendingRequest files one Human.Approve request through the landed handler.
func pendingRequest(t *testing.T, s approvalStore, scope, requester string, now int64) hashref.HashRef {
	t.Helper()
	out, err := newHumanHandler(s).Execute(boundedTestContext(t), EffectRequest{Effect: EffectHumanApprove, Scope: scope, Cost: 1, Now: now},
		mustApprovalJSON(approvalInputWire{Requester: requester}))
	if err != nil {
		t.Fatal(err)
	}
	ref, err := ParsePendingApproval(out)
	if err != nil {
		t.Fatal(err)
	}
	return ref
}

// BR1: deny, then approve → ErrApprovalAlreadyDecided, the head unchanged, the poll still deny (the
// design's V9 probe as a test); a deny by another decider or at another time is refused too; the
// byte-identical decision again answers its existing ref and writes nothing (plan §1.4 P2).
func TestDecideApprovalRefusesASecondDecision(t *testing.T) {
	s := openDecideStore(t)
	ctx := boundedTestContext(t)
	ref := pendingRequest(t, s, WorkspaceScope, "agent", 100)
	first, err := DecideApproval(ctx, s, ref, "deny", "operator:terminal", 101)
	if err != nil {
		t.Fatal(err)
	}
	head, _, err := s.GetRegistryHead(ctx, ApprovalsV1)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		verb, by string
		now      int64
	}{{"approve", "ssh:mallory", 102}, {"approve", "operator:terminal", 101}, {"deny", "operator:terminal", 102}, {"deny", "ssh:mallory", 101}} {
		if got, err := DecideApproval(ctx, s, ref, tc.verb, tc.by, tc.now); !errors.Is(err, ErrApprovalAlreadyDecided) || !got.IsZero() {
			t.Fatalf("second decision %+v: %v, %v", tc, got, err)
		}
	}
	if again, err := DecideApproval(ctx, s, ref, "deny", "operator:terminal", 101); err != nil || again != first {
		t.Fatalf("the identical decision again: %v, %v; want %v", again, err, first)
	}
	if after, _, err := s.GetRegistryHead(ctx, ApprovalsV1); err != nil || after != head {
		t.Fatalf("head moved: %v -> %v (%v)", head, after, err)
	}
	out, err := newHumanHandler(s).Execute(ctx, EffectRequest{Effect: EffectHumanPollApproval, Scope: WorkspaceScope, Now: 103},
		mustApprovalJSON(approvalInputWire{RequestRef: ref.String()}))
	if err != nil || !strings.Contains(string(out), `"decision":"deny","decidedBy":"operator:terminal"`) {
		t.Fatalf("poll = %s, %v", out, err)
	}
	if on, decided, verb, by, err := ApprovalChainState(ctx, s, ref); err != nil || !on || !decided || verb != "deny" || by != "operator:terminal" {
		t.Fatalf("chain state = %v %v %q %q %v", on, decided, verb, by, err)
	}
}

// BR2: an undecided request decides (the control), and the chain state reads it before and after;
// a ref that is not on the chain reads as absent.
func TestDecideApprovalFirstDecisionMustLive(t *testing.T) {
	s := openDecideStore(t)
	ctx := boundedTestContext(t)
	ref := pendingRequest(t, s, WorkspaceScope, "agent", 100)
	if on, decided, _, _, err := ApprovalChainState(ctx, s, ref); err != nil || !on || decided {
		t.Fatalf("pending chain state = %v %v %v", on, decided, err)
	}
	other := pendingRequest(t, s, WorkspaceScope, "other", 100)
	if _, err := DecideApproval(ctx, s, other, "deny", "operator:terminal", 101); err != nil {
		t.Fatal(err)
	}
	dec, err := DecideApproval(ctx, s, ref, "approve", "operator:terminal", 102)
	if err != nil || dec.IsZero() {
		t.Fatalf("first decision: %v, %v", dec, err)
	}
	if on, decided, verb, by, err := ApprovalChainState(ctx, s, ref); err != nil || !on || !decided || verb != "approve" || by != "operator:terminal" {
		t.Fatalf("decided chain state = %v %v %q %q %v", on, decided, verb, by, err)
	}
	if on, decided, _, _, err := ApprovalChainState(ctx, s, hashref.SumSHA256([]byte("nowhere"))); err != nil || on || decided {
		t.Fatalf("absent chain state = %v %v %v", on, decided, err)
	}
	if _, _, _, _, err := ApprovalChainState(ctx, openDecideStore(t), ref); err != nil {
		t.Fatalf("an empty chain: %v", err)
	}
}

// landPublishDecidedBy lands a publish stamp for f decided by decidedBy, through the landed request
// handler and decideApproval, and returns the payload naming it.
func landPublishDecidedBy(t *testing.T, base approvalStore, f publishFixture, decidedBy string) []byte {
	t.Helper()
	request := pendingRequest(t, base, PublishApprovalScopeFor(f.identity, f.hashes, 100), "r4-"+decidedBy, 10)
	decision, err := decideApproval(boundedTestContext(t), base, request, "approve", decidedBy, 11)
	if err != nil {
		t.Fatal(err)
	}
	f.identity.ApprovalRef = decision
	return EncodePublishPayload(f.identity, f.hashes)
}

// R4b: a publish-shaped request decided by a key (ssh:, webauthn:) never stamps a publish.
func TestPublishApprovalRefusesAKeyDecision(t *testing.T) {
	base := openPublishStore(t)
	f := newPublishFixture(t, "http://127.0.0.1:1", "r4b")
	req := EffectRequest{Effect: EffectRegistryPublish, Scope: f.scope, Cost: PublishCost, Now: 50}
	for _, by := range []string{"ssh:mark", "webauthn:mark-mac"} {
		_, err := validatePublishApproval(boundedTestContext(t), base, landPublishDecidedBy(t, base, f, by), req)
		if !errors.Is(err, ErrPublishApprovalDenied) || !strings.Contains(err.Error(), "publishing stays typed-only (D-WORLD-89 R4)") || !strings.Contains(err.Error(), "was made by "+by) {
			t.Errorf("%s: %v", by, err)
		}
	}
}

// R4c: the same stamp decided at a terminal, or by a free-text name, is accepted (the control).
func TestPublishApprovalTypedStampMustLive(t *testing.T) {
	base := openPublishStore(t)
	f := newPublishFixture(t, "http://127.0.0.1:1", "r4c")
	req := EffectRequest{Effect: EffectRegistryPublish, Scope: f.scope, Cost: PublishCost, Now: 50}
	for _, by := range []string{"operator:terminal", "mark", "attended-operator", "ssh", "sshmark", "webauthn"} {
		if _, err := validatePublishApproval(boundedTestContext(t), base, landPublishDecidedBy(t, base, f, by), req); err != nil {
			t.Errorf("%s: %v", by, err)
		}
	}
}

// decideCallers is the census of non-test calls of DecideApproval and decideApproval in the module,
// as "<file>:<enclosing func>" (design §5.1).
var decideCallers = []string{
	"host/broker/approve.go:DecideApproval",
	"host/broker/approve.go:RecordOperatorApproval",
	"host/broker/publish_op.go:mintAttendedApproval",
	"host/daemon/surface_decide.go:handleSurfaceDecide",
}

// XG1: the only non-test callers of the decide entry point are the ones named above: no effect
// handler, MCP, AG-UI or A2A route can decide a request.
func TestDecideApprovalCallerCensus(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	fset := token.NewFileSet()
	var got []string
	scanned := 0
	err := filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			if path != root && (strings.HasPrefix(e.Name(), ".") || e.Name() == "node_modules" || e.Name() == "vendor" || e.Name() == "testdata") {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return err
		}
		scanned++
		rel, _ := filepath.Rel(root, path)
		for _, decl := range f.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				name := ""
				switch fun := call.Fun.(type) {
				case *ast.Ident:
					name = fun.Name
				case *ast.SelectorExpr:
					name = fun.Sel.Name
				}
				if name == "DecideApproval" || name == "decideApproval" {
					got = append(got, filepath.ToSlash(rel)+":"+fn.Name.Name)
				}
				return true
			})
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if scanned < 100 {
		t.Fatalf("scanned %d Go files under %s; the census did not see the module", scanned, root)
	}
	sort.Strings(got)
	want := append([]string(nil), decideCallers...)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("decide callers = %q, want %q", got, want)
	}
}
