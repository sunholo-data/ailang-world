package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// Tests supply a private read-write file; production requires a controlling tty.
var openLaunchTerminal = func() (*os.File, error) {
	return os.OpenFile("/dev/tty", os.O_RDWR, 0)
}

func launchLauncher(agent, hostsDir, addr string) (string, error) {
	dir, err := filepath.Abs(hostsDir)
	if err != nil {
		return "", err
	}
	next := "ailang-worldd hosts init --out-dir " + shellQuote(dir) + " --addr http://" + addr + " --replace"
	var output bytes.Buffer
	doctorHost(&doctorReport{w: &output}, agent, dir)
	var failures []string
	for _, line := range strings.Split(output.String(), "\n") {
		if strings.HasPrefix(line, string(markWarn)) || strings.HasPrefix(line, string(markFail)) {
			failures = append(failures, line)
		}
	}
	if len(failures) > 0 {
		return "", presetRefusal(strings.Join(failures, "\n")+" (F10)", "interactive launch requires a current confined launcher", next)
	}
	data, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
	var manifest hostsManifest
	if err != nil {
		return "", err
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return "", err
	}
	got := strings.TrimRight(strings.TrimPrefix(manifest.Addr, "http://"), "/")
	if got != addr {
		return "", presetRefusal(fmt.Sprintf("launcher address %s differs from preset address %s (F10)", manifest.Addr, addr), "the launcher must use the preset daemon", next)
	}
	return filepath.Join(dir, "bin", agent+"-world"), nil
}

func launchSessionUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:]), nil
}

func interactiveLaunchArgv(host, agent, model, sessionID, context, prompt string, grants []broker.Capability) []string {
	a := []string{host}
	if agent == "claude" {
		a = append(a, "--session-id", sessionID, "--allowedTools", strings.Join(allowedTools(grants), ","))
		if model != "" {
			a = append(a, "--model", model)
		}
		a = append(a, "--append-system-prompt", context)
	} else {
		if model != "" {
			a = append(a, "-m", model)
		}
		a = append(a, "-c", "developer_instructions="+pythonJSONString(context))
	}
	if prompt != "" {
		a = append(a, prompt)
	}
	return a
}

func claudeTranscriptSlug(dir string) string {
	b := []byte(dir)
	for i, c := range b {
		if !(c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9') {
			b[i] = '-'
		}
	}
	return string(b)
}

func launchEnvValue(env []string, name string) string {
	for i := len(env) - 1; i >= 0; i-- {
		if v, ok := strings.CutPrefix(env[i], name+"="); ok {
			return v
		}
	}
	return ""
}

// hostTranscript identifies one source without modifying any host-owned file.
type hostTranscript struct {
	path, id, reason string
}

func regularHostTranscript(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode().IsRegular()
}

func claudeTranscriptPath(dir, id string, env []string) hostTranscript {
	root := launchEnvValue(env, "CLAUDE_CONFIG_DIR")
	if root == "" {
		root = filepath.Join(launchEnvValue(env, "HOME"), ".claude")
	}
	real, err := filepath.EvalSymlinks(dir)
	if err != nil {
		return hostTranscript{reason: fmt.Sprintf("resolve launch directory %s: %v", dir, err)}
	}
	slug := claudeTranscriptSlug(real)
	projects := filepath.Join(root, "projects")
	path := filepath.Join(projects, slug, id+".jsonl")
	if len(slug) <= 200 {
		if regularHostTranscript(path) {
			return hostTranscript{path: path, id: id}
		}
		if _, err := os.Lstat(path); err == nil {
			return hostTranscript{reason: "not a regular file: " + path}
		}
		return hostTranscript{reason: "no " + path}
	}
	prefix := slug[:200] + "-"
	entries, _ := os.ReadDir(projects)
	var matches []string
	var rejected string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), prefix) {
			p := filepath.Join(projects, entry.Name(), id+".jsonl")
			if regularHostTranscript(p) {
				matches = append(matches, p)
			} else if _, err := os.Lstat(p); err == nil {
				rejected = "not a regular file: " + p
			}
		}
	}
	if len(matches) == 1 {
		return hostTranscript{path: matches[0], id: id}
	}
	if len(matches) == 0 && rejected != "" {
		return hostTranscript{reason: rejected}
	}
	return hostTranscript{reason: fmt.Sprintf("%d transcripts under %s/%s* for %s", len(matches), projects, prefix, id)}
}

