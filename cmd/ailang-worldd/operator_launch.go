package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/agentrun"
	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const operatorLaunchBudget = 5 * time.Second

var errOperatorLaunchHandled = errors.New("operator launch tick handled")

type operatorLaunchJob struct {
	config runnerConfig
	row    store.OperatorActionRow
	wall   time.Duration
	secret [32]byte
	locks  [3]*os.File // runner, episode, lifecycle
}

type operatorLaunchInputs struct {
	token, dir, sessions, binary, digest, prompt string
	argv                                         []string
	proof                                        *approval.OperatorMintProof
	wall                                         time.Duration
}

func (e serviceEnv) launchInputs(ctx context.Context, c runnerConfig, root string, row store.OperatorActionRow,
	files store.OperatorFileState, a approval.OperatorAction, ro *store.Credentials) (operatorLaunchInputs, error) {
	in := operatorLaunchInputs{binary: operatorFileFact(a, "launch_binary"), sessions: filepath.Dir(files.Mint.TokenPath)}
	digest, err := authority.OperatorStartLaunchDigest(a)
	if err != nil {
		return in, err
	}
	in.digest = digest
	if !a.Options.Headless || a.Store != c.Credentials || c.Backend != "headless" ||
		filepath.Base(in.binary) != "codex" || c.Executables[in.binary] != in.digest ||
		!slices.Contains(c.PresetNames, a.Options.Preset) || operatorFileFact(a, "served_address") != c.Address {
		return in, store.ErrOperatorConflict
	}
	now := e.terminal.now()
	if now >= row.ExpiresAt || now >= files.Session.ExpiresAt {
		return in, store.ErrOperatorExpired
	}
	if err := e.checkStartPair(c, root); err != nil {
		return in, err
	}
	h, err := e.health(c.Address)
	build, _ := json.Marshal(h.Build)
	if err != nil || string(build) != operatorFileFact(a, "served_build") || h.ToolRef != operatorFileFact(a, "served_interpreter") {
		return in, store.ErrOperatorConflict
	}
	project := operatorFileFact(a, "mint_project")
	if !slices.ContainsFunc(h.ExecProfiles, func(p daemon.HealthExecProfile) bool {
		return p.Project == project && p.Digest == operatorFileFact(a, "mint_profile")
	}) {
		return in, store.ErrOperatorConflict
	}
	commands, err := launchCommands(h, project)
	if err != nil {
		return in, err
	}
	if _, err := recoveryWorktree(ctx, store.OperatorStartRecovery{Worktree: files.Worktree}, a); err != nil {
		return in, err
	}
	set, start, err := operatorFilePaths(files, a)
	proofs, proofErr := readStartRecoveryFiles(c, files, a)
	if err != nil || proofErr != nil || start.Agent != "codex" || proofs.Files[0].Status != "exact" || proofs.Files[1].Status != "exact" {
		return in, store.ErrOperatorConflict
	}
	set.Device, set.Inode = proofs.Dir.ID.Dev, proofs.Dir.ID.Inode
	d, _ := openOperatorFilesDir(set)
	if d == nil {
		return in, store.ErrOperatorConflict
	}
	defer d.Close()
	tokenProof, token := readOperatorFile(d, set.Names[1], set.Hashes[1], set.UID)
	if tokenProof.Status != "exact" {
		return in, store.ErrOperatorConflict
	}
	in.token = string(token)
	in.wall, err = time.ParseDuration(operatorFileFact(a, "launch_wall"))
	if err != nil || in.wall <= 0 || in.wall > time.Duration(1<<63-1)-time.Minute {
		return in, store.ErrOperatorConflict
	}
	in.dir = filepath.Join(filepath.Dir(in.sessions), "launch", row.Episode)
	if physical, err := operatorPhysicalTarget(in.dir); err != nil || physical != in.dir {
		return in, store.ErrOperatorConflict
	}
	if err := validateLaunchDir(in.dir, root, operatorWorktreeFact(a), in.sessions); err != nil {
		return in, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://"+c.Address+"/v1/session/self", nil)
	if err != nil {
		return in, err
	}
	req.Header.Set("Authorization", "Bearer "+in.token)
	res, err := e.client.Do(req)
	if err != nil {
		return in, store.ErrOperatorConflict
	}
	defer res.Body.Close()
	var self daemon.SessionSelf
	body, err := io.ReadAll(io.LimitReader(res.Body, 1<<20))
	got, gotErr := filepath.EvalSymlinks(operatorWorktreeFact(a))
	if err != nil || res.StatusCode != http.StatusOK || json.Unmarshal(body, &self) != nil ||
		self.SecondsLeft <= 0 || self.ExpiresAt != files.Session.ExpiresAt || self.Episode != row.Episode {
		return in, store.ErrOperatorConflict
	}
	peer, peerErr := filepath.EvalSymlinks(self.Worktree)
	if gotErr != nil || peerErr != nil || got != peer {
		return in, store.ErrOperatorConflict
	}
	contextText, err := launchContext(body, c.Address, in.token)
	if err != nil {
		return in, err
	}
	in.prompt = launchPrompt(a.Task, commands)
	in.argv = launchArgv(start.Agent, start.Model, in.prompt, in.dir, "", c.Address, contextText, "", nil)
	in.argv[0] = in.binary
	for _, arg := range in.argv {
		if strings.Contains(arg, in.token) {
			return in, store.ErrOperatorConflict
		}
	}
	proof, err := authority.RecheckClaimedOperatorAction(ctx, ro, a, [32]byte(row.Challenge), row.ExpiresAt, e.terminal.now(), row)
	if err != nil {
		return in, err
	}
	in.proof = proof.OperatorConfirmation().MintProof()
	return in, nil
}

type operatorAgentRunner func(context.Context, []string, string, []string, *os.File, *os.File, time.Duration, agentrun.LaunchHooks) (agentrun.Result, error)

func (e serviceEnv) launchPreparedOperatorStart(outer context.Context, c runnerConfig, h store.OperatorHeartbeat,
	src operatorStartSource, stops <-chan os.Signal, process *os.File, log io.Writer, ro *store.Credentials,
	row store.OperatorActionRow, current store.OperatorStartRecovery, episode *store.OperatorEpisodeLock) error {
	ctx, cancel := e.operatorContext(operatorLaunchBudget)
	defer cancel()
	files, a, err := ro.ReadOperatorFiles(ctx, row.ID)
	if err != nil {
		return err
	}
	if err := e.proveOldTerminalQuiet(); err != nil {
		return err
	}
	in, err := e.launchInputs(ctx, c, src.workspace, row, files, a, ro)
	if err != nil {
		return err
	}
	if err := launchLockCheck(filepath.Join(in.sessions, row.Episode+".launch.lock")); err != nil {
		return err
	}
	if outer.Err() != nil {
		return outer.Err()
	}
	fd, err := episode.Verify(c.Credentials, row.Episode)
	if err != nil {
		return err
	}
	lifecycle, err := runnerLifecycleLock(filepath.Dir(e.paths.Record))
	if err != nil {
		return errOperatorDraining
	}
	defer lifecycle.Close()
	select {
	case <-stops:
		return errOperatorDraining
	default:
	}
	job := operatorLaunchJob{config: c, row: row, wall: in.wall, locks: [3]*os.File{process, fd, lifecycle}}
	if _, err := rand.Read(job.secret[:]); err != nil {
		return err
	}
	if err := runnerWrite(c.Credentials, func(db *store.Credentials) error {
		live, err := db.OperatorHeartbeat(ctx)
		catalogue, catErr := strconv.ParseInt(operatorFileFact(a, "mint_catalogue"), 10, 64)
		if err != nil || catErr != nil || !preparationHeartbeat(live, h.Owner, h.Build, catalogue, src.clock()) {
			return store.ErrOperatorRunner
		}
		return db.ReserveOperatorLaunch(ctx, current, job.secret, in.proof, src.clock)
	}); err != nil {
		return err
	}
	if e.operatorSupervisor != nil {
		return e.operatorSupervisor(job, log)
	}
	if ctx.Err() != nil || serviceFileHash(e.paths.Binary) != c.Executables[e.paths.Binary] {
		return e.operatorSpawnRefusal(c, row.ID, in.digest, log)
	}
	return e.superviseOperatorLaunch(job, stops, log, h, in.digest)
}

func (e serviceEnv) superviseOperatorLaunch(job operatorLaunchJob, stops <-chan os.Signal, log io.Writer,
	h store.OperatorHeartbeat, digest string) error {
	cmd := exec.Command(e.paths.Binary, "operator-launch-worker", job.row.ID)
	cmd.Stdin, cmd.Stderr = bytes.NewReader(job.secret[:]), log
	cmd.ExtraFiles = job.locks[:]
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return e.operatorSpawnRefusal(job.config, job.row.ID, digest, log)
	}
	job.locks[2].Close() // Worker holds the same description through agent Start.
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	ticker := time.NewTicker(runnerPoll)
	defer ticker.Stop()
	cap := time.NewTimer(job.wall + time.Minute)
	defer cap.Stop()
	for {
		select {
		case err := <-done:
			return err
		case <-ticker.C:
			_ = e.runnerTick(job.config, &h)
		case <-stops:
			_ = cmd.Process.Signal(syscall.SIGTERM)
			stops = nil
		case <-cap.C:
			_ = cmd.Process.Signal(syscall.SIGTERM)
			return errors.New("launch supervisor did not report completion")
		}
	}
}

