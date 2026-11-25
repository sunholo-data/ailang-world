package daemon

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// Row 269 Phase 1 (design_docs/planned/w-ui-profile-preset-onboarding.md):
// read-only surface endpoints and address-bar verbs listing active/proposed
// exec profiles (joined with the daemon's live probe status) and configured
// launch presets, plus the standalone /surface/repos UI page.
const (
	SurfaceProfilesSchema = "world/surface-profiles/v1"
	SurfacePresetsSchema  = "world/surface-presets/v1"
)

type surfaceProfileEntry struct {
	Project       string   `json:"project"`
	State         string   `json:"state"`
	ProbeStatus   string   `json:"probe_status"`
	ExcludeReason string   `json:"exclude_reason,omitempty"`
	Root          string   `json:"root,omitempty"`
	Commands      []string `json:"commands,omitempty"`
	Digest        string   `json:"digest,omitempty"`
	Error         string   `json:"error,omitempty"`
}

type surfaceProfilesResponse struct {
	Schema   string                `json:"schema"`
	Profiles []surfaceProfileEntry `json:"profiles"`
}

type surfacePresetEntry struct {
	Name        string   `json:"name"`
	Repo        string   `json:"repo,omitempty"`
	Remote      string   `json:"remote,omitempty"`
	Base        string   `json:"base,omitempty"`
	Agent       string   `json:"agent,omitempty"`
	Model       string   `json:"model,omitempty"`
	Profile     string   `json:"profile,omitempty"`
	Mode        string   `json:"mode,omitempty"`
	TTLS        int64    `json:"ttl_s,omitempty"`
	GrantsCount int      `json:"grants_count"`
	Checks      []string `json:"checks,omitempty"`
	Digest      string   `json:"digest,omitempty"`
	Error       string   `json:"error,omitempty"`
}

type surfacePresetsResponse struct {
	Schema  string               `json:"schema"`
	Presets []surfacePresetEntry `json:"presets"`
}

func (d *Daemon) readSurfaceProfiles() []surfaceProfileEntry {
	loaded, excluded := d.workspace.healthProfiles()
	loadedByProj := make(map[string]HealthExecProfile, len(loaded))
	for _, p := range loaded {
		loadedByProj[p.Project] = p
	}
	excludedByProj := make(map[string]HealthExecExclusion, len(excluded))
	for _, ex := range excluded {
		excludedByProj[ex.Project] = ex
	}

	rows := []surfaceProfileEntry{}
	seenActive := map[string]bool{}
	if d.execProfilesDir != "" {
		if ents, err := os.ReadDir(d.execProfilesDir); err == nil {
			for _, e := range ents {
				if e.IsDir() {
					continue
				}
				name := e.Name()
				var project, state string
				switch {
				case strings.HasSuffix(name, broker.ExecProposedSuffix):
					project = strings.TrimSuffix(name, broker.ExecProposedSuffix)
					state = "proposed"
				case strings.HasSuffix(name, broker.ExecProfileSuffix):
					project = strings.TrimSuffix(name, broker.ExecProfileSuffix)
					state = "active"
					seenActive[project] = true
				default:
					continue
				}
				row := surfaceProfileEntry{Project: project, State: state, ProbeStatus: "unloaded"}
				if state == "active" {
					if _, ok := loadedByProj[project]; ok {
						row.ProbeStatus = "loaded"
					} else if ex, ok := excludedByProj[project]; ok {
						row.ProbeStatus = "excluded"
						row.ExcludeReason = ex.Reason
					}
				}
				raw, err := os.ReadFile(filepath.Join(d.execProfilesDir, name))
				if err != nil {
					row.Error = err.Error()
				} else if parsed, err := broker.ParseExecProfile(raw); err != nil {
					row.Error = err.Error()
				} else {
					row.Root = parsed.Root
					row.Digest = parsed.Digest
					for cmdName := range parsed.Commands {
						row.Commands = append(row.Commands, cmdName)
					}
					sort.Strings(row.Commands)
				}
				rows = append(rows, row)
			}
		}
	}
	// Include any --exec-profile loaded outside ExecProfilesDir.
	for _, lp := range loaded {
		if seenActive[lp.Project] {
			continue
		}
		rows = append(rows, surfaceProfileEntry{
			Project:     lp.Project,
			State:       "active",
			ProbeStatus: "loaded",
			Digest:      lp.Digest,
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Project != rows[j].Project {
			return rows[i].Project < rows[j].Project
		}
		return rows[i].State < rows[j].State
	})
	return rows
}

func (d *Daemon) readSurfacePresets() []surfacePresetEntry {
	rows := []surfacePresetEntry{}
	if d.presetsDir != "" {
		if ents, err := os.ReadDir(d.presetsDir); err == nil {
			for _, e := range ents {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
					continue
				}
				name := strings.TrimSuffix(e.Name(), ".json")
				row := surfacePresetEntry{Name: name}
				raw, err := os.ReadFile(filepath.Join(d.presetsDir, e.Name()))
				if err != nil {
					row.Error = err.Error()
					rows = append(rows, row)
					continue
				}
				sum := sha256.Sum256(raw)
				row.Digest = "sha256:" + hex.EncodeToString(sum[:])
				var doc struct {
					Name    string            `json:"name"`
					Repo    string            `json:"repo"`
					Remote  string            `json:"remote"`
					Base    string            `json:"base"`
					Agent   string            `json:"agent"`
					Model   string            `json:"model"`
					Profile string            `json:"profile"`
					TTLS    int64             `json:"ttl_s"`
					Grants  []json.RawMessage `json:"grants"`
					Launch  struct {
						Mode string `json:"mode"`
					} `json:"launch"`
					Checks []struct {
						Name string `json:"name"`
					} `json:"checks"`
				}
				if err := json.Unmarshal(raw, &doc); err != nil {
					row.Error = err.Error()
				} else {
					if doc.Name != "" {
						row.Name = doc.Name
					}
					row.Repo = doc.Repo
					row.Remote = doc.Remote
					row.Base = doc.Base
					row.Agent = doc.Agent
					row.Model = doc.Model
					row.Profile = doc.Profile
					row.Mode = doc.Launch.Mode
					row.TTLS = doc.TTLS
					row.GrantsCount = len(doc.Grants)
					for _, c := range doc.Checks {
						if c.Name != "" {
							row.Checks = append(row.Checks, c.Name)
						}
					}
				}
				rows = append(rows, row)
			}
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].Name < rows[j].Name })
	return rows
}

