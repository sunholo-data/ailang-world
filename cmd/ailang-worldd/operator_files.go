package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
	"golang.org/x/sys/unix"
)

type operatorFileID struct {
	Dev, Inode, Mode, UID, Links uint64
	Size                         int64
	Modified, Changed            unix.Timespec
}

func operatorFileStat(s unix.Stat_t) operatorFileID {
	return operatorFileID{uint64(s.Dev), s.Ino, uint64(s.Mode), uint64(s.Uid), uint64(s.Nlink), s.Size, s.Mtim, s.Ctim}
}

type operatorFileSet struct {
	Dir           string
	Device, Inode uint64
	UID           uint64
	Names, Hashes [3]string
}

type operatorFileProof struct {
	ID     operatorFileID
	Status string
	Digest string
}

type operatorFilesProof struct {
	Dir   operatorFileProof
	Files [3]operatorFileProof
}

func operatorFileFact(a approval.OperatorAction, name string) string {
	for _, fact := range a.Resolution {
		if fact.Name == name {
			return fact.Value
		}
	}
	return ""
}

func operatorFilePaths(s store.OperatorFileState, a approval.OperatorAction) (operatorFileSet, store.EpisodeStart, error) {
	var set operatorFileSet
	var start store.EpisodeStart
	mcp := operatorFileFact(a, "mcp_path")
	if err := json.Unmarshal([]byte(s.StartJSON), &start); err != nil || start.Task == nil || start.Allowance != nil {
		return set, start, store.ErrOperatorConflict
	}
	token, err := operatorPhysicalTarget(s.Mint.TokenPath)
	if err != nil {
		return set, start, err
	}
	set.Dir = filepath.Dir(token)
	set.Names = [3]string{a.Episode + ".task.md", a.Episode + ".session", a.Episode + ".mcp.json"}
	if filepath.Base(token) != set.Names[1] || filepath.Dir(s.Mint.TokenPath) != set.Dir ||
		start.Task.Path != filepath.Join(filepath.Dir(s.Mint.TokenPath), set.Names[0]) ||
		start.Task.SHA256 != serviceHash([]byte(a.Task)) ||
		(mcp != "") != (start.Agent == "claude") ||
		mcp != "" && mcp != filepath.Join(filepath.Dir(s.Mint.TokenPath), set.Names[2]) {
		return set, start, store.ErrOperatorConflict
	}
	if mcp == "" {
		set.Names[2] = ""
	}
	set.Hashes[0], set.Hashes[1] = start.Task.SHA256, "sha256:"+s.Mint.CredentialID
	set.UID = uint64(os.Geteuid())
	return set, start, nil
}

type operatorOpenat func(int, string, int, uint32) (int, error)

func openOperatorFilesDir(set operatorFileSet) (*os.File, operatorFileProof) {
	return openOperatorFilesDirAt(set, unix.Openat)
}

func openOperatorFilesDirAt(set operatorFileSet, open operatorOpenat) (*os.File, operatorFileProof) {
	var st unix.Stat_t
	var p operatorFileProof
	if err := unix.Lstat(set.Dir, &st); err != nil {
		p.Status = err.Error()
		return nil, p
	}
	p.ID = operatorFileStat(st)
	p.Status = "unattributed directory"
	if p.ID.Mode != unix.S_IFDIR|0700 || p.ID.UID != set.UID ||
		set.Inode != 0 && (p.ID.Dev != set.Device || p.ID.Inode != set.Inode) {
		return nil, p
	}
	const flags = unix.O_RDONLY | unix.O_DIRECTORY | unix.O_NOFOLLOW | unix.O_CLOEXEC
	fd, err := unix.Open("/", flags, 0)
	if err != nil {
		return nil, p
	}
	for _, part := range strings.Split(strings.TrimPrefix(set.Dir, "/"), "/") {
		next, err := open(fd, part, flags, 0)
		unix.Close(fd)
		if err != nil {
			return nil, p
		}
		fd = next
	}
	f := os.NewFile(uintptr(fd), set.Dir)
	if err := unix.Fstat(int(f.Fd()), &st); err != nil || operatorFileStat(st) != p.ID {
		f.Close()
		return nil, p
	}
	p.Status = "private directory"
	return f, p
}

func readOperatorFile(dir *os.File, name, want string, uid uint64) (operatorFileProof, []byte) {
	return readOperatorFileAt(dir, name, want, uid, unix.Openat, unix.Fstatat)
}