func (e serviceEnv) operatorSpawnRefusal(c runnerConfig, id, digest string, log io.Writer) error {
	r := store.OperatorLaunchReceipt{Run: id, Executable: digest, Status: "spawn-refused"}
	if err := e.writeOperatorLaunch(c, r, nil); err != nil {
		return err
	}
	fmt.Fprintf(log, "operator launch %s: spawn-refused\n", id)
	return nil
}

func (e serviceEnv) writeOperatorLaunch(c runnerConfig, r store.OperatorLaunchReceipt, history *store.EpisodeLaunch) error {
	ctx, cancel := e.operatorContext(operatorLaunchBudget)
	defer cancel()
	return runnerWrite(c.Credentials, func(db *store.Credentials) error {
		read, err := db.ReadOperatorLaunch(ctx, r.Run)
		if err != nil {
			return err
		}
		return db.RecordOperatorLaunch(ctx, read, r, history, e.terminal.now())
	})
}

func (e serviceEnv) runOperatorLaunchWorker(args []string, out io.Writer) int {
	proc, stop := serveSignalContext()
	defer stop()
	if len(args) != 1 || !verbs.ValidOperatorActionID(args[0]) {
		return exitUsage
	}
	c, b, err := e.installedRunner()
	if err != nil || os.Getenv("WORLD_OPERATOR_CONFIG_SHA256") != serviceHash(b) ||
		daemon.CurrentBuild().ExecutableSHA256 != c.Executables[e.paths.Binary] {
		return exitUsage
	}
	ctx, cancel := e.operatorContext(operatorLaunchBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
	if err != nil {
		return exitUsage
	}
	defer ro.Close()
	row, err := ro.OperatorAction(ctx, args[0], 0)
	if err != nil || !verbs.EpisodeIDPattern.MatchString(row.Episode) {
		return exitUsage
	}
	job := operatorLaunchJob{config: c, row: row}
	names := [3]string{c.Credentials + ".operator-runner.lock",
		store.OperatorEpisodeLockPath(c.Credentials, row.Episode), filepath.Dir(e.paths.Record)}
	for i, name := range names {
		fd := os.NewFile(uintptr(3+i), name)
		defer fd.Close()
		st, err := fd.Stat()
		path, pathErr := os.Lstat(name)
		if err != nil || pathErr != nil || !os.SameFile(st, path) || syscall.Flock(int(fd.Fd()), syscall.LOCK_EX|syscall.LOCK_NB) != nil {
			return exitUsage
		}
		syscall.CloseOnExec(int(fd.Fd()))
		job.locks[i] = fd
	}
	if _, err := io.ReadFull(os.Stdin, job.secret[:]); err != nil {
		return exitUsage
	}
	if err := e.executeOperatorLaunch(proc, job, out); err != nil {
		fmt.Fprintf(out, "operator launch %s: %s\n", row.ID, operatorLaunchFailure(err))
		return 1
	}
	return exitOK
}

func operatorLaunchFailure(err error) string {
	var verdict *approval.NotAdmittedError
	if errors.Is(err, store.ErrOperatorExpired) || errors.As(err, &verdict) && verdict.Verdict == approval.RefuseExpired {
		return "expired"
	}
	return "blocked or uncertain; inspect service runner launch-status"
}

func operatorLaunchProcessStatus(s store.OperatorLaunchSnapshot, read func(int) (string, bool)) string {
	if s.NoSpawn {
		return "spawn-refused"
	}
	if s.Phase == "reserved" {
		return "uncertain"
	}
	if s.Receipt.Status == "launched" {
		if birth, live := read(s.Receipt.PID); !live || birth != s.Receipt.Birth {
			return "uncertain"
		}
	}
	return s.Receipt.Status
}

func (e serviceEnv) operatorLaunchStatus(id string, out io.Writer) error {
	c, _, err := e.installedRunner()
	if err != nil {
		return err
	}
	db, err := store.OpenCredentialsReadOnly(c.Credentials)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := e.operatorContext(operatorLaunchBudget)
	defer cancel()
	s, err := db.ReadOperatorLaunch(ctx, id)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "operator launch %s: %s\n", id, operatorLaunchProcessStatus(s, operatorProcessIdentity))
	return nil
}

