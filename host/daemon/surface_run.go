package daemon

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

func durationBetween(start, end string) (time.Duration, bool) {
	a, e1 := time.Parse(time.RFC3339, start)
	b, e2 := time.Parse(time.RFC3339, end)
	return b.Sub(a), e1 == nil && e2 == nil && !b.Before(a)
}

func (d *Daemon) surfaceRun(ctx context.Context, start *store.EpisodeStart) (*surfacepage.NodeRun, *EpisodeFinish) {
	if start == nil {
		return nil, nil
	}
	run := &surfacepage.NodeRun{
		Episode: start.Episode, Harness: surface.SanitizeOpenText(start.Agent, 64),
		Model: surface.SanitizeOpenText(start.Model, 96), Launches: len(start.Launches),
	}
	run.Turns, run.CostUSD, run.Tokens = d.surfaceRunUsage(start)
	run.Permissions = d.surfacePermissions(ctx, start, time.Now().Unix())
	if len(start.Launches) > 0 {
		var total time.Duration
		for i, launch := range start.Launches {
			delta, ok := durationBetween(launch.StartedAt, launch.EndedAt)
			if !ok {
				break
			}
			total += delta
			if i == len(start.Launches)-1 {
				seconds := int64(total / time.Second)
				run.Runtime = &seconds
			}
		}
	}
	decision := func(action, at string) {
		if _, err := time.Parse(time.RFC3339, at); err == nil {
			run.Decisions = append(run.Decisions, surfacepage.NodeDecision{Action: action, At: at})
		}
	}
	decision("Authorized episode start", start.StartedAt)
	for _, resume := range start.Resumes {
		if !resume.Failed {
			decision("Renewed session", resume.ResumedAt)
		}
	}
	if d.store == nil {
		return run, nil
	}
	ref, ok, err := d.store.GetRegistryHead(ctx, finishHead(start.Episode))
	if err != nil || !ok {
		return run, nil
	}
	obj, ok, err := d.store.GetObject(ctx, ref)
	var finish EpisodeFinish
	if err != nil || !ok || obj.SemanticID != EpisodeFinishV1 || obj.Provenance != finishProvenance ||
		json.Unmarshal(obj.Payload, &finish) != nil || finish.Episode != start.Episode {
		return run, nil
	}
	if delta, valid := durationBetween(start.StartedAt, finish.FinishedAt); valid {
		seconds := int64(delta / time.Second)
		run.Elapsed = &seconds
	}
	if approval, err := hashref.Parse(finish.Approval); err == nil {
		record, found, err := d.store.GetObject(ctx, approval)
		var answer struct {
			Decision string
			Now      int64
		}
		if err == nil && found && record.SemanticID == broker.ApprovalDecisionV1 &&
			json.Unmarshal(record.Payload, &answer) == nil && answer.Decision == "approve" && answer.Now > 0 {
			decision("Approved pull request creation", time.Unix(answer.Now, 0).UTC().Format(time.RFC3339))
		}
	}
	for _, id := range finish.Checks {
		if ctx.Err() != nil {
			break
		}
		hash, err := hashref.Parse(id)
		if err != nil {
			continue
		}
		obj, ok, err := d.store.GetObject(ctx, hash)
		if !ok || err != nil || obj.SemanticID != "world/task-outcome/v1" || obj.Provenance != "harness:floor" {
			continue
		}
		var check struct {
			EpisodeID, Label, Task string
			Outcome                struct{ Class string }
		}
		if json.Unmarshal(obj.Payload, &check) != nil || check.EpisodeID != start.Episode || check.Label != "episode-finish" ||
			(check.Outcome.Class != "PASS" && check.Outcome.Class != "FAIL") {
			continue
		}
		run.Checks = append(run.Checks, surfacepage.NodeCheck{Name: surface.SanitizeOpenText(check.Task, 96), Result: check.Outcome.Class})
	}
	return run, &finish
}

const (
	nodePermissionSessionLimit = 32
	nodePermissionGrantLimit   = 128
	nodePermissionReceiptLimit = 512
)

type surfacePermissionUsage struct {
	debited, used int64
	budget        *int64
}

func (d *Daemon) surfacePermissionUsage(ctx context.Context, episode string) (map[store.EffectKey]surfacePermissionUsage, bool) {
	if d.store == nil {
		return nil, false
	}
	receipts, err := d.store.EpisodeEffectReceipts(ctx, episode, nodePermissionReceiptLimit+1)
	if err != nil || len(receipts) > nodePermissionReceiptLimit {
		return nil, false
	}
	type observedReceipt struct {
		intent  *store.EffectIntent
		outcome *store.EffectOutcome
		debited int64
	}
	totals := make(map[store.EffectKey]surfacePermissionUsage)
	last := make(map[store.EffectKey]observedReceipt)
	for _, receipt := range receipts {
		intent, outcome := receipt.EffectIntent, receipt.EffectOutcome
		key := store.EffectKey{Effect: intent.Effect, Scope: intent.Scope}
		total := totals[key]
		if intent.Cost < 0 || intent.Cost > math.MaxInt64-total.debited {
			return nil, false
		}
		total.debited += intent.Cost
		if outcome != nil && (outcome.Status == "succeeded" || outcome.Status == "failed") {
			total.used += intent.Cost
			last[key] = observedReceipt{intent, outcome, total.debited}
		}
		totals[key] = total
	}
	for key, evidence := range last {
		obj, ok, err := d.store.GetObject(ctx, evidence.outcome.RecordRef)
		if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 {
			continue
		}
		record, err := broker.DecodeRecord(obj.Payload)
		if err != nil || !record.Allowed || record.Failed != (evidence.outcome.Status == "failed") ||
			record.Effect != key.Effect || record.Scope != key.Scope || record.Cost != evidence.intent.Cost ||
			record.RequestRef != evidence.intent.RequestRef || !broker.RecordConsistent(record) ||
			record.BudgetAfter < 0 || record.BudgetAfter > math.MaxInt64-evidence.debited {
			continue
		}
		total := totals[key]
		budget := record.BudgetAfter + evidence.debited
		total.budget = &budget
		totals[key] = total
	}
	return totals, true
}

