package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/store"
)

const allowanceHelp = `usage: ailang-worldd allowance <grant|list|revoke> [flags]

Manage standing per-preset session allowances (row 224, D-WORLD-96):
  ailang-worldd allowance grant --preset <name> [--max-per-day 10] [--valid-days 30]
      [--presets-dir D] [--exec-profiles-dir D] [--credentials-db P]
  ailang-worldd allowance list [--all] [--json] [--credentials-db P]
  ailang-worldd allowance revoke <allowance-id|preset> [--reason TEXT] [--credentials-db P]

D-WORLD-96 bounds:
  - At most 10 sessions/day per preset, 12 h max session TTL, 30 days max validity
  - Bound to exact preset file SHA-256 and exec profile digest (any edit invalidates)
  - Mission loop (row 196) is always excluded
  - Revocable at any time; over-bound draws refuse and require explicit approval
`

func runAllowance(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, allowanceHelp)
		return exitUsage
	}
	switch args[0] {
	case "grant":
		return runAllowanceGrant(args[1:], stdout, stderr, env)
	case "list":
		return runAllowanceList(args[1:], stdout, stderr, env)
	case "revoke":
		return runAllowanceRevoke(args[1:], stdout, stderr, env)
	case "--help", "-help", "-h", "help":
		fmt.Fprint(stdout, allowanceHelp)
		return exitOK
	default:
		fmt.Fprintf(stderr, "ailang-worldd allowance: unknown subcommand %q\n\n", args[0])
		fmt.Fprint(stderr, allowanceHelp)
		return exitUsage
	}
}

