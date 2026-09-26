package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/sunholo-data/ailang-world/host/archive"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// daemonFakeRelease is the release line the daemon-fixture interpreter
// reports; the daemon's own bootstrap derives it from the manifest.
const daemonFakeRelease = "DAEMON-FAKE v1"

// daemonFakeInterpreter writes the house-pattern shell-script fake: it
// prints the release for --version and exits 0 for everything else, so
// `check` accepts every source.
func daemonFakeInterpreter(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "daemon-fake-ailang")
	script := "#!/bin/sh\n" +
		"case \"$1\" in\n" +
		"  --version) echo \"" + daemonFakeRelease + "\";;\n" +
		"esac\n" +
		"exit 0\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path
}

// newPublisherDaemon starts a REAL daemon on a temp store with its interpreter
// configured — the production startup shape: daemon.New archives the binary
// and bootstraps the epoch registry with its release (daemon.go:494-512), so
// publishing through the daemon's own archived interpreter and bootstrapped
// epoch registry is exactly what the verb does against a served store.
func newPublisherDaemon(t *testing.T) (*Daemon, string) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "world.db")
	d, err := New(context.Background(), Config{
		DBPath: dbPath, BindHost: DefaultBindHost, AilangBin: daemonFakeInterpreter(t),
	})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() {
		if err := d.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	})
	return d, dbPath
}

// daemonInterpreterRef is the daemon's own archived interpreter pin.
func daemonInterpreterRef(t *testing.T, d *Daemon) hashref.HashRef {
	t.Helper()
	if d.interpreterRef == "" {
		t.Fatal("daemon was started without AilangBin; no interpreter is pinned")
	}
	ref, err := hashref.Parse(d.interpreterRef)
	if err != nil {
		t.Fatalf("parse daemon interpreter ref %q: %v", d.interpreterRef, err)
	}
	return ref
}

// publishProductionRevision publishes one honest descriptor through the
// PRODUCTION publisher core — transitionreg.NewPublisher(db, archive).PublishSet,
// the exact function the `world-publish transitions` operator verb calls after
// its attended fences — into the live daemon's own store, pinning the
// DAEMON'S archived interpreter and the epoch its bootstrap derived. The
// transition source is stored as an object exactly as the verb's
// transitionFnFile step stores it.
func publishProductionRevision(t *testing.T, d *Daemon, dbPath, id, effect string, epoch int64) transitionreg.SetResult {
	t.Helper()
	payload := []byte("transition source for " + id)
	src := store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("daemon-test/transition-source")),
		SemanticID: "daemon-test/transition-source", Provenance: "registry_publisher_test", Payload: payload,
	}
	if err := d.store.PutObject(src); err != nil {
		t.Fatalf("put transition source: %v", err)
	}
	desc := transitionreg.Descriptor{
		ID: id, TransitionFn: src.Hash,
		Interpreter:    daemonInterpreterRef(t, d),
		SemanticsEpoch: epoch,
		InputSchema:    []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: effect, Scope: "world", Cost: 1},
		DeclaredEffects: []transitionreg.EffectRequirement{{Effect: effect, Scope: "world", Cost: 1}},
		Title:           "title-" + id, Description: "description-" + id,
	}
	res, err := transitionreg.NewPublisher(d.store, archive.New(dbPath)).PublishSet(context.Background(),
		[]transitionreg.Change{{ID: id, Descriptor: &desc}})
	if err != nil {
		t.Fatalf("publish through the production path: %v", err)
	}
	if res.Unchanged || res.Revision < 1 {
		t.Fatalf("publish result = %+v, want a changed revision", res)
	}
	return res
}

