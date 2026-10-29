// Package surfacefollow implements the read-only M1 source boundary. It does
// not write a cache or install a daemon route. LoadConfig enrolls sources;
// NewCollector and Follower.Run consume a caller's context and deliver immutable
// projections and status to a prompt, caller-owned Sink.
//
// A configuration file (regular, owned by the current user, mode 0600 in an
// operator-owned directory without group/other write access) looks like:
//
//	{
//	  "schema": "world/surface-follow-config/v1",
//	  "repos": [
//	    {"id":"world", "repo":"~/dev/world", "charter":"docs/charter.md",
//	     "branch":"dev", "limit":200}
//	  ]
//	}
//
// LoadOptions.Home explicitly supplies tilde expansion; ForbiddenRoots must
// include every workspace root. Repository paths become canonical absolute
// directories. The default limit is 200, with a 2000 summed limit. Config errors
// name a field and repo, suppress untrusted values, and direct the operator to
// correct that field and run ailang-worldd serve --help.
//
// With cfg returned by LoadConfig, a caller can construct:
//
//	f := Follower{Config: cfg, Collector: NewCollector(nil, nil),
//	    Sink: func(s Status) { /* retain the accepted partition and status */ }}
//	err := f.Run(ctx)
//
// A cold collection issues rev-parse on refs/heads/BRANCH, then log and cat-file
// on its pinned full hash, then merge-base on immutable parents. An unchanged
// poll issues only rev-parse. All calls carry -C CANONICAL_REPO
// --no-optional-locks and pass through worktree.Output. ReadGit never receives a
// repo-qualified graph identifier. Real source bytes and refs are compared in
// TestRealGitLeavesSourceUnchanged.
//
// The scheduler starts immediately, polls at start-relative 15-second ticks,
// bounds builds to two workers and one latest pending job per repo, and joins
// all work on cancellation. Queue time counts toward 20 seconds; each git call
// gets two seconds. Sixteen of the rolling 1200 launches are reserved for probes.
// Failures retain accepted partitions and back off 15, 30, then 60 seconds;
// observed head changes reset backoff. Budget refusals retry on the next poll.
// Sinks must return promptly and must not mutate delivered indexes.
package surfacefollow
