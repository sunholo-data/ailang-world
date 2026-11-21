package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const episodeFileProofBytes = 1 << 20

func (d *Daemon) handleSurfaceEpisodeFile(w http.ResponseWriter, r *http.Request, query url.Values, path, source string) {
	refuse := func(status int, reason string) {
		d.writeSurfaceJSON(w, r, status, surfaceFileSchema, surfaceReason{Schema: surfaceFileSchema, Reason: reason})
	}
	if d.follow == nil {
		refuse(http.StatusNotFound, "no recorded episode file for this selection")
		return
	}
	if d.follow.fileReads.Add(1) > 2 {
		d.follow.fileReads.Add(-1)
		refuse(http.StatusServiceUnavailable, "episode file preview is busy")
		return
	}
	defer d.follow.fileReads.Add(-1)
	ctx, cancel := d.readCtx(r)
	defer cancel()
	env := d.followSnapshot().Envelope
	selection, through, err := selectSurfaceEpisode(env, query.Get("brane"), query.Get("at"), time.Now())
	if errors.Is(err, errSurfaceEpisodeAddress) {
		refuse(http.StatusBadRequest, errSurfaceEpisodeAddress.Error())
		return
	}
	if err != nil {
		refuse(http.StatusNotFound, "no recorded episode file for this selection")
		return
	}
	start := d.readFollowEpisodeStart(ctx, env, selection)
	if start == nil {
		refuse(http.StatusNotFound, "no recorded episode file for this selection")
		return
	}
	evidence := buildEpisodeFileActivity(ctx, d.intents, d.reads, start.Episode, through)
	rows := evidence.Reads
	if source == "episode-write" {
		rows = evidence.Writes
	}
	for _, row := range rows {
		if row.Path != path {
			continue
		}
		for _, repo := range d.follow.config.Repos {
			if repo.ID == selection.Repo {
				file := d.episodeFilePreview(ctx, env, selection, repo, start.Episode, start.Base.SHA, row, source)
				d.writeSurfaceJSON(w, r, http.StatusOK, surfaceFileSchema, file)
				return
			}
		}
	}
	refuse(http.StatusNotFound, "no recorded episode file for this selection")
}

// episodeFilePreview only emits bytes read from immutable Git objects. The
// stored effect supplies a comparison digest, never browser-visible content.
func (d *Daemon) episodeFilePreview(ctx context.Context, env surface.FollowIndex, selection surfacepage.FollowStack,
	repo surfacefollow.Repo, episode, base string, row surfaceFileUse, source string) surfacefollow.FilePreview {
	ctx, cancel := context.WithTimeout(ctx, surfacefollow.CallTimeout)
	defer cancel()
	unavailable := surfacefollow.FilePreview{Path: row.Path, Kind: "unavailable", Status: "read", SourceLabel: "Read version unavailable"}
	if source == "episode-read" {
		digest := d.episodeFileDigest(ctx, episode, row.Path, "read", row.FirstEntry)
		snapshot, err := surfacefollow.ReadCommittedFileVersion(ctx, repo, base, row.Path, hashref.HashRef{})
		if err == nil {
			snapshot.Status, snapshot.SourceLabel = "read", "At episode start"
			if snapshot.Text != nil && hashref.SumSHA256([]byte(*snapshot.Text)) == digest {
				snapshot.SourceLabel = "Verified recorded read"
				return snapshot
			}
		}
		if !digest.IsZero() {
			finish := d.episodeFileFinishCommit(ctx, selection, episode)
			for _, commit := range episodeFileCommitCandidates(env, selection, finish) {
				if commit == base {
					continue
				}
				file, err := surfacefollow.ReadCommittedFileVersion(ctx, repo, commit, row.Path, digest)
				if err == nil {
					file.Status, file.SourceLabel = "read", "Verified recorded read"
					return file
				}
			}
		}
		if err == nil {
			return snapshot
		}
		return unavailable
	}

	unavailable.Status, unavailable.SourceLabel = "write", "Recorded write version unavailable"
	finishCommit := d.episodeFileFinishCommit(ctx, selection, episode)
	var digest hashref.HashRef
	if row.Op == "write" {
		digest = d.episodeFileDigest(ctx, episode, row.Path, "write", row.lastEntry)
	}
	if !digest.IsZero() {
		for _, commit := range episodeFileCommitCandidates(env, selection, finishCommit) {
			file, err := surfacefollow.ReadCommittedFileVersion(ctx, repo, commit, row.Path, digest)
			if err == nil {
				file.Status, file.SourceLabel = "write", "Verified written version"
				return file
			}
		}
	}
	if finishCommit != "" {
		file, err := surfacefollow.ReadCommittedFileVersion(ctx, repo, finishCommit, row.Path, hashref.HashRef{})
		if err == nil {
			file.Status, file.SourceLabel = "write", "At episode finish"
			return file
		}
	}
	return unavailable
}