type cardSkills struct {
	Skills []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"skills"`
}

func cardSkillIDs(t *testing.T, d *Daemon, auth string) cardSkills {
	t.Helper()
	rec := requestRecorderAuth(t, d, auth, http.MethodGet, "/.well-known/agent.json", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("card for an authorized session = %d %s, want 200", rec.Code, rec.Body)
	}
	var card cardSkills
	if err := json.Unmarshal(rec.Body.Bytes(), &card); err != nil {
		t.Fatalf("card body is not JSON: %v\n%s", err, rec.Body)
	}
	return card
}

// TestPublishedTransitionsAppearOnAgentCard is the row-107 acceptance test:
// a revision published through the production path lands in the live daemon's
// store, and the daemon's OWN reader (projection.go) then lists the published
// transition as a card skill for a session whose capabilities allow the
// descriptor's Access requirement — while a session without that capability
// sees a zero-skills card. No wall-clock assertions; the recorder harness
// binds no sockets.
func TestPublishedTransitionsAppearOnAgentCard(t *testing.T) {
	d, dbPath := newPublisherDaemon(t)
	publishProductionRevision(t, d, dbPath, "tools.echo", "world.apply", 1)

	// Positive arm: the session's capabilities include the descriptor's
	// Access requirement, so the card lists exactly the published skill.
	auth := mintSessionGrants(t, d, "ep-publisher", "world.apply")
	card := cardSkillIDs(t, d, "Bearer "+auth)
	if len(card.Skills) != 1 {
		t.Fatalf("card skills = %+v, want exactly the published transition", card.Skills)
	}
	if card.Skills[0].ID != "tools.echo" || card.Skills[0].Name != "title-tools.echo" {
		t.Fatalf("card skill = %+v, want the published descriptor's ID and title verbatim", card.Skills[0])
	}

	// Negative arm: a session whose capabilities do NOT allow the descriptor's
	// Access requirement sees a legitimate ZERO-skills card at 200.
	other := mintSessionGrants(t, d, "ep-unrelated", "fs.read")
	cardOther := cardSkillIDs(t, d, "Bearer "+other)
	if len(cardOther.Skills) != 0 {
		t.Fatalf("card skills for an unrelated session = %+v, want zero (the capability filter must apply to published entries)", cardOther.Skills)
	}
}

// TestPublisherRefusalKeepsCardEmpty is the negative production path: a
// descriptor whose TransitionFn object is absent from the store is REFUSED by
// the publisher core, so the daemon's card still has zero skills — a card
// skill must always name a source execution and replay can resolve.
func TestPublisherRefusalKeepsCardEmpty(t *testing.T) {
	d, dbPath := newPublisherDaemon(t)
	ghost := transitionreg.Descriptor{
		ID: "tools.ghost", TransitionFn: hashref.SumSHA256([]byte("never stored")),
		Interpreter:    daemonInterpreterRef(t, d),
		SemanticsEpoch: 1, InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1},
		DeclaredEffects: []transitionreg.EffectRequirement{{Effect: "world.apply", Scope: "world", Cost: 1}},
	}
	_, err := transitionreg.NewPublisher(d.store, archive.New(dbPath)).PublishSet(context.Background(),
		[]transitionreg.Change{{ID: ghost.ID, Descriptor: &ghost}})
	var absent *transitionreg.TransitionSourceAbsentError
	if !errors.As(err, &absent) {
		t.Fatalf("absent-source publish error = %v, want *TransitionSourceAbsentError", err)
	}
	auth := mintSessionGrants(t, d, "ep-empty", "world.apply")
	card := cardSkillIDs(t, d, "Bearer "+auth)
	if len(card.Skills) != 0 {
		t.Fatalf("card skills after a refused publish = %+v, want zero", card.Skills)
	}
}

// TestEpochRefusalKeepsCardUnchanged is the objection-A daemon arm: after one
// honest publish populates the card, a second descriptor pinning an epoch the
// registry does NOT nominate for the interpreter's release (epoch 2 against
// the daemon's bootstrapped epoch 1) is refused with the typed mismatch error
// — and BOTH the head and the card stay exactly as they were.
func TestEpochRefusalKeepsCardUnchanged(t *testing.T) {
	d, dbPath := newPublisherDaemon(t)
	publishProductionRevision(t, d, dbPath, "tools.echo", "world.apply", 1)
	headBefore, revBefore, ok, err := transitionreg.NewReader(d.store).CurrentRevision(context.Background())
	if err != nil || !ok {
		t.Fatalf("read head after the honest publish: ok=%v err=%v", ok, err)
	}

	blocked := transitionreg.Descriptor{
		ID: "tools.blocked", TransitionFn: hashref.SumSHA256([]byte("epoch mismatch source")),
		Interpreter:    daemonInterpreterRef(t, d),
		SemanticsEpoch: 2, InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1},
		DeclaredEffects: []transitionreg.EffectRequirement{{Effect: "world.apply", Scope: "world", Cost: 1}},
		Title:           "title-tools.blocked", Description: "must never appear",
	}
	if err := d.store.PutObject(store.Object{
		Hash: blocked.TransitionFn, InterfaceHash: hashref.SumSHA256([]byte("daemon-test/transition-source")),
		SemanticID: "daemon-test/transition-source", Provenance: "registry_publisher_test",
		Payload: []byte("epoch mismatch source"),
	}); err != nil {
		t.Fatal(err)
	}
	_, err = transitionreg.NewPublisher(d.store, archive.New(dbPath)).PublishSet(context.Background(),
		[]transitionreg.Change{{ID: blocked.ID, Descriptor: &blocked}})
	var mismatch *transitionreg.InterpreterEpochMismatchError
	if !errors.As(err, &mismatch) {
		t.Fatalf("mismatched-epoch publish error = %v, want *InterpreterEpochMismatchError", err)
	}

	// The head did not move.
	headAfter, revAfter, ok, err := transitionreg.NewReader(d.store).CurrentRevision(context.Background())
	if err != nil || !ok || headAfter != headBefore || revAfter.Revision != revBefore.Revision {
		t.Fatalf("head after the refused publish = (%q, rev %d, ok=%v), want it unchanged at (%q, rev %d)",
			headAfter, revAfter.Revision, ok, headBefore, revBefore.Revision)
	}
	// The card did not change: exactly the first skill, nothing else.
	auth := mintSessionGrants(t, d, "ep-epoch", "world.apply")
	card := cardSkillIDs(t, d, "Bearer "+auth)
	if len(card.Skills) != 1 || card.Skills[0].ID != "tools.echo" {
		t.Fatalf("card after the refused publish = %+v, want exactly the first published skill", card.Skills)
	}
}

// TestEpochTwinInterpretersBothListedOnCard is AC-EPOCH-TWIN at daemon
// level: a second archived interpreter whose bytes differ from the daemon's
// own but whose first --version line is the same release is epoch-eligible
// under the daemon's bootstrapped nomination; both published skills appear on
// the card and each stored descriptor keeps its own Interpreter pin. A
// non-nominated release is refused and head + card stay unchanged.
func TestEpochTwinInterpretersBothListedOnCard(t *testing.T) {
	d, dbPath := newPublisherDaemon(t)
	arch := archive.New(dbPath)
	own := daemonInterpreterRef(t, d)
	writeFake := func(name, versionOut string) hashref.HashRef {
		t.Helper()
		path := filepath.Join(t.TempDir(), name)
		script := "#!/bin/sh\ncase \"$1\" in\n  --version) printf '" + versionOut + "';;\nesac\nexit 0\n"
		if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
			t.Fatal(err)
		}
		ref, err := arch.Archive(path)
		if err != nil {
			t.Fatalf("archive %s: %v", name, err)
		}
		return ref
	}
	twin := writeFake("daemon-twin-ailang", daemonFakeRelease+"\\nCommit: twin\\n")
	if twin == own {
		t.Fatal("premise: the twin interpreter must hash differently from the daemon's own")
	}
	// pinned is one honest descriptor pinning interp; its source object is stored.
	pinned := func(id string, interp hashref.HashRef) transitionreg.Change {
		payload := []byte("transition source for " + id)
		if err := d.store.PutObject(store.Object{
			Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte("daemon-test/transition-source")),
			SemanticID: "daemon-test/transition-source", Provenance: "registry_publisher_test", Payload: payload,
		}); err != nil {
			t.Fatal(err)
		}
		desc := transitionreg.Descriptor{
			ID: id, TransitionFn: hashref.SumSHA256(payload), Interpreter: interp, SemanticsEpoch: 1,
			InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
			Access:          transitionreg.EffectRequirement{Effect: "world.apply", Scope: "world", Cost: 1},
			DeclaredEffects: []transitionreg.EffectRequirement{{Effect: "world.apply", Scope: "world", Cost: 1}},
			Title:           "title-" + id,
		}
		return transitionreg.Change{ID: id, Descriptor: &desc}
	}
	publish := func(changes ...transitionreg.Change) error {
		_, err := transitionreg.NewPublisher(d.store, arch).PublishSet(context.Background(), changes)
		return err
	}
	// Both twins in ONE publish: the set is where a release-level stand-in
	// for the descriptor's own pin would conflate them.
	if err := publish(pinned("tools.own", own), pinned("tools.twin", twin)); err != nil {
		t.Fatalf("publish on the daemon's own interpreter and its same-release twin: %v", err)
	}
	snap, err := transitionreg.NewReader(d.store).ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	entries := snap.List()
	if len(entries) != 2 || entries[0].Interpreter != own || entries[1].Interpreter != twin {
		t.Fatalf("stored entries = %+v, want tools.own pinned to %s and tools.twin pinned to %s", entries, own, twin)
	}
	auth := mintSessionGrants(t, d, "ep-twin", "world.apply")
	if card := cardSkillIDs(t, d, "Bearer "+auth); len(card.Skills) != 2 {
		t.Fatalf("card = %+v, want both twin-pinned skills", card.Skills)
	}

	other := writeFake("daemon-other-ailang", "DAEMON-OTHER v2\\n")
	var mismatch *transitionreg.InterpreterEpochMismatchError
	if err := publish(pinned("tools.other", other)); !errors.As(err, &mismatch) {
		t.Fatalf("non-nominated release = %v, want *InterpreterEpochMismatchError", err)
	}
	after, err := transitionreg.NewReader(d.store).ReadSnapshot(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if after.Head != snap.Head || len(after.List()) != 2 {
		t.Fatalf("head after refusal = %s (%d entries), want unchanged %s (2)", after.Head, len(after.List()), snap.Head)
	}
	if card := cardSkillIDs(t, d, "Bearer "+auth); len(card.Skills) != 2 {
		t.Fatalf("card after refusal = %+v, want unchanged (2 skills)", card.Skills)
	}
}
