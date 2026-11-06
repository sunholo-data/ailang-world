// Package checkrun runs finish's attended subprocesses: checks and gh.
// It owns one bounded subprocess site, driven by the AC10 census. Checks
// execute agent code, so credentials are always removed from their environment.
package checkrun

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/childenv"
	"github.com/sunholo-data/ailang-world/host/procbound"
)

func checkEnv(inherited []string, env map[string]string, pathPrefix []string) []string {
	values := map[string]string{}
	for _, kv := range childenv.Scrubbed(inherited) {
		key, value, ok := strings.Cut(kv, "=")
		if ok && key != "WORLD_SESSION" {
			values[key] = value
		}
	}
	for key, value := range env {
		if key != "WORLD_SESSION" {
			values[key] = value
		}
	}
	if len(pathPrefix) > 0 {
		parts := append([]string{}, pathPrefix...)
		if values["PATH"] != "" {
			parts = append(parts, values["PATH"])
		}
		values["PATH"] = strings.Join(parts, string(os.PathListSeparator))
	}
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+"="+values[key])
	}
	return result
}

const HeadBytes = 400

type Result struct {
	ExitCode int
	Stdout   string
	Stderr   string
}

type head struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (h *head) Write(p []byte) (int, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	n := len(p)
	if room := HeadBytes - h.buf.Len(); room > 0 {
		h.buf.Write(p[:min(room, n)])
	}
	return n, nil
}

func (h *head) String() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.buf.String()
}

// Run enforces the wall cap across the whole process group. Output heads are
// bounded while the rest is drained; stdin is never inherited.
func Run(ctx context.Context, argv []string, dir string, timeout time.Duration) (Result, error) {
	return RunEnv(ctx, argv, dir, timeout, nil, nil)
}

// RunEnv overlays the pinned profile on the scrubbed inherited environment.
func RunEnv(ctx context.Context, argv []string, dir string, timeout time.Duration, env map[string]string, pathPrefix []string) (Result, error) {
	r := Result{ExitCode: -1}
	if len(argv) == 0 || timeout <= 0 {
		return r, fmt.Errorf("checkrun: argv and a positive wall cap are required")
	}
	runCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	cmd := exec.CommandContext(runCtx, argv[0], argv[1:]...)
	cmd.Dir = dir
	cmd.Env = checkEnv(os.Environ(), env, pathPrefix)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error { return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL) }
	cmd.WaitDelay = 100 * time.Millisecond
	var out, errout head
	cmd.Stdout, cmd.Stderr = &out, &errout
	release, err := procbound.Admit()
	if err != nil {
		return r, err
	}
	if err := cmd.Start(); err != nil {
		release()
		return r, err
	}
	err = procbound.Wait(cmd.Wait, timeout+time.Second, release)
	if errors.Is(err, procbound.ErrCleanupIncomplete) {
		return r, err
	}
	// Also kill descendants after the direct child exits successfully. Otherwise
	// a background process could outlive the check and mutate after attribution.
	killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	if killErr != nil && !procbound.NothingToKill(killErr, func() bool { return procbound.GroupEmpty(cmd.Process.Pid) }) {
		err = errors.Join(err, killErr)
	}
	r.Stdout, r.Stderr = out.String(), errout.String()
	if cmd.ProcessState != nil {
		r.ExitCode = cmd.ProcessState.ExitCode()
	}
	if runCtx.Err() != nil {
		return r, runCtx.Err()
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && (killErr == nil || procbound.NothingToKill(killErr, func() bool { return procbound.GroupEmpty(cmd.Process.Pid) })) {
		return r, nil
	}
	return r, err
}