func launchOperatorEnvironment(c runnerConfig, token string) []string {
	path := strings.Join(append(slices.Clone(c.ToolDirs), runnerSystemPath), ":")
	return []string{"HOME=" + c.Home, "PATH=" + path, "GIT_TERMINAL_PROMPT=0", "GH_PROMPT_DISABLED=1", "WORLD_SESSION=" + token}
}

func operatorLaunchStage(dir string) (*os.File, error) {
	f, err := os.CreateTemp(dir, ".operator-output-*")
	if err != nil {
		return nil, err
	}
	if err := os.Remove(f.Name()); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func (e serviceEnv) executeOperatorLaunch(proc context.Context, job operatorLaunchJob, log io.Writer) error {
	ctx, cancel := e.operatorContext(operatorLaunchBudget)
	defer cancel()
	c, row := job.config, job.row
	record, err := readServiceRecord(e.paths)
	if err != nil {
		return err
	}
	ro, err := store.OpenCredentialsReadOnly(c.Credentials)
	if err != nil {
		return err
	}
	defer ro.Close()
	files, a, err := ro.ReadOperatorSpawnInput(ctx, row.ID)
	if err != nil {
		return err
	}
	in, err := e.launchInputs(ctx, c, record.Workspace, row, files, a, ro)
	if err != nil {
		return err
	}
	lock := filepath.Join(in.sessions, row.Episode+".launch.lock")
	if err := launchLockCheck(lock); err != nil {
		return err
	}
	f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	identity, statErr := f.Stat()
	defer func() {
		if st, err := os.Lstat(lock); err == nil && statErr == nil && os.SameFile(identity, st) {
			_ = os.Remove(lock)
		}
	}()
	_, err = fmt.Fprintln(f, os.Getpid())
	if err := errors.Join(err, f.Close()); err != nil {
		return err
	}
	for _, dir := range []string{filepath.Dir(in.dir), in.dir} {
		if err := os.Mkdir(dir, 0700); err == nil {
			defer os.Remove(dir) // Remove an empty directory created by this run.
		} else if !errors.Is(err, os.ErrExist) {
			return err
		}
	}
	if err := runnerPath(c.Home, in.dir, c.UID, true, 0700); err != nil {
		return err
	}
	stdout, err := operatorLaunchStage(in.sessions)
	if err != nil {
		return err
	}
	defer stdout.Close()
	stderr, err := operatorLaunchStage(in.sessions)
	if err != nil {
		return err
	}
	defer stderr.Close()
	r := store.OperatorLaunchReceipt{Run: row.ID, Executable: in.digest, Status: "launched"}
	var writer *store.Credentials
	granted := false
	started := time.Time{}
	hooks := agentrun.LaunchHooks{Files: job.locks[:2],
		Before: func() error {
			var err error
			writer, err = store.OpenCredentials(runnerDB(c.Credentials))
			if err != nil {
				return err
			}
			if err := proc.Err(); err != nil {
				return err
			}
			if serviceFileHash(in.binary) != in.digest {
				return &approval.BindingError{Field: "launch_binary_sha256"}
			}
			if err := validateLaunchDir(in.dir, record.Workspace, operatorWorktreeFact(a), in.sessions); err != nil {
				return err
			}
			err = writer.AuthorizeOperatorSpawn(ctx, row.ID, job.secret, in.proof, e.terminal.now)
			granted = err == nil
			return err
		},
		Started: func(pid int) error {
			defer func() { writer.Close(); writer = nil; job.locks[2].Close() }()
			r.PID = pid
			r.Birth, _ = operatorProcessIdentity(pid)
			if r.Birth == "" {
				return errors.New("process birth unavailable")
			}
			s, err := writer.ReadOperatorLaunch(ctx, row.ID)
			if err != nil {
				return err
			}
			started = time.Now()
			if err := writer.RecordOperatorLaunch(ctx, s, r, nil, e.terminal.now()); err != nil {
				return err
			}
			fmt.Fprintf(log, "operator launch %s: launched\n", row.ID)
			return nil
		},
	}
	run := e.operatorAgentRun
	if run == nil {
		run = agentrun.RunObserved
	}
	result, runErr := run(proc, in.argv, in.dir, launchOperatorEnvironment(c, in.token), stdout, stderr, in.wall, hooks)
	if writer != nil {
		writer.Close()
	}
	if r.PID == 0 {
		if granted {
			return e.operatorSpawnRefusal(c, row.ID, in.digest, log)
		}
		return runErr
	}
	if runErr != nil {
		return runErr // the active reservation or observed receipt stays uncertain
	}
	r.Exit = result.ExitCode
	r.Status = operatorLaunchExitStatus(proc, result)
	history := operatorLaunchHistory(in, r, row.Episode, stdout, stderr, started, time.Now())
	if history.FinishReason == "output-failed" {
		r.Status = "failed"
	}
	if err := e.writeOperatorLaunch(c, r, &history); err != nil {
		return err
	}
	fmt.Fprintf(log, "operator launch %s: %s\n", row.ID, r.Status)
	return nil
}

func operatorLaunchHistory(in operatorLaunchInputs, r store.OperatorLaunchReceipt, ep string, stdout, stderr *os.File, started, ended time.Time) store.EpisodeLaunch {
	transcript := filepath.Join(in.sessions, ep+".transcript.jsonl")
	errPath := filepath.Join(in.sessions, ep+".agent.err")
	state, digest, outErr := keepLaunchStream(stdout, transcript, in.token, 1)
	_, _, errErr := keepLaunchStream(stderr, errPath, in.token, 1)
	finish := r.Status
	if outErr != nil || errErr != nil {
		finish = "output-failed"
	}
	return store.EpisodeLaunch{Mode: "headless", StartedAt: started.UTC().Format(time.RFC3339),
		EndedAt: ended.UTC().Format(time.RFC3339), Exit: r.Exit, FinishReason: finish,
		Host: in.binary, ArgvRedacted: redactLaunchArgv(in.argv, in.prompt, in.token),
		Transcript: transcript, TranscriptSHA256: digest, TranscriptState: state}
}

func operatorLaunchExitStatus(proc context.Context, result agentrun.Result) string {
	if result.Reason == "timeout" && proc.Err() == nil {
		return "capped"
	}
	if result.ExitCode != 0 && result.Reason != "signal" && proc.Err() == nil {
		return "failed"
	}
	return "stopped"
}
