package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"sort"
	"time"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const (
	trajectorySchema     = "world/surface-trajectory/v1"
	maxTrajectoryEvents  = 120
	maxTrajectoryScan    = 999
	maxTrajectoryLookups = 512
)

type surfaceTrajectory struct {
	Schema        string                   `json:"schema"`
	Episode       string                   `json:"episode"`
	SourceState   string                   `json:"sourceState"`
	Harness       trajectoryHarness        `json:"harness"`
	World         trajectoryWorld          `json:"world"`
	RawTranscript *trajectoryRawTranscript `json:"rawTranscript,omitempty"`
}

type trajectoryRawTranscript struct {
	Filename string `json:"filename"`
	SHA256   string `json:"sha256"`
}

type trajectoryHarness struct {
	Available bool              `json:"available"`
	Events    []trajectoryEvent `json:"events"`
	Truncated bool              `json:"truncated,omitempty"`
}

type trajectoryEvent struct {
	Launch int    `json:"launch"`
	Kind   string `json:"kind"`
	At     string `json:"at,omitempty"`
}

type trajectoryWorld struct {
	Available bool             `json:"available"`
	Total     int              `json:"total,omitempty"`
	Calls     []trajectoryCall `json:"calls"`
	Truncated bool             `json:"truncated,omitempty"`
}

type trajectoryCall struct {
	Entry   int64  `json:"entry"`
	At      string `json:"at"`
	Tool    string `json:"tool"`
	Surface string `json:"surface"`
	Target  string `json:"target,omitempty"`
	Status  string `json:"status,omitempty"`
}

