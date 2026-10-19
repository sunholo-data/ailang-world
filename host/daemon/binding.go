package daemon

// Row 189 M2 (w-world-onboarding D-189-2 A): which exec profile an episode's
// workspace-exec runs under is authority, so the human confirms it at the
// mint: `session new --profile` / `episode start --profile` record the
// project and the sha256 of the exact profile bytes on the credential. The
// daemon reads that binding at request time, from the credentials store it
// already reads per request, so a newly minted episode needs no restart.
//
// The coordinator opens an episode's registry through BinderFor(episode,
// grants), which names the episode, not the credential. So the binding read
// here is the episode's: every LIVE credential of the episode that carries a
// binding. They must agree — two live credentials of one episode bound to
// different profiles (or to two digests of one) are a refusal naming both,
// never a guess. Unbound credentials (all of M1's, every mint without
// --profile) add nothing, so a bound and an unbound credential of one episode
// resolve to the binding.

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

// execBinding is an episode's minted profile binding; the zero value is
// unbound.
type execBinding struct {
	Project string
	Digest  string
}

func (b execBinding) bound() bool { return b.Project != "" }

// bindingConflictError refuses an episode whose live credentials disagree on
// the binding. It is a refusal (the operator revokes one), not a failure.
type bindingConflictError struct {
	episode  string
	bindings []string
}

func (e *bindingConflictError) Error() string {
	return fmt.Sprintf("episode %q has live session credentials bound to different exec profiles (%s); "+
		"revoke the ones you do not want (ailang-worldd session list --episode %s, then session revoke <credential_id>)",
		e.episode, strings.Join(e.bindings, ", "), e.episode)
}

// bindingOf is the episode's binding from its credential rows at now: the
// one binding every live bound row carries, unbound when none carries one.
func bindingOf(episodeID string, rows []store.SessionRow, now int64) (execBinding, error) {
	seen := map[execBinding]bool{}
	for _, r := range rows {
		if r.EpisodeID != episodeID || now > r.ExpiresAt || !r.Bound() {
			continue
		}
		seen[execBinding{Project: r.Project, Digest: r.ProfileDigest}] = true
	}
	switch len(seen) {
	case 0:
		return execBinding{}, nil
	case 1:
		for b := range seen {
			return b, nil
		}
	}
	names := make([]string, 0, len(seen))
	for b := range seen {
		names = append(names, b.Project+" "+b.Digest)
	}
	sort.Strings(names)
	return execBinding{}, &bindingConflictError{episode: episodeID, bindings: names}
}

// episodeBinding reads the episode's credentials through a fresh read-only
// handle, bounded by the credential budget (B4), and resolves its binding.
// The root context is here because BinderFor carries no ctx (as for
// episodeHandler); the read is one indexed scan of one episode's rows.
func (d *Daemon) episodeBinding(episodeID string) (execBinding, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.credentialBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return execBinding{}, fmt.Errorf("read the exec profile binding of episode %q: %w", episodeID, err)
	}
	defer func() { _ = ro.Close() }()
	rows, err := ro.ListSessions(ctx, episodeID, store.MaxListSessions)
	if err != nil {
		return execBinding{}, fmt.Errorf("read the exec profile binding of episode %q: %w", episodeID, err)
	}
	return bindingOf(episodeID, rows, time.Now().Unix())
}
