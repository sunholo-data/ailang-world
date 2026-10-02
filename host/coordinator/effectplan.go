package coordinator

// Effect-plan law (queue row 134, M1; design
// design_docs/planned/w-software-engineering-domain.md §4.1). A transition
// whose descriptor declares effects answers its plan phase with a
// world/effect-plan/v1 object; parsePlan admits it only under plan laws
// L1–L4, L6 and L7 (L5 is reserved for multi-effect ordering, R-SE-10), and
// parseFinish admits a finish-phase output under L6.
//
// The laws are specified, with Z3-proven contracts, in
// design_docs/sketches/effectplan.ail. Every lower-case predicate below
// mirrors the sketch function of the same name, and
// effectplan_drift_test.go evaluates every inline test row of the sketch
// through these mirrors, so the Go law and the proven law cannot drift.
//
// parsePlan and parseFinish are called from effectful dispatch and replay
// (effectful.go, row 134 M2a/M2b).

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

const (
	// EffectPlanV1 is the plan object's own discriminator value.
	EffectPlanV1 = "world/effect-plan/v1"

	maxPlanEffectsV1  = 1       // L1
	maxEffectIDBytes  = 16      // L2
	maxPayloadBytes   = 1 << 20 // L4: canonical payload object, 1 MiB
	reservedWorldKey  = "world" // L6
	finishEffectsKey  = "effects"
	maxPhaseOutputRaw = 4 << 20 // parse bound; Dispatch's MaxOutput cap applies first
)

// PlanLaw names the law a refused phase output violates.
type PlanLaw string

const (
	// LawShape: not a well-formed world/effect-plan/v1 object (JSON syntax,
	// discriminator, closed key sets, member types). A codec property, not a
	// sketch law.
	LawShape PlanLaw = "shape"
	LawL1    PlanLaw = "L1" // effect count
	LawL2    PlanLaw = "L2" // effect id grammar and uniqueness
	LawL3    PlanLaw = "L3" // (effect, scope, cost) ∈ DeclaredEffects
	LawL4    PlanLaw = "L4" // payload is a JSON object, canonical ≤ 1 MiB
	LawL6    PlanLaw = "L6" // reserved `world` key; finish cannot request effects
	LawL7    PlanLaw = "L7" // zero effects ⇔ non-null result; finish needs an effect
)

// PlanLawError refuses a plan- or finish-phase output. Law is the first law,
// in parse order, that the output violates.
type PlanLawError struct {
	Law    PlanLaw
	Reason string
}

func (e *PlanLawError) Error() string {
	return fmt.Sprintf("coordinator: effect plan violates %s: %s", e.Law, e.Reason)
}

func lawErr(law PlanLaw, format string, args ...any) error {
	return &PlanLawError{Law: law, Reason: fmt.Sprintf(format, args...)}
}

// PlannedEffect is one admitted effect: its plan-local id, the declared
// requirement it matched exactly, and its canonically re-encoded payload.
type PlannedEffect struct {
	ID          string
	Requirement transitionreg.EffectRequirement
	Payload     []byte
}

// Plan is an admitted world/effect-plan/v1 object.
type Plan struct {
	Effects []PlannedEffect
	Finish  bool
	// Result is the canonical result object; non-nil exactly when Effects is
	// empty (L7).
	Result []byte
	// Canonical is the canonical encoding of the whole plan object — the bytes
	// M2 stores and references as the record's `plan`.
	Canonical []byte
}

var (
	planKeys   = []string{"effects", "finish", "plan", "result"}
	effectKeys = []string{"cost", "effect", "id", "payload", "scope"}
)