func runAllowanceGrant(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	fs := flag.NewFlagSet("ailang-worldd allowance grant", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	presetName := fs.String("preset", "", "launch preset name (required)")
	maxPerDay := fs.Int("max-per-day", approval.AllowanceMaxPerDay, "max sessions per UTC day (1..10)")
	validDays := fs.Int("valid-days", approval.AllowanceMaxValidDays, "validity in days (1..30)")
	presetsDir := fs.String("presets-dir", defaultPresetsDir(), "launch presets directory")
	profilesDir := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "exec profiles directory")
	credFlag := fs.String("credentials-db", "", "credentials store path (default from preset daemon)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, allowanceHelp)
			return exitOK
		}
		fmt.Fprint(stderr, allowanceHelp)
		return exitUsage
	}
	fail := func(format string, a ...any) int {
		fmt.Fprintf(stderr, "ailang-worldd allowance grant: "+format+"\n", a...)
		return exitUsage
	}
	if *presetName == "" {
		return fail("--preset <name> is required")
	}

	def := defaultWorldPaths()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	p, err := loadLaunchPreset(ctx, *presetName, *presetsDir, *profilesDir, def.workspaceRoot)
	cancel()
	if err != nil {
		return fail("%v", err)
	}
	ttl := p.Preset.TTLS
	if ttl > approval.AllowanceMaxTTLSeconds {
		return fail("preset ttl_s %d exceeds D-WORLD-96 ceiling %d (12h)", ttl, approval.AllowanceMaxTTLSeconds)
	}

	credPath := *credFlag
	if credPath == "" {
		credPath = p.Daemon.CredentialsDB
	}
	if credPath == "" {
		var err error
		credPath, err = store.DefaultCredentialsPath(def.db)
		if err != nil {
			return fail("default credentials path: %v", err)
		}
	}
	worldDB := p.Daemon.DB
	if worldDB == "" {
		worldDB = def.db
	}
	wsRoot := p.Daemon.WorkspaceRoot
	if wsRoot == "" {
		wsRoot = def.workspaceRoot
	}

	apGrants := make([]approval.Grant, len(p.Grants))
	for i, g := range p.Grants {
		apGrants[i] = approval.Grant{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget}
	}
	presetHex := strings.TrimPrefix(p.SHA256, "sha256:")
	al := approval.StandingAllowance{
		Schema:             approval.AllowanceSchema,
		PresetName:         *presetName,
		PresetSHA256:       presetHex,
		Profile:            approval.ProfileRef{Project: p.Preset.Profile, Digest: p.ProfileDigest},
		Repos:              []string{filepath.Clean(p.Preset.Repo)},
		Grants:             apGrants,
		TTLSeconds:         ttl,
		MaxPerDay:          *maxPerDay,
		ValidDays:          *validDays,
		WorldDB:            filepath.Clean(worldDB),
		CredentialsDB:      filepath.Clean(credPath),
		WorkspaceRoot:      filepath.Clean(wsRoot),
		ExcludeMissionLoop: true,
	}
	payload, err := approval.CanonicalAllowance(al)
	if err != nil {
		return fail("%v", err)
	}
	digest, err := approval.AllowanceDigest(al)
	if err != nil {
		return fail("%v", err)
	}

	// Attended TTY confirmation before touching credentials.db.
	now := env.now()
	expiresAt := now + int64(*validDays)*86400
	var prompt strings.Builder
	fmt.Fprintf(&prompt, "Grant standing session allowance for preset %s (%s)\n", *presetName, p.SHA256)
	fmt.Fprintf(&prompt, "  repo         %s\n", p.Preset.Repo)
	if p.Preset.Profile != "" {
		fmt.Fprintf(&prompt, "  profile      %s %s\n", p.Preset.Profile, p.ProfileDigest)
	}
	fmt.Fprintf(&prompt, "  grants       %s\n", grantList(p.Grants))
	fmt.Fprintf(&prompt, "  bounds       max %d sessions/day · TTL %ds · valid %d days (until %s)\n",
		*maxPerDay, ttl, *validDays, time.Unix(expiresAt, 0).UTC().Format(time.RFC3339))
	fmt.Fprintf(&prompt, "  mission-loop excluded (D-WORLD-96)\n")
	fmt.Fprintf(&prompt, "  credentials  %s\n", credPath)
	fmt.Fprintf(&prompt, "Proceed? [y/N] ")
	if code := confirmApproverAct(env, stderr, "ailang-worldd allowance grant: ", prompt.String(), "no standing allowance was written"); code != exitOK {
		return code
	}

	idBytes := make([]byte, 8)
	chalBytes := make([]byte, 32)
	if _, err := rand.Read(idBytes); err != nil {
		return fail("random id: %v", err)
	}
	if _, err := rand.Read(chalBytes); err != nil {
		return fail("random challenge: %v", err)
	}
	allowanceID := "sa-" + hex.EncodeToString(idBytes)

	writeCtx, writeCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer writeCancel()

	st, err := store.OpenCredentials(credPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance grant: open credentials %q: %v\n", credPath, err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	row := store.StandingAllowanceRow{
		AllowanceID:    allowanceID,
		PresetName:     *presetName,
		PresetSHA256:   presetHex,
		ProfileProject: p.Preset.Profile,
		ProfileDigest:  p.ProfileDigest,
		Payload:        payload,
		Digest:         hex.EncodeToString(digest[:]),
		Challenge:      hex.EncodeToString(chalBytes),
		MaxPerDay:      *maxPerDay,
		TTLS:           ttl,
		ApprovedAt:     now,
		ExpiresAt:      expiresAt,
		ApprovedBy:     "tty",
	}
	if err := st.GrantAllowance(writeCtx, row); err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance grant: %v\n", err)
		return exitUsage
	}
	fmt.Fprintf(stdout, "granted standing allowance %s for preset %s (max %d/day, TTL %ds, expires %s)\n",
		allowanceID, *presetName, *maxPerDay, ttl, time.Unix(expiresAt, 0).UTC().Format(time.RFC3339))
	return exitOK
}

type allowanceListItem struct {
	AllowanceID   string `json:"allowance_id"`
	PresetName    string `json:"preset_name"`
	PresetSHA256  string `json:"preset_sha256"`
	Profile       string `json:"profile,omitempty"`
	MaxPerDay     int    `json:"max_per_day"`
	DrawsToday    int    `json:"draws_today"`
	DayKey        string `json:"day_key"`
	TTLS          int64  `json:"ttl_s"`
	ApprovedBy    string `json:"approved_by"`
	ApprovedAt    string `json:"approved_at"`
	ExpiresAt     string `json:"expires_at"`
	Status        string `json:"status"`
	RevokedAt     string `json:"revoked_at,omitempty"`
	RevokedReason string `json:"revoked_reason,omitempty"`
}