func (d *Daemon) episodeFileFinishCommit(ctx context.Context, selection surfacepage.FollowStack, episode string) string {
	if d.store == nil {
		return ""
	}
	finish, err := newFollowEpisodes(d.credentialsPath, d.store).trustedFinish(ctx, episode)
	if err != nil || finish == nil {
		return ""
	}
	_, _, selectedCommit, _ := surface.ParseQualifiedID(selection.Front.ID)
	if selection.Front.Kind == surface.KindStep && selectedCommit != finish.Commit {
		return ""
	}
	return finish.Commit
}

func episodeFileCommitCandidates(env surface.FollowIndex, selection surfacepage.FollowStack, finish string) []string {
	commits := []string{}
	if finish != "" {
		commits = append(commits, finish)
	}
	if selection.Front.Kind != surface.KindStep {
		return commits
	}
	_, _, id, _ := surface.ParseQualifiedID(selection.Front.ID)
	if id != finish {
		commits = append(commits, id)
	}
	for _, partition := range env.Repos {
		if partition.ID == selection.Repo {
			for _, state := range partition.Index.States {
				if state.ID == id && state.Merge != "" && state.Merge != finish {
					return append(commits, state.Merge)
				}
			}
		}
	}
	return commits
}

// episodeFileDigest rechecks the exact successful invocation behind a row.
// For an explicit write it hashes the actual bytes passed to Workspace.Write,
// not the outer tool arguments. Ranged reads do not prove full-file identity.
func (d *Daemon) episodeFileDigest(ctx context.Context, episode, path, op string, entry int64) hashref.HashRef {
	var zero hashref.HashRef
	if d.reads == nil {
		return zero
	}
	log, ok, err := d.reads.GetLogEntry(ctx, entry)
	if err != nil || !ok {
		return zero
	}
	obj, ok, err := d.reads.GetObject(ctx, log.TransitionRef)
	var call reportRecord
	if err != nil || !ok || obj.SemanticID != coordinator.RecordV2 || len(obj.Payload) > episodeFileProofBytes ||
		json.Unmarshal(obj.Payload, &call) != nil || call.EpisodeID != episode || call.SkillID != "ailang-"+op || len(call.Effects) != 1 {
		return zero
	}
	ref, err := hashref.Parse(call.Effects[0])
	if err != nil {
		return zero
	}
	obj, ok, err = d.reads.GetObject(ctx, ref)
	if err != nil || !ok || obj.SemanticID != broker.EffectRecordV1 || obj.Provenance != "host/broker" {
		return zero
	}
	effect, err := broker.DecodeRecord(obj.Payload)
	want := broker.EffectWorkspaceRead
	if op == "write" {
		want = broker.EffectWorkspaceWrite
	}
	if err != nil || !effect.Allowed || effect.Failed || effect.Scope != broker.WorkspaceScope || effect.Effect != want {
		return zero
	}
	obj, ok, err = d.reads.GetObject(ctx, effect.RequestRef)
	if err != nil || !ok || obj.SemanticID != broker.EffectRequestV1 || obj.Provenance != "host/broker" || len(obj.Payload) > episodeFileProofBytes {
		return zero
	}
	request := parseEpisodeFileRequest(effect, obj.Payload)
	if request.Op != op || request.Path != path {
		return zero
	}
	obj, ok, err = d.reads.GetObject(ctx, effect.ResultRef)
	var result struct {
		OK         bool   `json:"ok"`
		Content    string `json:"content"`
		TotalLines *int   `json:"total_lines"`
	}
	if err != nil || !ok || obj.SemanticID != broker.EffectResultV1 || obj.Provenance != "host/broker" ||
		len(obj.Payload) > episodeFileProofBytes || json.Unmarshal(obj.Payload, &result) != nil || !result.OK {
		return zero
	}
	if op == "write" && request.Content != nil {
		return hashref.SumSHA256([]byte(*request.Content))
	}
	if op == "read" && result.TotalLines == nil && request.Offset == nil && request.Limit == nil {
		return hashref.SumSHA256([]byte(result.Content))
	}
	return zero
}

type episodeFileRequest struct {
	Path    string  `json:"path"`
	Op      string  `json:"op"`
	Content *string `json:"content"`
	Offset  *int    `json:"offset"`
	Limit   *int    `json:"limit"`
}

func parseEpisodeFileRequest(effect broker.EffectRecord, raw []byte) episodeFileRequest {
	prefix := fmt.Sprintf("%d:%s%d:%s%d:", len(effect.Effect), effect.Effect, len(effect.Scope), effect.Scope, effect.Cost)
	body, ok := strings.CutPrefix(string(raw), prefix)
	if !ok {
		return episodeFileRequest{}
	}
	timestamp, payload, ok := strings.Cut(body, ":")
	if _, err := strconv.ParseInt(timestamp, 10, 64); err != nil || !ok {
		return episodeFileRequest{}
	}
	var request episodeFileRequest
	if json.Unmarshal([]byte(payload), &request) != nil {
		return episodeFileRequest{}
	}
	return request
}
