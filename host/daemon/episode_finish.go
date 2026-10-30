package daemon

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

const EpisodeFinishV1 = "world/episode-finish/v1"

// EpisodeFinish is the attended dispatch result, linked to its previous record.
type EpisodeFinish struct {
	Episode      string   `json:"episode"`
	Branch       string   `json:"branch"`
	Base         string   `json:"base"`
	Commit       string   `json:"commit"`
	Remote       string   `json:"remote"`
	PR           FinishPR `json:"pr"`
	Approval     string   `json:"approval"`
	Scope        string   `json:"scope"`
	Checks       []string `json:"checks"`
	Unattributed []string `json:"unattributed"`
	ReportSHA256 string   `json:"report_sha256"`
	Trailer      string   `json:"trailer"`
	Outcome      string   `json:"outcome"`
	FinishedAt   string   `json:"finished_at"`
	Previous     string   `json:"previous"`
}

type FinishPR struct {
	URL    string `json:"url"`
	Number int    `json:"number"`
	Draft  bool   `json:"draft"`
}

type FinishRequest struct {
	Phase  string         `json:"phase"`
	Scope  string         `json:"scope,omitempty"`
	Record *EpisodeFinish `json:"record,omitempty"`
}

type FinishApproval struct {
	Request  string `json:"request"`
	Approval string `json:"approval"`
}

type operatorProof struct {
	SHA256    string `json:"sha256"`
	Phase     string `json:"phase"`
	ExpiresAt int64  `json:"expires_at"`
}

var finishEpisodeName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

func finishHead(episode string) string { return "world/episode-finish/" + episode }

