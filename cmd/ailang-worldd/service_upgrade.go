package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

type serviceUpgrade struct {
	From    string `json:"from"`
	To      string `json:"to"`
	At      int64  `json:"at"`
	Outcome string `json:"outcome"`
}

type serviceVersion struct {
	Build              daemon.BuildIdentity `json:"build"`
	Platform           string               `json:"platform"`
	ToolRelease        string               `json:"tool_release"`
	InterpreterRelease string               `json:"interpreter_release"`
}

// The environment is deliberately exact, rather than inherited and filtered.
func serviceVersionOutput(binary, home string) ([]byte, error) {
	ctx, cancel := budgetContext(5 * time.Second)
	defer cancel()
	return toolprobe.OutputEnv(ctx, binary, "", []string{"PATH=/usr/bin:/bin:/usr/sbin:/sbin", "HOME=" + home}, "version", "--json")
}

func (e serviceEnv) version(binary, sum string) (serviceVersion, error) {
	probe := e.probe
	if probe == nil {
		probe = serviceVersionOutput
	}
	absolute, err := filepath.Abs(binary)
	if err != nil {
		return serviceVersion{}, err
	}
	b, err := probe(absolute, e.paths.Home)
	var v serviceVersion
	if err != nil {
		return v, err
	}
	d := json.NewDecoder(bytes.NewReader(b))
	if err = d.Decode(&v); err != nil {
		return v, err
	}
	var trailing any
	if d.Decode(&trailing) != io.EOF || v.Platform != runtime.GOOS+"/"+runtime.GOARCH ||
		v.Build.ExecutableSHA256 != sum || v.ToolRelease == "" || v.InterpreterRelease == "" {
		return v, errors.New("invalid version identity")
	}
	return v, nil
}

func serviceShort(sum string) string {
	sum = strings.TrimPrefix(sum, "sha256:")
	if len(sum) > 12 {
		return sum[:12]
	}
	return sum
}

// Health exposes recent episodes and a total in-flight count, but does not
// attribute that count to individual episodes. Do not invent that attribution.
func serviceActiveText(a daemon.HealthActive, now int64) string {
	var episodes []string
	for _, s := range a.Sessions {
		age := now - s.LastCallUnix
		if age < 0 {
			age = 0
		}
		episodes = append(episodes, fmt.Sprintf("%s (last call %d s ago)", s.Episode, age))
	}
	text := fmt.Sprintf("%d episodes are working: %s", len(episodes), strings.Join(episodes, ", "))
	if len(episodes) == 0 {
		text = "calls are working"
	}
	if a.InFlight > 0 {
		text += fmt.Sprintf("; %d calls in flight", a.InFlight)
	}
	return text
}

func (e serviceEnv) waitForBuild(addr, sum string) bool {
	end := time.Now().Add(e.wait)
	for {
		h, err := e.health(addr)
		if err == nil && h.Build.ExecutableSHA256 == sum {
			return true
		}
		remaining := time.Until(end)
		if remaining <= 0 {
			return false
		}
		if remaining > 100*time.Millisecond {
			remaining = 100 * time.Millisecond
		}
		time.Sleep(remaining)
	}
}

func (e serviceEnv) saveUpgrade(r serviceRecord, from, to, outcome string) error {
	r.Upgrades = append(r.Upgrades, serviceUpgrade{From: from, To: to, At: e.terminal.now(), Outcome: outcome})
	if len(r.Upgrades) > 10 {
		r.Upgrades = r.Upgrades[len(r.Upgrades)-10:]
	}
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return e.atomic(e.paths.Record, b, 0600)
}

