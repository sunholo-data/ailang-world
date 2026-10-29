package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/toolprobe"
)

func doctorClass(name string) (string, int) {
	switch name {
	case "daemon":
		return "daemon", 10
	case "session":
		return "session", 11
	case "host":
		return "host", 12
	case "exec tools":
		return "exec", 13
	default:
		return "other", 1
	}
}

// doctorSession never writes credentials or prints their token or hash.
// Revocation deletes the mapping; an absent mapping is indistinguishable from unknown.
func doctorSession(r *doctorReport, db, creds, file string, now time.Time) string {
	fix := "start a new session: ailang-worldd session new <episode> … (or episode start)"
	info, err := os.Stat(file)
	if err != nil {
		r.add(markFail, "session", file+": session file is missing or unreadable", fix)
		return ""
	}
	if !info.Mode().IsRegular() || info.Mode().Perm() != 0600 {
		r.add(markFail, "session", file+": session file must be a regular file with mode 0600", "chmod 0600 "+shellQuote(file))
		return ""
	}
	raw, err := os.ReadFile(file)
	if err != nil {
		r.add(markFail, "session", file+": cannot read session file", "make "+file+" readable by its owner")
		return ""
	}
	id := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.TrimSpace(string(raw)))))
	credPath, err := credentialsFor(db, creds)
	if err != nil {
		r.add(markFail, "session", "cannot locate the credentials store", "pass --db <world store> and --credentials-db <credentials store>")
		return ""
	}
	ctx, cancel := budgetContext(doctorStoreBudget)
	defer cancel()
	// Indexed lookup avoids session list's 500-row display limit.
	ro, err := store.OpenCredentialsReadOnly(credPath)
	var row store.SessionRow
	found := false
	if err == nil {
		defer ro.Close()
		row, found, err = ro.ResolveSession(ctx, id)
	} else {
		// Match session list's read-only compatibility path for schema 4.
		var rows []store.SessionRow
		rows, err = listSessionRows(ctx, db, credPath, "", io.Discard)
		for _, candidate := range rows {
			if candidate.CredentialID == id {
				row, found = candidate, true
				break
			}
		}
	}
	if err != nil {
		r.add(markFail, "session", "cannot read the credentials store "+credPath, "pass --db and --credentials-db naming the store used to mint this session")
		return ""
	}
	if !found {
		r.add(markFail, "session", "session is revoked or unknown", fix)
		return ""
	}
	fix = "start a new session: ailang-worldd session new " + row.EpisodeID + " … (or episode start)"
	if now.Unix() >= row.ExpiresAt {
		r.add(markFail, "session", row.EpisodeID+"'s session expired at "+time.Unix(row.ExpiresAt, 0).In(time.Local).Format("15:04 MST"), fix)
		return row.Project
	}
	left := time.Unix(row.ExpiresAt, 0).Sub(now)
	detail := fmt.Sprintf("%s: %d h %d m left", row.EpisodeID, int(left.Hours()), int(left.Minutes())%60)
	var grants []broker.Capability
	if json.Unmarshal([]byte(row.GrantsJSON), &grants) == nil && len(grants) > 0 {
		detail += "; grants " + grantList(grants)
	}
	r.add(markOK, "session", detail, "")
	return row.Project
}