// parsePlan admits a plan-phase output against the descriptor's declared
// effects, or returns a *PlanLawError naming the violated law.
func parsePlan(output []byte, declared []transitionreg.EffectRequirement) (Plan, error) {
	canonical, err := transitionreg.CanonicalJSON(output, maxPhaseOutputRaw)
	if err != nil {
		return Plan{}, lawErr(LawShape, "not well-formed JSON: %v", err)
	}
	top, err := objectMembers(canonical)
	if err != nil {
		return Plan{}, lawErr(LawShape, "plan is not a JSON object")
	}
	if _, ok := top[reservedWorldKey]; ok {
		return Plan{}, lawErr(LawL6, "the plan carries the reserved key %q", reservedWorldKey)
	}
	if err := closedKeys(top, planKeys); err != nil {
		return Plan{}, lawErr(LawShape, "plan %v", err)
	}
	if tag, ok := jsonString(top["plan"]); !ok || tag != EffectPlanV1 {
		return Plan{}, lawErr(LawShape, "plan discriminator is not %q", EffectPlanV1)
	}
	finish, ok := jsonBool(top["finish"])
	if !ok {
		return Plan{}, lawErr(LawShape, "finish is not a boolean")
	}
	var rawEffects []json.RawMessage
	if kindOf(top["effects"]) != '[' || json.Unmarshal(top["effects"], &rawEffects) != nil {
		return Plan{}, lawErr(LawShape, "effects is not a JSON array")
	}
	n := len(rawEffects)
	if !effectCountOk(n) {
		return Plan{}, lawErr(LawL1, "%d effects; a v1 plan carries at most %d", n, maxPlanEffectsV1)
	}

	plan := Plan{Finish: finish, Canonical: canonical}
	for i, raw := range rawEffects {
		e, err := parsePlannedEffect(i, raw, declared)
		if err != nil {
			return Plan{}, err
		}
		plan.Effects = append(plan.Effects, e)
	}
	for i := range plan.Effects {
		for j := i + 1; j < len(plan.Effects); j++ {
			if !idsDistinct(plan.Effects[i].ID, plan.Effects[j].ID) {
				return Plan{}, lawErr(LawL2, "effect id %q is not unique", plan.Effects[i].ID)
			}
		}
	}

	result := top["result"]
	hasResult := kindOf(result) != 'n'
	if hasResult && kindOf(result) != '{' {
		return Plan{}, lawErr(LawL7, "result is neither null nor a JSON object")
	}
	if !resultPresenceOk(n, hasResult) {
		if n == 0 {
			return Plan{}, lawErr(LawL7, "a plan with zero effects must carry a result object")
		}
		return Plan{}, lawErr(LawL7, "a plan with effects must carry a null result")
	}
	if !finishNeedsEffect(n, finish) {
		return Plan{}, lawErr(LawL7, "finish requested by a plan with zero effects")
	}
	if hasResult {
		members, _ := objectMembers(result) // kind checked above
		for k := range members {
			if reservedOutputKey(k) {
				return Plan{}, lawErr(LawL6, "the result carries the reserved key %q", k)
			}
		}
		plan.Result = append([]byte(nil), result...)
	}
	return plan, nil
}

func parsePlannedEffect(i int, raw json.RawMessage, declared []transitionreg.EffectRequirement) (PlannedEffect, error) {
	members, err := objectMembers(raw)
	if err != nil {
		return PlannedEffect{}, lawErr(LawShape, "effects[%d] is not a JSON object", i)
	}
	if err := closedKeys(members, effectKeys); err != nil {
		return PlannedEffect{}, lawErr(LawShape, "effects[%d] %v", i, err)
	}
	id, okID := jsonString(members["id"])
	effect, okEffect := jsonString(members["effect"])
	scope, okScope := jsonString(members["scope"])
	if !okID || !okEffect || !okScope {
		return PlannedEffect{}, lawErr(LawShape, "effects[%d] id, effect and scope must be strings", i)
	}
	cost, err := strconv.ParseInt(string(members["cost"]), 10, 64)
	if kindOf(members["cost"]) != '0' || err != nil {
		return PlannedEffect{}, lawErr(LawShape, "effects[%d] cost is not an integer", i)
	}
	if !idOk(id) {
		return PlannedEffect{}, lawErr(LawL2, "effects[%d] id %q does not match [a-z0-9]{1,%d}", i, id, maxEffectIDBytes)
	}
	payload := members["payload"]
	if kindOf(payload) != '{' {
		return PlannedEffect{}, lawErr(LawL4, "effects[%d] payload is not a JSON object", i)
	}
	// payload is a member of the canonical plan, so its bytes are already the
	// canonical re-encoding; re-running the codec makes that a checked fact.
	canonicalPayload, err := transitionreg.CanonicalJSON(payload, maxPhaseOutputRaw)
	if err != nil || string(canonicalPayload) != string(payload) {
		return PlannedEffect{}, lawErr(LawL4, "effects[%d] payload is not canonical", i)
	}
	law := lawEffect{ID: id, Effect: effect, Scope: scope, Cost: cost, PayloadBytes: len(canonicalPayload)}
	if !declaredContains(declared, law) {
		return PlannedEffect{}, lawErr(LawL3, "effects[%d] (%s, %s, %d) is not a declared effect", i, effect, scope, cost)
	}
	if !payloadSizeOk(law.PayloadBytes) {
		return PlannedEffect{}, lawErr(LawL4, "effects[%d] payload is %d canonical bytes; limit is %d", i, law.PayloadBytes, maxPayloadBytes)
	}
	return PlannedEffect{
		ID:          id,
		Requirement: transitionreg.EffectRequirement{Effect: effect, Scope: scope, Cost: cost},
		Payload:     canonicalPayload,
	}, nil
}

