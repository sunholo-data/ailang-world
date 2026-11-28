package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/checkrun"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

var errOperatorFinishEligibility = errors.New("finish needs an eligible episode with a stopped launch and its original live session")

type operatorFinishBusy string

func (id operatorFinishBusy) Error() string {
	return fmt.Sprintf("active action %s holds the episode turn; run ailang-worldd service runner status", string(id))
}

func (operatorFinishBusy) Unwrap() error { return store.ErrOperatorConflict }

func (e serviceEnv) finishInputs(ctx context.Context, c runnerConfig, src operatorStartSource, ro *store.Credentials,
	ep string, lock *os.File) (operatorFinishInputs, error) {
	var out operatorFinishInputs
	if !verbs.EpisodeIDPattern.MatchString(ep) {
		return out, errOperatorFinishEligibility
	}
	start, err := ro.EpisodeStart(ctx, ep)
	if err != nil || start == nil || len(start.Resumes) != 0 || !slices.Contains(c.PresetNames, start.Preset.Name) ||
		!slices.Contains(c.FinishRemotes, start.Remote.URL) {
		return out, errOperatorFinishEligibility
	}
	launch, err := ro.OperatorFinishLaunch(ctx, ep)
	if err != nil || launch.Key.Version != 4 {
		return out, errOperatorFinishEligibility
	}
	if birth, live := operatorProcessIdentity(launch.Receipt.PID); live && birth == launch.Receipt.Birth {
		return out, errOperatorFinishEligibility
	}
	health, err := e.health(c.Address)
	if err != nil || health.Status != "ok" || slices.Contains(health.DegradedEpisodes, ep) ||
		health.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
		return out, errOperatorFinishEligibility
	}
	_, rep, err := fetchReport(ctx, finishClient(c.Address), ep)
	if err != nil || rep.Episode != ep || rep.Uncommitted != 0 || rep.Finish != nil || !reflect.DeepEqual(rep.Start, start) {
		return out, errOperatorFinishEligibility
	}
	presetDir, err := filepath.EvalSymlinks(src.presets)
	if err != nil {
		return out, errOperatorFinishReview
	}
	raw, _, err := operatorFinishRead(presetDir, start.Preset.Name+".json", 128<<10)
	if err != nil || serviceHash(raw) != start.Preset.SHA256 || string(raw) != start.Preset.Bytes {
		return out, errOperatorFinishReview
	}
	var preset launchPreset
	if strictPresetJSON(raw, &preset) != nil || preset.Name != start.Preset.Name {
		return out, errOperatorFinishReview
	}
	p, err := finishPreset(start.Preset.Name, src.presets)
	physical, physicalErr := filepath.EvalSymlinks(p.WorkspaceRoot)
	root, rootErr := filepath.EvalSymlinks(src.workspace)
	if err != nil || physicalErr != nil || rootErr != nil || physical != root || p.Addr != c.Address {
		return out, errOperatorFinishReview
	}
	if equal, _ := store.SameDatabase(p.DB, c.World); !equal {
		return out, errOperatorFinishReview
	}
	creds, err := credentialsFor(p.DB, p.CredentialsDB)
	if equal, _ := store.SameDatabase(creds, c.Credentials); err != nil || !equal {
		return out, errOperatorFinishReview
	}
	profile, _, err := loadMintProfile(src.profiles, preset.Profile)
	if err != nil {
		return out, errOperatorFinishReview
	}
	for _, check := range start.Checks {
		if check.Unsandboxed {
			return out, errOperatorFinishTerminal
		}
	}
	checks, err := validateGraderChecks(preset.Checks, preset.Name)
	if err != nil || pinPresetChecks(checks, profile) != nil || !reflect.DeepEqual(checks, start.Checks) {
		return out, errOperatorFinishReview
	}
	sessions, err := ro.ListSessions(ctx, ep, 2)
	if err != nil || len(sessions) != 1 || sessions[0].ExpiresAt <= src.clock() ||
		sessions[0].Project != preset.Profile {
		return out, errOperatorFinishEligibility
	}
	if sessions[0].ProfileDigest != profile.Digest {
		return out, errOperatorFinishReview
	}
	sessionDir, err := filepath.EvalSymlinks(p.SessionsDir)
	if err != nil || runnerPath(c.Home, sessionDir, c.UID, true, 0700) != nil {
		return out, errOperatorFinishEligibility
	}
	token, mode, err := operatorFinishRead(sessionDir, ep+".session", maxSessionFileBytes)
	if err != nil || mode.Perm() != 0600 {
		return out, errOperatorFinishEligibility
	}
	out.token = strings.TrimSpace(string(token))
	resolved, err := authority.New(ro).ResolveContext(ctx, "Bearer "+out.token, src.clock())
	if err != nil || resolved.Success == nil || resolved.Success.EpisodeID != ep || resolved.Success.ExpiresAt != sessions[0].ExpiresAt {
		return out, errOperatorFinishEligibility
	}

	auth, err := planOperatorFinishAuth(c, start.Remote.URL)
	if err != nil {
		return out, err
	}
	g, err := openOperatorFinishGit(auth, filepath.Join(root, ep), start.Repo, lock)
	if err != nil || g.dir != filepath.Join(root, ep) {
		return out, errOperatorFinishFence
	}
	name, email, err := worktree.ParseCommitter(preset.Committer)
	if err != nil || name == "" {
		return out, errOperatorFinishFence
	}
	title := "World episode " + ep
	if start.Task != nil && start.Task.Title != "" {
		title = start.Task.Title
	}
	rep.ExecCache = nil
	source, _ := json.Marshal(rep)
	sessionWire, _ := json.Marshal(sessions)
	launchWire, _ := json.Marshal(launch.Receipt)
	out.git, out.report, out.profile, out.world = g, rep, profile, c.World
	out.manifest = store.OperatorFinishReview{Store: c.Credentials, Episode: ep, Preset: start.Preset.Name, Address: c.Address,
		Worktree: g.dir, Report: finishSHA(source), PresetDigest: start.Preset.SHA256, Profile: profile.Digest,
		Sessions: finishSHA(sessionWire), Launch: finishSHA(launchWire), Remote: start.Remote, Base: start.Base, Checks: start.Checks,
		Commit: store.OperatorLocalCommit{Committer: name + " <" + email + ">", Message: title}}
	return out, nil
}

