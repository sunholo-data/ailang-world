package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/pinfetch"
	"github.com/sunholo-data/ailang-world/host/store"
)

const worldServiceLabel = "com.sunholo.ailang-world.worldd"
const serviceHelp = `usage: ailang-worldd service install|status|upgrade|stop|start|uninstall [flags]
  install [--bind 127.0.0.1:7644] [--db PATH] [--credentials-db PATH]
    [--workspace-root DIR] [--interpreter-dir DIR] [--tools-dir DIR]
    [--exec-sandbox DIR] [--exec-profiles-dir DIR] [--exec-node PATH]
    [--examples-dir DIR] [--surface-index PATH] [--surface-follow PATH] [--binary PATH] [--replace]
  status [--json] [--verbose]
  upgrade [--binary FILE] [--pause]
Install, stop and uninstall require a typed y at the controlling terminal.
Start uses the approved plist; edited plists require confirmation.
Upgrade requires a typed y; --pause permits interruption of active episodes.
Uninstall keeps all data, binaries and logs. macOS preview.
`

type servicePaths struct {
	Home, Plist, Binary, Record, Lock, Log string
}

func homeServicePaths(home string) servicePaths {
	return servicePaths{home, filepath.Join(home, "Library", "LaunchAgents", worldServiceLabel+".plist"),
		filepath.Join(home, ".ailang-world", "bin", "ailang-worldd"), filepath.Join(home, ".ailang-world", "service.json"),
		filepath.Join(home, ".ailang-world", "service.lock"), filepath.Join(home, "Library", "Logs", "ailang-world", "worldd.log")}
}

type serviceRecord struct {
	Schema         string           `json:"schema"`
	Label          string           `json:"label"`
	Plist          string           `json:"plist"`
	PlistSHA256    string           `json:"plist_sha256"`
	PlistText      string           `json:"plist_text"`
	Binary         string           `json:"binary"`
	BinarySHA256   string           `json:"binary_sha256"`
	PreviousSHA256 string           `json:"previous_sha256"`
	Address        string           `json:"address"`
	Store          string           `json:"store"`
	Credentials    string           `json:"credentials"`
	Workspace      string           `json:"workspace"`
	Log            string           `json:"log"`
	InstalledAt    int64            `json:"installed_at"`
	Uninstalled    bool             `json:"uninstalled"`
	Upgrades       []serviceUpgrade `json:"upgrades"`
}
type serviceEnv struct {
	paths               servicePaths
	controller          serviceController
	terminal            sessionEnv
	goos, domain, label string
	wait                time.Duration
	client              *http.Client
	portFree            func(string) error
	// observe records every intended filesystem mutation, including temp paths.
	observe func(string)
	probe   func(string, string) ([]byte, error)
}