// parseFinish admits a finish-phase output under L6: a JSON object that
// neither requests effects nor carries the reserved `world` key. It returns
// the canonical encoding.
func parseFinish(output []byte) ([]byte, error) {
	canonical, err := transitionreg.CanonicalJSON(output, maxPhaseOutputRaw)
	if err != nil {
		return nil, lawErr(LawShape, "finish output is not well-formed JSON: %v", err)
	}
	members, err := objectMembers(canonical)
	if err != nil {
		return nil, lawErr(LawL6, "finish output is not a JSON object")
	}
	keys := make([]string, 0, len(members))
	for k := range members {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if !finishKeyAllowed(k) {
			return nil, lawErr(LawL6, "finish output carries the key %q", k)
		}
	}
	return canonical, nil
}

// ── JSON helpers over canonical bytes ─────────────────────────────────────

// kindOf classifies a canonical JSON value by its first byte: '{', '[', '"',
// 't'/'f' (bool), 'n' (null), or '0' (number).
func kindOf(raw json.RawMessage) byte {
	if len(raw) == 0 {
		return 0
	}
	switch c := raw[0]; c {
	case '{', '[', '"', 't', 'f', 'n':
		return c
	default:
		return '0'
	}
}

func objectMembers(raw []byte) (map[string]json.RawMessage, error) {
	if kindOf(raw) != '{' {
		return nil, fmt.Errorf("not an object")
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func closedKeys(m map[string]json.RawMessage, want []string) error {
	var missing, extra []string
	for _, k := range want {
		if _, ok := m[k]; !ok {
			missing = append(missing, k)
		}
	}
	for k := range m {
		if !containsString(want, k) {
			extra = append(extra, k)
		}
	}
	sort.Strings(extra)
	if len(missing) > 0 || len(extra) > 0 {
		return fmt.Errorf("keys: missing %v, unknown %v", missing, extra)
	}
	return nil
}

func containsString(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func jsonString(raw json.RawMessage) (string, bool) {
	var s string
	if kindOf(raw) != '"' || json.Unmarshal(raw, &s) != nil {
		return "", false
	}
	return s, true
}

func jsonBool(raw json.RawMessage) (bool, bool) {
	switch string(raw) {
	case "true":
		return true, true
	case "false":
		return false, true
	}
	return false, false
}

// ── Mirrors of design_docs/sketches/effectplan.ail ────────────────────────

// lawEffect is the sketch's PlannedEffect: the law's view of one effect.
type lawEffect struct {
	ID, Effect, Scope string
	Cost              int64
	PayloadBytes      int
}

func effectCountOk(n int) bool     { return n >= 0 && n <= maxPlanEffectsV1 }
func idLengthOk(id string) bool    { return len(id) >= 1 && len(id) <= maxEffectIDBytes }
func idByteOk(c int) bool          { return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') }
func idsDistinct(a, b string) bool { return a != b }

// idOk applies L2's grammar: the length law and the byte law over every byte.
func idOk(id string) bool {
	if !idLengthOk(id) {
		return false
	}
	for i := 0; i < len(id); i++ {
		if !idByteOk(int(id[i])) {
			return false
		}
	}
	return true
}

func requirementMatches(d transitionreg.EffectRequirement, e lawEffect) bool {
	return d.Effect == e.Effect && d.Scope == e.Scope && d.Cost == e.Cost
}

func declaredContains(ds []transitionreg.EffectRequirement, e lawEffect) bool {
	for _, d := range ds {
		if requirementMatches(d, e) {
			return true
		}
	}
	return false
}

func payloadSizeOk(n int) bool        { return n >= 0 && n <= maxPayloadBytes }
func reservedOutputKey(k string) bool { return k == reservedWorldKey }
func finishKeyAllowed(k string) bool  { return k != reservedWorldKey && k != finishEffectsKey }

func resultPresenceOk(n int, hasResult bool) bool { return (n == 0) == hasResult }
func finishNeedsEffect(n int, finish bool) bool   { return !finish || n > 0 }

func effectLawfulAgainst(d transitionreg.EffectRequirement, e lawEffect) bool {
	return idLengthOk(e.ID) && requirementMatches(d, e) && payloadSizeOk(e.PayloadBytes)
}

func planShapeLawful(n int, hasResult, finish bool) bool {
	return effectCountOk(n) && resultPresenceOk(n, hasResult) && finishNeedsEffect(n, finish)
}
