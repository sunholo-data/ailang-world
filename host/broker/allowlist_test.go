package broker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var allowedDepModules = []string{
	"github.com/sunholo-data/ailang-world",
	"modernc.org/sqlite",
	"modernc.org/libc",
	"modernc.org/mathutil",
	"modernc.org/memory",
	"golang.org/x/sys",
	"github.com/dustin/go-humanize",
	"github.com/google/uuid",
	"github.com/mattn/go-isatty",
	"github.com/ncruces/go-strftime",
	"github.com/remyoudompheng/bigfft",
}

var brokerCorePatterns = []string{"./host/broker/...", "./host/capsule/..."}

func isStdlibImportPath(p string) bool {
	first, _, _ := strings.Cut(p, "/")
	return !strings.Contains(first, ".")
}

func disallowedDeps(deps []string) ([]string, error) {
	if len(deps) == 0 {
		return nil, errors.New("dependency list is empty: the allowlist would pass vacuously")
	}
	var bad []string
	for _, dep := range deps {
		if isStdlibImportPath(dep) {
			continue
		}
		allowed := false
		for _, module := range allowedDepModules {
			if dep == module || strings.HasPrefix(dep, module+"/") {
				allowed = true
				break
			}
		}
		if !allowed {
			bad = append(bad, dep)
		}
	}
	return bad, nil
}

func goListDeps(repoRoot string, patterns ...string) ([]string, error) {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return nil, fmt.Errorf("cannot locate the `go` toolchain on PATH: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, goBin, append([]string{"list", "-deps"}, patterns...)...)
	cmd.Dir = repoRoot
	var stdout, stderr strings.Builder
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("%s list -deps %s: %w\nstderr: %s",
			goBin, strings.Join(patterns, " "), err, stderr.String())
	}
	var deps []string
	scanner := bufio.NewScanner(strings.NewReader(stdout.String()))
	for scanner.Scan() {
		if line := strings.TrimSpace(scanner.Text()); line != "" {
			deps = append(deps, line)
		}
	}
	return deps, scanner.Err()
}

func TestBrokerDependencyAllowlist(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	deps, err := goListDeps(root, brokerCorePatterns...)
	if err != nil {
		t.Fatal(err)
	}
	bad, err := disallowedDeps(deps)
	if err != nil {
		t.Fatal(err)
	}
	if len(bad) != 0 {
		t.Fatalf("broker dependency allowlist rejected: %s", strings.Join(bad, ", "))
	}
}

func TestBrokerDependencyAllowlistNullCase(t *testing.T) {
	if _, err := disallowedDeps(nil); err == nil {
		t.Fatal("empty dependency list passed vacuously")
	}
}
