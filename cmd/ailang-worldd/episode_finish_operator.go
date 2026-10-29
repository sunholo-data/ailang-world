package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
)

func finishSHA(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// writeOperatorProof never puts the proof itself on disk. Its hash is beside
// the independently protected credentials store, not in the episode worktree.
func writeOperatorProof(credentials, episode, phase string) (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	proof := hex.EncodeToString(raw)
	dir := filepath.Join(filepath.Dir(credentials), "operator-proofs")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	info, err := os.Lstat(dir)
	if err != nil || !info.IsDir() || info.Mode().Perm() != 0o700 {
		return "", fmt.Errorf("operator-proofs must be a private 0700 directory")
	}
	data, err := json.Marshal(struct {
		SHA256    string `json:"sha256"`
		Phase     string `json:"phase"`
		ExpiresAt int64  `json:"expires_at"`
	}{finishSHA([]byte(proof)), phase, time.Now().Unix() + 300})
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp(dir, ".proof-")
	if err != nil {
		return "", err
	}
	defer os.Remove(f.Name())
	if err := f.Chmod(0o600); err != nil {
		f.Close()
		return "", err
	}
	_, writeErr := f.Write(data)
	syncErr := f.Sync()
	closeErr := f.Close()
	for _, err := range []error{writeErr, syncErr, closeErr} {
		if err != nil {
			return "", err
		}
	}
	if err := os.Rename(f.Name(), filepath.Join(dir, episode+".json")); err != nil {
		return "", err
	}
	return proof, nil
}

func finishOperatorPost(ctx context.Context, c *client, credentials, ep string, input daemon.FinishRequest, out any) error {
	proof, err := writeOperatorProof(credentials, ep, input.Phase)
	if err != nil {
		return err
	}
	data, err := json.Marshal(input)
	if err != nil {
		return err
	}
	status, _, body, err := c.doRequest(ctx, http.MethodPost, "/v1/episodes/"+ep+"/finish", bytes.NewReader(data),
		map[string]string{"X-World-Operator-Proof": proof, "Content-Type": "application/json"})
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("finish operator route: HTTP %d: %s", status, body)
	}
	if out != nil {
		return json.Unmarshal(body, out)
	}
	return nil
}

func finishProtected(ctx context.Context, start *store.EpisodeStart, credentials, profilesDir string, changed []string) ([]string, error) {
	var pinned launchPreset
	if start.Preset.Bytes == "" || json.Unmarshal([]byte(start.Preset.Bytes), &pinned) != nil {
		return nil, fmt.Errorf("the start record has no valid pinned preset bytes")
	}
	if pinned.Profile == "" {
		return nil, nil
	}
	profile, _, err := loadMintProfile(profilesDir, pinned.Profile)
	if err != nil {
		return nil, err
	}
	creds, err := store.OpenCredentialsReadOnly(credentials)
	if err != nil {
		return nil, err
	}
	defer creds.Close()
	rows, err := creds.ListSessions(ctx, start.Episode, store.MaxListSessions)
	if err != nil {
		return nil, err
	}
	found := false
	for _, row := range rows {
		if row.ExpiresAt > time.Now().Unix() {
			found = true
			if row.Project != pinned.Profile || row.ProfileDigest != profile.Digest {
				return nil, fmt.Errorf("the exec profile %s changed since episode %s started", pinned.Profile, start.Episode)
			}
		}
	}
	if !found {
		return nil, fmt.Errorf("episode %s has no live profile-bound credential", start.Episode)
	}
	var paths []string
	for _, path := range changed {
		if broker.MatchProtected(profile.Protected, path) {
			paths = append(paths, path)
		}
	}
	return paths, nil
}

func finishPrompt(tty io.Writer, scope broker.FinishApprovalScope, title, body string, protected []string) {
	lines := strings.Split(body, "\n")
	preview := strings.Join(lines[:min(6, len(lines))], "\n")
	state := "draft"
	if !scope.Draft {
		state = "ready"
	}
	label := "ailang-world"
	if scope.CreateLabel {
		label += " (will be created)"
	}
	fmt.Fprintf(tty, "repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\nbody (%d bytes), first lines:\n%s\nlabel: %s\n%s\nprotected paths changed: %s\nOpen this pull request? [y/N] ",
		scope.Repository, scope.Head, scope.SHA, scope.Base, title, len(body), preview, label, state, strings.Join(protected, ", "))
}
