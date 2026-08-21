package evidence_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/evidence"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

var (
	testKey      = [32]byte{1, 2, 3, 4}
	testSubject  = hashref.SumSHA256([]byte("subject"))
	testCompiler = hashref.SumSHA256([]byte("compiler"))
)

type fakeReader struct {
	meta    store.ObjectMeta
	payload []byte
	err     error
	busy    time.Duration
}

func (f *fakeReader) ReadObject(context.Context, hashref.HashRef, int64) (store.ObjectMeta, []byte, error) {
	return f.meta, f.payload, f.err
}
func (f *fakeReader) BusyTimeout() time.Duration { return f.busy }

func reportFor(subject, compiler hashref.HashRef) evidence.ProofReportV1 {
	return evidence.ProofReportV1{Schema: evidence.ProofReportSchemaV1, Subject: subject, Compiler: compiler,
		CompilerVersion: "AILANG v0.30.0", Verified: []string{"world/types.ail"}, CheckPassed: true, ProofSucceeded: true}
}

func envelopeFor(t *testing.T, key [32]byte, report evidence.ProofReportV1, tagMode string) []byte {
	t.Helper()
	raw, err := evidence.EncodeProofReportV1(report)
	if err != nil {
		t.Fatal(err)
	}
	m := hmac.New(sha256.New, key[:])
	_, _ = m.Write(raw)
	tag := m.Sum(nil)
	switch tagMode {
	case "absent":
		return []byte(fmt.Sprintf(`{"report":"%s"}`, base64.RawURLEncoding.EncodeToString(raw)))
	case "wrong":
		tag = make([]byte, sha256.Size)
	}
	b, err := evidence.EncodeAuthenticatedEnvelope(evidence.AuthenticatedEnvelope{Report: raw, MAC: tag, MACValid: len(tag) == sha256.Size})
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func newFixture(t *testing.T) (*evidence.Validator, *fakeReader, hashref.HashRef) {
	t.Helper()
	payload := envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "good")
	reader := &fakeReader{meta: store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID}, payload: payload}
	v, err := evidence.NewValidator(testKey, reader, evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}, []string{"world/types.ail"})
	if err != nil {
		t.Fatal(err)
	}
	return v, reader, hashref.SumSHA256(payload)
}

func requireProvenControl(t *testing.T, v *evidence.Validator, ref hashref.HashRef) evidence.ValidatedEvidence {
	t.Helper()
	r := v.ValidateProof(context.Background(), ref, testSubject)
	seal, ok := r.Validated()
	if !ok || r.Err() != nil {
		t.Fatalf("shared control did not mint: unsupported=%v err=%v", unsupportedOf(r), r.Err())
	}
	grade, ok := v.Resolve(seal).Proven()
	if !ok || grade != evidence.ResolvedGradeProven {
		t.Fatalf("shared control did not resolve PROVEN: %v %v", grade, ok)
	}
	return seal
}

func unsupportedOf(r evidence.ValidationResult) evidence.UnsupportedReason {
	u, _ := r.Unsupported()
	return u
}
func requireReason(t *testing.T, r evidence.ValidationResult, want evidence.UnsupportedReason, mechanism string) {
	t.Helper()
	got, ok := r.Unsupported()
	if !ok || got != want {
		t.Fatalf("%s: got %s; want %s", mechanism, got, want)
	}
	if _, sealed := r.Validated(); sealed {
		t.Fatalf("%s: refusal minted a seal", mechanism)
	}
}

func TestInvalidProofRefIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.payload = nil
	requireReason(t, v.ValidateProof(context.Background(), hashref.HashRef{}, testSubject), evidence.UnsupportedInvalidRef, "invalid-ref guard")
}

func TestMissingProofReportIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.payload = nil
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256([]byte("absent")), testSubject), evidence.UnsupportedMissing, "missing-object guard")
}

func TestPayloadHashMismatchIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.payload = envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "wrong")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256([]byte("different")), testSubject), evidence.UnsupportedHashMismatch, "recomputed payload-hash guard")
}

func TestWrongSemanticIDIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.meta.SemanticID = "wrong"
	r.payload = envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "wrong")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedWrongSemanticID, "semantic-ID guard")
}

func TestWrongInterfaceIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.meta.InterfaceHash = hashref.SumSHA256([]byte("wrong"))
	r.payload = envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "wrong")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedWrongInterface, "interface-hash guard")
}

func TestMalformedProofReportIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	badReport := []byte(`{"not":"a report"}`)
	tag := make([]byte, 32)
	r.payload = []byte(fmt.Sprintf(`{"report":"%s","mac":"%s"}`, base64.RawURLEncoding.EncodeToString(badReport), base64.RawURLEncoding.EncodeToString(tag)))
	ref := hashref.SumSHA256(r.payload)
	requireReason(t, v.ValidateProof(context.Background(), ref, testSubject), evidence.UnsupportedMalformed, "strict report-decode guard")
}

func TestOtherwisePerfectReportWithoutMACIsUnauthenticated(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.payload = envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "absent")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedUnauthenticatedReport, "absent-MAC authentication guard")
}

func TestOtherwisePerfectReportWithWrongMACIsUnauthenticated(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	r.payload = envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "wrong")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedUnauthenticatedReport, "wrong-MAC authentication guard")
}

func TestMismatchedProofSubjectIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	report := reportFor(hashref.SumSHA256([]byte("other subject")), hashref.SumSHA256([]byte("other compiler")))
	r.payload = envelopeFor(t, testKey, report, "good")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedSubjectMismatch, "subject-binding guard")
}

func TestMismatchedProofToolIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	report := reportFor(testSubject, hashref.SumSHA256([]byte("other compiler")))
	report.Verified = []string{"other"}
	r.payload = envelopeFor(t, testKey, report, "good")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedToolMismatch, "compiler-identity guard")
}

func TestFailedProofReportIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	report := reportFor(testSubject, testCompiler)
	report.CheckPassed = false
	report.Verified = []string{"other"}
	r.payload = envelopeFor(t, testKey, report, "good")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedProofFailed, "proof-success guard")
}

func TestIncompleteProofReportIsRefused(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	report := reportFor(testSubject, testCompiler)
	report.Verified = []string{"other"}
	r.payload = envelopeFor(t, testKey, report, "good")
	requireReason(t, v.ValidateProof(context.Background(), hashref.SumSHA256(r.payload), testSubject), evidence.UnsupportedProofIncomplete, "required-identity guard")
}

func TestAttackerChosenValidatorCannotMintForHostAuthority(t *testing.T) {
	v1, _, ref1 := newFixture(t)
	own := requireProvenControl(t, v1, ref1)
	v1copy := *v1
	if grade, ok := v1copy.Resolve(own).Proven(); !ok || grade != evidence.ResolvedGradeProven {
		t.Fatal("validator copy lost its mint identity")
	}
	key2 := [32]byte{9}
	payload2 := envelopeFor(t, key2, reportFor(testSubject, testCompiler), "good")
	r2 := &fakeReader{meta: store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID}, payload: payload2}
	v2, err := evidence.NewValidator(key2, r2, evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}, []string{"world/types.ail"})
	if err != nil {
		t.Fatal(err)
	}
	foreign, ok := v2.ValidateProof(context.Background(), hashref.SumSHA256(payload2), testSubject).Validated()
	if !ok {
		t.Fatal("caller-constructed validator did not self-mint")
	}
	result := v1.Resolve(foreign)
	if !errors.Is(result.Err(), evidence.ErrForeignSeal) {
		t.Fatalf("binding check: got %v; want ErrForeignSeal", result.Err())
	}
	if grade, ok := result.Proven(); ok {
		t.Fatalf("foreign seal resolved %v; want ErrForeignSeal", grade)
	}
}

func TestValidatorMintIdentitiesAreDistinct(t *testing.T) {
	v1, _, ref1 := newFixture(t)
	requireProvenControl(t, v1, ref1)
	payload2 := envelopeFor(t, testKey, reportFor(testSubject, testCompiler), "good")
	r2 := &fakeReader{meta: store.ObjectMeta{InterfaceHash: evidence.InterfaceHashV1, SemanticID: evidence.ProofSemanticID}, payload: payload2}
	v2, err := evidence.NewValidator(testKey, r2, evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}, []string{"world/types.ail"})
	if err != nil {
		t.Fatal(err)
	}
	seal2, ok := v2.ValidateProof(context.Background(), hashref.SumSHA256(payload2), testSubject).Validated()
	if !ok {
		t.Fatal("second same-key validator did not mint")
	}
	if !errors.Is(v1.Resolve(seal2).Err(), evidence.ErrForeignSeal) {
		t.Fatal("distinct NewValidator calls with the same key shared a mint identity")
	}
}