func (e serviceEnv) upgrade(candidate string, pause bool, out, errOut io.Writer) int {
	refuse := func(what, next string) int {
		return serviceMessage(errOut, what, "nothing was swapped", next, 1)
	}
	st, err := os.Lstat(candidate)
	if err != nil || !st.Mode().IsRegular() {
		return refuse(candidate+" is not a regular file", "choose a regular ailang-worldd binary with --binary")
	}
	source, err := os.ReadFile(candidate)
	if err != nil {
		return refuse("cannot read "+candidate, "check the candidate's permissions")
	}
	to := serviceHash(source)
	v, err := e.version(candidate, to)
	if err != nil {
		return refuse(candidate+" is not an ailang-worldd for "+runtime.GOOS+"/"+runtime.GOARCH, "ailang-worldd setup with that binary")
	}
	r, err := readServiceRecord(e.paths)
	if err != nil || r.Uninstalled {
		return serviceMessage(errOut, "no World service is installed", "no valid service.json and plist at "+e.paths.Plist, "ailang-worldd service install", 1)
	}
	if serviceFileHash(e.paths.Plist) != r.PlistSHA256 {
		return serviceMessage(errOut, "the plist differs from the one installed", "its sha256 differs from service.json", "service start shows the diff, or service install --replace", 1)
	}
	from := serviceFileHash(e.paths.Binary)
	st, err = os.Lstat(e.paths.Binary)
	if err != nil || !st.Mode().IsRegular() || from == "" || from != r.BinarySHA256 {
		return refuse("installed file sha256 differs from service.json", "ailang-worldd service status")
	}
	h, err := e.health(r.Address)
	if err != nil {
		return refuse("installed service does not answer /v1/health", "ailang-worldd service status")
	}
	if h.Build.ExecutableSHA256 != from {
		return refuse("restart pending: health build.executable_sha256 differs from the installed file", "ailang-worldd service status")
	}
	if to == from {
		fmt.Fprintf(out, "already installed (%s)\n", to)
		return 0
	}
	installed, err := e.version(e.paths.Binary, from)
	if err != nil {
		return refuse("cannot read the installed build's version identity", "ailang-worldd service status")
	}
	if v.ToolRelease != installed.ToolRelease {
		return refuse(fmt.Sprintf("it pins tool AILANG %s (installed %s)", v.ToolRelease, installed.ToolRelease), candidate+" setup (ailang-worldd setup with that binary)")
	}
	if v.InterpreterRelease != installed.InterpreterRelease {
		return refuse(fmt.Sprintf("it serves interpreter %s, which the store's epoch registry does not admit (installed %s)", v.InterpreterRelease, installed.InterpreterRelease), "world-publish epoch nominate (row 139)")
	}
	active := h.Active.InFlight > 0 || len(h.Active.Sessions) > 0
	episodes := serviceActiveText(h.Active, e.terminal.now())
	if active && !pause {
		return serviceMessage(errOut, episodes, "a restart fails their calls for a few seconds", "wait, session revoke --episode <ep>, or service upgrade --pause", 1)
	}
	if err := serviceLockFree(e.paths.Lock); err != nil {
		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
	}
	if e.goos != "darwin" {
		return refuse("service: not supported on "+e.goos+" yet", "run serve by hand (QUICKSTART §1)")
	}
	// Sibling paths must not redirect writes outside the service home.
	for _, suffix := range []string{".new", ".prev", ".failed"} {
		if err := serviceSafePath(e.paths.Home, e.paths.Binary+suffix); err != nil {
			return refuse("unsafe service path: "+err.Error(), "remove symlinks from the home service paths")
		}
	}
	recordHash := serviceFileHash(e.paths.Record)
	change := fmt.Sprintf("from %s → to %s\ntool AILANG: %s (unchanged)\ninterpreter: %s (unchanged)\nplist: unchanged\nexpected downtime: 10 s drain plus a restart", serviceShort(from), serviceShort(to), v.ToolRelease, v.InterpreterRelease)
	if active {
		change += "\n--pause will interrupt: " + episodes
	}
	if !e.confirm(errOut, change) {
		return 1
	}
	lock, err := e.lock()
	if err != nil {
		return serviceMessage(errOut, "a service change is in progress", err.Error(), "wait and retry", 1)
	}
	defer lock.Close()
	// Revalidate the approved inputs after acquiring the lock. Another attended
	// change may have completed while the confirmation prompt was open.
	if serviceFileHash(e.paths.Binary) != from || serviceFileHash(e.paths.Plist) != r.PlistSHA256 || serviceFileHash(e.paths.Record) != recordHash {
		return serviceMessage(errOut, "service changed during confirmation", "no binary was swapped", "ailang-worldd service status; retry upgrade", 1)
	}
	fresh, previous := e.paths.Binary+".new", e.paths.Binary+".prev"
	if err = e.atomic(fresh, source, 0755); err != nil {
		return serviceMessage(errOut, "cannot stage upgrade", err.Error(), "inspect "+fresh+"; retry upgrade", 2)
	}
	e.note(previous)
	if err = os.Remove(previous); err != nil && !errors.Is(err, os.ErrNotExist) {
		return serviceMessage(errOut, "cannot remove previous build", err.Error(), "inspect "+previous, 2)
	}
	e.note(previous)
	if err = os.Link(e.paths.Binary, previous); err != nil {
		return serviceMessage(errOut, "cannot retain previous build", err.Error(), "inspect "+previous, 2)
	}
	e.note(fresh)
	e.note(e.paths.Binary)
	if err = os.Rename(fresh, e.paths.Binary); err != nil {
		return serviceMessage(errOut, "cannot swap build", err.Error(), "inspect "+fresh, 2)
	}
	err = e.rotate()
	if err == nil {
		err = e.controller.Kickstart(e.domain, e.label)
	}
	if err == nil && e.waitForBuild(r.Address, to) {
		updated := r
		updated.BinarySHA256, updated.PreviousSHA256 = to, from
		if err = e.saveUpgrade(updated, from, to, "ok"); err == nil {
			fmt.Fprintf(out, "upgraded %s → %s\n", serviceShort(from), serviceShort(to))
			return 0
		}
		// A healthy but unrecorded swap is not a successful upgrade.
	}
	return e.rollbackUpgrade(r, from, to, candidate, errOut)
}