func realServiceEnv() serviceEnv {
	home, _ := os.UserHomeDir()
	return serviceEnv{paths: homeServicePaths(home), controller: systemServiceController{}, terminal: realSessionEnv(),
		goos: runtime.GOOS, domain: "gui/" + strconv.Itoa(os.Getuid()), label: worldServiceLabel,
		wait: 30 * time.Second, client: &http.Client{Timeout: time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
			return errors.New("health redirects are refused")
		}}}
}
func runService(args []string, out, errOut io.Writer) int {
	return runServiceWith(args, out, errOut, realServiceEnv())
}
func serviceMessage(w io.Writer, what, why, next string, code int) int {
	fmt.Fprintf(w, "%s\n  why: %s\n  next: %s\n", what, why, next)
	return code
}
func serviceHash(b []byte) string { h := sha256.Sum256(b); return "sha256:" + hex.EncodeToString(h[:]) }
func serviceFileHash(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return serviceHash(b)
}
func (e serviceEnv) note(path string) {
	if e.observe != nil {
		e.observe(path)
	}
}
func (e serviceEnv) atomic(path string, b []byte, mode os.FileMode) error {
	e.note(filepath.Dir(path))
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	tempDir := filepath.Dir(e.paths.Binary)
	if err := os.MkdirAll(tempDir, 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(tempDir, ".service-*")
	if err != nil {
		return err
	}
	e.note(f.Name())
	defer os.Remove(f.Name())
	if err = f.Chmod(mode); err == nil {
		_, err = f.Write(b)
	}
	if err == nil {
		err = f.Sync()
	}
	closeErr := f.Close()
	if err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	e.note(path)
	return os.Rename(f.Name(), path)
}
func (e serviceEnv) lock() (*os.File, error) {
	e.note(filepath.Dir(e.paths.Lock))
	if err := os.MkdirAll(filepath.Dir(e.paths.Lock), 0700); err != nil {
		return nil, err
	}
	e.note(e.paths.Lock)
	f, err := os.OpenFile(e.paths.Lock, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		f.Close()
		return nil, err
	}
	if err = f.Truncate(0); err == nil {
		_, err = fmt.Fprintf(f, "%d\n", os.Getpid())
	}
	if err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}
func (e serviceEnv) confirm(w io.Writer, change string) bool {
	term, err := e.terminal.openTerminal()
	if err != nil {
		serviceMessage(w, "refusing: no controlling terminal", "this change requires one human act at a terminal", sessionTTYFix, 1)
		return false
	}
	defer term.Close()
	fmt.Fprintf(term, "%s\nProceed? [y/N] ", change)
	if strings.TrimSpace(readTermLine(term)) != "y" {
		serviceMessage(w, "aborted (not confirmed)", "the answer was not y", "run again and type y", 1)
		return false
	}
	return true
}
func serviceXML(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}
func renderServicePlist(label string, p servicePaths, args []string, working string) string {
	var b strings.Builder
	b.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\"><dict>\n")
	fmt.Fprintf(&b, "<key>Label</key><string>%s</string>\n<key>ProgramArguments</key><array>\n", serviceXML(label))
	for _, a := range args {
		fmt.Fprintf(&b, "<string>%s</string>\n", serviceXML(a))
	}
	b.WriteString("</array>\n<key>RunAtLoad</key><true/>\n<key>KeepAlive</key><true/>\n<key>ThrottleInterval</key><integer>30</integer>\n<key>ExitTimeOut</key><integer>30</integer>\n<key>ProcessType</key><string>Standard</string>\n<key>Umask</key><integer>63</integer>\n")
	fmt.Fprintf(&b, "<key>WorkingDirectory</key><string>%s</string>\n", serviceXML(working))
	b.WriteString("<key>EnvironmentVariables</key><dict><key>PATH</key><string>/usr/bin:/bin:/usr/sbin:/sbin</string></dict>\n")
	fmt.Fprintf(&b, "<key>StandardOutPath</key><string>%s</string>\n<key>StandardErrorPath</key><string>%s</string>\n</dict></plist>\n", serviceXML(p.Log), serviceXML(p.Log))
	return b.String()
}
func (e serviceEnv) health(addr string) (daemon.HealthResponse, error) {
	var h daemon.HealthResponse
	res, err := e.client.Get("http://" + addr + "/v1/health")
	if err != nil {
		return h, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return h, fmt.Errorf("health status %d", res.StatusCode)
	}
	err = json.NewDecoder(io.LimitReader(res.Body, 1<<20)).Decode(&h)
	return h, err
}
func serviceTail(path string, n int) string {
	b, _ := os.ReadFile(path)
	lines := strings.Split(strings.TrimSpace(string(b)), "\n")
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.Join(lines, "\n")
}
func (e serviceEnv) rotate() error {
	e.note(filepath.Dir(e.paths.Log))
	if err := os.MkdirAll(filepath.Dir(e.paths.Log), 0700); err != nil {
		return err
	}
	if st, err := os.Stat(e.paths.Log); err == nil && st.Size() > 10<<20 {
		e.note(e.paths.Log + ".1")
		e.note(e.paths.Log)
		if err := os.Rename(e.paths.Log, e.paths.Log+".1"); err != nil {
			return err
		}
	}
	return nil
}
func (e serviceEnv) boot(r serviceRecord, w io.Writer) int {
	expected := serviceFileHash(e.paths.Binary)
	if err := e.rotate(); err != nil {
		return serviceMessage(w, "cannot prepare logs", err.Error(), "check "+e.paths.Log, 2)
	}
	if err := e.controller.Bootstrap(e.domain, e.paths.Plist); err != nil {
		return serviceMessage(w, "could not start World", err.Error(), "inspect "+e.paths.Log+"; service start", 2)
	}
	end := time.Now().Add(e.wait)
	for {
		h, err := e.health(r.Address)
		if err == nil && h.Build.ExecutableSHA256 == expected && h.Build.ExecutableSHA256 != "" {
			return 0
		}
		if !time.Now().Before(end) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	err := e.controller.Bootout(e.domain, e.label)
	fmt.Fprintln(w, serviceTail(e.paths.Log, 20))
	if err != nil {
		return serviceMessage(w, "health timed out; stopping the service also failed", err.Error(), "inspect "+e.paths.Log+" and stop the service manually", 3)
	}
	return serviceMessage(w, "build did not answer /v1/health within 30 s; service stopped", "installed files were kept", "read "+e.paths.Log+"; fix and run service start", 2)
}
func readServiceRecord(p servicePaths) (serviceRecord, error) {
	var r serviceRecord
	b, err := os.ReadFile(p.Record)
	if err != nil {
		return r, err
	}
	if err = json.Unmarshal(b, &r); err != nil {
		return r, err
	}
	if !serviceLoopback(r.Address) {
		return r, errors.New("service address is not loopback")
	}
	if r.Schema != "world/service/v1" || r.Label != worldServiceLabel || r.Plist != p.Plist || r.Binary != p.Binary {
		return r, errors.New("service record does not match the home service paths")
	}
	return r, nil
}
func runServiceWith(args []string, out, errOut io.Writer, e serviceEnv) int {
	if len(args) == 0 {
		fmt.Fprint(errOut, serviceHelp)
		return 1
	}
	verb := args[0]
	if verb == "--help" || verb == "-h" {
		fmt.Fprint(out, serviceHelp)
		return 0
	}
	fs := flag.NewFlagSet("service "+verb, flag.ContinueOnError)
	fs.SetOutput(errOut)
	asJSON := fs.Bool("json", false, "JSON status")
	verbose := fs.Bool("verbose", false, "full supervisor status")
	pause := fs.Bool("pause", false, "interrupt active episodes during upgrade")
	def := defaultWorldPaths()
	bind := fs.String("bind", daemon.DefaultBind, "loopback address")
	db := fs.String("db", def.db, "store")
	creds := fs.String("credentials-db", "", "credentials")
	root := fs.String("workspace-root", def.workspaceRoot, "workspaces")
	interp := fs.String("interpreter-dir", def.interpDir, "interpreter pin")
	tools := fs.String("tools-dir", def.toolsDir, "tool pin")
	binaryDefault, _ := os.Executable()
	binary := fs.String("binary", binaryDefault, "source binary")
	replace := fs.Bool("replace", false, "replace plist")
	extra := map[string]*string{}
	for _, name := range []string{"exec-sandbox", "exec-profiles-dir", "exec-node", "examples-dir", "surface-index", "surface-follow"} {
		extra[name] = fs.String(name, "", "absolute path")
	}
	if code, done := parseVerbFlags(fs, args[1:], serviceHelp, out, errOut); done {
		return code
	}
	if fs.NArg() != 0 {
		return serviceMessage(errOut, "unexpected argument", "service verbs take flags only", serviceHelp, 1)
	}
	if verb != "install" && verb != "status" && verb != "start" && verb != "stop" && verb != "uninstall" && verb != "upgrade" {
		fmt.Fprint(errOut, serviceHelp)
		return 1
	}
	badFlag := false
	fs.Visit(func(f *flag.Flag) {
		allowed := verb == "install" && f.Name != "json" && f.Name != "verbose" && f.Name != "pause" ||
			verb == "status" && (f.Name == "json" || f.Name == "verbose") ||
			verb == "upgrade" && (f.Name == "binary" || f.Name == "pause")
		if !allowed {
			badFlag = true
		}
	})
	if badFlag || (*asJSON && verb != "status") {
		return serviceMessage(errOut, "invalid service flags", "install flags apply only to install", serviceHelp, 1)
	}
	if strings.HasPrefix(e.label, "dev.ailang") || e.label != worldServiceLabel {
		return serviceMessage(errOut, "refusing service label", "only World's own label is permitted", "use "+worldServiceLabel, 1)
	}
	if verb != "status" {
		for _, path := range []string{e.paths.Plist, e.paths.Binary, e.paths.Record, e.paths.Lock, e.paths.Log} {
			if err := serviceSafePath(e.paths.Home, path); err != nil {
				return serviceMessage(errOut, "unsafe service path", err.Error(), "remove symlinks from the home service paths", 1)
			}
		}
	}
	if verb == "upgrade" {
		return e.upgrade(*binary, *pause, out, errOut)
	}
	if verb == "install" {
		if *extra["surface-index"] != "" && *extra["surface-follow"] != "" {
			return serviceMessage(errOut, "surface follow refused", "--surface-follow and --surface-index are mutually exclusive", "ailang-worldd serve --help", 1)
		}
		if e.goos != "darwin" {
			return serviceMessage(errOut, "service install: not supported on "+e.goos+" yet", "the systemd unit is M6", "run serve by hand (QUICKSTART §1)", 1)
		}
		// Use the existing session terminal seam before any filesystem mutation.
		term, err := e.terminal.openTerminal()
		if err != nil {
			return serviceMessage(errOut, "refusing: no controlling terminal", "install requires an attended terminal", sessionTTYFix, 1)
		}
		defer term.Close()
		for _, path := range []string{e.paths.Home, e.paths.Plist, e.paths.Binary, e.paths.Record, e.paths.Lock, e.paths.Log, *db, *root, *interp, *tools, *binary} {
			if !filepath.IsAbs(path) {
				return serviceMessage(errOut, "refusing relative path", path, "use absolute paths", 1)
			}
		}
		if *creds == "" {
			*creds = filepath.Join(filepath.Dir(*db), "credentials.db")
		}
		if !filepath.IsAbs(*creds) {
			return serviceMessage(errOut, "refusing relative credentials path", *creds, "use an absolute path", 1)
		}
		for _, dir := range []string{filepath.Dir(*db), *root} {
			if st, err := os.Lstat(dir); err != nil || !st.IsDir() {
				return serviceMessage(errOut, "service install: "+dir+" does not exist", "the store directory and the workspace root must exist before install; ailang-worldd setup creates both", "ailang-worldd setup, then ailang-worldd service install", 1)
			}
		}
		if within(*root, filepath.Dir(*db)) || within(*root, e.paths.Binary) || within(filepath.Dir(*db), e.paths.Binary) {
			return serviceMessage(errOut, "unsafe service layout", "binary and store must be outside agent workspaces", "use the default home layout", 1)
		}
		roles := pinRoles(*interp, *tools)
		for _, role := range roles {
			pin, err := lookupPin(role.release, setupPlatform)
			sum, hashErr := pinfetch.HashFile(role.dest)
			st, statErr := os.Lstat(role.dest)
			if err != nil || hashErr != nil || statErr != nil || !st.Mode().IsRegular() || sum != pin.BinarySHA256 {
				return serviceMessage(errOut, "the "+role.name+" pin is missing or drifted at "+role.dest, "serve archives both pins at startup", "ailang-worldd setup", 1)
			}
		}
		host, port, err := net.SplitHostPort(*bind)
		ip := net.ParseIP(host)
		pn, pe := strconv.Atoi(port)
		if err != nil || pe != nil || pn < 1 || pn > 65535 || (host != "localhost" && (ip == nil || !ip.IsLoopback())) {
			return serviceMessage(errOut, "invalid --bind", "only a loopback host and nonzero port are allowed", "use 127.0.0.1:7644", 1)
		}
		err = servicePortFree(e, *bind)
		if err != nil {
			return serviceMessage(errOut, *bind+" is in use", "another process listens there", "lsof -nP -iTCP:"+port+" -sTCP:LISTEN; stop it or choose --bind another port", 1)
		}
		held, err := store.WriterLockHeld(*db)
		if err != nil || held {
			return serviceMessage(errOut, "another process holds the writer lock on "+*db, "serve would fail at store-open", "stop that daemon; run doctor", 1)
		}
		key, err := e.controller.Getenv("AILANG_REGISTRY_API_KEY")
		if err != nil || key != "" {
			return serviceMessage(errOut, "launchd's environment carries AILANG_REGISTRY_API_KEY", "World must not inherit publish authority (or the environment could not be checked)", serviceRegistryKeyFix, 1)
		}
		if _, err := e.controller.Print(e.domain, e.label); err == nil {
			return serviceMessage(errOut, "World service is already loaded", "install cannot replace a loaded job", "service stop first", 1)
		} else if !errors.Is(err, errServiceNotLoaded) {
			return serviceMessage(errOut, "cannot verify the service is unloaded", err.Error(), "inspect service status and retry", 1)
		}
		if err := serviceLockFree(e.paths.Lock); err != nil {
			return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
		}
		old, oldErr := os.ReadFile(e.paths.Plist)
		if oldErr == nil && !*replace {
			return serviceMessage(errOut, "World plist already exists", e.paths.Plist, "service install --replace", 1)
		}
		if oldErr != nil && !errors.Is(oldErr, os.ErrNotExist) {
			return serviceMessage(errOut, "cannot read plist", oldErr.Error(), "check permissions", 1)
		}
		for _, p := range []string{e.paths.Binary, e.paths.Record, e.paths.Lock} {
			if st, err := os.Lstat(p); err == nil && !st.Mode().IsRegular() {
				return serviceMessage(errOut, "refusing nonregular service file", p, "restore a regular file", 1)
			}
		}
		source, err := os.ReadFile(*binary)
		st, se := os.Lstat(*binary)
		if err != nil || se != nil || !st.Mode().IsRegular() {
			return serviceMessage(errOut, "invalid source binary", *binary, "choose a regular binary file", 1)
		}
		argv := []string{e.paths.Binary, "serve", "--db", *db, "--credentials-db", *creds, "--bind", *bind, "--ailang-bin", roles[0].dest, "--workspace-root", *root, "--tool-ailang-bin", roles[1].dest}
		for _, name := range []string{"exec-sandbox", "exec-profiles-dir", "exec-node", "examples-dir", "surface-index", "surface-follow"} {
			if *extra[name] != "" {
				if !filepath.IsAbs(*extra[name]) {
					return serviceMessage(errOut, "relative --"+name, "all service paths must be absolute", "supply an absolute path", 1)
				}
				argv = append(argv, "--"+name, *extra[name])
			}
		}
		plist := renderServicePlist(e.label, e.paths, argv, filepath.Dir(*db))
		r := serviceRecord{Schema: "world/service/v1", Label: e.label, Plist: e.paths.Plist, PlistSHA256: serviceHash([]byte(plist)), PlistText: plist, Binary: e.paths.Binary, BinarySHA256: serviceHash(source), Address: *bind, Store: *db, Credentials: *creds, Workspace: *root, Log: e.paths.Log, InstalledAt: e.terminal.now(), Upgrades: []serviceUpgrade{}}
		record, _ := json.MarshalIndent(r, "", "  ")
		if len(old) > 0 {
			fmt.Fprint(term, servicePlistDiff(string(old), plist))
		}
		fmt.Fprintf(term, "%s\n%s → %s (%s)\n%s\n%s\nlog: %s\n", plist, *binary, e.paths.Binary, r.BinarySHA256, e.paths.Record, record, e.paths.Log)
		var warnings bytes.Buffer
		report := doctorReport{w: &warnings}
		doctorStore(&report, *db, "")
		fmt.Fprint(term, warnings.String())
		if strings.Contains(warnings.String(), "no transition registry") || strings.Contains(warnings.String(), "does not exist yet") {
			serviceMessage(term, "warning: the home store has no transition registry", "no tool will serve until transitions are published", "service stop; world-publish transitions --store "+*db+" …; service start (QUICKSTART §9)", 0)
		}
		if strings.Contains(warnings.String(), "no world head") || strings.Contains(warnings.String(), "does not exist yet") {
			serviceMessage(term, "warning: the home store has no world head", "tool calls need a genesis world", "commit genesis after publishing transitions (QUICKSTART §9)", 0)
		}
		fmt.Fprintf(term, "Publish transitions with world-publish transitions --store %s …; then commit genesis (QUICKSTART §9).\nProceed? [y/N] ", *db)
		if strings.TrimSpace(readTermLine(term)) != "y" {
			return serviceMessage(errOut, "aborted (not confirmed)", "the answer was not y", "run again and type y", 1)
		}
		lock, err := e.lock()
		if err != nil {
			return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
		}
		defer lock.Close()
		if err = e.atomic(e.paths.Binary, source, 0755); err == nil {
			err = e.atomic(e.paths.Plist, []byte(plist), 0644)
		}
		if err == nil {
			err = e.atomic(e.paths.Record, record, 0600)
		}
		if err != nil {
			return serviceMessage(errOut, "service files could not be installed", err.Error(), "inspect the retained files and retry", 3)
		}
		code := e.boot(r, errOut)
		if code == 0 {
			fmt.Fprintf(out, "World service installed at %s\nexport PATH=%s:$PATH\n", r.Address, filepath.Dir(e.paths.Binary))
		}
		return code
	}
	r, err := readServiceRecord(e.paths)
	if err != nil || r.Uninstalled {
		code := 1
		if verb == "status" {
			code = 10
			if *asJSON {
				_ = json.NewEncoder(out).Encode(map[string]any{"state": "not installed", "installed": false, "next": "ailang-worldd service install"})
				return code
			}
		}
		return serviceMessage(errOut, "no World service is installed", "no valid service.json and plist at "+e.paths.Plist, "ailang-worldd service install", code)
	}
	if verb == "status" {
		return e.status(r, *asJSON, *verbose, out, errOut)
	}
	if err := serviceLockFree(e.paths.Lock); err != nil {
		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
	}
	if e.goos != "darwin" {
		return serviceMessage(errOut, "service: not supported on "+e.goos+" yet", "the systemd unit is M6", "run serve by hand (QUICKSTART §1)", 1)
	}
	plist, err := os.ReadFile(e.paths.Plist)
	if err != nil {
		return serviceMessage(errOut, "cannot read service plist", err.Error(), "service install --replace", 1)
	}
	if verb == "start" {
		if err := serviceLockFree(e.paths.Lock); err != nil {
			return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
		}
		if st, err := os.Lstat(e.paths.Binary); err != nil || !st.Mode().IsRegular() {
			return serviceMessage(errOut, "installed binary is not a regular file", e.paths.Binary, "service install --replace", 1)
		}
	}
	if verb == "start" {
		if err := validateServicePlist(plist); err != nil {
			return serviceMessage(errOut, "invalid service plist", err.Error(), "service install --replace", 1)
		}
	}
	if verb == "start" && serviceHash(plist) != r.PlistSHA256 {
		if !e.confirm(errOut, "the plist differs from the one installed\n"+servicePlistDiff(r.PlistText, string(plist))) {
			return 1
		}
	}
	if verb == "stop" || verb == "uninstall" {
		if _, err := e.controller.Print(e.domain, e.label); err != nil && !errors.Is(err, errServiceNotLoaded) {
			return serviceMessage(errOut, "cannot inspect the service", err.Error(), "inspect service status and retry", 1)
		}
		h, _ := e.health(r.Address)
		active, _ := json.Marshal(h.Active)
		change := fmt.Sprintf("Stop %s at %s. Active: %s. Plist stays: %s", e.label, r.Address, active, e.paths.Plist)
		if verb == "uninstall" {
			change = fmt.Sprintf("Stop %s; remove only %s; mark %s uninstalled.\nKeep binary %s and %s.prev, store %s, credentials %s, sessions %s, worktrees %s, logs %s. Delete data by hand only.", e.label, e.paths.Plist, e.paths.Record, e.paths.Binary, e.paths.Binary, r.Store, r.Credentials, filepath.Join(filepath.Dir(r.Store), "sessions"), r.Workspace, e.paths.Log)
		}
		if !e.confirm(errOut, change) {
			return 1
		}
		fmt.Fprintln(out, change)
	}
	lock, err := e.lock()
	if err != nil {
		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
	}
	defer lock.Close()
	if verb == "start" {
		return e.boot(r, errOut)
	}
	if _, err := e.controller.Print(e.domain, e.label); err == nil {
		if err = e.controller.Bootout(e.domain, e.label); err != nil {
			return serviceMessage(errOut, "could not stop World", err.Error(), "inspect service status and retry", 3)
		}
	}
	if verb == "uninstall" {
		e.note(e.paths.Plist)
		if err = os.Remove(e.paths.Plist); err != nil && !errors.Is(err, os.ErrNotExist) {
			return serviceMessage(errOut, "cannot remove plist", err.Error(), "check permissions and retry", 3)
		}
		r.Uninstalled = true
		b, _ := json.MarshalIndent(r, "", "  ")
		if err = e.atomic(e.paths.Record, b, 0600); err != nil {
			return serviceMessage(errOut, "cannot mark service uninstalled", err.Error(), "inspect "+e.paths.Record, 3)
		}
	}
	return 0
}
func (e serviceEnv) status(r serviceRecord, asJSON, verbose bool, out, errOut io.Writer) int {
	print, printErr := e.controller.Print(e.domain, e.label)
	h, healthErr := e.health(r.Address)
	installed := serviceFileHash(e.paths.Binary)
	edited := serviceFileHash(e.paths.Plist) != r.PlistSHA256
	uptime := int64(0)
	if h.StartedAt > 0 {
		uptime = time.Now().Unix() - h.StartedAt
	}
	state := "serving"
	code := 0
	if healthErr != nil || serviceFileHash(e.paths.Plist) == "" {
		state = "not serving"
		code = 10
	} else if h.Build.ExecutableSHA256 != installed || installed == "" {
		state = "restart pending"
		code = 10
	}
	logSize := int64(0)
	if st, err := os.Stat(e.paths.Log); err == nil {
		logSize = st.Size()
	}
	report := struct {
		PID        int64                 `json:"pid"`
		Uptime     int64                 `json:"uptime_seconds"`
		Runs       int64                 `json:"runs"`
		LastExit   int64                 `json:"last_exit_status"`
		State      string                `json:"state"`
		Edited     bool                  `json:"plist_edited"`
		Loaded     bool                  `json:"loaded"`
		Supervisor string                `json:"supervisor_state"`
		Record     serviceRecord         `json:"service"`
		Health     daemon.HealthResponse `json:"health"`
		LogBytes   int64                 `json:"log_bytes"`
	}{servicePrintNumber(print, "pid"), uptime, servicePrintNumber(print, "runs"), servicePrintNumber(print, "last exit code"), state, edited, printErr == nil, print, r, h, logSize}
	if asJSON {
		_ = json.NewEncoder(out).Encode(report)
	} else {
		fmt.Fprintf(out, "World service: %s\n", state)
		fmt.Fprintf(out, "pid: %d; uptime: %ds; runs: %d; last exit: %d\n", report.PID, uptime, report.Runs, report.LastExit)
		fmt.Fprintf(out, "address: %s\ninstalled: %s\nrunning: %s\nstore: %s\nlog: %s (%d bytes)\n", r.Address, installed, h.Build.ExecutableSHA256, r.Store, e.paths.Log, logSize)
		if verbose {
			fmt.Fprintln(out, print)
		}
		if edited {
			serviceMessage(out, "! the plist differs from the one installed", "its sha256 differs from service.json", "service start shows the diff, or service install --replace", 0)
		}
		if code == 10 {
			serviceMessage(errOut, "the World service is not serving and matching at "+r.Address, fmt.Sprintf("loaded: %t; last exit: %d; %s", report.Loaded, report.LastExit, serviceTail(e.paths.Log, 10)), "ailang-worldd service start", 10)
		}
		if logSize > 100<<20 {
			fmt.Fprintln(out, "warning: log exceeds 100 MiB; service stop then service start rotates it")
		}
	}
	return code
}

func servicePortFree(e serviceEnv, addr string) error {
	if e.portFree != nil {
		return e.portFree(addr)
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return l.Close()
}
func serviceLockFree(path string) error {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	defer f.Close()
	if err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		pid, _ := io.ReadAll(io.LimitReader(f, 64))
		return fmt.Errorf("a service change is in progress (pid %s): %w", strings.TrimSpace(string(pid)), err)
	}
	return syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
}

func serviceSafePath(home, path string) error {
	if !filepath.IsAbs(home) || !within(home, path) {
		return fmt.Errorf("%s is outside home", path)
	}
	for p := path; p != home; p = filepath.Dir(p) {
		st, err := os.Lstat(p)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return err
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink", p)
		}
	}
	return nil
}

