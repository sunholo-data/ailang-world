package main

// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.10, Q2): approve --list
// also prints the pending Human.Approve requests an operator decides with a request ref. It reads
// the world store --db read-only, offline, through the daemon's own listing (daemon.PendingDecisions,
// the rows GET /surface/requests.json serves), so it never dials a daemon.

import (
	"context"
	"fmt"
	"io"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// approveListDecisions prints the decisions section of approve --list: nothing when there is none;
// a world store that cannot be read is one line on stderr and changes neither the mint list nor the
// exit code.
func approveListDecisions(ctx context.Context, db string, now int64, stdout, stderr io.Writer) {
	ro, err := store.OpenReadOnly(db)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd approve: decisions: cannot read the world store %s: %v\n", db, err)
		return
	}
	rows, truncated, err := daemon.PendingDecisions(ctx, ro, now)
	_ = ro.Close()
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd approve: decisions: cannot list the decisions in %s: %v\n", db, err)
		return
	}
	printDecisions(stdout, rows, truncated)
}

// printDecisions is the section's text: nothing for no row; one line per request, its refusal or
// its deadline and the command that decides it; the bound when the list was cut.
func printDecisions(stdout io.Writer, rows []surfacepage.PendingDecision, truncated bool) {
	if len(rows) == 0 && !truncated {
		return
	}
	fmt.Fprintln(stdout, "pending decisions (Human.Approve requests):")
	for _, d := range rows {
		line := fmt.Sprintf("%s  %s  requested by %s", d.Request, d.Scope, d.Requester)
		if d.Refusal != "" {
			line += "  refused: " + d.Refusal
		} else {
			line += fmt.Sprintf("  deadline %s  decide: ailang-worldd approve %s [--deny]", approverTime(d.Deadline), d.Request)
		}
		fmt.Fprintln(stdout, line)
	}
	if truncated {
		fmt.Fprintln(stdout, "(more requests than listed: the approvals chain or the packet scan reached its bound)")
	}
}
