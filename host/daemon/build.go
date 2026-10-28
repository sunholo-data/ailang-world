package daemon

// Row 189 M3 (w-world-onboarding §3.1 step 1, AC3): the build identity a
// binary reports, so `ailang-worldd episode start` can refuse a daemon built
// from different source than itself (P4: a stale rig build silently served an
// old tool pin). It is read from the binary's own debug.BuildInfo: the
// vcs.revision and vcs.modified settings `go build` stamps from the checkout,
// plus the sha256 of the running executable's bytes.
//
// What a build without VCS stamping reports: `go run`, a `go test` binary,
// a build with -buildvcs=false and a build from outside a git checkout carry no
// vcs.revision, so VCSRevision is "". Nor does a build inside a git worktree
// (.git is a file there), which is where every World episode builds. So
// episode start matches on executable bytes first, and on a revision only
// when both sides carry the same clean one (SameBuild).

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"runtime/debug"
	"sync"
)

// BuildIdentity is a binary's build, as /v1/health reports it.
type BuildIdentity struct {
	// VCSRevision is the vcs.revision the build was stamped with; "" when
	// the build carries none (see the file comment).
	VCSRevision string `json:"vcs_revision"`
	// VCSModified is vcs.modified: the checkout had uncommitted changes.
	VCSModified bool `json:"vcs_modified"`
	// ExecutableSHA256 is "sha256:<hex>" of the running binary's bytes, ""
	// when they could not be read.
	ExecutableSHA256 string `json:"executable_sha256,omitempty"`
	// GoVersion is the toolchain that built it.
	GoVersion string `json:"go_version,omitempty"`
}

// Known reports whether the build carries a VCS identity at all.
func (b BuildIdentity) Known() bool { return b.VCSRevision != "" }

// String is the revision, "+dirty" when modified, or "none".
func (b BuildIdentity) String() string {
	switch {
	case !b.Known():
		return "none (no vcs.revision: a go run, go test or non-VCS build)"
	case b.VCSModified:
		return b.VCSRevision + "+dirty"
	default:
		return b.VCSRevision
	}
}

// SameBuild reports whether a and b are one build: the same executable bytes
// are the same build, stamped or not (a build in a git worktree, where .git is
// a file, carries no vcs.revision); failing that, both carry the same clean
// revision. Anything else — no stamp, a dirty tree, another revision — is a
// mismatch.
func SameBuild(a, b BuildIdentity) bool {
	if a.ExecutableSHA256 != "" && a.ExecutableSHA256 == b.ExecutableSHA256 {
		return true
	}
	return a.Known() && b.Known() && a.VCSRevision == b.VCSRevision && !a.VCSModified && !b.VCSModified
}

// buildIdentityOf reads the identity out of build info and an executable
// digest.
func buildIdentityOf(info *debug.BuildInfo, exeDigest string) BuildIdentity {
	b := BuildIdentity{ExecutableSHA256: exeDigest}
	if info == nil {
		return b
	}
	b.GoVersion = info.GoVersion
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			b.VCSRevision = s.Value
		case "vcs.modified":
			b.VCSModified = s.Value == "true"
		}
	}
	return b
}

// currentBuild is computed once per process. Run eagerly freezes it before listening.
var currentBuild = sync.OnceValue(func() BuildIdentity {
	info, _ := debug.ReadBuildInfo()
	return buildIdentityOf(info, executableDigest())
})

// CurrentBuild is this process's build identity.
func CurrentBuild() BuildIdentity { return currentBuild() }

func executableDigest() string {
	path, err := os.Executable()
	if err != nil {
		return ""
	}
	return executableDigestAt(path)
}

func executableDigestAt(path string) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer func() { _ = f.Close() }()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}
	return "sha256:" + hex.EncodeToString(h.Sum(nil))
}