func servicePlistDiff(old, current string) string {
	return "--- approved plist\n+++ current plist\n@@\n-" + strings.ReplaceAll(strings.TrimSuffix(old, "\n"), "\n", "\n-") + "\n+" + strings.ReplaceAll(strings.TrimSuffix(current, "\n"), "\n", "\n+") + "\n"
}

func serviceLoopback(addr string) bool {
	host, port, err := net.SplitHostPort(addr)
	n, parseErr := strconv.Atoi(port)
	ip := net.ParseIP(host)
	return err == nil && parseErr == nil && n > 0 && n <= 65535 && (host == "localhost" || ip != nil && ip.IsLoopback())
}

// Parse the argument array rather than accepting an edited non-loopback bind.
func validateServicePlist(data []byte) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	key := ""
	inArguments := false
	var args []string
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "key" {
				if err := decoder.DecodeElement(&key, &token); err != nil {
					return err
				}
			} else if token.Name.Local == "array" && key == "ProgramArguments" {
				inArguments = true
			} else if token.Name.Local == "string" && inArguments {
				var arg string
				if err := decoder.DecodeElement(&arg, &token); err != nil {
					return err
				}
				args = append(args, arg)
			}
		case xml.EndElement:
			if token.Name.Local == "array" {
				inArguments = false
			}
		}
	}
	for i, arg := range args {
		if arg == "--bind" && i+1 < len(args) && serviceLoopback(args[i+1]) {
			return nil
		}
	}
	return errors.New("ProgramArguments must explicitly bind to loopback")
}

func servicePrintNumber(text, key string) int64 {
	for _, line := range strings.Split(text, "\n") {
		parts := strings.SplitN(strings.TrimSpace(line), " = ", 2)
		if len(parts) == 2 && parts[0] == key {
			n, _ := strconv.ParseInt(parts[1], 10, 64)
			return n
		}
	}
	return 0
}