func TestZeroValueForgeryCannotResolve(t *testing.T) {
	v, _, ref := newFixture(t)
	requireProvenControl(t, v, ref)
	var zeroV evidence.Validator
	var zeroS evidence.ValidatedEvidence
	result := zeroV.Resolve(zeroS)
	if !errors.Is(result.Err(), evidence.ErrUnmintedAuthority) {
		t.Fatalf("mint-validity check: got %v; want ErrUnmintedAuthority", result.Err())
	}
	if grade, ok := result.Proven(); ok {
		t.Fatalf("zero-value seal resolved %v; want ErrUnmintedAuthority", grade)
	}
	var zeroR evidence.ResolutionResult
	if grade, ok := zeroR.Proven(); ok {
		t.Fatalf("zero ResolutionResult resolved %v", grade)
	}
}

// frozenPublicSurface is the EXACT inventory of exported declarations in package
// evidence. It is a MANIFEST, not a pattern: any new exported symbol — of any
// shape, under any name, whatever it returns — reds this test until it is added
// here deliberately.
//
// The completeness is the whole point, and it was bought with a measurement.
// This test previously enumerated only package-level funcs whose *result type
// was the bare identifier* ResolvedGrade or ResolutionResult. That fires on the
// spelling §5's M16 row happens to name, and iteration 106 drilled three other
// natural spellings of the very package-level PROVEN resolver AC1 forbids, each
// minting ResolvedGradeProven from a raw HashRef with no seal — and the whole
// package stayed GREEN for all three:
//
//	func GradeOfRef(hashref.HashRef) *ResolutionResult   // *ast.StarExpr, not *ast.Ident
//	func (RawGrader) Grade(hashref.HashRef) ResolvedGrade // fn.Recv != nil, skipped
//	type Grade = ResolvedGrade; func GradeOfClaim(...) Grade // alias, name differs
//
// A removal proves a check FIRES; only an addition proves it LOOKS. This list is
// the addition-proof form.
var frozenPublicSurface = []string{
	"func DecodeAuthenticatedEnvelope",
	"func DecodeProofReportV1",
	"func DecodeProposal",
	"func EncodeAuthenticatedEnvelope",
	"func EncodeProofReportV1",
	"func NewValidator",
	"method ClaimedEvidence.IsZero",
	"method DecodeRefusal.Error",
	"method ResolutionResult.Err",
	"method ResolutionResult.Proven",
	"method ValidationResult.Err",
	"method ValidationResult.Unsupported",
	"method ValidationResult.Validated",
	"method Validator.Resolve",
	"method Validator.ValidateProof",
	"type AuthenticatedEnvelope",
	"type ClaimedEvidence",
	"type CompilerConfig",
	"type DecodeRefusal",
	"type ObjectReader",
	"type ProofReportV1",
	"type ResolutionResult",
	"type ResolvedGrade",
	"type UnsupportedReason",
	"type ValidatedEvidence",
	"type ValidationResult",
	"type Validator",
	"value ErrForeignSeal",
	"value ErrInvalidValidatorConfig",
	"value ErrUnmintedAuthority",
	"value ErrUnorderedTimeouts",
	"value InterfaceHashV1",
	"value MaxBytes",
	"value MaxStringBytes",
	"value MaxVerifiedIdentities",
	"value ProofReportSchemaV1",
	"value ProofSemanticID",
	"value RefusalInvalidUTF8",
	"value RefusalLimit",
	"value RefusalMalformed",
	"value RefusalNonCanonical",
	"value RefusalOversize",
	"value ResolvedGradeProven",
	"value UnsupportedHashMismatch",
	"value UnsupportedInvalidRef",
	"value UnsupportedMalformed",
	"value UnsupportedMissing",
	"value UnsupportedOversize",
	"value UnsupportedProofFailed",
	"value UnsupportedProofIncomplete",
	"value UnsupportedSubjectMismatch",
	"value UnsupportedToolMismatch",
	"value UnsupportedUnauthenticatedReport",
	"value UnsupportedWrongInterface",
	"value UnsupportedWrongSemanticID",
}

