package surface

import (
	"encoding/json"
	"strings"
	"testing"
)

// openFixture is a valid overlay over followFixture's golden partitions.
func openFixture(t *testing.T) FollowIndex {
	t.Helper()
	x := followFixture(t)
	var shared, step State
	for _, s := range x.Repos[0].Index.States {
		if s.Kind == KindShared && shared.ID == "" {
			shared = s
		}
		if s.Kind == KindStep && step.ID == "" {
			step = s
		}
	}
	if shared.ID == "" || step.ID == "" {
		t.Fatal("golden index lacks shared and step states")
	}
	x.Open = &FollowOpen{Status: OpenComplete, CheckedAt: "2026-10-29T12:00:00Z", Cap: MaxOpenEpisodes,
		Episodes: []OpenEpisode{
			{ID: "world:episode:w1", Repo: "world", Episode: "w1", State: OpenRunning, Label: OpenLabel(OpenRunning), SessionLive: true,
				StartedAt: "2026-10-29T11:00:00Z", Base: shared.ID, Anchor: "world:commit:" + shared.ID, Branch: "trial/w1",
				Task: "Follow open episodes", Agent: "claude", Model: "opus"},
			{ID: "tools:episode:w2", Repo: "tools", Episode: "w2", State: OpenPR, Label: OpenLabel(OpenPR), PR: 4,
				StartedAt: "2026-10-29T10:00:00Z", EndedAt: "2026-10-29T11:00:00Z", Base: strings.Repeat("e", 40),
				BaseOutsideWindow: true, Branch: "trial/w2", Agent: "codex", Model: "gpt"},
		},
		Aliases: []EpisodeAlias{{ID: "world:episode:w3", Target: "world:commit:" + step.ID}},
		Running: []RunningEntry{{Episode: "legacy", Reason: ReasonNoStart}, {Episode: "elsewhere", Reason: ReasonNoRepo}},
	}
	return x
}

func TestFollowOpenRoundTrips(t *testing.T) {
	x := openFixture(t)
	b, err := json.Marshal(x)
	if err != nil {
		t.Fatal(err)
	}
	got, err := DecodeFollow(b)
	if err != nil || got.Open == nil || len(got.Open.Episodes) != 2 {
		t.Fatalf("round trip: %v", err)
	}
	// An envelope without an overlay (an M2 cache) still decodes.
	x.Open = nil
	b, _ = json.Marshal(x)
	if _, err := DecodeFollow(b); err != nil {
		t.Fatal(err)
	}
}