func (d *Daemon) readEpisodeFinish(ctx context.Context, episode string) (*EpisodeFinish, error) {
	ref, ok, err := d.store.GetRegistryHead(ctx, finishHead(episode))
	if err != nil || !ok {
		return nil, err
	}
	obj, ok, err := d.store.GetObject(ctx, ref)
	if err != nil {
		return nil, err
	}
	if !ok || obj.SemanticID != EpisodeFinishV1 {
		return nil, fmt.Errorf("finish head names no finish record")
	}
	var record EpisodeFinish
	if err := json.Unmarshal(obj.Payload, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// consumeOperatorProof runs under humanMu. Lstat refuses links and permissions
// before opening; unlink precedes every store mutation, including approval.
func (d *Daemon) consumeOperatorProof(episode, phase, proof string) error {
	dir := filepath.Join(filepath.Dir(d.credentialsPath), "operator-proofs")
	dirInfo, err := os.Lstat(dir)
	if err != nil {
		return fmt.Errorf("operator proof directory is missing")
	}
	dirStat, owned := dirInfo.Sys().(*syscall.Stat_t)
	if !dirInfo.IsDir() || dirInfo.Mode().Perm() != 0o700 || !owned || dirStat.Uid != uint32(os.Getuid()) {
		return fmt.Errorf("operator proof directory must be owned and 0700")
	}
	path := filepath.Join(dir, episode+".json")
	info, err := os.Lstat(path)
	if err != nil {
		return fmt.Errorf("operator proof is missing")
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !info.Mode().IsRegular() || info.Mode() != 0o600 || !ok || stat.Uid != uint32(os.Getuid()) {
		return fmt.Errorf("operator proof must be an owned regular 0600 file")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("operator proof cannot be read")
	}
	var p operatorProof
	now := time.Now().Unix()
	if json.Unmarshal(data, &p) != nil || p.Phase != phase || p.ExpiresAt <= now || p.ExpiresAt > now+300 {
		return fmt.Errorf("operator proof phase or expiry is invalid")
	}
	raw, err := hex.DecodeString(proof)
	if err != nil || len(raw) != 32 {
		return fmt.Errorf("operator proof is invalid")
	}
	want, err := hex.DecodeString(p.SHA256)
	sum := sha256.Sum256([]byte(proof))
	if err != nil || len(want) != len(sum) || subtle.ConstantTimeCompare(want, sum[:]) != 1 {
		return fmt.Errorf("operator proof does not match")
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("operator proof cannot be consumed")
	}
	return nil
}

func (d *Daemon) handleEpisodeFinish(w http.ResponseWriter, r *http.Request) {
	refuse := func(what string) {
		writeAPIError(w, "Forbidden", what+"\nwhy: the agent holds the episode's session, so a session-gated write is one the agent could make\nnext: run ailang-worldd episode finish at a terminal", http.StatusForbidden)
	}
	for name := range r.Header {
		if strings.EqualFold(name, "Authorization") {
			refuse("operator route: session credentials are refused")
			return
		}
	}
	episode := r.PathValue("episode")
	if !finishEpisodeName.MatchString(episode) || len(episode) > maxReportEpisodeBytes {
		refuse("operator route: invalid episode")
		return
	}
	// No body is read before session credentials have been refused.
	var input FinishRequest
	dec := json.NewDecoder(io.LimitReader(r.Body, 128<<10))
	dec.DisallowUnknownFields()
	if dec.Decode(&input) != nil || dec.Decode(new(any)) != io.EOF || (input.Phase != "approve" && input.Phase != "record") {
		refuse("operator route: invalid finish request")
		return
	}
	d.humanMu.Lock()
	defer d.humanMu.Unlock()
	if err := d.consumeOperatorProof(episode, input.Phase, r.Header.Get("X-World-Operator-Proof")); err != nil {
		refuse("operator route: " + err.Error())
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
	defer cancel()
	if input.Phase == "approve" {
		scope, err := broker.ParseFinishApprovalScope(input.Scope)
		now := time.Now().Unix()
		if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {
			refuse("operator route: invalid or expired finish scope")
			return
		}
		request, decision, err := broker.RecordOperatorApproval(ctx, d.store, input.Scope, "operator:episode-finish", "operator:terminal", now)
		if err != nil {
			d.writeInternalError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, FinishApproval{request.String(), decision.String()})
		return
	}
	record := input.Record
	if record == nil || input.Scope != "" || record.Episode != episode {
		refuse("operator route: record episode does not match")
		return
	}
	scope, scopeText, err := d.finishDecisionScope(ctx, record.Approval)
	if err != nil || scope.Head != record.Branch || scope.SHA != record.Commit || scope.Remote != record.Remote {
		refuse("operator route: record does not match an approved finish scope")
		return
	}
	if record.Scope != "" && record.Scope != scopeText {
		refuse("operator route: record scope differs from approval")
		return
	}
	record.Scope = scopeText
	switch record.Outcome {
	case "done", "push-failed", "pr-failed", "ambiguous":
	default:
		refuse("operator route: invalid outcome")
		return
	}
	if record.Outcome == "done" && (record.PR.URL == "" || record.PR.Number <= 0) {
		refuse("operator route: done requires an open PR")
		return
	}
	usedHead := "world/episode-finish-used/" + strings.TrimPrefix(record.Approval, "sha256:")
	if _, used, err := d.store.GetRegistryHead(ctx, usedHead); err != nil {
		d.writeInternalError(w, r, err)
		return
	} else if used {
		writeAPIError(w, "Conflict", "finish approval already consumed\nwhy: done is single use\nnext: inspect the episode report", http.StatusConflict)
		return
	}
	previous, ok, err := d.store.GetRegistryHead(ctx, finishHead(episode))
	if err != nil {
		d.writeInternalError(w, r, err)
		return
	}
	record.Previous = ""
	if ok {
		record.Previous = previous.String()
	}
	payload, err := json.Marshal(record)
	if err != nil {
		d.writeInternalError(w, r, err)
		return
	}
	obj := store.Object{Hash: hashref.SumSHA256(payload), InterfaceHash: hashref.SumSHA256([]byte(EpisodeFinishV1)),
		SemanticID: EpisodeFinishV1, Provenance: "operator:episode-finish", Payload: payload}
	if err = d.store.PutObject(ctx, obj); err == nil && record.Outcome == "done" {
		err = d.store.SetRegistryHead(ctx, usedHead, obj.Hash)
	}
	if err == nil {
		err = d.store.SetRegistryHead(ctx, finishHead(episode), obj.Hash)
	}
	if err != nil {
		d.writeInternalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"record": obj.Hash.String()})
}

func (d *Daemon) finishDecisionScope(ctx context.Context, text string) (broker.FinishApprovalScope, string, error) {
	bad := func() (broker.FinishApprovalScope, string, error) {
		return broker.FinishApprovalScope{}, "", fmt.Errorf("invalid finish approval")
	}
	ref, err := hashref.Parse(text)
	if err != nil {
		return bad()
	}
	obj, ok, err := d.store.GetObject(ctx, ref)
	if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 {
		return bad()
	}
	var decision struct {
		RequestRef string `json:"requestRef"`
		Decision   string `json:"decision"`
	}
	if json.Unmarshal(obj.Payload, &decision) != nil || decision.Decision != "approve" {
		return bad()
	}
	ref, err = hashref.Parse(decision.RequestRef)
	if err != nil {
		return bad()
	}
	obj, ok, err = d.store.GetObject(ctx, ref)
	if err != nil || !ok || obj.SemanticID != broker.ApprovalRequestV1 {
		return bad()
	}
	request, err := broker.DecodeApprovalRequest(obj.Payload)
	if err != nil {
		return bad()
	}
	scope, err := broker.ParseFinishApprovalScope(request.Scope)
	now := time.Now().Unix()
	if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 {
		return bad()
	}
	return scope, request.Scope, nil
}