// exportedSurfaceOf enumerates package evidence's exported declarations from its
// own source, counting METHODS and TYPES and VALUES, not only package-level
// funcs — the categories the previous scan could not see.
func exportedSurfaceOf(t *testing.T, dir string) []string {
	t.Helper()
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	pkg, ok := pkgs["evidence"]
	if !ok || pkg == nil || len(pkg.Files) == 0 {
		t.Fatal("instrument failure: parsed no non-test files for package evidence")
	}
	var got []string
	for _, file := range pkg.Files {
		for _, decl := range file.Decls {
			switch node := decl.(type) {
			case *ast.FuncDecl:
				if !node.Name.IsExported() {
					continue
				}
				if node.Recv != nil && len(node.Recv.List) > 0 {
					got = append(got, "method "+receiverName(node.Recv.List[0].Type)+"."+node.Name.Name)
					continue
				}
				got = append(got, "func "+node.Name.Name)
			case *ast.GenDecl:
				for _, spec := range node.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						if s.Name.IsExported() {
							got = append(got, "type "+s.Name.Name)
						}
					case *ast.ValueSpec:
						for _, id := range s.Names {
							if id.IsExported() {
								got = append(got, "value "+id.Name)
							}
						}
					}
				}
			}
		}
	}
	if len(got) == 0 {
		t.Fatal("instrument failure: exported-surface enumeration is empty")
	}
	sort.Strings(got)
	return got
}

func receiverName(e ast.Expr) string {
	switch t := e.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.IndexExpr:
		return receiverName(t.X)
	}
	return "?"
}

func TestPublicAuthoritySurfaceIsFrozen(t *testing.T) {
	for _, typ := range []reflect.Type{reflect.TypeOf(evidence.ValidatedEvidence{}), reflect.TypeOf(evidence.ResolutionResult{})} {
		for i := 0; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				t.Fatalf("public authority surface exposes field %s.%s", typ.Name(), typ.Field(i).Name)
			}
		}
	}
	rt := reflect.TypeOf(evidence.ResolutionResult{})
	if rt.NumMethod() != 2 || rt.Method(0).Name != "Err" || rt.Method(1).Name != "Proven" {
		t.Fatalf("public authority surface exposes non-sealed PROVEN ingress: methods=%v", methodNames(rt))
	}
	want := append([]string(nil), frozenPublicSurface...)
	sort.Strings(want)
	if got := exportedSurfaceOf(t, "."); !reflect.DeepEqual(got, want) {
		t.Fatalf("public authority surface changed:\n  added:   %v\n  removed: %v\n(any new exported declaration must be reviewed against AC1 and then added to frozenPublicSurface)",
			missingFrom(got, want), missingFrom(want, got))
	}

	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", func(info fs.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range pkgs["evidence"].Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() || fn.Type.Results == nil {
				continue
			}
			for _, result := range fn.Type.Results.List {
				name, ok := result.Type.(*ast.Ident)
				if ok && (name.Name == "ResolvedGrade" || name.Name == "ResolutionResult") {
					t.Fatalf("public authority surface exposes non-sealed PROVEN ingress: function %s returns %s", fn.Name.Name, name.Name)
				}
			}
		}
	}
}

func methodNames(t reflect.Type) []string {
	out := make([]string, t.NumMethod())
	for i := range out {
		out[i] = t.Method(i).Name
	}
	return out
}

func TestConstructorRefusesNonPositiveObjectReadTimeout(t *testing.T) {
	r := &fakeReader{}
	cfg := evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0"}
	for _, d := range []time.Duration{0, -time.Second} {
		cfg.ObjectReadTimeout = d
		if _, err := evidence.NewValidator(testKey, r, cfg, []string{"id"}); !errors.Is(err, evidence.ErrInvalidValidatorConfig) {
			t.Fatalf("NewValidator accepted ObjectReadTimeout %s: %v", d, err)
		}
	}
}

// The three constructor refusals below were unpinned as delivered: iteration 106
// measured each of them neutered independently with the ENTIRE host/evidence suite
// green, because a single non-positive-timeout test stood in for a four-disjunct
// compound guard. Each arm now asserts the branch's OWN message, not merely the
// shared ErrInvalidValidatorConfig sentinel — the sentinel is produced by six
// different refusals, so matching on it alone is an observable whose value set is
// wider than the mechanism's.
func TestConstructorRefusesNilReader(t *testing.T) {
	cfg := evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}
	_, err := evidence.NewValidator(testKey, nil, cfg, []string{"id"})
	if !errors.Is(err, evidence.ErrInvalidValidatorConfig) || !strings.Contains(err.Error(), "reader is nil") {
		t.Fatalf("nil-reader guard: got %v; want ErrInvalidValidatorConfig naming the reader", err)
	}
}

