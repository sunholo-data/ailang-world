package main

// Row 243 V2a (design_docs/planned/w-world-address-bar-verbs-v2.md §2.1): why's CLI front
// (help, flags, the --result and - reads, --expect-store) over host/why's walk, then the
// names the other CLI files use, kept so that none of them is edited.

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/why"
)

const (
	whyBudget      = 60 * time.Second
	whyScanDefault = why.ScanDefault
	whyScanMax     = why.ScanMax
	logPageMax     = why.LogPageMax
	frozenRigPort  = 47_330 // historical trailer recovery, never a default
)

const whyHelp = `usage: ailang-worldd [--addr <url>] why <target> [--scan N] [--json] [--expect-store <path>]
       ailang-worldd [--addr <url>] why --result <file> [--scan N] [--json]

Walks one committed invocation back to its whole provenance chain, checking
every link (✓ or ✗):

  world    the world ref, recomputed from the entry and the output, and
           confirmed by GET /v1/worlds/<ref>
  entry    index, writtenBy, prev, transitionFn and the interpreter (with its
           pin name); the entry hash is recomputed
  record   invocation id, episode, skill (its hash and its fields agree with
           the entry)
  input    your arguments, decoded and elided; when elided, the short
           top-level fields (path, command, op, …) are printed first
  plan     the effects planned, or the refusal result
  effect   per effect record: allowed/failed/denial, budget before→after,
           request and result sizes; a toolchain call that changed the
           worktree adds "worktree: N created, M modified, K deleted" and
           its first paths (--json: the whole worktree_delta)
  output   the committed output; it must equal world.stateRoot

Targets:
  <index> | head               a log entry
  sha256:<ref>                 any object (dispatched on its semanticId: a
                               record, input, output, plan, effect record,
                               effect request or result), or a world ref
  a2a:<id> | mcp:<id> | rest:<id>
                               an invocation id
  - | --result <file>          the output of 'call --json-out': its sha256 is
                               the output ref (the world.plan is the fallback)

A target that is not an index is found by scanning the log backwards from
head: --scan N entries (default 500, max 5000) within a 60 s budget; past
either it prints "not found in the last N entries" and exits 1.

Exit: 0 every link ✓; 3 any link ✗ (an integrity refusal); 1 not found or
usage. --json prints the chain as JSON.
`

// parseWhyTarget makes the CLI's two reads (--result's file, -'s stdin) and
// hands every other target to why.ParseTarget.
func parseWhyTarget(arg, resultFile string, stdin io.Reader) (why.Target, error) {
	if resultFile != "" {
		if arg != "" {
			return why.Target{}, errors.New("give either a target or --result, not both")
		}
		b, err := os.ReadFile(resultFile)
		if err != nil {
			return why.Target{}, fmt.Errorf("read --result %s: %v", resultFile, errors.Unwrap(err))
		}
		return why.Target{Kind: "result", Result: trimOneNewline(b)}, nil
	}
	if arg == "-" {
		b, err := io.ReadAll(io.LimitReader(stdin, maxClientResponseBytes+1))
		if err != nil {
			return why.Target{}, fmt.Errorf("read stdin: %v", err)
		}
		return why.Target{Kind: "result", Result: trimOneNewline(b)}, nil
	}
	return why.ParseTarget(arg)
}

// trimOneNewline drops the one newline `call --json-out` appends.
func trimOneNewline(b []byte) []byte {
	if n := len(b); n > 0 && b[n-1] == '\n' {
		return b[:n-1]
	}
	return b
}

