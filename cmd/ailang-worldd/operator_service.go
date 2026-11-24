package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

const runnerLabel = "com.sunholo.ailang-world.operator-runner"
const runnerSystemPath = "/usr/bin:/bin:/usr/sbin:/sbin"
const runnerHelp = `usage: ailang-worldd service runner install --config FILE | start | stop | status | cancel ACTION-ID | uninstall
Install, stop, cancel and uninstall require typed y. Start uses the approved configuration.
Config JSON keys: home, uid, address, world, credentials, executables (absolute path to sha256 digest), preset_names, tool_dirs, backend ("headless").
The World service must already be running. Use absolute, canonical paths and include its installed binary in executables.
`

type runnerConfig struct {
	Home        string            `json:"home"`
	UID         int               `json:"uid"`
	Address     string            `json:"address"`
	World       string            `json:"world"`
	Credentials string            `json:"credentials"`
	Executables map[string]string `json:"executables"`
	PresetNames []string          `json:"preset_names"`
	ToolDirs    []string          `json:"tool_dirs"`
	Backend     string            `json:"backend"`
}

func (e serviceEnv) runnerPaths() servicePaths {
	p := e.paths
	root := filepath.Dir(p.Record)
	p.Plist = filepath.Join(filepath.Dir(p.Plist), runnerLabel+".plist")
	p.Record = filepath.Join(root, "operator-runner.json")
	p.Log = filepath.Join(root, "operator-runner", "runner.log")
	return p
}

func actualRunnerLogin() (int, string, error) {
	u, err := user.LookupId(strconv.Itoa(os.Getuid()))
	if err != nil {
		return 0, "", err
	}
	return os.Getuid(), u.HomeDir, nil
}

// The injected home is the test trust anchor. Production additionally gets
// its login home from the account database, independently of the environment.
func runnerPath(home, path string, uid int, directory bool, private os.FileMode) error {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return errors.New("path is not absolute and clean")
	}
	for p := path; ; p = filepath.Dir(p) {
		st, err := os.Lstat(p)
		if err != nil {
			return err
		}
		owner, ok := st.Sys().(*syscall.Stat_t)
		if !ok || (int(owner.Uid) != uid && (within(home, p) || p == home || owner.Uid != 0)) {
			return errors.New("path ownership mismatch")
		}
		wantDir := p != path || directory
		if wantDir && !st.IsDir() || !wantDir && !st.Mode().IsRegular() {
			return errors.New("path has wrong type or is a symlink")
		}
		if st.Mode().Perm()&0022 != 0 {
			return errors.New("path writable by others")
		}
		if p == path && private != 0 && st.Mode().Perm() != private {
			return errors.New("path permissions are not private")
		}
		if p == home || p == "/" {
			return nil
		}
	}
}

