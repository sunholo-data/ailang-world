package surface

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// The open-episode overlay (row 220 M3) is validated separately from the v1
// merge partitions it sits beside. "Running" means a session credential is
// live; it never claims an agent process is alive.

// MaxOpenEpisodes caps the open states in one envelope. With MaxStates it
// leaves 4000 merged states and 1000 open ones.
const MaxOpenEpisodes = 1000

// Overlay completeness. Only OpenComplete may be read as "this is all of it".
const (
	OpenComplete    = "complete"
	OpenIncomplete  = "incomplete"
	OpenUnavailable = "unavailable"
)

// Open state kinds. Finish outcomes other than done keep their own state and
// never become OpenPR.
const (
	OpenRunning     = "running"
	OpenEnded       = "ended"
	OpenPR          = "pr-open"
	OpenPushFailed  = "push-failed"
	OpenPRFailed    = "pr-failed"
	OpenAmbiguous   = "ambiguous"
	ReasonNoStart   = "start record unavailable"
	ReasonNoRepo    = "repo not configured"
	openTextBound   = 120
	openReasonBound = 200
)

var openLabels = map[string]string{
	OpenRunning:    "session live",
	OpenEnded:      "session ended; outcome unknown",
	OpenPR:         "PR open",
	OpenPushFailed: "finish: push failed; no PR is implied",
	OpenPRFailed:   "finish: PR creation failed; no PR is implied",
	OpenAmbiguous:  "finish: outcome ambiguous; no PR is implied",
}

// OpenLabel is the fixed browser label of a state kind.
func OpenLabel(state string) string { return openLabels[state] }

// EpisodeGrammar is the daemon's episode id grammar.
var EpisodeGrammar = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)

// FollowOpen is the open-episode collection of one publication.
type FollowOpen struct {
	Status string `json:"status"`
	// Reason is a fixed, path-free sentence; empty only when complete.
	Reason string `json:"reason,omitempty"`
	// CheckedAt is the last successful read; empty if none has succeeded.
	CheckedAt string `json:"checked_at"`
	Cap       int    `json:"cap"`
	// Omitted counts open states dropped by the cap. RemainderUnknown says a
	// read bound was reached, so more may exist than any count shows.
	Omitted          int            `json:"omitted"`
	RemainderUnknown bool           `json:"remainder_unknown"`
	Unconfigured     int            `json:"unconfigured"`
	Episodes         []OpenEpisode  `json:"episodes"`
	Aliases          []EpisodeAlias `json:"aliases"`
	Running          []RunningEntry `json:"running"`
}

// OpenEpisode holds only whitelisted start facts. It never carries preset
// bytes, paths, remote URLs, checks, tokens or e-mail addresses.
type OpenEpisode struct {
	ID          string `json:"id"` // repo:episode:<ep>
	Repo        string `json:"repo"`
	Episode     string `json:"episode"`
	State       string `json:"state"`
	Label       string `json:"label"`
	SessionLive bool   `json:"session_live"`
	StartedAt   string `json:"started_at"`
	EndedAt     string `json:"ended_at,omitempty"`
	Base        string `json:"base"`
	// Anchor is the qualified commit id of Base when Base is a state in the
	// repo's retained window; otherwise BaseOutsideWindow is true.
	Anchor            string `json:"anchor,omitempty"`
	BaseOutsideWindow bool   `json:"base_outside_window"`
	Branch            string `json:"branch"`
	Task              string `json:"task,omitempty"`
	Agent             string `json:"agent"`
	Model             string `json:"model"`
	PR                int    `json:"pr,omitempty"`
}

// EpisodeAlias keeps an episode address on its landed step after a merge.
type EpisodeAlias struct {
	ID          string `json:"id"`     // repo:episode:<ep>
	Target      string `json:"target"` // repo:commit:<step>
	SessionLive bool   `json:"session_live"`
}

// RunningEntry is a live session without an open state, and why.
type RunningEntry struct {
	Episode string `json:"episode"`
	Reason  string `json:"reason"`
}

var (
	openEmail = regexp.MustCompile(`[^\s@]+@[^\s@]+\.[^\s@]+`)
	openURL   = regexp.MustCompile(`(?i)\b[a-z][a-z0-9+.-]*://\S*`)
	openPath  = regexp.MustCompile(`(^|[\s"'(=])(~|\.{1,2})?/\S+`)
	openToken = regexp.MustCompile(`\b[0-9a-fA-F]{32,}\b|\b(?:ghp|gho|ghs|github_pat|sk|xox[a-z])[-_][A-Za-z0-9_-]{8,}`)
)

// SanitizeOpenText reduces untrusted start text to one short line with
// e-mail addresses, URLs, paths and token-shaped strings redacted.
func SanitizeOpenText(s string, bound int) string {
	s = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.IsSpace(r) {
			return ' '
		}
		if !unicode.IsPrint(r) {
			return -1
		}
		return r
	}, s)
	s = strings.Join(strings.Fields(s), " ")
	s = openURL.ReplaceAllString(s, "[url]")
	s = openEmail.ReplaceAllString(s, "[email]")
	s = openPath.ReplaceAllString(s, "${1}[path]")
	s = openToken.ReplaceAllString(s, "[redacted]")
	if r := []rune(s); len(r) > bound {
		s = string(r[:bound-1]) + "…"
	}
	return s
}

func safeOpenText(s string, bound int, empty bool) bool {
	if s == "" {
		return empty
	}
	return SanitizeOpenText(s, bound) == s
}

var openReasons = regexp.MustCompile(`^[A-Za-z0-9 ,;:.()'-]+$`)

