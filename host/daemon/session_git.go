package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/worktree"
)

// SessionGit is an additive, credential-free projection of local repository
// facts. An error replaces the entire object, so partial facts never escape.
type SessionGit struct {
	Branch      string `json:"branch"`
	Head        string `json:"head"`
	Base        string `json:"base"`
	BaseSubject string `json:"base_subject"`
	BaseDate    string `json:"base_date"`
	Ref         string `json:"ref"`
	RefHead     string `json:"ref_head"`
	Behind      int    `json:"behind"`
	Ahead       int    `json:"ahead"`
	Dirty       bool   `json:"dirty"`
	Error       string `json:"error,omitempty"`
}

// MarshalJSON keeps zero counts and false dirty on success, but only error on failure.
func (g SessionGit) MarshalJSON() ([]byte, error) {
	if g.Error != "" {
		return json.Marshal(map[string]string{"error": g.Error})
	}
	type facts SessionGit
	return json.Marshal(facts(g))
}

func readSessionGit(parent context.Context, dir, git string) *SessionGit {
	if !worktree.IsLinked(dir) {
		return nil
	}
	// Reserve 100 ms of the two-second phase for Output's bounded pipe drain
	// (runGit caps WaitDelay at one twentieth of this deadline).
	ctx, cancel := context.WithTimeout(parent, 1900*time.Millisecond)
	defer cancel()
	g := &SessionGit{Ref: "dev"}
	run := func(args ...string) (string, error) {
		out, err := worktree.Output(ctx, git, append([]string{"-C", dir, "-c", "core.fsmonitor=false"}, args...))
		return strings.TrimSpace(out), err
	}
	fail := func(err error) *SessionGit {
		msg := "local git metadata unavailable"
		if ctx.Err() != nil {
			msg = "local git metadata timed out (2s budget)"
		}
		// Do not echo child stderr, paths or environment into the context.
		_ = err
		return &SessionGit{Error: msg}
	}
	for _, step := range []struct {
		dst  *string
		args []string
	}{
		{&g.Branch, []string{"rev-parse", "--abbrev-ref", "HEAD"}},
		{&g.Head, []string{"rev-parse", "HEAD"}},
		{&g.Base, []string{"merge-base", "HEAD", "dev"}},
		{&g.RefHead, []string{"rev-parse", "dev"}},
	} {
		v, err := run(step.args...)
		if err != nil {
			return fail(err)
		}
		*step.dst = v
	}
	meta, err := run("show", "-s", "--format=%cI%n%s", g.Base)
	if err != nil {
		return fail(err)
	}
	parts := strings.SplitN(meta, "\n", 2)
	if len(parts) != 2 {
		return fail(fmt.Errorf("missing base metadata"))
	}
	g.BaseDate, g.BaseSubject = parts[0], parts[1]
	for _, step := range []struct {
		dst *int
		rev string
	}{
		{&g.Behind, g.Base + "..dev"}, {&g.Ahead, g.Base + "..HEAD"},
	} {
		v, err := run("rev-list", "--count", step.rev)
		if err != nil {
			return fail(err)
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			return fail(err)
		}
		*step.dst = n
	}
	status, err := run("status", "--porcelain")
	if err != nil {
		return fail(err)
	}
	g.Dirty = status != ""
	return g
}

func (g *SessionGit) BaseLine() string {
	if g == nil {
		return ""
	}
	if g.Error != "" {
		return "Base unknown: " + g.Error
	}
	base := g.Base
	if len(base) > 7 {
		base = base[:7]
	}
	line := fmt.Sprintf("Base %s (%s, %q) · dev is %d commits ahead of it", base, g.BaseDate, g.BaseSubject, g.Behind)
	if g.Dirty {
		line += " · worktree has uncommitted changes"
	}
	return line
}