func (e serviceEnv) finishCatalogue(ctx context.Context, c runnerConfig, src operatorStartSource) map[string]store.OperatorFinishReview {
	catalogue := make(map[string]store.OperatorFinishReview)
	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return catalogue
	}
	defer ro.Close()
	ids, err := ro.OperatorCompletedLaunches(ctx)
	if err != nil {
		return catalogue
	}
	for _, id := range ids {
		row, err := ro.OperatorAction(ctx, id, src.clock())
		if err != nil {
			continue
		}
		in, err := e.finishInputs(ctx, c, src, ro, row.Episode, nil)
		if err != nil {
			continue
		}
		v, err := in.frozen(ctx)
		if err == nil {
			if _, err := v.Action(1, false); err == nil {
				catalogue[row.Episode] = v
			}
		}
	}
	return catalogue
}

func (e serviceEnv) consumeFinish(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource) error {
	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return err
	}
	defer ro.Close()
	row, err := ro.NextOperatorFinish(ctx, src.clock())
	if err != nil {
		return err
	}
	a, m, generation, err := finishManifestForRow(c, row)
	if err != nil {
		return err
	}
	cas := operatorPreparationCAS(row)
	journal, err := ro.ReadOperatorFinishJournal(ctx, cas)
	if err != nil {
		return err
	}
	observation := row.Owner != "" && journal.CommitPhase == "reserved"
	verifyAt := src.clock()
	if observation {
		// Expiry stops new effects. Use the signed deadline for proof of what
		// was approved, then observe only the already reserved commit.
		verifyAt = min(verifyAt, row.ExpiresAt-1)
	}
	_, err = verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner)
	if err != nil {
		return err
	}
	for _, check := range m.Checks {
		if check.Unsandboxed {
			return errOperatorFinishTerminal
		}
	}
	if err := prepareOperatorFinish(ctx, c, m.Remote.URL); err != nil {
		return err
	}

	lock, err := store.OpenOperatorEpisodeLock(c.Credentials, row.Episode)
	if errors.Is(err, syscall.EWOULDBLOCK) {
		return operatorFinishBusy(row.ID)
	}
	if err != nil {
		return err
	}
	defer lock.Close()
	fd, err := lock.Verify(c.Credentials, row.Episode)
	if err != nil {
		return err
	}
	w, err := e.finishWork(ctx, c, h, src, ro, fd, row, cas, a, m, generation, observation, verifyAt)
	if err != nil {
		return err
	}
	return runEpisodeFinishWithPrep(ctx, w)
}