func readRunnerConfig(path string) (runnerConfig, []byte, error) {
	var c runnerConfig
	f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return c, nil, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil || !st.Mode().IsRegular() {
		return c, nil, errors.New("expected a regular runner configuration")
	}
	b, err := io.ReadAll(io.LimitReader(f, 16385))
	if err != nil || len(b) > 16384 {
		return c, nil, errors.New("cannot read bounded runner configuration")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if d.Decode(&c) != nil || d.Decode(new(any)) != io.EOF {
		return c, nil, errors.New("invalid runner configuration JSON")
	}
	return c, b, nil
}

func (e serviceEnv) checkRunnerConfig(c runnerConfig) error {
	login := e.runnerLogin
	if login == nil {
		login = actualRunnerLogin
	}
	uid, home, err := login()
	switch {
	case err != nil || uid == 0 || uid != c.UID:
		return errors.New("runner requires the approved non-root login UID")
	case home != c.Home || home != e.paths.Home:
		return errors.New("runner HOME differs from the login home")
	case c.Backend != "headless":
		return errors.New("runner backend must be headless")
	case len(c.PresetNames) == 0 || len(c.PresetNames) > 16 || len(c.Executables) > 16 || len(c.ToolDirs) > 16:
		return errors.New("runner configuration list size is invalid")
	}
	host, _, _ := net.SplitHostPort(c.Address)
	if !serviceLoopback(c.Address) || !net.ParseIP(host).IsLoopback() {
		return errors.New("runner daemon must be numeric loopback")
	}
	anchor := home
	if e.runnerLogin == nil {
		anchor = "" // verify all login-home ancestors in production
	}
	if err := runnerPath(anchor, home, uid, true, 0); err != nil {
		return err
	}
	for _, preset := range c.PresetNames {
		if preset == "se-tools" || len(preset) > verbs.MaxPresetName || !verbs.EpisodeIDPattern.MatchString(preset) {
			return errors.New("invalid runner preset name")
		}
	}
	if c.Executables[e.paths.Binary] == "" {
		return errors.New("missing installed World executable")
	}
	for path, sum := range c.Executables {
		if err := runnerPath(home, path, uid, false, 0); err != nil {
			return err
		}
		st, err := os.Stat(path)
		if err != nil || st.Mode().Perm()&0100 == 0 || sum == "" || serviceFileHash(path) != sum {
			return errors.New("runner executable changed or is not executable")
		}
	}
	for _, dir := range c.ToolDirs {
		if strings.ContainsAny(dir, ":\n") {
			return errors.New("invalid tool directory")
		}
		if err := runnerPath(home, dir, uid, true, 0); err != nil {
			return err
		}
	}
	world, worldErr := store.CanonicalPath(c.World)
	creds, credsErr := store.CanonicalPath(c.Credentials)
	if worldErr != nil || credsErr != nil || world != c.World || creds != c.Credentials || world == creds {
		return errors.New("runner stores must be distinct canonical files")
	}
	for _, path := range []string{c.World, c.Credentials} {
		if err := runnerPath(home, path, uid, false, 0); err != nil {
			return err
		}
	}
	return nil
}

func (e serviceEnv) checkRunnerPair(c runnerConfig) error {
	r, err := readServiceRecord(e.paths)
	if err != nil || r.Uninstalled || r.Address != c.Address || r.BinarySHA256 != c.Executables[e.paths.Binary] {
		return errors.New("runner differs from the installed World service")
	}
	worldRecord, worldErr := store.CanonicalPath(r.Store)
	credsRecord, credsErr := store.CanonicalPath(r.Credentials)
	if worldErr != nil || credsErr != nil || worldRecord != c.World || credsRecord != c.Credentials {
		return errors.New("runner store pair differs from the installed World service")
	}
	h, err := e.health(c.Address)
	if err != nil || h.Status != "ok" {
		return errors.New("daemon unavailable")
	}
	world, err := store.CanonicalPath(h.DBPath)
	if err != nil || world != c.World {
		return errors.New("daemon World store mismatch")
	}
	creds, err := store.CanonicalPath(h.CredentialsDB)
	if err != nil || creds != c.Credentials {
		return errors.New("daemon credentials store mismatch")
	}
	if h.Build.ExecutableSHA256 != c.Executables[e.paths.Binary] {
		return errors.New("daemon build mismatch; restart pending")
	}
	db, err := store.OpenCredentialsReadOnly(runnerDB(c.Credentials))
	if err != nil {
		return errors.New("cannot read daemon credentials schema")
	}
	defer db.Close()
	if db.Version() != store.CredentialsSchemaVersion {
		return errors.New("daemon credentials schema mismatch")
	}
	return nil
}

func runnerPlist(p servicePaths, c runnerConfig, configBytes []byte) string {
	paths := append(append([]string(nil), c.ToolDirs...), runnerSystemPath)
	return renderServicePlist(runnerLabel, p, []string{p.Binary, "operator-runner"}, c.Home,
		[2]string{"HOME", c.Home}, [2]string{"PATH", strings.Join(paths, ":")},
		[2]string{"WORLD_OPERATOR_CONFIG_SHA256", serviceHash(configBytes)},
		[2]string{"GIT_TERMINAL_PROMPT", "0"}, [2]string{"GH_PROMPT_DISABLED", "1"})
}

func (e serviceEnv) installedRunner() (runnerConfig, []byte, error) {
	p := e.runnerPaths()
	c, b, err := readRunnerConfig(p.Record)
	if err != nil {
		return c, nil, err
	}
	if err := e.checkRunnerConfig(c); err != nil {
		return c, nil, err
	}
	for _, path := range []string{p.Record, p.Plist, p.Log} {
		if err := runnerPath(c.Home, path, c.UID, false, 0600); err != nil {
			return c, nil, err
		}
	}
	if err := runnerPath(c.Home, filepath.Dir(p.Log), c.UID, true, 0700); err != nil {
		return c, nil, err
	}
	if serviceFileHash(p.Plist) != serviceHash([]byte(runnerPlist(p, c, b))) {
		return c, nil, errors.New("runner plist differs from approved configuration")
	}
	return c, b, nil
}

// Lock the service root inode, without rewriting any base service file.
func runnerLifecycleLock(root string) (*os.File, error) {
	f, err := os.OpenFile(root, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_DIRECTORY, 0)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func runServiceRunner(args []string, out, errOut io.Writer, e serviceEnv) int {
	refuse := func(err error) int {
		return serviceMessage(errOut, "runner refused", runnerDiagnostic(err), "ailang-worldd service runner --help; inspect the approved configuration", 1)
	}
	if len(args) == 0 || isHelpToken(args[len(args)-1]) {
		fmt.Fprint(out, runnerHelp)
		return 0
	}
	verb := args[0]
	if e.goos != "darwin" || e.label != worldServiceLabel {
		return refuse(errors.New("runner requires the macOS World service"))
	}
	switch verb {
	case "install":
		if len(args) != 3 || args[1] != "--config" {
			return refuse(errors.New("expected install --config FILE"))
		}
	case "cancel":
		if len(args) != 2 || !verbs.ValidOperatorActionID(args[1]) {
			return refuse(errors.New("expected cancel ACTION-ID"))
		}
	case "start", "stop", "status", "uninstall":
		if len(args) != 1 {
			return refuse(errors.New("unexpected runner argument"))
		}
	default:
		return refuse(errors.New("unknown runner operation"))
	}
	p := e.runnerPaths()
	if verb == "status" {
		if err := e.runnerStatus(out); err != nil {
			return refuse(err)
		}
		return 0
	}

	if verb == "install" || verb == "stop" || verb == "cancel" || verb == "uninstall" {
		change := "operator runner " + strings.Join(args, " ")
		if verb == "install" {
			c, b, err := readRunnerConfig(args[2])
			if err != nil {
				return refuse(errors.New("invalid runner setup file"))
			}
			if err := e.checkRunnerConfig(c); err != nil {
				return refuse(err)
			}
			// Display only typed, validated fields from the approved bytes.
			pretty, _ := json.MarshalIndent(c, "", "  ")
			change = "install operator runner " + serviceHash(b) + "\n" + string(pretty)
			if !e.confirm(errOut, change) {
				return 1
			}
			lock, err := runnerLifecycleLock(filepath.Dir(p.Record))
			if err != nil {
				return refuse(errors.New("runner lifecycle is busy or unavailable"))
			}
			defer lock.Close()
			if err := e.installRunner(c, b); err != nil {
				return refuse(err)
			}
			fmt.Fprintln(out, "runner installed; next: ailang-worldd service runner start")
			return 0
		}
		if !e.confirm(errOut, change) {
			return 1
		}
	}
	lock, err := runnerLifecycleLock(filepath.Dir(p.Record))
	if err != nil {
		return refuse(errors.New("runner lifecycle is busy or unavailable"))
	}
	defer lock.Close()
	if err := e.changeRunner(verb, args[1:]); err != nil {
		return refuse(err)
	}
	fmt.Fprintln(out, "runner operation complete")
	return 0
}

func (e serviceEnv) installRunner(c runnerConfig, b []byte) error {
	p := e.runnerPaths()
	if err := e.checkRunnerConfig(c); err != nil {
		return err
	}
	if err := e.checkRunnerPair(c); err != nil {
		return err
	}
	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
		return errors.New("runner is loaded or its supervisor is unavailable")
	}
	for _, path := range []string{p.Plist, p.Record} {
		if err := serviceSafePath(p.Home, path); err != nil {
			return err
		}
		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
			return errors.New("runner installation already exists; stop and uninstall it before replacing")
		}
	}
	if err := runnerPath(p.Home, filepath.Dir(p.Plist), c.UID, true, 0); err != nil {
		return err
	}
	if err := runnerPath(p.Home, filepath.Dir(p.Record), c.UID, true, 0700); err != nil {
		return err
	}
	// A pre-existing spool may contain durable evidence. Keep it, including its log.
	spool := filepath.Dir(p.Log)
	e.note(spool)
	if err := os.Mkdir(spool, 0700); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	if err := runnerPath(c.Home, spool, c.UID, true, 0700); err != nil {
		return err
	}
	if _, err := os.Lstat(p.Log); errors.Is(err, os.ErrNotExist) {
		if err := e.atomic(p.Log, nil, 0600); err != nil {
			return err
		}
	}
	if err := runnerPath(c.Home, p.Log, c.UID, false, 0600); err != nil {
		return err
	}
	if err := e.atomic(p.Record, b, 0600); err != nil {
		return err
	}
	if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, b)), 0600); err != nil {
		e.note(p.Record)
		return errors.Join(err, os.Remove(p.Record))
	}
	return nil
}