func doctorHost(r *doctorReport, host, dir string) {
	if absolute, err := filepath.Abs(dir); err == nil {
		dir = absolute
	}
	if host != "claude" && host != "codex" && host != "pi" {
		r.add(markFail, "host", "unknown host "+host, "pass --host claude, --host codex or --host pi")
		return
	}
	bin, err := exec.LookPath(host)
	if err != nil {
		r.add(markFail, "host", host+" is missing from PATH", "install "+host+" and add its binary directory to PATH")
	} else {
		// AC10: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential
		// drives toolprobe's shared site; TestProductionContextRoots stays unchanged.
		ctx, cancel := budgetContext(5 * time.Second)
		version, e := toolprobe.OutputEnv(ctx, bin, "", []string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}, "--version")
		cancel()
		if e != nil {
			r.add(markFail, "host", host+" at "+bin+": --version failed or timed out", "repair or reinstall "+host+", then run "+host+" --version")
		} else {
			r.add(markOK, "host", host+" at "+bin+": "+firstLine(string(version)), "")
		}
	}
	data, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
	var manifest hostsManifest
	fix := "ailang-worldd hosts init --out-dir " + shellQuote(dir) + " --replace"
	if err != nil || json.Unmarshal(data, &manifest) != nil || manifest.Schema != "world/hosts/v1" {
		r.add(markWarn, "host", "launch manifest missing or invalid in "+dir+"; launch would be observed, not confined", fix)
		return
	}
	// Require the entire generated set, not merely the entries a modified manifest lists.
	expected := hostFiles(dir, manifest.Addr)
	names := make([]string, 0, len(expected))
	for name := range expected {
		names = append(names, name)
	}
	sort.Strings(names)
	drift := false
	if host == "codex" {
		config, _ := os.ReadFile(filepath.Join(dir, "codex", "config.toml"))
		if !strings.Contains(string(config), "approval_policy = \"never\"\n") {
			drift = true
			r.add(markWarn, "host", "codex/config.toml: approval_policy = \"never\" is missing; launch config drift", fix)
		}
	}
	for _, name := range names {
		b, e := os.ReadFile(filepath.Join(dir, name))
		sum := fmt.Sprintf("sha256:%x", sha256.Sum256(b))
		expectedSum := fmt.Sprintf("sha256:%x", sha256.Sum256(expected[name]))
		if e != nil || manifest.Files[name] != sum || sum != expectedSum {
			drift = true
			r.add(markWarn, "host", name+": launch config missing or digest differs; launch would be observed, not confined", fix)
		}
	}
	if !drift {
		r.add(markOK, "host", "launch files match hosts.json in "+dir, "")
	}
}

// resolveDoctorTool uses only the profile's PATH; it never executes the binary.
func resolveDoctorTool(name string, paths []string) string {
	candidates := []string{name}
	if !filepath.IsAbs(name) {
		if strings.ContainsRune(name, filepath.Separator) {
			return ""
		}
		candidates = nil
		for _, dir := range paths {
			candidates = append(candidates, filepath.Join(dir, name))
		}
	}
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0 {
			p, err := filepath.EvalSymlinks(candidate)
			if err == nil {
				return p
			}
		}
	}
	return ""
}

func doctorExecTools(r *doctorReport, profile string, startup map[string]string) {
	paths := filepath.SplitList(os.Getenv("PATH"))
	names := map[string]bool{}
	if profile != "" {
		b, err := os.ReadFile(profile)
		var p struct {
			Path     []string `json:"path"`
			Commands map[string]struct {
				Argv []string `json:"argv"`
			} `json:"commands"`
		}
		if err != nil || json.Unmarshal(b, &p) != nil || len(p.Commands) == 0 {
			r.add(markFail, "exec tools", "cannot read exec profile "+profile, "pass --profile <valid exec profile JSON>")
			return
		}
		paths = p.Path
		for _, command := range p.Commands {
			if len(command.Argv) > 0 {
				names[command.Argv[0]] = true
			}
		}
	}
	for name := range startup {
		names[name] = true
	}
	ordered := make([]string, 0, len(names))
	for name := range names {
		ordered = append(ordered, name)
	}
	sort.Strings(ordered)
	for _, name := range ordered {
		if status, ok := r.execToolStatus[name]; ok {
			switch {
			case status.Probe == "pending":
				r.add(markOK, "exec tools", name+": resolves to "+status.Target+"; the next exec call re-checks it, no restart needed", "")
			case strings.HasPrefix(status.Probe, "failed:") || strings.HasPrefix(status.Probe, "does-not-start:"):
				next := status.Next
				if next == "" {
					next = "reinstall the package that provides " + status.Configured + ", then retry (no restart needed)"
				}
				r.add(markFail, "exec tools", name+": "+status.Probe, next)
			default:
				r.add(markOK, "exec tools", name+": "+status.Accepted, "")
			}
			continue
		}
		current := resolveDoctorTool(name, paths)
		old := startup[name]
		if old != "" {
			_, err := os.Stat(old)
			if err != nil || current != "" && current != old {
				why := "it now resolves elsewhere"
				if os.IsNotExist(err) {
					why = "that file is gone"
				}
				if current != "" {
					why += " (now " + current + ")"
				}
				r.add(markFail, "exec tools", name+": the daemon resolved "+old+" at startup; "+why, "restart World so it resolves the tools again")
				continue
			}
		}
		if current == "" {
			r.add(markFail, "exec tools", name+": binary is missing from the profile path", "install "+name+" in the profile path, then restart World")
		} else {
			r.add(markOK, "exec tools", name+": "+current, "")
		}
	}
}