func readOperatorFileAt(dir *os.File, name, want string, uid uint64, open operatorOpenat, stat func(int, string, *unix.Stat_t, int) error) (operatorFileProof, []byte) {
	var st unix.Stat_t
	fd := int(dir.Fd())
	err := stat(fd, name, &st, unix.AT_SYMLINK_NOFOLLOW)
	if errors.Is(err, os.ErrNotExist) {
		return operatorFileProof{Status: "absent"}, nil
	}
	if err != nil {
		return operatorFileProof{Status: err.Error()}, nil
	}
	p := operatorFileProof{ID: operatorFileStat(st), Status: "unattributed file"}
	if p.ID.Mode != unix.S_IFREG|0600 || p.ID.UID != uid || p.ID.Links != 1 {
		return p, nil
	}
	raw, err := open(fd, name, unix.O_RDONLY|unix.O_NOFOLLOW|unix.O_NONBLOCK|unix.O_CLOEXEC, 0)
	if err != nil {
		return p, nil
	}
	f := os.NewFile(uintptr(raw), name)
	defer f.Close()
	if unix.Fstat(raw, &st) != nil || operatorFileStat(st) != p.ID {
		return p, nil
	}
	const limit = approval.MaxOperatorTaskBytes + 4096
	body, err := io.ReadAll(io.LimitReader(f, limit+1))
	if err != nil || len(body) > limit || stat(fd, name, &st, unix.AT_SYMLINK_NOFOLLOW) != nil || operatorFileStat(st) != p.ID {
		return p, nil
	}
	p.Digest = serviceHash(body)
	if p.Digest == want {
		p.Status = "exact"
		return p, body
	}
	return p, nil
}

func readOperatorFileSet(set operatorFileSet, addr string) operatorFilesProof {
	return readOperatorFileSetWithOpener(set, addr, openOperatorFilesDir)
}

func readOperatorFileSetWithOpener(set operatorFileSet, addr string,
	open func(operatorFileSet) (*os.File, operatorFileProof)) operatorFilesProof {
	f, dir := open(set)
	p := operatorFilesProof{Dir: dir}
	if f == nil {
		return p
	}
	defer f.Close()
	var token []byte
	for i, name := range set.Names {
		if name != "" {
			var body []byte
			p.Files[i], body = readOperatorFile(f, name, set.Hashes[i], set.UID)
			if i == 1 {
				token = body
			}
		}
	}
	if p.Files[2].Status == "exact" {
		body, err := mcpConfigBytes(addr, string(token))
		if token == nil || err != nil || serviceHash(body) != set.Hashes[2] {
			p.Files[2].Status = "unattributed file"
		}
	}
	current, reopened := open(set)
	if current != nil {
		current.Close()
	}
	if reopened != dir {
		p.Dir.Status = "directory changed during readback"
	}
	return p
}

func installOperatorFile(dir *os.File, name string, data []byte) (operatorFileID, error) {
	return installOperatorFileAt(dir, name, data, unix.Openat)
}

func installOperatorFileAt(dir *os.File, name string, data []byte, open operatorOpenat) (operatorFileID, error) {
	raw, err := open(int(dir.Fd()), name, unix.O_WRONLY|unix.O_CREAT|unix.O_EXCL|unix.O_NOFOLLOW|unix.O_CLOEXEC, 0600)
	if err != nil {
		return operatorFileID{}, err
	}
	f := os.NewFile(uintptr(raw), name)
	defer f.Close()
	if err := f.Chmod(0600); err != nil {
		return operatorFileID{}, err
	}
	if _, err := f.Write(data); err != nil {
		return operatorFileID{}, err
	}
	if err := f.Sync(); err != nil {
		return operatorFileID{}, err
	}
	var st unix.Stat_t
	err = unix.Fstat(raw, &st)
	return operatorFileStat(st), err
}

func operatorFileAddress(start store.EpisodeStart, signed, approved string) (string, error) {
	var preset launchPreset
	if err := json.Unmarshal([]byte(start.Preset.Bytes), &preset); err != nil {
		return "", store.ErrOperatorConflict
	}
	d, err := resolvePresetDaemon(preset.Daemon)
	if err != nil || d.Addr != signed || signed != approved {
		return "", store.ErrOperatorConflict
	}
	return d.Addr, nil
}

func recheckFileApproval(ctx context.Context, credentials string, a approval.OperatorAction, now int64) (*approval.OperatorMintProof, error) {
	db, err := store.OpenCredentialsReadOnly(runnerDB(credentials))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row, err := db.OperatorAction(ctx, a.ID, now)
	if err != nil {
		return nil, err
	}
	if now >= row.ExpiresAt {
		return nil, store.ErrOperatorExpired
	}
	if len(row.Challenge) != 32 {
		return nil, approval.ErrNotVerified
	}
	proof, err := authority.RecheckClaimedOperatorAction(ctx, db, a, [32]byte(row.Challenge), row.ExpiresAt, now, row)
	if err != nil {
		return nil, err
	}
	return proof.OperatorConfirmation().MintProof(), nil
}