func (e serviceEnv) finishWork(ctx context.Context, c runnerConfig, h store.OperatorHeartbeat, src operatorStartSource,
	ro *store.Credentials, fd *os.File, row store.OperatorActionRow, cas store.OperatorCAS, a approval.OperatorAction,
	m store.OperatorFinishReview, generation int64, observation bool, verifyAt int64) (operatorFinishWork, error) {
	w := operatorFinishWork{generation: generation, manifest: m, spool: filepath.Dir(e.runnerPaths().Log),
		read: ro, now: src.clock, budget: e.operatorContext}
	if observation {
		var err error
		w.prep, err = verifyOperatorFinishPrep(ctx, ro, a, verifyAt, cas, row.Owner)
		if err != nil {
			return w, err
		}
		start, err := ro.EpisodeStart(ctx, m.Episode)
		if err != nil || start == nil || start.Remote != m.Remote || start.Base != m.Base {
			return w, errOperatorFinishUncertain
		}
		auth, err := planOperatorFinishAuth(c, m.Remote.URL)
		if err != nil {
			return w, err
		}
		w.git, err = openOperatorFinishGit(auth, filepath.Join(src.workspace, m.Episode), start.Repo, fd)
		if err != nil || w.git.dir != m.Worktree {
			return w, errOperatorFinishUncertain
		}
		w.report.Start = start
		return w, nil
	}
	if _, err := verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner); err != nil {
		return w, err
	}
	in, err := e.finishInputs(ctx, c, src, ro, row.Episode, fd)
	if err != nil {
		return w, err
	}
	if !finishReviewInputsMatch(in.manifest, m) {
		return w, errOperatorFinishReview
	}
	if row.Owner == "" {
		reviewed, err := in.frozen(ctx)
		if err != nil {
			return w, errors.Join(errOperatorFinishReview, err)
		}
		if !reflect.DeepEqual(reviewed, m) || len(m.Protected) != 0 || len(m.Unattributed) != 0 {
			return w, errOperatorFinishReview
		}
		admit, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
		if err != nil {
			return w, errOperatorDraining
		}
		if ctx.Err() == nil && e.checkStartPair(c, src.workspace) == nil && e.proveOldTerminalQuiet() == nil {
			err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
				write, cancel := e.operatorContext(time.Second)
				defer cancel()
				if _, err := verifyOperatorFinishPrep(write, db, a, src.clock(), cas, ""); err != nil {
					return err
				}
				live, err := db.OperatorHeartbeat(write)
				if err != nil || live.Owner != h.Owner || live.CatalogueGeneration != generation || !live.Available(h.Build, "finish", src.clock()) {
					return store.ErrOperatorRunner
				}
				return db.ClaimOperatorAction(write, cas, h.Owner, src.clock())
			})
		} else {
			err = store.ErrOperatorRunner
		}
		admit.Close()
		if err != nil {
			return w, err
		}
		cas.Generation++
		row.Owner = h.Owner.Nonce
	}
	w.prep, err = verifyOperatorFinishPrep(ctx, ro, a, src.clock(), cas, row.Owner)
	if err != nil {
		return w, err
	}

	client := finishClient(c.Address)
	w.report, w.git = in.report, in.git
	w.fresh = func(bound context.Context) error {
		if err := bound.Err(); err != nil {
			return err
		}
		live, err := ro.OperatorHeartbeat(bound)
		if err != nil || live.Owner != h.Owner || row.Owner != h.Owner.Nonce || !live.Available(h.Build, "finish", src.clock()) ||
			e.checkStartPair(c, src.workspace) != nil {
			return store.ErrOperatorRunner
		}
		current, err := e.finishInputs(bound, c, src, ro, row.Episode, fd)
		if err != nil {
			return err
		}
		if !finishReviewInputsMatch(current.manifest, m) || !current.baseMatches(bound) {
			return errOperatorFinishReview
		}
		_, err = verifyOperatorFinishPrep(bound, ro, a, src.clock(), cas, row.Owner)
		return err
	}
	w.run = func(bound context.Context, check store.StartCheck) (checkrun.Result, error) {
		return finishSandboxCheck(bound, client, in.token, check)
	}
	w.record = func(bound context.Context, outcome demoObject) (int64, error) {
		return commitOneObject(bound, client, in.token, outcome, taskOutcomeWrittenBy, "", "outcome record")
	}
	w.evidence = func(bound context.Context, receipt store.OperatorFinishCheck) bool {
		return finishWorldEvidence(bound, c.World, receipt)
	}
	w.commit = func(bound context.Context, spec store.OperatorLocalCommit, contents [][]byte, tree string) error {
		return in.git.commit(bound, w.spool, spec, contents, tree, in.report.Start)
	}
	return w, nil
}