// handleSurfaceProfiles serves GET /surface/profiles.json behind approvalRoute.
func (d *Daemon) handleSurfaceProfiles(w http.ResponseWriter, r *http.Request, _ []string) {
	if r.URL.RawQuery != "" {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", unknownSurfaceKeyMessage)
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfaceProfilesResponse{
		Schema:   SurfaceProfilesSchema,
		Profiles: d.readSurfaceProfiles(),
	})
}

// handleSurfacePresets serves GET /surface/presets.json behind approvalRoute.
func (d *Daemon) handleSurfacePresets(w http.ResponseWriter, r *http.Request, _ []string) {
	if r.URL.RawQuery != "" {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", unknownSurfaceKeyMessage)
		return
	}
	writeApprovalJSON(w, http.StatusOK, surfacePresetsResponse{
		Schema:  SurfacePresetsSchema,
		Presets: d.readSurfacePresets(),
	})
}

// handleSurfaceReposPage serves GET /surface/repos behind approvalRoute.
func (d *Daemon) handleSurfaceReposPage(w http.ResponseWriter, r *http.Request, _ []string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; script-src 'self'; style-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; form-action 'none'; frame-ancestors 'none'")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(surfacepage.ReposPage)
}

// handleSurfaceReposScript serves GET /surface/repos.js behind approvalRoute.
func (d *Daemon) handleSurfaceReposScript(w http.ResponseWriter, r *http.Request, _ []string) {
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(surfacepage.ReposScript)
}

func parseListJSONFlag(verb string, args []string) (asJSON bool, reason string) {
	for _, arg := range args {
		if arg == "--json" {
			asJSON = true
			continue
		}
		return false, fmt.Sprintf("%s accepts only --json (unexpected %q)", verb, arg)
	}
	return asJSON, ""
}

func (d *Daemon) verbProfileList(args []string, a *verbAnswer) {
	asJSON, reason := parseListJSONFlag("profile list", args)
	if reason != "" {
		a.Reason = reason
		return
	}
	rows := d.readSurfaceProfiles()
	a.Outcome = "result"
	if asJSON {
		raw, _ := json.MarshalIndent(rows, "", "  ")
		a.JSON = string(raw) + "\n"
		return
	}
	if len(rows) == 0 {
		a.Text = "no exec profiles configured\n"
		return
	}
	var sb strings.Builder
	for _, r := range rows {
		if r.Error != "" {
			fmt.Fprintf(&sb, "%-20s %-8s %-8s ERROR: %s\n", r.Project, r.State, r.ProbeStatus, r.Error)
			continue
		}
		fmt.Fprintf(&sb, "%-20s %-8s %-8s root=%-10s cmds=[%s]\n",
			r.Project, r.State, r.ProbeStatus, r.Root, strings.Join(r.Commands, ", "))
	}
	a.Text = sb.String()
}

func (d *Daemon) verbPresetList(args []string, a *verbAnswer) {
	asJSON, reason := parseListJSONFlag("preset list", args)
	if reason != "" {
		a.Reason = reason
		return
	}
	rows := d.readSurfacePresets()
	a.Outcome = "result"
	if asJSON {
		raw, _ := json.MarshalIndent(rows, "", "  ")
		a.JSON = string(raw) + "\n"
		return
	}
	if len(rows) == 0 {
		a.Text = "no launch presets configured\n"
		return
	}
	var sb strings.Builder
	for _, p := range rows {
		if p.Error != "" {
			fmt.Fprintf(&sb, "%-20s ERROR: %s\n", p.Name, p.Error)
			continue
		}
		fmt.Fprintf(&sb, "%-20s %-6s %-11s profile=%-16s base=%-6s %s\n",
			p.Name, p.Agent, p.Mode, p.Profile, p.Base, p.Repo)
	}
	a.Text = sb.String()
}
