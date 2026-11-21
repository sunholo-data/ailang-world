package surfacefollow

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surface"
)

const filePreviewBytes = 128 << 10

// MaxPreviewPath bounds the decoded path in a preview request.
const MaxPreviewPath = 4096

type FilePreview struct {
	Path        string   `json:"path"`
	Status      string   `json:"status"`
	Merge       string   `json:"merge,omitempty"`
	Commit      string   `json:"commit,omitempty"`
	SourceLabel string   `json:"sourceLabel,omitempty"`
	Kind        string   `json:"kind"`
	Size        FileSize `json:"size"`
	Text        *string  `json:"text,omitempty"`
	Diff        string   `json:"diff,omitempty"`
	Truncated   bool     `json:"truncated,omitempty"`
}

type FileSize struct {
	Before *int64 `json:"before,omitempty"`
	After  *int64 `json:"after,omitempty"`
}

type previewEntry struct {
	mode string
	oid  string
	size *int64
}

var errFilePreview = errors.New("committed file unavailable")

func previewGit(ctx context.Context, repo Repo) func(...string) (string, error) {
	return func(args ...string) (string, error) {
		prefix := []string{"-C", repo.Repo, "--no-optional-locks", "--no-pager", "--no-lazy-fetch",
			"--no-replace-objects", "--literal-pathspecs"}
		return (ReadGit{}).Git(ctx, append(prefix, args...))
	}
}

// ReadCommittedFileVersion reads an exact path in one repository commit.
// If want is set, only content with that digest can be presented as evidence.
func ReadCommittedFileVersion(ctx context.Context, repo Repo, commit, path string, want hashref.HashRef) (FilePreview, error) {
	qualified, err := surface.QualifyID(repo.ID, "commit", commit)
	if err != nil || !filepath.IsLocal(path) || len(path) > MaxPreviewPath || strings.ContainsRune(path, 0) {
		return FilePreview{}, errFilePreview
	}
	ctx, cancel := context.WithTimeout(ctx, CallTimeout)
	defer cancel()
	git := previewGit(ctx, repo)
	out, err := git("ls-tree", "-lz", "--full-tree", commit, "--", path)
	entry, ok := parsePreviewEntry(out, path)
	if err != nil || !ok || entry.oid == "" {
		return FilePreview{}, errFilePreview
	}
	p := FilePreview{Path: path, Commit: qualified, Size: FileSize{After: entry.size}}
	switch {
	case entry.mode == "120000":
		p.Kind = "symlink"
	case entry.mode == "160000":
		p.Kind = "submodule"
	case entry.size == nil:
		return FilePreview{}, errFilePreview
	case *entry.size > filePreviewBytes:
		p.Kind = "too-large"
	default:
		text, err := git("cat-file", "blob", entry.oid)
		if err != nil || int64(len(text)) != *entry.size || !want.IsZero() && hashref.SumSHA256([]byte(text)) != want {
			return FilePreview{}, errFilePreview
		}
		if previewText(text) {
			p.Kind, p.Text = "text", &text
		} else {
			p.Kind = "binary"
		}
		return p, nil
	}
	if !want.IsZero() {
		return FilePreview{}, errFilePreview
	}
	return p, nil
}