func TestConstructorRefusesUnsetCompilerIdentity(t *testing.T) {
	r := &fakeReader{}
	cfg := evidence.CompilerConfig{CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}
	_, err := evidence.NewValidator(testKey, r, cfg, []string{"id"})
	if !errors.Is(err, evidence.ErrInvalidValidatorConfig) || !strings.Contains(err.Error(), "compiler identity is unset") {
		t.Fatalf("compiler-identity guard: got %v; want ErrInvalidValidatorConfig naming the compiler identity", err)
	}
}

func TestConstructorRefusesEmptyCompilerVersion(t *testing.T) {
	r := &fakeReader{}
	cfg := evidence.CompilerConfig{Compiler: testCompiler, ObjectReadTimeout: time.Second}
	_, err := evidence.NewValidator(testKey, r, cfg, []string{"id"})
	if !errors.Is(err, evidence.ErrInvalidValidatorConfig) || !strings.Contains(err.Error(), "compiler version is empty") {
		t.Fatalf("compiler-version guard: got %v; want ErrInvalidValidatorConfig naming the compiler version", err)
	}
}

// The validator reads the report the ENVELOPE decoder already decoded rather than
// decoding a second time, because the second decode's error branch is unreachable
// and therefore unpinnable. That shortcut is only safe while the envelope really
// does carry a faithfully decoded report, so pin the invariant rather than assume
// it: a validated envelope must mint, and the same bytes must decode identically
// through the public report decoder.
func TestEnvelopeCarriesTheReportItAlreadyDecoded(t *testing.T) {
	v, r, good := newFixture(t)
	requireProvenControl(t, v, good)
	env, err := evidence.DecodeAuthenticatedEnvelope(r.payload)
	if err != nil {
		t.Fatal(err)
	}
	fresh, err := evidence.DecodeProofReportV1(env.Report)
	if err != nil {
		t.Fatalf("envelope accepted a report the public decoder refuses: %v", err)
	}
	want := reportFor(testSubject, testCompiler)
	if !reflect.DeepEqual(fresh, want) {
		t.Fatalf("envelope report bytes decode to %+v; want %+v", fresh, want)
	}
}

func TestConstructorRefusesEmptyRequiredIdentities(t *testing.T) {
	v, _, ref := newFixture(t)
	requireProvenControl(t, v, ref)
	r := &fakeReader{}
	cfg := evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}
	for _, ids := range [][]string{nil, {}} {
		if _, err := evidence.NewValidator(testKey, r, cfg, ids); !errors.Is(err, evidence.ErrInvalidValidatorConfig) {
			t.Fatalf("NewValidator accepted empty required identities; want ErrInvalidValidatorConfig")
		}
	}
}

func TestConstructorRefusesUnknownBusyTimeout(t *testing.T) {
	r := &fakeReader{busy: -1}
	cfg := evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}
	if _, err := evidence.NewValidator(testKey, r, cfg, []string{"id"}); !errors.Is(err, evidence.ErrInvalidValidatorConfig) {
		t.Fatalf("NewValidator accepted unknown BusyTimeout: %v", err)
	}
}

func TestConstructorNamesActuallyUsedUnorderedTimeouts(t *testing.T) {
	r := &fakeReader{busy: 2 * time.Second}
	cfg := evidence.CompilerConfig{Compiler: testCompiler, CompilerVersion: "AILANG v0.30.0", ObjectReadTimeout: time.Second}
	_, err := evidence.NewValidator(testKey, r, cfg, []string{"id"})
	if !errors.Is(err, evidence.ErrUnorderedTimeouts) || !strings.Contains(err.Error(), "1s") || !strings.Contains(err.Error(), "2s") {
		t.Fatalf("ordering refusal did not name runtime values: %v", err)
	}
}

func missingFrom(have, other []string) []string {
	set := make(map[string]bool, len(other))
	for _, s := range other {
		set[s] = true
	}
	var out []string
	for _, s := range have {
		if !set[s] {
			out = append(out, s)
		}
	}
	return out
}
