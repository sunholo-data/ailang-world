package broker

import (
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// protectedBaseline belongs to one handler. ignored is private and read-only
// after capture; neither phase checks nor record building may mutate it.
type protectedBaseline struct {
	base, digest, ignoredDigest string
	files                       int
	ignored                     map[string]struct{}
	err                         error
}

var protectedTreeRecord = regexp.MustCompile(`^(100644|100755|120000|160000) (blob|commit) (` + protectedHash + `)$`)
var protectedBaseHash = regexp.MustCompile(`^` + protectedHash + `$`)

// captureProtectedBaseline shares ctx's one construction deadline across every
// leg; the caller supplies it (NewExecHandler: protectedBaselineTimeout). On
// any failure, discard all partially captured identities and the snapshot.
func captureProtectedBaseline(ctx context.Context, git, root string, canonical []string) (b protectedBaseline) {
	defer func() {
		if b.err != nil {
			b = protectedBaseline{err: b.err}
		}
	}()
	if git == "" {
		if git, b.err = worktreeGit(); b.err != nil {
			return
		}
	}
	if b.err = protectedGitPreflight(ctx, git, root); b.err != nil {
		return
	}
	if _, err := os.Lstat(filepath.Join(root, ".git")); err != nil {
		b.err = fmt.Errorf("worktree is not a git checkout")
		return
	}
	run := func(args ...string) []byte {
		stderr := newBoundedSink(4096)
		out, err := runBounded(ctx, handlerBounds{execTimeout: protectedBaselineTimeout, maxOutputBytes: protectedCheckMaxOutput},
			handlerCommand{path: git, args: append([]string{"-c", "core.fsmonitor=false"}, args...), dir: root, env: gitEnv(), stderr: stderr})
		if err != nil {
			b.err = fmt.Errorf("protected baseline Git %s: %w (stderr %q)", args[0], err, stderr.Bytes())
		}
		return out
	}
	b.base = strings.TrimSuffix(string(run("rev-parse", "HEAD")), "\n")
	if b.err != nil {
		return
	}
	if !protectedBaseHash.MatchString(b.base) {
		b.err = fmt.Errorf("protected baseline: malformed HEAD")
		return
	}
	// Validate even unprotected records: corrupt output never means intact.
	records := func(out []byte) []string {
		if b.err != nil || len(out) == 0 {
			return nil
		}
		if out[len(out)-1] != 0 {
			b.err = fmt.Errorf("protected baseline: missing NUL terminator")
			return nil
		}
		return strings.Split(string(out[:len(out)-1]), "\x00")
	}
	h := sha256.New()
	for _, record := range records(run("ls-tree", "-r", "-z", "--full-tree", b.base)) {
		header, path, ok := strings.Cut(record, "\t")
		fields := strings.Split(header, " ")
		if !ok || path == "" || !protectedTreeRecord.MatchString(header) ||
			(fields[0] == "160000") != (fields[1] == "commit") {
			b.err = fmt.Errorf("protected baseline: malformed ls-tree record")
			return
		}
		if matchProtected(canonical, path) {
			fmt.Fprintf(h, "%s %s\t%s\n", fields[0], fields[2], path)
			b.files++
		}
	}
	if b.err != nil {
		return
	}
	b.digest = fmt.Sprintf("sha256:%x", h.Sum(nil))
	b.ignored = map[string]struct{}{}
	for _, path := range records(run(append(strings.Fields(protectedIgnoredCommand), gitCandidatePathspecs(canonical)...)...)) {
		if path == "" {
			b.err = fmt.Errorf("protected baseline: empty ignored path")
			return
		}
		if matchProtected(canonical, path) {
			b.ignored[foldProtectedKey(path)] = struct{}{}
		}
	}
	if b.err != nil {
		return
	}
	b.ignoredDigest = protectedKeysDigest(b.ignored)
	b.err = ctx.Err()
	return
}

// protectedKeysDigest is sha256 over the sorted fold keys, each LF-terminated;
// the empty set hashes the empty byte string.
func protectedKeysDigest(set map[string]struct{}) string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, key := range keys {
		fmt.Fprintln(h, key)
	}
	return fmt.Sprintf("sha256:%x", h.Sum(nil))
}

// protectedRecord is data only in M6a; execResult starts emitting it in M6b.
type protectedRecord struct {
	Digest            string                  `json:"digest"`
	Patterns          []string                `json:"patterns"`
	Enforced          string                  `json:"enforced"`
	Base              string                  `json:"base"`
	BaseDigest        string                  `json:"base_digest"`
	BaseFiles         int                     `json:"base_files"`
	BaseIgnoredDigest string                  `json:"base_ignored_digest"`
	BaseIgnoredFiles  int                     `json:"base_ignored_files"`
	IntactBefore      *bool                   `json:"intact_before"`
	IntactAfter       *bool                   `json:"intact_after"`
	Changes           []protectedRecordChange `json:"changes"`
	Truncated         bool                    `json:"truncated"`
	Total             int                     `json:"total"`
	Error             string                  `json:"error"`
}

type protectedRecordChange struct {
	protectedChange
	Phase         string `json:"phase"`
	PathTruncated bool   `json:"path_truncated,omitempty"`
}

func buildProtectedRecord(b protectedBaseline, canonical []string, enforced string, before, after protectedCheck) protectedRecord {
	r := protectedRecord{Digest: protectedDigest(canonical), Patterns: append([]string(nil), canonical...), Enforced: enforced,
		Base: b.base, BaseDigest: b.digest, BaseFiles: b.files, BaseIgnoredDigest: b.ignoredDigest, BaseIgnoredFiles: len(b.ignored),
		Changes: []protectedRecordChange{}}
	for _, err := range []error{b.err, before.err, after.err} {
		if err != nil {
			r.Error = err.Error()
			if len(r.Error) > 4096 {
				r.Error = r.Error[:4093] + "..."
			}
			break
		}
	}
	for i, phase := range []protectedCheck{before, after} {
		if b.err == nil && phase.err == nil && phase.intact != nil {
			v := *phase.intact
			if i == 0 {
				r.IntactBefore = &v
			} else {
				r.IntactAfter = &v
			}
		}
		r.Total += len(phase.changes)
		for _, change := range phase.changes {
			if len(r.Changes) == 32 {
				break
			}
			c := protectedRecordChange{protectedChange: change, Phase: []string{"before", "after"}[i]}
			c.Path, c.PathTruncated = boundPathJSON(c.Path)
			r.Changes = append(r.Changes, c)
		}
	}
	r.Truncated = r.Total > 32
	return r
}