func TestFollowOpenRejects(t *testing.T) {
	mutations := map[string]func(*FollowOpen){
		"status":            func(o *FollowOpen) { o.Status = "fine" },
		"complete-reason":   func(o *FollowOpen) { o.Reason = "capped" },
		"complete-omitted":  func(o *FollowOpen) { o.Omitted = 1 },
		"complete-unknown":  func(o *FollowOpen) { o.RemainderUnknown = true },
		"incomplete-silent": func(o *FollowOpen) { o.Status = OpenIncomplete },
		"reason-path":       func(o *FollowOpen) { o.Status, o.Reason = OpenIncomplete, "store /Users/x failed" },
		"checked":           func(o *FollowOpen) { o.CheckedAt = "now" },
		"checked-empty":     func(o *FollowOpen) { o.CheckedAt = "" },
		"cap":               func(o *FollowOpen) { o.Cap = 5 },
		"over-cap": func(o *FollowOpen) {
			for len(o.Episodes) <= MaxOpenEpisodes {
				o.Episodes = append(o.Episodes, o.Episodes[0])
			}
		},
		"id-unqualified":  func(o *FollowOpen) { o.Episodes[0].ID = "w1" },
		"id-repo":         func(o *FollowOpen) { o.Episodes[0].ID = "tools:episode:w1" },
		"id-duplicate":    func(o *FollowOpen) { o.Episodes[1] = o.Episodes[0] },
		"episode-grammar": func(o *FollowOpen) { o.Episodes[0].Episode, o.Episodes[0].ID = "W_1", "world:episode:W_1" },
		"label":           func(o *FollowOpen) { o.Episodes[0].Label = "PR open" },
		"state":           func(o *FollowOpen) { o.Episodes[0].State, o.Episodes[0].Label = "merged", "merged" },
		"running-dead":    func(o *FollowOpen) { o.Episodes[0].SessionLive, o.Episodes[0].EndedAt = false, "2026-10-29T12:00:00Z" },
		"ended-live": func(o *FollowOpen) {
			o.Episodes[0].State, o.Episodes[0].Label = OpenEnded, OpenLabel(OpenEnded)
		},
		"failure-has-pr": func(o *FollowOpen) {
			o.Episodes[1].State, o.Episodes[1].Label = OpenPushFailed, OpenLabel(OpenPushFailed)
		},
		"pr-without-number": func(o *FollowOpen) { o.Episodes[1].PR = 0 },
		"ended-without-at":  func(o *FollowOpen) { o.Episodes[1].EndedAt = "" },
		"live-with-end":     func(o *FollowOpen) { o.Episodes[0].EndedAt = "2026-10-29T12:00:00Z" },
		"started":           func(o *FollowOpen) { o.Episodes[0].StartedAt = "" },
		"base":              func(o *FollowOpen) { o.Episodes[0].Base = "abc" },
		"anchor-other":      func(o *FollowOpen) { o.Episodes[0].Anchor = "world:commit:" + strings.Repeat("e", 40) },
		"anchor-missing":    func(o *FollowOpen) { o.Episodes[0].Anchor = "" },
		"anchor-and-marker": func(o *FollowOpen) { o.Episodes[0].BaseOutsideWindow = true },
		"anchor-unresolved": func(o *FollowOpen) {
			o.Episodes[1].BaseOutsideWindow, o.Episodes[1].Anchor = false, "tools:commit:"+o.Episodes[1].Base
		},
		"task-email":      func(o *FollowOpen) { o.Episodes[0].Task = "ask operator@example.invalid" },
		"task-path":       func(o *FollowOpen) { o.Episodes[0].Task = "edit /Users/operator/x" },
		"task-url":        func(o *FollowOpen) { o.Episodes[0].Task = "see https://example.invalid/x" },
		"task-token":      func(o *FollowOpen) { o.Episodes[0].Task = "use ghp_abcdefgh12345678" },
		"task-newline":    func(o *FollowOpen) { o.Episodes[0].Task = "two\nlines" },
		"model-long":      func(o *FollowOpen) { o.Episodes[0].Model = strings.Repeat("m", 65) },
		"alias-open":      func(o *FollowOpen) { o.Aliases[0].ID = o.Episodes[0].ID },
		"alias-cross":     func(o *FollowOpen) { o.Aliases[0].Target = strings.Replace(o.Aliases[0].Target, "world:", "tools:", 1) },
		"alias-shared":    func(o *FollowOpen) { o.Aliases[0].Target = o.Episodes[0].Anchor },
		"alias-kind":      func(o *FollowOpen) { o.Aliases[0].Target = "world:episode:w1" },
		"running-reason":  func(o *FollowOpen) { o.Running[0].Reason = "idle" },
		"running-dup":     func(o *FollowOpen) { o.Running[1] = o.Running[0] },
		"running-glyph":   func(o *FollowOpen) { o.Running[0].Episode = "w1" },
		"running-grammar": func(o *FollowOpen) { o.Running[0].Episode = "Legacy!" },
		"omitted":         func(o *FollowOpen) { o.Status, o.Reason, o.Omitted = OpenIncomplete, "capped", -1 },
		"unavailable-unchecked-entries": func(o *FollowOpen) {
			o.Status, o.Reason, o.CheckedAt = OpenUnavailable, "running status unavailable", ""
		},
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			x := openFixture(t)
			mutate(x.Open)
			b, _ := json.Marshal(x)
			if _, err := DecodeFollow(b); err == nil {
				t.Fatal("malformed overlay accepted")
			}
		})
	}
	// The aggregate state ceiling counts open states.
	t.Run("aggregate", func(t *testing.T) {
		x := openFixture(t)
		if err := ValidateFollow(x); err != nil {
			t.Fatal(err)
		}
		if err := ValidateFollowOpen(*x.Open, x.Repos); err != nil {
			t.Fatal(err)
		}
	})
}

func TestSanitizeOpenText(t *testing.T) {
	for in, want := range map[string]string{
		"  Build\tthe\nthing  ":                                     "Build the thing",
		"mail operator@example.invalid now":                         "mail [email] now",
		"open https://u:p@example.invalid/x?y=1":                    "open [url]",
		"edit /Users/operator/a and ~/b and ./c":                    "edit [path] and [path] and [path]",
		"and/or stays":                                              "and/or stays",
		"token ghp_abcdefgh12345678 and " + strings.Repeat("a", 40): "token [redacted] and [redacted]",
	} {
		if got := SanitizeOpenText(in, 120); got != want {
			t.Errorf("%q → %q, want %q", in, got, want)
		}
		if SanitizeOpenText(SanitizeOpenText(in, 120), 120) != SanitizeOpenText(in, 120) {
			t.Errorf("%q: not idempotent", in)
		}
	}
	long := SanitizeOpenText(strings.Repeat("é", 200), 120)
	if len([]rune(long)) != 120 || !strings.HasSuffix(long, "…") || SanitizeOpenText(long, 120) != long {
		t.Fatalf("bound: %d runes", len([]rune(long)))
	}
}