func (d *Daemon) surfacePermissions(ctx context.Context, start *store.EpisodeStart, now int64) *surfacepage.NodePermissions {
	if start == nil || d.credentialsPath == "" || !surface.EpisodeGrammar.MatchString(start.Episode) {
		return nil
	}
	credentials, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return nil
	}
	defer credentials.Close()
	rows, err := credentials.ListSessions(ctx, start.Episode, nodePermissionSessionLimit+1)
	if err != nil {
		return nil
	}
	out := &surfacepage.NodePermissions{Sessions: []surfacepage.NodePermissionSession{}}
	sessionsTruncated := len(rows) > nodePermissionSessionLimit
	if sessionsTruncated {
		out.Truncated = true
		rows = rows[:nodePermissionSessionLimit]
	}
	current := store.CurrentEpisodeCredential(start, rows)
	totals, complete := d.surfacePermissionUsage(ctx, start.Episode)
	if !complete {
		out.Truncated = true
	}
	ids := map[store.EffectKey]int{}
	retained := map[store.EffectKey]bool{}
	project := func(key store.EffectKey, state string, budget *int64) surfacepage.NodeGrant {
		if ids[key] == 0 {
			ids[key] = len(ids) + 1
		}
		grant := surfacepage.NodeGrant{Effect: surface.SanitizeOpenText(key.Effect, 96), Scope: surfaceGrantScope(key.Scope),
			ScopeKey: ids[key], State: state, Budget: budget}
		if complete {
			total := totals[key]
			grant.Used = &total.used
			reserved := total.debited - total.used
			grant.Reserved = &reserved
			if budget != nil {
				left := max(0, *budget-total.debited)
				grant.LedgerRemaining = &left
				if state == "active" {
					grant.Remaining = &left
				}
			}
		}
		if state == "expired" {
			zero := int64(0)
			grant.Remaining = &zero
		}
		return grant
	}
	stamp := func(t int64) string { return time.Unix(t, 0).UTC().Format(time.RFC3339) }
	known := map[string]bool{}
	for _, row := range rows {
		known[row.CredentialID] = true
		state := "active"
		if now > row.ExpiresAt {
			state = "expired"
		}
		session := surfacepage.NodePermissionSession{State: state, Current: row.CredentialID == current,
			CreatedAt: stamp(row.CreatedAt), ExpiresAt: stamp(row.ExpiresAt), Grants: []surfacepage.NodeGrant{}}
		var caps []broker.Capability
		if len(row.GrantsJSON) > 64<<10 || json.Unmarshal([]byte(row.GrantsJSON), &caps) != nil {
			out.Truncated = true
			session.State = "unavailable"
			caps = nil
		}
		for i, cap := range caps {
			if i == nodePermissionGrantLimit {
				out.Truncated = true
				break
			}
			if cap.Budget < 0 || cap.Effect == "" || cap.Scope == "" {
				out.Truncated = true
				continue
			}
			key := store.EffectKey{Effect: cap.Effect, Scope: cap.Scope}
			retained[key] = true
			grantState := state
			if now >= cap.ExpiresAt {
				grantState = "expired"
			}
			grant := project(key, grantState, &cap.Budget)
			grant.ExpiresAt = stamp(cap.ExpiresAt)
			session.Grants = append(session.Grants, grant)
		}
		out.Sessions = append(out.Sessions, session)
	}
	if !sessionsTruncated {
		// Revocation deletes credential rows. Only the recorded resume chain
		// can establish a missing ID; no grant values are inferred from it.
		missing := func(id string, at string) {
			if id == "" || known[id] {
				return
			}
			if len(out.Sessions) == nodePermissionSessionLimit {
				out.Truncated = true
				return
			}
			known[id] = true
			out.Sessions = append(out.Sessions, surfacepage.NodePermissionSession{
				State: "revoked-or-unavailable", Current: id == current, CreatedAt: at, Grants: []surfacepage.NodeGrant{}})
		}
		for _, resume := range start.Resumes {
			missing(resume.ResumedFrom, "")
			if !resume.Failed {
				missing(resume.CredentialID, resume.ResumedAt)
			}
		}
		if len(rows) == 0 && len(out.Sessions) == 0 {
			out.Sessions = append(out.Sessions, surfacepage.NodePermissionSession{
				State: "revoked-or-unavailable", Current: true, CreatedAt: start.StartedAt, Grants: []surfacepage.NodeGrant{}})
		}
	}
	keys := make([]store.EffectKey, 0, len(totals))
	for key := range totals {
		if !retained[key] {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Effect != keys[j].Effect {
			return keys[i].Effect < keys[j].Effect
		}
		return keys[i].Scope < keys[j].Scope
	})
	for _, key := range keys {
		out.Observed = append(out.Observed, project(key, "observed", totals[key].budget))
	}
	return out
}

func surfaceGrantScope(scope string) string {
	if strings.ContainsAny(scope, "/\\;:@ \t\r\n") {
		return "[private scope]"
	}
	if safe := surface.SanitizeOpenText(scope, 120); safe != "" {
		return safe
	}
	return "[private scope]"
}