// ReadFilePreview reads only the given merge and its real first parent. The
// caller must authorize change against its immutable, published partition.
func ReadFilePreview(ctx context.Context, repo Repo, merge string, change surface.Change) (FilePreview, error) {
	qualified, err := surface.QualifyID(repo.ID, "commit", merge)
	if err != nil || !filepath.IsLocal(change.Path) || len(change.Path) > MaxPreviewPath ||
		strings.ContainsRune(change.Path, 0) {
		return FilePreview{}, errFilePreview
	}
	ctx, cancel := context.WithTimeout(ctx, CallTimeout)
	defer cancel()
	git := previewGit(ctx, repo)
	parent := merge + "^1"
	entries := [2]previewEntry{}
	for i, rev := range []string{parent, merge} {
		out, err := git("ls-tree", "-lz", "--full-tree", rev, "--", change.Path)
		if err != nil {
			return FilePreview{}, errFilePreview
		}
		entry, ok := parsePreviewEntry(out, change.Path)
		if !ok {
			return FilePreview{}, errFilePreview
		}
		entries[i] = entry
	}
	before, after := entries[0], entries[1]
	var effective previewEntry
	switch {
	case change.Status == "D" && before.oid != "" && after.oid == "":
		effective = before
	case change.Status == "A" && before.oid == "" && after.oid != "":
		effective = after
	case change.Status == "M" && before.oid != "" && after.oid != "":
		effective = after
	default:
		return FilePreview{}, errFilePreview
	}
	p := FilePreview{Path: change.Path, Status: change.Status, Merge: qualified,
		Size: FileSize{Before: before.size, After: after.size}}
	switch effective.mode {
	case "120000":
		p.Kind = "symlink"
		return p, nil
	case "160000":
		p.Kind = "submodule"
		return p, nil
	}
	if effective.size == nil {
		return FilePreview{}, errFilePreview
	}
	if *effective.size > filePreviewBytes {
		p.Kind = "too-large"
		return p, nil
	}
	text, err := git("cat-file", "blob", effective.oid)
	if err != nil || int64(len(text)) != *effective.size {
		return FilePreview{}, errFilePreview
	}
	if !previewText(text) {
		p.Kind = "binary"
		return p, nil
	}
	p.Kind, p.Text = "text", &text
	// Diff only bounded regular blobs. Type changes and an oversized or binary
	// prior copy can still show the effective text independently.
	for _, entry := range entries {
		if entry.oid == "" {
			continue
		}
		if (entry.mode != "100644" && entry.mode != "100755") || entry.size == nil || *entry.size > filePreviewBytes {
			return p, nil
		}
		if entry.oid != effective.oid {
			other, err := git("cat-file", "blob", entry.oid)
			if err != nil || int64(len(other)) != *entry.size || !previewText(other) {
				return p, nil
			}
		}
	}
	diff, err := git("diff", "--no-ext-diff", "--no-textconv", "--no-renames", "--no-color",
		"--no-relative", "--text", "--ignore-submodules=none", "--unified=3", parent, merge, "--", change.Path)
	if err == nil && previewText(diff) {
		p.Truncated = len(diff) > filePreviewBytes
		if p.Truncated {
			end := filePreviewBytes
			for !utf8.RuneStart(diff[end]) {
				end--
			}
			diff = diff[:end]
		}
		p.Diff = diff
	}
	return p, nil
}

func parsePreviewEntry(out, path string) (previewEntry, bool) {
	if out == "" {
		return previewEntry{}, true
	}
	row, tail, ok := strings.Cut(out, "\x00")
	if !ok || tail != "" {
		return previewEntry{}, false
	}
	header, name, ok := strings.Cut(row, "\t")
	fields := strings.Fields(header)
	if !ok || name != path || len(fields) != 4 || !validHead(fields[2]) {
		return previewEntry{}, false
	}
	mode, kind := fields[0], fields[1]
	entry := previewEntry{mode: mode, oid: fields[2]}
	switch {
	case mode == "040000" && kind == "tree":
		return previewEntry{}, true // a directory can replace a deleted file or vice versa
	case mode == "160000" && kind == "commit" && fields[3] == "-":
		return entry, true
	case (mode == "100644" || mode == "100755" || mode == "120000") && kind == "blob":
		size, err := strconv.ParseInt(fields[3], 10, 64)
		entry.size = &size
		return entry, err == nil && size >= 0
	default:
		return previewEntry{}, false
	}
}

func previewText(s string) bool {
	if !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if unicode.IsControl(r) && r != '\t' && r != '\n' && r != '\r' {
			return false
		}
	}
	return true
}