func runWhy(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	target, args := verbs.WhyLeadingTarget(args)
	fs := flag.NewFlagSet("ailang-worldd why", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	a := verbs.WhyFlags(fs)
	if code, done := parseVerbFlags(fs, args, whyHelp, stdout, stderr); done {
		return code
	}
	target, err := verbs.WhyTrailingTarget(fs, target)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
		return exitUsage
	}
	if !verbs.WhyScanOK(*a.Scan) {
		fmt.Fprintf(stderr, "ailang-worldd why: --scan must be 1..%d\n", whyScanMax)
		return exitUsage
	}
	expectStore, resultFile, scan, asJSON := a.ExpectStore, a.Result, a.Scan, a.JSON
	t, err := parseWhyTarget(target, *resultFile, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd why: %v\n%s", err, whyHelp)
		return exitUsage
	}
	ctx, cancel := budgetContext(whyBudget)
	defer cancel()
	r := worldReader{ctx: ctx, c: newClient(addr)}
	if *expectStore != "" {
		h, err := r.health()
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
			return exitUsage
		}
		r.healthCache = &h
		if !sameStorePath(h.DBPath, *expectStore) {
			return serviceMessage(stderr,
				fmt.Sprintf("ailang-worldd why: the daemon at %s serves store %s, not %s", addr, h.DBPath, *expectStore),
				"entry indexes are per store, so the same index in another store is another entry (a World-Provenance trailer names its store as store=<path>)",
				fmt.Sprintf("pass the --addr of the daemon serving %s (the frozen rig: start it from daemon-cmdline.txt on 127.0.0.1:%d), or drop --expect-store", *expectStore, frozenRigPort), exitUsage)
		}
	}
	ans := why.Run(r.why(), t, why.Options{Scan: *scan, JSON: *asJSON})
	fmt.Fprint(stdout, ans.Stdout)
	fmt.Fprint(stderr, ans.Stderr)
	return ans.Exit
}

// The compatibility layer: the names other CLI files use, on host/why.

type (
	wireHeader = why.Header
	wireEntry  = why.Entry
	wireObject = why.Object
	wireWorld  = why.World
	invRecord  = why.Record
	healthInfo = why.Health
	whyChain   = why.Chain
	whyLink    = why.Link
)

// entryWire and worldWire are defined (not aliased) because demo_seed.go and
// demo_test.go build them positionally, which vet refuses for another
// package's struct; the encoding stays host/why's.
type (
	entryWire why.EntryWire
	worldWire why.WorldWire
)

func sumJSON(v any) string { return why.SumJSON(v) }
func sum(b []byte) string  { return why.Sum(b) }

func contains(xs []string, x string) bool {
	for _, v := range xs {
		if v == x {
			return true
		}
	}
	return false
}

// worldReader is the GET side of the daemon, under one command budget.
type worldReader struct {
	ctx         context.Context
	c           *client
	healthCache *healthInfo
}

// why is r as host/why's Reader, reading through the CLI's HTTP client.
func (r worldReader) why() why.Reader {
	return why.Reader{Ctx: r.ctx, Get: clientGetter(r.c), HealthCache: r.healthCache}
}

// clientGetter is c's GET as a why.Getter (c.get, cli.go, is another shape).
func clientGetter(c *client) why.Getter {
	return func(ctx context.Context, path string) (int, []byte, error) {
		return c.do(ctx, http.MethodGet, path, nil)
	}
}

type whyTarget struct { // a target as demo_seed builds it (kind index)
	kind   string // index, head, ref, invocation, result
	index  int64
	text   string
	result []byte
}

// The reader's methods other CLI files call, each host/why's.
func (r worldReader) getJSON(path string, v any) (bool, error) { return r.why().GetJSON(path, v) }
func (r worldReader) entry(i int64) (wireEntry, bool, error)   { return r.why().Entry(i) }
func (r worldReader) page(from int64, n int) ([]wireEntry, error) {
	return r.why().Page(from, n)
}
func (r worldReader) object(ref string) (wireObject, bool, error) { return r.why().Object(ref) }
func (r worldReader) world(ref string) (wireWorld, bool, error)   { return r.why().World(ref) }
func (r worldReader) headIndex() (int64, bool, error)             { return r.why().HeadIndex() }
func (r worldReader) health() (healthInfo, error)                 { return r.why().Health() }
func (r worldReader) record(e wireEntry) (*invRecord, wireObject, error) {
	return r.why().Record(e)
}
func (r worldReader) scanEntries(from, to int64, reverse bool, max int,
	visit func(e wireEntry, rec *invRecord) (bool, error)) (int, error) {
	return r.why().ScanEntries(from, to, reverse, max, visit)
}
func (r worldReader) locate(t whyTarget, scan int) (wireEntry, string, error) {
	return r.why().Locate(why.Target{Kind: t.kind, Index: t.index, Text: t.text, Result: t.result}, scan)
}
func (r worldReader) effectRecord(ref string) (broker.EffectRecord, bool) {
	return r.why().EffectRecord(ref)
}
func (r worldReader) checkedObject(ref, semantic string) (wireObject, bool, string) {
	return r.why().CheckedObject(ref, semantic)
}
func (r worldReader) buildChain(e wireEntry, result []byte) (whyChain, error) {
	return r.why().BuildChain(e, result)
}