func (e serviceEnv) rollbackUpgrade(r serviceRecord, from, to, candidate string, w io.Writer) int {
	tail := serviceTail(e.paths.Log, 20)
	failed, previous := e.paths.Binary+".failed", e.paths.Binary+".prev"
	restored := false
	e.note(failed)
	err := os.Remove(failed)
	if err == nil || errors.Is(err, os.ErrNotExist) {
		e.note(failed)
		err = os.Link(e.paths.Binary, failed)
	}
	if err == nil {
		e.note(previous)
		e.note(e.paths.Binary)
		err = os.Rename(previous, e.paths.Binary)
		restored = err == nil
	}
	if err == nil {
		err = e.controller.Kickstart(e.domain, e.label)
	}
	if err == nil && e.waitForBuild(r.Address, from) {
		recordErr := e.saveUpgrade(r, from, to, "rolled back")
		fmt.Fprintln(w, tail)
		if recordErr != nil {
			return serviceMessage(w, "old build is serving but upgrade history could not be saved", recordErr.Error(), "inspect "+e.paths.Record, 3)
		}
		return serviceMessage(w, fmt.Sprintf("build %s did not answer /v1/health within 30 s; rolled back to %s, which is serving", serviceShort(to), serviceShort(from)), "new build log: "+strings.ReplaceAll(tail, "\n", "; "), "read "+e.paths.Log+"; fix and re-run service upgrade --binary "+candidate, 2)
	}
	stopErr := e.controller.Bootout(e.domain, e.label)
	historyErr := e.saveUpgrade(r, from, to, "rollback failed")
	why := fmt.Sprintf("both builds kept: %s (old) and %s (new)", e.paths.Binary, failed)
	next := e.paths.Binary + " service start"
	if !restored {
		why = fmt.Sprintf("swap recovery incomplete: %s (new) and %s (old)", e.paths.Binary, previous)
		next = "restore " + previous + " over " + e.paths.Binary + ", then " + next
	}
	if err != nil {
		why += "; recovery error: " + err.Error()
	}
	if stopErr != nil {
		why += "; stop failed: " + stopErr.Error()
	}
	if historyErr != nil {
		why += "; history: " + historyErr.Error()
	}
	what := "rollback failed too: the service is stopped"
	if stopErr != nil {
		what = "rollback failed too; could not stop the service"
	}
	return serviceMessage(w, what, why, next, 3)
}