// ValidateFollowOpen checks an overlay against the partitions it is
// published with. Each refusal names the field and the entry.
func ValidateFollowOpen(o FollowOpen, parts []FollowPartition) error {
	fail := func(f string, a ...any) error { return fmt.Errorf("surface follow open: "+f, a...) }
	switch o.Status {
	case OpenComplete:
		if o.Reason != "" || o.Omitted != 0 || o.RemainderUnknown {
			return fail("complete overlay carries a reason, omission or unknown remainder")
		}
	case OpenIncomplete, OpenUnavailable:
		if o.Reason == "" {
			return fail("%s overlay must give a reason", o.Status)
		}
	default:
		return fail("status %q is not complete, incomplete or unavailable", o.Status)
	}
	if len(o.Reason) > openReasonBound || (o.Reason != "" && !openReasons.MatchString(o.Reason)) {
		return fail("reason must be a short, path-free sentence")
	}
	if o.CheckedAt != "" {
		if _, err := time.Parse(time.RFC3339, o.CheckedAt); err != nil {
			return fail("checked_at must be RFC 3339")
		}
	} else if o.Status != OpenUnavailable || len(o.Episodes)+len(o.Aliases)+len(o.Running) > 0 {
		return fail("checked_at is empty without an empty unavailable overlay")
	}
	if o.Cap != MaxOpenEpisodes || len(o.Episodes) > o.Cap || len(o.Running) > o.Cap || len(o.Aliases) > MaxStates {
		return fail("cap must be %d and bound episodes and running entries", MaxOpenEpisodes)
	}
	if o.Omitted < 0 || o.Unconfigured < 0 || (o.Omitted > 0 && o.Status == OpenComplete) {
		return fail("omitted and unconfigured must be non-negative, and omissions are incomplete")
	}
	byID := map[string]map[string]State{}
	for _, p := range parts {
		states := map[string]State{}
		for _, s := range p.Index.States {
			states[s.ID] = s
		}
		byID[p.ID] = states
	}
	seen := map[string]bool{}
	for i, e := range o.Episodes {
		id, err := QualifyID(e.Repo, "episode", e.Episode)
		if err != nil || id != e.ID || !EpisodeGrammar.MatchString(e.Episode) || seen[e.ID] {
			return fail("episodes[%d].id must be a unique repo:episode:<ep>", i)
		}
		seen[e.ID] = true
		label, ok := openLabels[e.State]
		if !ok || label != e.Label {
			return fail("episodes[%d] %s: state or label is not a fixed pair", i, e.ID)
		}
		if (e.State == OpenRunning && !e.SessionLive) || (e.State == OpenEnded && e.SessionLive) {
			return fail("episodes[%d] %s: running needs a live session and ended needs none", i, e.ID)
		}
		if (e.PR > 0) != (e.State == OpenPR) || e.PR < 0 {
			return fail("episodes[%d] %s: only pr-open carries a PR number", i, e.ID)
		}
		if _, err := time.Parse(time.RFC3339, e.StartedAt); err != nil {
			return fail("episodes[%d] %s: started_at must be RFC 3339", i, e.ID)
		}
		if e.EndedAt != "" {
			if _, err := time.Parse(time.RFC3339, e.EndedAt); err != nil || e.SessionLive {
				return fail("episodes[%d] %s: ended_at must be RFC 3339 and only without a live session", i, e.ID)
			}
		} else if !e.SessionLive {
			return fail("episodes[%d] %s: an ended session needs ended_at", i, e.ID)
		}
		if e.Base != "" && !commitHash.MatchString(e.Base) {
			return fail("episodes[%d] %s: base must be a full commit hash", i, e.ID)
		}
		if e.Anchor == "" {
			if !e.BaseOutsideWindow {
				return fail("episodes[%d] %s: an unanchored base must be marked outside the window", i, e.ID)
			}
		} else if e.BaseOutsideWindow || e.Base == "" || e.Anchor != e.Repo+":commit:"+e.Base {
			return fail("episodes[%d] %s: anchor must qualify the recorded base", i, e.ID)
		} else if _, ok := byID[e.Repo][e.Base]; !ok {
			return fail("episodes[%d] %s: anchor does not resolve in the repo's window", i, e.ID)
		}
		if !safeOpenText(e.Task, openTextBound, true) || !safeOpenText(e.Agent, 64, true) ||
			!safeOpenText(e.Model, 64, true) || !safeOpenText(e.Branch, openTextBound, true) {
			return fail("episodes[%d] %s: text fields must be sanitized", i, e.ID)
		}
	}
	for i, a := range o.Aliases {
		repo, kind, local, err := ParseQualifiedID(a.ID)
		if err != nil || kind != "episode" || !EpisodeGrammar.MatchString(local) || seen[a.ID] {
			return fail("aliases[%d].id must be a unique episode id without an open state", i)
		}
		seen[a.ID] = true
		trepo, tkind, tlocal, err := ParseQualifiedID(a.Target)
		if err != nil || trepo != repo || tkind != "commit" {
			return fail("aliases[%d] %s: target must be a commit in the same repo", i, a.ID)
		}
		if s, ok := byID[repo][tlocal]; !ok || s.Kind != KindStep {
			return fail("aliases[%d] %s: target is not a landed step in the window", i, a.ID)
		}
	}
	running := map[string]bool{}
	for i, r := range o.Running {
		if !EpisodeGrammar.MatchString(r.Episode) || running[r.Episode] || (r.Reason != ReasonNoStart && r.Reason != ReasonNoRepo) {
			return fail("running[%d] must be a unique episode with a fixed reason", i)
		}
		running[r.Episode] = true
	}
	for _, e := range o.Episodes {
		if running[e.Episode] {
			return fail("running entry %s also has an open state", e.Episode)
		}
	}
	return nil
}