func finishReviewInputsMatch(current, frozen store.OperatorFinishReview) bool {
	return current.Remote == frozen.Remote && current.Base == frozen.Base && current.PresetDigest == frozen.PresetDigest &&
		current.Profile == frozen.Profile && current.Sessions == frozen.Sessions && current.Launch == frozen.Launch &&
		reflect.DeepEqual(current.Checks, frozen.Checks)
}

func (e serviceEnv) printFinishStatus(id string, out io.Writer) error {
	c, _, err := e.installedRunner()
	if err != nil {
		return err
	}
	ro, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return errOperatorFinishEligibility
	}
	defer ro.Close()
	ctx, cancel := e.operatorContext(time.Second)
	defer cancel()
	row, err := ro.OperatorAction(ctx, id, e.terminal.now())
	if err != nil || row.Kind != "finish" {
		return store.ErrOperatorMissing
	}
	j, err := ro.ReadOperatorFinishJournal(ctx, operatorPreparationCAS(row))
	if err != nil {
		return errOperatorFinishUncertain
	}
	fmt.Fprintln(out, "finish action", row.ID, row.State)
	for _, check := range j.Checks {
		fmt.Fprintf(out, "check %q: %s\n", check.Name, check.Result)
	}
	if j.OpenCheck || j.CommitPhase == "reserved" {
		fmt.Fprintln(out, errOperatorFinishUncertain)
	}
	if j.Commit.State == store.OperatorLocalCommitDone {
		fmt.Fprintln(out, "local commit", j.Commit.SHA, "retained; PR confirmation has not been issued")
	} else if j.Commit.State == "commit-failed" {
		fmt.Fprintln(out, errOperatorCommitFailed)
	}
	return nil
}

// An admitted check may outlast the heartbeat's admission window. The runner
// keeps the same pinned catalogue alive while the exclusive finish turn runs.
func (e serviceEnv) runFinishTurn(c runnerConfig, h *store.OperatorHeartbeat, src operatorStartSource,
	stops <-chan os.Signal, ticks <-chan time.Time) error {
	ctx, cancel := e.operatorContext(2 * time.Hour)
	defer cancel()
	reply := make(chan error, 1)
	look := *h
	go func() { reply <- e.consumeFinish(ctx, c, look, src) }()
	for {
		select {
		case err := <-reply:
			return err
		case <-stops:
			cancel()
			<-reply
			return errOperatorDraining
		case _, live := <-ticks:
			if !live {
				ticks = nil
				continue
			}
			_ = e.touchRunner(c, h) // A busy writer leaves the real heartbeat deadline in force.
		}
	}
}