func codexTranscriptPath(hostsDir, dir string, start, end time.Time) hostTranscript {
	root := filepath.Join(hostsDir, "codex", "sessions")
	lo, hi := start.Add(-time.Second), end.Add(time.Second)
	first := lo.In(time.Local)
	day := time.Date(first.Year(), first.Month(), first.Day(), 0, 0, 0, 0, time.Local)
	last := hi.In(time.Local).Format("2006/01/02")
	real, _ := filepath.EvalSymlinks(dir)
	var matches []hostTranscript
	var rejected string
	for {
		date := day.Format("2006/01/02")
		entries, _ := os.ReadDir(filepath.Join(root, date))
		for _, entry := range entries {
			name := entry.Name()
			// rollout-YYYY-MM-DDThh-mm-ss-<uuid>.jsonl
			if len(name) <= 28 || !strings.HasPrefix(name, "rollout-") || !strings.HasSuffix(name, ".jsonl") || name[27] != '-' {
				continue
			}
			stamp, err := time.ParseInLocation("2006-01-02T15-04-05", name[8:27], time.Local)
			if err != nil || stamp.Before(lo) || stamp.After(hi) {
				continue
			}
			path := filepath.Join(root, date, name)
			if !regularHostTranscript(path) {
				rejected = "not a regular file: " + path
				continue
			}
			f, err := os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
			if err != nil {
				rejected = fmt.Sprintf("not a regular file: %s: %v", path, err)
				continue
			}
			info, statErr := f.Stat()
			if statErr != nil || !info.Mode().IsRegular() {
				_ = f.Close()
				rejected = "not a regular file: " + path
				continue
			}
			var meta struct {
				Type    string `json:"type"`
				Payload struct {
					ID  string `json:"id"`
					Cwd string `json:"cwd"`
				} `json:"payload"`
			}
			err = json.NewDecoder(io.LimitReader(f, 4<<20)).Decode(&meta)
			_ = f.Close()
			if err == nil && meta.Type == "session_meta" && (meta.Payload.Cwd == dir || real != "" && meta.Payload.Cwd == real) {
				matches = append(matches, hostTranscript{path: path, id: meta.Payload.ID})
			}
		}
		if date == last || day.After(hi) {
			break
		}
		day = day.AddDate(0, 0, 1)
	}
	if len(matches) == 1 {
		return matches[0]
	}
	if len(matches) == 0 && rejected != "" {
		return hostTranscript{reason: rejected}
	}
	return hostTranscript{reason: fmt.Sprintf("%d rollouts under %s for %s", len(matches), root, dir)}
}

func copyHostTranscript(source hostTranscript, stage *os.File) error {
	if !regularHostTranscript(source.path) {
		return fmt.Errorf("host transcript is not a regular file: %s", source.path)
	}
	f, err := os.OpenFile(source.path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return fmt.Errorf("host transcript is not a regular file: %s: %w", source.path, err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("host transcript is not a regular file: %s", source.path)
	}
	const capBytes = 64 << 20
	if info.Size() > capBytes {
		return fmt.Errorf("host transcript larger than 64 MiB: %s", source.path)
	}
	n, err := io.Copy(stage, io.LimitReader(f, capBytes+1))
	if n > capBytes {
		return fmt.Errorf("host transcript larger than 64 MiB: %s", source.path)
	}
	return err
}
