package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

type runnerUpgradePair struct {
	config        runnerConfig
	before, after []byte
	loaded        bool
}

func (e serviceEnv) runnerUpgradeReview() [2]string {
	p := e.runnerPaths()
	return [2]string{serviceFileHash(p.Record), serviceFileHash(p.Plist)}
}

func (e serviceEnv) reviewRunnerUpgrade(approved [2]string, from, to string, schemas ...int) (*runnerUpgradePair, error) {
	if approved != e.runnerUpgradeReview() {
		return nil, errors.New("operator runner changed during confirmation")
	}
	p := e.runnerPaths()
	_, recordErr := os.Lstat(p.Record)
	_, plistErr := os.Lstat(p.Plist)
	_, loadedErr := e.controller.Print(e.domain, runnerLabel)
	if errors.Is(recordErr, os.ErrNotExist) && errors.Is(plistErr, os.ErrNotExist) {
		if !errors.Is(loadedErr, errServiceNotLoaded) {
			return nil, errors.New("cannot verify that the operator runner is unloaded")
		}
		return nil, nil
	}
	c, seen, err := e.installedRunner()
	if err != nil {
		return nil, errors.New("operator runner is still installed but its configuration and plist cannot be verified")
	}
	if approved[0] != serviceHash(seen) || c.Executables[e.paths.Binary] != from {
		return nil, errors.New("operator runner changed during confirmation")
	}
	if loadedErr != nil && !errors.Is(loadedErr, errServiceNotLoaded) {
		return nil, errors.New("cannot verify that the operator runner is unloaded")
	}
	if err := e.checkRunnerPair(c); err != nil {
		return nil, err
	}
	pair := &runnerUpgradePair{config: c, before: seen, loaded: loadedErr == nil}
	if err := pair.checkSchema(schemas...); err != nil {
		return nil, err
	}
	c.Executables[e.paths.Binary] = to
	pair.after, err = json.MarshalIndent(c, "", "  ")
	c.Executables[e.paths.Binary] = from
	return pair, err
}

func (p *runnerUpgradePair) checkSchema(schemas ...int) error {
	db, err := store.OpenCredentialsReadOnly(runnerDB(p.config.Credentials))
	if err != nil {
		return err
	}
	defer db.Close()
	for _, schema := range schemas {
		if schema != db.Version() {
			return errors.New("paired upgrade requires both builds to match the actual credentials schema")
		}
	}
	return nil
}

func (e serviceEnv) waitRunnerExit(c runnerConfig) (*os.File, error) {
	return e.waitRunnerExitUntil(c, e.runnerNow().Add(min(e.wait, serviceExitTimeout)))
}

func (e serviceEnv) waitRunnerExitUntil(c runnerConfig, deadline time.Time) (*os.File, error) {
	for {
		lock, err := runnerKernelLock(c.Credentials + ".operator-runner.lock")
		if err == nil || !errors.Is(err, syscall.EWOULDBLOCK) {
			return lock, err
		}
		remaining := deadline.Sub(e.runnerNow())
		if remaining <= 0 {
			return nil, errors.New("companion or retained child still holds the runner lock")
		}
		e.runnerSleep(min(remaining, 100*time.Millisecond))
	}
}

func (e serviceEnv) stopRunnerJob(budget time.Duration) error {
	_, err := e.controller.Print(e.domain, runnerLabel)
	if errors.Is(err, errServiceNotLoaded) {
		return nil
	}
	if err != nil {
		return errors.New("cannot verify that the operator runner is unloaded")
	}
	if controller, ok := e.controller.(systemServiceController); ok {
		var output string
		output, err = controller.commandWithin(budget, "bootout", e.domain+"/"+runnerLabel)
		if err != nil && strings.Contains(output, "Could not find service") {
			err = errServiceNotLoaded
		}
	} else {
		err = e.controller.Bootout(e.domain, runnerLabel)
	}
	if err != nil && !errors.Is(err, errServiceNotLoaded) {
		return errors.New("could not stop the companion")
	}
	if _, err := e.controller.Print(e.domain, runnerLabel); !errors.Is(err, errServiceNotLoaded) {
		return errors.New("cannot verify that the operator runner is unloaded")
	}
	return nil
}

func (e serviceEnv) drainRunner(c runnerConfig) (*os.File, error) {
	budget := min(e.wait, serviceExitTimeout)
	deadline := e.runnerNow().Add(budget)
	if err := e.stopRunnerJob(budget); err != nil {
		return nil, err
	}
	if !e.runnerNow().Before(deadline) {
		return nil, errors.New("companion drain budget expired")
	}
	return e.waitRunnerExitUntil(c, deadline)
}

func (e serviceEnv) proveOldTerminalQuiet() error {
	allowed := map[int]bool{os.Getpid(): true}
	for _, label := range []string{e.label, runnerLabel} {
		job, err := e.controller.Print(e.domain, label)
		if errors.Is(err, errServiceNotLoaded) {
			continue
		}
		if err != nil {
			return errors.New("cannot inventory supervised World processes")
		}
		pid := 0
		for _, line := range strings.Split(job, "\n") {
			if raw, ok := strings.CutPrefix(strings.TrimSpace(line), "pid = "); ok {
				pid, _ = strconv.Atoi(raw)
				break
			}
		}
		if pid <= 0 {
			return errors.New("cannot identify a loaded World process")
		}
		allowed[pid] = true
	}
	read := e.runnerProcesses
	if read == nil {
		read = listWorldProcesses
	}
	rows, err := read()
	if err != nil {
		return errors.New("cannot verify old terminal process retirement")
	}
	self := false
	for _, process := range rows {
		self = self || process.pid == os.Getpid()
		if !allowed[process.pid] && strings.HasPrefix(process.name, "ailang-worldd") {
			return fmt.Errorf("old terminal process retirement not proved (pid %d)", process.pid)
		}
	}
	if !self {
		return errors.New("World process inventory did not include this process")
	}
	return nil
}

type worldProcess struct {
	pid  int
	name string
}

func (e serviceEnv) writeRunnerPair(pair *runnerUpgradePair, before bool) error {
	data := pair.after
	if before {
		data = pair.before
	}
	var c runnerConfig
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}
	p := e.runnerPaths()
	if err := e.checkRunnerConfig(c); err != nil {
		return err
	}
	if err := e.checkRunnerPair(c); err != nil {
		return err
	}
	if err := e.atomic(p.Record, data, 0600); err != nil {
		return err
	}
	if err := e.atomic(p.Plist, []byte(runnerPlist(p, c, data)), 0600); err != nil {
		return err
	}
	if pair.loaded {
		return e.changeRunner("start", nil)
	}
	return nil
}