func (e serviceEnv) changeRunner(verb string, args []string) error {
	p := e.runnerPaths()
	login := e.runnerLogin
	if login == nil {
		login = actualRunnerLogin
	}
	uid, home, err := login()
	if err != nil || uid == 0 || home != p.Home {
		return errors.New("runner requires the non-root login home")
	}
	if err := runnerPath(home, filepath.Dir(p.Record), uid, true, 0700); err != nil {
		return err
	}
	switch verb {
	case "stop", "uninstall":
		for _, path := range []string{p.Plist, p.Record} {
			if err := serviceSafePath(p.Home, path); err != nil {
				return err
			}
		}
		if err := e.controller.Bootout(e.domain, runnerLabel); err != nil && !errors.Is(err, errServiceNotLoaded) {
			return errors.New("could not stop the companion")
		}
		if verb == "uninstall" {
			for _, path := range []string{p.Plist, p.Record} {
				e.note(path)
				if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
					return err
				}
			}
		}
		return nil
	}
	c, _, err := e.installedRunner()
	if err != nil {
		return err
	}
	if verb == "cancel" {
		return runnerCancel(c, args[0], e.terminal.now())
	}
	if err := e.checkRunnerPair(c); err != nil {
		return err
	}
	_, err = e.controller.Print(e.domain, runnerLabel)
	switch {
	case errors.Is(err, errServiceNotLoaded):
		if err := e.controller.Bootstrap(e.domain, p.Plist); err != nil {
			_ = e.controller.Bootout(e.domain, runnerLabel)

			return errors.New("could not bootstrap the companion")
		}
	case err != nil:
		return errors.New("could not inspect the companion")
	}
	return nil
}

func runnerDiagnostic(err error) string {
	var pathErr *os.PathError
	if errors.As(err, &pathErr) {
		return "runner file is absent or inaccessible"
	}
	return err.Error()
}
