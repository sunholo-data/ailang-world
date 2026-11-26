package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sunholo-data/ailang-world/host/store"
)

func reconcileFilesCommand(id string) string {
	return "ailang-worldd service runner reconcile-files " + id
}

func (e serviceEnv) reconcileRunnerFiles(id string, out, stderr io.Writer) int {
	remedy := reconcileFilesCommand(id)
	refuse := func(err error) int {
		return serviceMessage(stderr, "file recovery refused", err.Error(), remedy, exitUsage)
	}
	c, _, err := e.installedRunner()
	if err != nil {
		return refuse(err)
	}
	ctx, cancel := budgetContext(5 * time.Second)
	defer cancel()
	s, a, err := readRunnerFiles(ctx, c.Credentials, id)
	if err != nil {
		return refuse(err)
	}
	if s.Phase != "reserved" {
		return refuse(store.ErrOperatorMissing)
	}
	lock, err := openRunnerLock(operatorEpisodeLockPath(c.Credentials, a.Episode), 0)
	if err != nil {
		return refuse(errors.New("episode recovery lock held or unavailable"))
	}
	defer lock.Close()
	current, a, err := readRunnerFiles(ctx, c.Credentials, id)
	if err != nil || current != s {
		return refuse(store.ErrOperatorConflict)
	}
	want, start, err := operatorFilePaths(s, a)
	var set operatorFileSet
	if err != nil || json.Unmarshal([]byte(s.Step.Identity), &set) != nil || set.Dir != want.Dir ||
		set.UID != want.UID || set.Inode == 0 || set.Names != want.Names ||
		set.Hashes[0] != want.Hashes[0] || set.Hashes[1] != want.Hashes[1] ||
		(set.Names[2] != "" && !operatorFileDigest(set.Hashes[2])) ||
		(set.Names[2] == "" && set.Hashes[2] != "") {
		return refuse(store.ErrOperatorConflict)
	}
	addr, err := operatorFileAddress(start, operatorFileFact(a, "served_address"), c.Address)
	if err != nil {
		return refuse(err)
	}
	before := readOperatorFileSet(set, addr)
	absent, exact := before.Dir.Status == "private directory", before.Dir.Status == "private directory"
	fmt.Fprintln(stderr, "File directory:", before.Dir.Status)
	for i, name := range set.Names {
		if name != "" {
			fmt.Fprintln(stderr, name, before.Files[i].Status)
			absent = absent && before.Files[i].Status == "absent"
			exact = exact && before.Files[i].Status == "exact"
		}
	}
	if !absent && !exact {
		return refuse(errors.New("file evidence remains uncertain"))
	}
	cancel()
	disposition := store.OperatorFilesStoppedRetained
	if absent {
		disposition = store.OperatorFilesStoppedAbsent
	}
	if !e.confirm(stderr, fmt.Sprintf("Stop action %s and keep its files, credential and worktree as found?", id)) {
		fmt.Fprintln(stderr, "Retry at a terminal:", remedy)
		return exitUsage
	}
	if after := readOperatorFileSet(set, addr); after != before {
		return refuse(errors.New("file evidence changed after review"))
	}
	err = runnerWrite(c.Credentials, func(db *store.Credentials) error {
		ctx, cancel := budgetContext(time.Second)
		defer cancel()
		return db.RecordOperatorFiles(ctx, s, disposition, nil, e.terminal.now)
	})
	if err != nil {
		return refuse(err)
	}
	fmt.Fprintln(out, "File reservation disposed for", id)
	return exitOK
}