// publishOperatorStartFiles is dormant until the paired-consumer gates open.
// rawToken is the existing mint output. Neither it nor the MCP bytes enter SQL.
func publishOperatorStartFiles(ctx context.Context, c runnerConfig, held *store.OperatorEpisodeLock, id, rawToken string, clock func() int64, checkpoint func(string) error) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	s, a, err := readRunnerFiles(ctx, c.Credentials, id)
	if err != nil {
		return err
	}
	if _, err := held.Verify(c.Credentials, a.Episode); err != nil {
		return err
	}
	current, a, err := readRunnerFiles(ctx, c.Credentials, id)
	if err != nil || current != s || s.Phase != "" {
		return store.ErrOperatorConflict
	}
	set, start, err := operatorFilePaths(s, a)
	if err != nil || serviceHash([]byte(rawToken)) != set.Hashes[1] {
		return store.ErrOperatorConflict
	}
	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
	if err != nil {
		return err
	}
	data := [3][]byte{[]byte(a.Task), []byte(rawToken), nil}
	if set.Names[2] != "" {
		data[2], err = mcpConfigBytes(addr, rawToken)
		if err != nil {
			return err
		}
		set.Hashes[2] = serviceHash(data[2])
	}
	var target operatorWorktreeTarget
	path, err := operatorPhysicalTarget(operatorWorktreeFact(a))
	if err != nil || json.Unmarshal([]byte(s.Worktree), &target) != nil ||
		target.Action != id || target.Path != path || target.Branch != start.Branch || target.SHA != start.Base.SHA {
		return store.ErrOperatorConflict
	}
	git, err := sessionGit()
	if err != nil {
		return err
	}
	repo, err := filepath.EvalSymlinks(start.Repo)
	if err != nil || repo != target.Repo {
		return store.ErrOperatorConflict
	}
	worktreeBefore, err := target.inspect(ctx, git)
	if err != nil || worktreeBefore.state != treeExact {
		return store.ErrOperatorConflict
	}
	dir, proof := openOperatorFilesDir(set)
	if dir == nil {
		return store.ErrOperatorConflict
	}
	defer dir.Close()
	set.Device, set.Inode = proof.ID.Dev, proof.ID.Inode
	for i, name := range set.Names {
		if name != "" {
			file, _ := readOperatorFile(dir, name, set.Hashes[i], set.UID)
			if file.Status != "absent" {
				return store.ErrOperatorConflict
			}
		}
	}
	if err := checkpoint("before reservation"); err != nil {
		return err
	}
	approvalProof, err := recheckFileApproval(ctx, c.Credentials, a, clock())
	if err != nil {
		return err
	}
	record := func(value string) error {
		return runnerWrite(c.Credentials, func(db *store.Credentials) error {
			return db.RecordOperatorFiles(ctx, s, value, approvalProof, clock)
		})
	}
	wire, _ := json.Marshal(set) // scalar identity, no raw credential
	if err := record(string(wire)); err != nil {
		return err
	}
	s.Step.Identity, s.Step.Version, s.Phase = string(wire), 1, "reserved"
	remedy := reconcileFilesCommand(id)
	outcome := func(err error) error { return fmt.Errorf("%w; %s", err, remedy) }
	uncertain := func() error {
		if err := record(store.OperatorFilesUncertain); err != nil {
			return outcome(err)
		}
		return outcome(errors.New("file publication uncertain"))
	}
	var installed [3]operatorFileID
	for i, name := range set.Names {
		if name == "" {
			continue
		}
		if err := checkpoint("before " + name); err != nil {
			return outcome(err)
		}
		current, action, err := readRunnerFiles(ctx, c.Credentials, id)
		now := clock()
		if err != nil || current != s || now >= current.Session.ExpiresAt {
			return uncertain()
		}
		if _, err := recheckFileApproval(ctx, c.Credentials, action, now); err != nil {
			return uncertain()
		}
		fresh, _ := openOperatorFilesDir(set)
		if fresh == nil {
			return uncertain()
		}
		fresh.Close()
		installed[i], err = installOperatorFile(dir, name, data[i])
		if err != nil {
			return uncertain()
		}
	}
	if err := dir.Sync(); err != nil {
		return uncertain()
	}
	if err := checkpoint("before readback"); err != nil {
		return outcome(err)
	}
	got := readOperatorFileSet(set, addr)
	if got.Dir.Status != "private directory" {
		return uncertain()
	}
	for i, name := range set.Names {
		if name != "" && (got.Files[i].Status != "exact" || got.Files[i].ID != installed[i]) {
			return uncertain()
		}
	}
	if err := checkpoint("before observe"); err != nil {
		return outcome(err)
	}
	worktreeAfter, err := target.inspect(ctx, git)
	if err != nil || worktreeAfter != worktreeBefore || readOperatorFileSet(set, addr) != got {
		return uncertain()
	}
	approvalProof, err = recheckFileApproval(ctx, c.Credentials, a, clock())
	if err != nil {
		return uncertain()
	}
	if err := record(store.OperatorFilesPublished); err != nil {
		return outcome(err)
	}
	return nil
}

func readRunnerFiles(ctx context.Context, creds, id string) (store.OperatorFileState, approval.OperatorAction, error) {
	db, err := store.OpenCredentialsReadOnly(runnerDB(creds))
	if err != nil {
		return store.OperatorFileState{}, approval.OperatorAction{}, err
	}
	defer db.Close()
	read, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	s, a, err := db.ReadOperatorFiles(read, id)
	if err == nil && (s.Mint.Store != creds || ctx.Err() != nil) {
		err = store.ErrOperatorConflict
	}
	return s, a, err
}

func operatorFileDigest(s string) bool {
	b, err := hex.DecodeString(strings.TrimPrefix(s, "sha256:"))
	return strings.HasPrefix(s, "sha256:") && err == nil && len(b) == 32
}
