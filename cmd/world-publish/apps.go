package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/appregistry"
	"github.com/sunholo-data/ailang-world/host/store"
)

// runApps is the local, attended, single-writer write of the app registry (queue row
// 274 M1b-2, design_docs/planned/w-app-registry.md §3). It publishes one app.json as
// its app's next version, or retires one app. Like `transitions` it is NOT the network
// package publish: --live is refused, and there is no credential and no origin.
//
// Everything that decides WHAT will be written happens before the gate: open the store,
// read the head, read and check the manifest, build the next revision. So the typed
// phrase names the exact write (app, version or retire, store), and a mistyped path or
// an invalid manifest is reported before anyone is asked to type. After the gate the
// revision is published by compare-and-set against the head read here; a lost race
// writes nothing.
func runApps(opts options, in io.Reader, out, errw io.Writer, env environment) int {
	if serr := refuseLiveOnReadOnlyVerb("apps", opts.live); serr != nil {
		return report(errw, serr)
	}
	if serr := requireStorePath(opts.store); serr != nil {
		return report(errw, serr)
	}
	manifestPath, retire := strings.TrimSpace(opts.manifest), strings.TrimSpace(opts.retire)
	if (manifestPath == "") == (retire == "") {
		fmt.Fprintln(errw, "world-publish apps: pass exactly one of --manifest <app.json> (publish) or --retire <name>")
		return exitUsage
	}

	db, err := store.Open(opts.store)
	if err != nil {
		if store.IsWriterAlreadyActive(err) {
			for _, line := range writerLockHelp("apps", opts.store) {
				fmt.Fprintln(errw, line)
			}
		}
		return report(errw, storeOpenFailed(opts.store, err))
	}
	defer func() { _ = db.Close() }()

	// Two bounded contexts, never one across the gate: the operator may take any time to
	// type the phrase, so the publish gets its own budget, created after the gate (as
	// `transitions` does). The compare-and-set against snap.Head keeps the two reads honest.
	readCtx, cancelRead := context.WithTimeout(context.Background(), appsStepBudget)
	defer cancelRead()
	snap, err := appregistry.Read(readCtx, db)
	if err != nil {
		fmt.Fprintln(errw, "world-publish apps: "+err.Error())
		return exitError
	}
	var current *appregistry.Revision
	if snap.Present {
		current = &snap.Revision
	}

	var (
		next     appregistry.Revision
		manifest *store.Object
		prompt   string
		phrase   string
		done     string
	)
	if retire != "" {
		next, err = appregistry.BuildRetire(current, snap.Head, retire)
		if err != nil {
			fmt.Fprintf(errw, "world-publish apps: %v\n  next: world-publish apps --manifest <app.json> --store %s publishes one\n", err, opts.store)
			return exitError
		}
		prompt = fmt.Sprintf("Type exactly, to RETIRE app %s from the LOCAL app registry (its history stays in the parent chain):", retire)
		phrase = appsRetirePhrase(retire, opts.store)
		done = fmt.Sprintf("retired app %s", retire)
	} else {
		raw, rerr := os.ReadFile(manifestPath)
		if rerr != nil {
			fmt.Fprintf(errw, "world-publish apps: read %s: %v\n", manifestPath, rerr)
			return exitError
		}
		obj, m, merr := appregistry.ManifestObject(raw)
		if merr != nil {
			fmt.Fprintf(errw, "world-publish apps: %v\n  next: ailang-worldd app check %s\n", merr, manifestPath)
			return exitError
		}
		if base := strings.TrimSuffix(filepath.Base(manifestPath), ".json"); base != m.Name {
			fmt.Fprintf(errw, "world-publish apps: %s names app %q but the file is %s.json\n  next: rename the file to %s.json, or fix the manifest's name\n",
				manifestPath, m.Name, base, m.Name)
			return exitError
		}
		next, err = appregistry.BuildPublish(current, snap.Head, m.Name, obj.Hash, nil, time.Now())
		if err != nil {
			fmt.Fprintf(errw, "world-publish apps: %v\n  why: publishing the bytes the app already pins would mint a version that changes nothing\n", err)
			return exitError
		}
		entry, _ := next.Lookup(m.Name)
		manifest = &obj
		prompt = fmt.Sprintf("Type exactly, to publish app %s version %d (%s) to the LOCAL app registry:", m.Name, entry.Version, obj.Hash)
		phrase = appsPublishPhrase(m.Name, entry.Version, opts.store)
		done = fmt.Sprintf("published app %s version %d (manifest %s)", m.Name, entry.Version, obj.Hash)
	}

	// THE GATE: the fence stack `transitions` uses — refuse an automation environment,
	// require a controlling terminal, require the typed phrase for THIS write. A
	// headless loop cannot publish or retire an app.
	if serr := requireAttendedLocalWrite(in, out, env.getenv, env.probe(), prompt, phrase); serr != nil {
		return report(errw, serr)
	}

	writeCtx, cancelWrite := context.WithTimeout(context.Background(), appsStepBudget)
	defer cancelWrite()
	head, err := appregistry.Publish(writeCtx, db, snap.Head, next, manifest)
	if err != nil {
		if store.IsRegistryCASConflict(err) {
			fmt.Fprintln(errw, "world-publish apps: the app registry moved while you confirmed; nothing was written — re-run the command")
			return exitError
		}
		fmt.Fprintln(errw, "world-publish apps: "+err.Error())
		return exitError
	}
	fmt.Fprintf(out, "%s: app registry revision %d (head %s)\n", done, next.Revision, head)
	return exitOK
}

// appsStepBudget bounds each of runApps' two store phases (the pre-gate read and the
// post-gate publish). A var only so TestAppsVerbSlowConfirmation can shrink it.
var appsStepBudget = 30 * time.Second

// appsPublishPhrase and appsRetirePhrase name the app, the version or the retire, and
// the store, so a confirmation cannot be reused for another app, another version or
// another store, and is never the package-publish phrase.
func appsPublishPhrase(name string, version int64, storePath string) string {
	return fmt.Sprintf("publish app %s version %d to %s", name, version, storePath)
}

func appsRetirePhrase(name, storePath string) string {
	return fmt.Sprintf("retire app %s from %s", name, storePath)
}
