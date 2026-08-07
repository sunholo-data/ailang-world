package broker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunholo-data/ailang-world/host/childenv"
)

// RegistryCredentialVariable is the environment variable the pinned publisher
// reads for its API key. It is re-exported from host/childenv so callers do
// not spell the name a second time.
const RegistryCredentialVariable = childenv.CredentialVariable

// redactedMarker replaces every occurrence of the credential in any byte
// stream this package hands to an error, a result object, or a log.
const redactedMarker = "[REDACTED:AILANG_REGISTRY_API_KEY]"

// ErrAmbientRegistryCredential reports that the irreversible-publish
// credential is present in the process environment. Decision 4 makes this a
// startup failure rather than a warning: the public registry is immutable, so
// an ambient credential means every subprocess World launches — every agent,
// every shell command, every `ailang` invocation — inherits the authority to
// perform an unrecallable public write merely by being a child.
var ErrAmbientRegistryCredential = errors.New(
	"broker: " + RegistryCredentialVariable + " is set in the process environment")

// AssertNoAmbientRegistryCredential fails when environ carries a non-empty
// registry credential. The returned error names the VARIABLE and never its
// value; a guard that prints the secret it is guarding has published it to
// every log that captures stderr.
func AssertNoAmbientRegistryCredential(environ []string) error {
	if childenv.Has(environ, RegistryCredentialVariable) {
		return fmt.Errorf(
			"%w: move it to a mode-0600 file outside the working tree and unset it "+
				"(see design_docs/planned/w-self-mod-vertical.md Decision 4)",
			ErrAmbientRegistryCredential)
	}
	return nil
}

// RegistryCredentialProvider yields the registry API key for exactly one live
// publish dispatch. Implementations return the bytes; they never retain them,
// never log them, and are never consulted on the dry-run path.
type RegistryCredentialProvider interface {
	Credential() ([]byte, error)
}

// FileRegistryCredentialProvider reads the credential from a mode-0600 regular
// file outside the working tree.
//
// Only the PATH is a field. The bytes are read on demand and returned to the
// caller, so nothing in this object, its String form, or any structure derived
// from it can carry the secret — which is what makes "never stored in an
// object" checkable rather than a promise.
type FileRegistryCredentialProvider struct {
	path string
}

// CredentialFileError reports a refused credential file. It names the path and
// the reason and, like every error in this file, never the contents.
type CredentialFileError struct {
	Path string
	Why  string
}

func (e *CredentialFileError) Error() string {
	return fmt.Sprintf("broker: registry credential file %q: %s", e.Path, e.Why)
}

// NewFileRegistryCredentialProvider validates the credential file's location
// and mode at construction, so a misplaced or world-readable secret is a
// startup failure rather than a surprise at the one moment an irreversible
// publish is already authorized.
//
// repoRoot is the working tree the secret must NOT live inside. It is required
// rather than optional: a tree-relative secret is one `git add -A` from being
// committed, and this constructor is the only place that can see the two paths
// together.
func NewFileRegistryCredentialProvider(path, repoRoot string) (*FileRegistryCredentialProvider, error) {
	if !filepath.IsAbs(path) {
		return nil, &CredentialFileError{Path: path, Why: "path is not absolute"}
	}
	if repoRoot == "" || !filepath.IsAbs(repoRoot) {
		return nil, &CredentialFileError{Path: path, Why: "an absolute repository root is required"}
	}
	resolvedRoot, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		return nil, &CredentialFileError{Path: path, Why: "resolve repository root: " + err.Error()}
	}
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, &CredentialFileError{Path: path, Why: "resolve path: " + err.Error()}
	}
	rel, err := filepath.Rel(resolvedRoot, resolved)
	if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return nil, &CredentialFileError{Path: path, Why: "file is inside the working tree " + repoRoot}
	}
	info, err := os.Lstat(resolved)
	if err != nil {
		return nil, &CredentialFileError{Path: path, Why: "stat: " + err.Error()}
	}
	if !info.Mode().IsRegular() {
		return nil, &CredentialFileError{Path: path, Why: "not a regular file"}
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		return nil, &CredentialFileError{
			Path: path,
			Why:  fmt.Sprintf("mode is %#o, want exactly 0600", perm),
		}
	}
	return &FileRegistryCredentialProvider{path: resolved}, nil
}

// Credential reads and returns the secret bytes. The mode is re-checked on
// every read because construction and dispatch are separated in time and a
// file can be chmod'ed between them.
func (p *FileRegistryCredentialProvider) Credential() ([]byte, error) {
	info, err := os.Lstat(p.path)
	if err != nil {
		return nil, &CredentialFileError{Path: p.path, Why: "stat: " + err.Error()}
	}
	if !info.Mode().IsRegular() {
		return nil, &CredentialFileError{Path: p.path, Why: "not a regular file"}
	}
	if perm := info.Mode().Perm(); perm != 0o600 {
		return nil, &CredentialFileError{
			Path: p.path,
			Why:  fmt.Sprintf("mode is %#o, want exactly 0600", perm),
		}
	}
	data, err := os.ReadFile(p.path)
	if err != nil {
		return nil, &CredentialFileError{Path: p.path, Why: "read: " + err.Error()}
	}
	secret := []byte(strings.TrimRight(string(data), "\r\n"))
	if len(secret) == 0 {
		return nil, &CredentialFileError{Path: p.path, Why: "file is empty"}
	}
	return secret, nil
}

// redactSecret removes every occurrence of secret from text.
//
// An empty secret is returned unchanged rather than replaced: strings.Replace
// with an empty `old` inserts the marker between every rune, which would turn
// a harmless dry-run transcript into unreadable noise and — worse — would make
// the redaction assertion pass for a reason unrelated to the secret.
func redactSecret(text string, secret []byte) string {
	if len(secret) == 0 {
		return text
	}
	return strings.ReplaceAll(text, string(secret), redactedMarker)
}