func runAllowanceList(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	fs := flag.NewFlagSet("ailang-worldd allowance list", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	all := fs.Bool("all", false, "include revoked allowances")
	asJSON := fs.Bool("json", false, "emit JSON")
	credFlag := fs.String("credentials-db", "", "credentials store path (default beside world.db)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, allowanceHelp)
			return exitOK
		}
		fmt.Fprint(stderr, allowanceHelp)
		return exitUsage
	}
	credPath := *credFlag
	if credPath == "" {
		var err error
		credPath, err = store.DefaultCredentialsPath(defaultWorldPaths().db)
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd allowance list: %v\n", err)
			return exitFatal
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	st, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance list: %v\n", err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	rows, err := st.ListAllowances(ctx, *all)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance list: %v\n", err)
		return exitFatal
	}
	now := env.now()
	dayKey := approval.AllowanceDayKey(now)
	items := make([]allowanceListItem, 0, len(rows))
	for _, r := range rows {
		draws, _ := st.CountAllowanceDraws(ctx, r.AllowanceID, dayKey)
		status := "active"
		var revAt string
		if r.RevokedAt != nil {
			status = "revoked"
			revAt = time.Unix(*r.RevokedAt, 0).UTC().Format(time.RFC3339)
		} else if now >= r.ExpiresAt {
			status = "expired"
		}
		items = append(items, allowanceListItem{
			AllowanceID:   r.AllowanceID,
			PresetName:    r.PresetName,
			PresetSHA256:  r.PresetSHA256,
			Profile:       r.ProfileProject,
			MaxPerDay:     r.MaxPerDay,
			DrawsToday:    draws,
			DayKey:        dayKey,
			TTLS:          r.TTLS,
			ApprovedBy:    r.ApprovedBy,
			ApprovedAt:    time.Unix(r.ApprovedAt, 0).UTC().Format(time.RFC3339),
			ExpiresAt:     time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339),
			Status:        status,
			RevokedAt:     revAt,
			RevokedReason: r.RevokedReason,
		})
	}
	if *asJSON {
		enc := json.NewEncoder(stdout)
		enc.SetIndent("", "  ")
		_ = enc.Encode(items)
		return exitOK
	}
	if len(items) == 0 {
		fmt.Fprintln(stdout, "no standing allowances")
		return exitOK
	}
	for _, it := range items {
		fmt.Fprintf(stdout, "%-19s %-16s %-8s today=%d/%-2d ttl=%ds expires=%s\n",
			it.AllowanceID, it.PresetName, it.Status, it.DrawsToday, it.MaxPerDay, it.TTLS, it.ExpiresAt)
	}
	return exitOK
}

func runAllowanceRevoke(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	fs := flag.NewFlagSet("ailang-worldd allowance revoke", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	reason := fs.String("reason", "", "optional revocation reason")
	credFlag := fs.String("credentials-db", "", "credentials store path (default beside world.db)")
	var target string
	var flagArgs []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if strings.HasPrefix(a, "-") {
			flagArgs = append(flagArgs, a)
			if (a == "--reason" || a == "-reason" || a == "--credentials-db" || a == "-credentials-db") && i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		} else if target == "" {
			target = a
		} else {
			flagArgs = append(flagArgs, a)
		}
	}
	if err := fs.Parse(flagArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, allowanceHelp)
			return exitOK
		}
		fmt.Fprint(stderr, allowanceHelp)
		return exitUsage
	}
	if target == "" && fs.NArg() > 0 {
		target = fs.Arg(0)
	}
	if target == "" {
		fmt.Fprintln(stderr, "ailang-worldd allowance revoke: <allowance-id|preset> is required")
		return exitUsage
	}

	credPath := *credFlag
	if credPath == "" {
		var err error
		credPath, err = store.DefaultCredentialsPath(defaultWorldPaths().db)
		if err != nil {
			fmt.Fprintf(stderr, "ailang-worldd allowance revoke: %v\n", err)
			return exitFatal
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	st, err := store.OpenCredentials(credPath)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance revoke: %v\n", err)
		return exitFatal
	}
	defer func() { _ = st.Close() }()

	rev, err := st.RevokeAllowance(ctx, target, env.now(), "tty", *reason)
	if err != nil {
		fmt.Fprintf(stderr, "ailang-worldd allowance revoke: %v\n", err)
		return exitUsage
	}
	fmt.Fprintf(stdout, "revoked standing allowance %s (preset %s)\n", rev.AllowanceID, rev.PresetName)
	return exitOK
}