func (d *Daemon) handleSurfaceTrajectory(w http.ResponseWriter, r *http.Request) {
	refuse := func(status int, reason string) {
		d.writeSurfaceJSON(w, r, status, trajectorySchema, surfaceReason{Schema: trajectorySchema, Reason: reason})
	}
	if len(r.URL.RawQuery) > 2048 {
		refuse(http.StatusBadRequest, "expected a brane and state")
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1 {
		refuse(http.StatusBadRequest, "expected a brane and state")
		return
	}
	if d.follow == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	snapshot := d.followSnapshot()
	env := snapshot.Envelope
	selection, through, err := selectSurfaceEpisode(env, q.Get("brane"), q.Get("at"), time.Now())
	if errors.Is(err, errSurfaceEpisodeAddress) {
		refuse(http.StatusBadRequest, errSurfaceEpisodeAddress.Error())
		return
	}
	if err != nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	if d.follow.trajectoryReads.Add(1) > 1 {
		d.follow.trajectoryReads.Add(-1)
		refuse(http.StatusServiceUnavailable, "episode trajectory is busy")
		return
	}
	defer d.follow.trajectoryReads.Add(-1)

	ctx, cancel := d.readCtx(r)
	defer cancel()
	start := d.readFollowEpisodeStart(ctx, env, selection)
	if start == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	harness, raw := d.readTrajectoryHarness(ctx, start, through)
	result := surfaceTrajectory{Schema: trajectorySchema, Episode: start.Episode, SourceState: selection.Front.ID,
		Harness: harness, RawTranscript: raw,
		World: readTrajectoryWorld(ctx, d.intents, d.reads, start.Episode, through),
	}
	d.writeSurfaceJSON(w, r, http.StatusOK, trajectorySchema, result)
}

// A shared merge has no launch of its own. For linked episode evidence, resolve
// its exact second-parent step; the episode binding is still verified against
// the follower alias and configured repository by readFollowEpisodeStart.
func trajectorySelection(env surface.FollowIndex, selection surfacepage.FollowStack) surfacepage.FollowStack {
	if selection.Front.Kind == surface.KindShared && selection.Front.PR > 0 {
		_, _, merge, _ := surface.ParseQualifiedID(selection.Front.ID)
		if head := surfacePRHead(env, selection.Repo, merge, selection.Front.PR); head != "" {
			selection.Front.ID, _ = surface.QualifyID(selection.Repo, "commit", head)
			selection.Front.Kind = surface.KindStep
		}
	}
	return selection
}

// readTrajectoryHarness opens only the daemon service's fixed transcript leaf.
// Decode only event discriminators and timestamps. Native content is private.
func (d *Daemon) readTrajectoryHarness(ctx context.Context, start *store.EpisodeStart, through time.Time) (trajectoryHarness, *trajectoryRawTranscript) {
	empty := trajectoryHarness{Events: []trajectoryEvent{}}
	if !surface.EpisodeGrammar.MatchString(start.Episode) || d.cfg.DBPath == "" ||
		len(start.Launches) == 0 || len(start.Launches) > maxSurfaceUsageLaunches {
		return empty, nil
	}
	last := start.Launches[len(start.Launches)-1]
	finalDigest, err := hashref.Parse(last.TranscriptSHA256)
	if err != nil || (last.TranscriptState != "kept" && last.TranscriptState != "scrubbed") {
		return empty, nil
	}
	root := openSurfaceDir(filepath.Join(filepath.Dir(d.cfg.DBPath), "sessions"))
	if root == nil {
		return empty, nil
	}
	defer root.Close()
	filename := start.Episode + ".transcript.jsonl"
	data := readSurfaceFile(root, filename, maxSurfaceUsageBytes)
	if len(data) == 0 || ctx.Err() != nil || hashref.SumSHA256(data) != finalDigest {
		return empty, nil
	}
	var raw *trajectoryRawTranscript
	if trajectoryThrough(through, last.EndedAt) {
		raw = &trajectoryRawTranscript{Filename: filename, SHA256: finalDigest.Digest()}
	}
	if !utf8.Valid(data) || (start.Agent != "claude" && start.Agent != "codex" && start.Agent != "agy") {
		return empty, raw
	}
	digests := make([]hashref.HashRef, 0, len(start.Launches))
	for _, launch := range start.Launches {
		digest, err := hashref.Parse(launch.TranscriptSHA256)
		if err != nil || launch.TranscriptState != "kept" {
			return empty, raw
		}
		digests = append(digests, digest)
	}
	return projectTrajectoryHarness(ctx, start, through, data, digests), raw
}

func trajectoryThrough(through time.Time, ts string) bool {
	if through.IsZero() {
		return true
	}
	stamp, err := time.Parse(time.RFC3339, ts)
	return err == nil && !stamp.After(through)
}

func projectTrajectoryHarness(ctx context.Context, start *store.EpisodeStart, through time.Time, data []byte, digests []hashref.HashRef) trajectoryHarness {
	out := trajectoryHarness{Events: []trajectoryEvent{}}
	empty := trajectoryHarness{Events: []trajectoryEvent{}}
	launch := 0
	allowed := trajectoryThrough(through, start.Launches[0].EndedAt)
	for offset := 0; offset < len(data); {
		if ctx.Err() != nil {
			return empty
		}
		line, rest, separated := bytes.Cut(data[offset:], []byte{'\n'})
		if bytes.HasPrefix(line, []byte(`{"world_launch":`)) {
			if !separated || launch+1 >= len(digests) || string(line) != fmt.Sprintf(`{"world_launch": %d}`, launch+2) ||
				!trajectoryPrefixMatches(data[:offset], digests[launch]) {
				return empty
			}
			launch++
			allowed = trajectoryThrough(through, start.Launches[launch].EndedAt)
		} else if allowed {
			if event := harnessTrajectoryEvent(line); event.Kind != "" {
				event.Launch = launch + 1
				if len(out.Events) < maxTrajectoryEvents {
					out.Events = append(out.Events, event)
				} else {
					out.Truncated = true
				}
			}
		}
		offset = len(data) - len(rest)
	}
	if launch+1 != len(digests) {
		return empty
	}
	out.Available = len(out.Events) > 0
	return out
}

// The launch writer inserts a newline before a separator if the previous
// transcript ended with an incomplete line. That newline was not in its digest.
func trajectoryPrefixMatches(prefix []byte, want hashref.HashRef) bool {
	return hashref.SumSHA256(prefix) == want ||
		bytes.HasSuffix(prefix, []byte{'\n'}) && hashref.SumSHA256(prefix[:len(prefix)-1]) == want
}

func harnessTrajectoryEvent(line []byte) trajectoryEvent {
	var e struct {
		Type      string            `json:"type"`
		Subtype   string            `json:"subtype"`
		Timestamp string            `json:"timestamp"`
		CreatedAt string            `json:"created_at"`
		ToolCalls []json.RawMessage `json:"tool_calls"`
		Payload   struct {
			Type string `json:"type"`
		} `json:"payload"`
	}
	if json.Unmarshal(line, &e) != nil {
		return trajectoryEvent{}
	}
	var kind string
	switch e.Type {
	case "thread.started", "session_meta":
		kind = "session-started"
	case "system":
		if e.Subtype == "init" {
			kind = "session-started"
		}
	case "turn.started", "USER_INPUT":
		kind = "turn-started"
	case "turn.completed":
		kind = "turn-completed"
	case "turn.failed":
		kind = "turn-failed"
	case "item.completed":
		kind = "item-completed"
	case "assistant":
		kind = "assistant-event"
	case "PLANNER_RESPONSE":
		if len(e.ToolCalls) > 0 {
			kind = "tool-called"
		} else {
			kind = "assistant-event"
		}
	case "GENERIC":
		kind = "tool-result"
	case "result":
		kind = "run-result"
	case "event_msg":
		switch e.Payload.Type {
		case "task_started":
			kind = "turn-started"
		case "task_complete":
			kind = "turn-completed"
		case "turn_aborted":
			kind = "turn-failed"
		}
	case "response_item":
		switch e.Payload.Type {
		case "function_call":
			kind = "tool-called"
		case "function_call_output":
			kind = "tool-result"
		}
	}
	event := trajectoryEvent{Kind: kind}
	ts := e.Timestamp
	if ts == "" {
		ts = e.CreatedAt
	}
	if stamp, err := time.Parse(time.RFC3339, ts); err == nil {
		event.At = stamp.UTC().Format(time.RFC3339)
	}
	return event
}

var (
	trajectoryToolName    = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.-]{0,95}$`)
	trajectoryExecCommand = regexp.MustCompile(`^[a-z][a-z0-9-]{0,31}$`)
)

func trajectoryCallMeta(ctx context.Context, reads objectWorldReader, lookups *int, tool, inputRef string, effects []string) (target, status string) {
	if inHash, err := hashref.Parse(inputRef); err == nil && *lookups > 0 {
		*lookups--
		if inObj, ok, err := reads.GetObject(ctx, inHash); err == nil && ok && len(inObj.Payload) <= 64<<10 {
			var args struct {
				Path    string `json:"path"`
				Command string `json:"command"`
				Offset  int    `json:"offset"`
				Limit   int    `json:"limit"`
			}
			if json.Unmarshal(inObj.Payload, &args) == nil {
				switch {
				case args.Path != "":
					if (tool == "workspace-list" || tool == "workspace-search") && args.Path == "." {
						target = "."
					} else if safeEpisodeFilePath(args.Path) {
						if tool == "ailang-read" && args.Offset > 0 && args.Limit > 0 {
							target = fmt.Sprintf("%s:%d-%d", args.Path, args.Offset, args.Offset+args.Limit-1)
						} else {
							target = args.Path
						}
					}
				case tool == "workspace-list":
					target = "."
				case tool == "workspace-exec" && trajectoryExecCommand.MatchString(args.Command):
					target = args.Command
				}
			}
		}
	}
	if len(effects) == 0 {
		return target, "refused"
	}
	resolved := false
	for _, effText := range effects {
		if *lookups <= 0 {
			break
		}
		effHash, err := hashref.Parse(effText)
		if err != nil {
			continue
		}
		*lookups--
		effObj, ok, err := reads.GetObject(ctx, effHash)
		if err != nil || !ok || effObj.SemanticID != broker.EffectRecordV1 {
			continue
		}
		rec, err := broker.DecodeRecord(effObj.Payload)
		if err != nil {
			continue
		}
		resolved = true
		if !rec.Allowed {
			return target, "denied"
		}
		if rec.Failed {
			return target, "failed"
		}
		if rec.ResultRef != (hashref.HashRef{}) && *lookups > 0 {
			*lookups--
			if resObj, ok, err := reads.GetObject(ctx, rec.ResultRef); err == nil && ok && len(resObj.Payload) <= 64<<10 {
				var out struct {
					OK *bool `json:"ok"`
				}
				if json.Unmarshal(resObj.Payload, &out) == nil && out.OK != nil && !*out.OK {
					return target, "refused"
				}
			}
		}
	}
	if resolved {
		return target, "ok"
	}
	return target, ""
}

// The report's journal ranges provide committed invocation receipts. Scan up to
// maxTrajectoryScan per namespace, verify exact episode and commit world, and
// project the tool identifier, committed entry, and sanitized target/status.
func readTrajectoryWorld(ctx context.Context, intents intentReader, reads objectWorldReader, episode string, through time.Time) trajectoryWorld {
	out := trajectoryWorld{Calls: []trajectoryCall{}}
	if intents == nil || reads == nil || !surface.EpisodeGrammar.MatchString(episode) {
		return out
	}
	lookups := maxTrajectoryLookups
	type pendingCall struct {
		call    trajectoryCall
		v2      bool
		input   string
		effects []string
	}
	var pending []pendingCall
	for _, via := range reportSurfaces {
		page, err := intents.InvocationIntents(ctx, string(via)+":"+episode+":", 0, maxTrajectoryScan+1)
		if err != nil {
			return trajectoryWorld{Calls: []trajectoryCall{}}
		}
		if len(page) > maxTrajectoryScan {
			page = page[:maxTrajectoryScan]
			out.Truncated = true
		}
		for _, it := range page {
			at := time.Unix(it.Intent.LogicalTime, 0).UTC()
			if it.Outcome == nil || it.Outcome.Status != "committed" || !through.IsZero() && at.After(through) {
				continue
			}
			obj, ok, err := reads.GetObject(ctx, it.Intent.TransitionRef)
			if err != nil {
				return trajectoryWorld{Calls: []trajectoryCall{}}
			}
			var rec struct {
				Episode    string   `json:"episodeId"`
				Invocation string   `json:"invocationId"`
				Tool       string   `json:"skillId"`
				Input      string   `json:"input"`
				Effects    []string `json:"effects"`
			}
			if !ok || (obj.SemanticID != coordinator.RecordV1 && obj.SemanticID != coordinator.RecordV2) ||
				json.Unmarshal(obj.Payload, &rec) != nil || rec.Episode != episode || rec.Invocation != it.Intent.InvocationID {
				continue
			}
			world, ok, err := reads.GetWorld(ctx, it.Outcome.ResultRef)
			if err != nil {
				return trajectoryWorld{Calls: []trajectoryCall{}}
			}
			if !ok {
				continue
			}
			tool := rec.Tool
			if !trajectoryToolName.MatchString(tool) || surface.SanitizeOpenText(tool, 96) != tool {
				tool = "Recorded tool"
			}
			pending = append(pending, pendingCall{
				call: trajectoryCall{
					Entry: world.Revision, At: at.Format(time.RFC3339),
					Tool: tool, Surface: string(via),
				},
				v2:      obj.SemanticID == coordinator.RecordV2,
				input:   rec.Input,
				effects: rec.Effects,
			})
		}
	}
	sort.SliceStable(pending, func(i, j int) bool { return pending[i].call.Entry < pending[j].call.Entry })
	out.Total = len(pending)
	if len(pending) > maxTrajectoryEvents {
		pending = pending[len(pending)-maxTrajectoryEvents:]
		out.Truncated = true
	}
	for _, p := range pending {
		if p.v2 {
			p.call.Target, p.call.Status = trajectoryCallMeta(ctx, reads, &lookups, p.call.Tool, p.input, p.effects)
		}
		out.Calls = append(out.Calls, p.call)
	}
	out.Available = true
	return out
}
