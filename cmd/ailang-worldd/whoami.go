package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

var whoamiTransport http.RoundTripper = http.DefaultTransport

const whoamiHelp = "usage: ailang-worldd [--addr <url>] whoami [--json] [--session <file>]\nSession comes from WORLD_SESSION or a private file, never an argv token.\n"

func runWhoami(addr string, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("whoami", flag.ContinueOnError)
	// Flag parsing must never echo a credential supplied as an unknown flag.
	fs.SetOutput(io.Discard)
	file := fs.String("session", "", "session file")
	asJSON := fs.Bool("json", false, "JSON output")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			fmt.Fprint(stdout, whoamiHelp)
			return 0
		}
		fmt.Fprint(stderr, whoamiHelp)
		return exitUsage
	}
	if fs.NArg() != 0 {
		fmt.Fprint(stderr, whoamiHelp)
		return exitUsage
	}
	fail := func(what, why, next string, code int) int {
		fmt.Fprintf(stderr, "what: %s\nwhy: %s\nnext: %s\n", what, why, next)
		return code
	}
	token := strings.TrimSpace(os.Getenv("WORLD_SESSION"))
	if *file != "" {
		b, err := os.ReadFile(*file)
		if err != nil {
			return fail("session unavailable", "cannot read session file", "set WORLD_SESSION or pass --session <file>", 11)
		}
		token = strings.TrimSpace(string(b))
	}
	valid := len(token) == 64
	for _, c := range token {
		valid = valid && ((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F'))
	}
	if !valid {
		return fail("session invalid or expired", "missing or malformed session credential", "set WORLD_SESSION or ask the operator to mint a new session", 11)
	}
	client := &http.Client{Transport: whoamiTransport, Timeout: 5 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(addr, "/")+"/v1/session/self", nil)
	if err != nil {
		return fail("daemon unavailable", "invalid daemon address", "pass --addr <World URL>", 10)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return fail("daemon unavailable", "World did not answer", "start ailang-worldd serve, then retry whoami", 10)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized {
		return fail("session invalid or expired", "World refused this credential", "ask the operator to mint a new session", 11)
	}
	if resp.StatusCode != http.StatusOK {
		return fail("daemon unavailable", "World could not read the session", "run ailang-worldd doctor", 10)
	}
	var self daemon.SessionSelf
	if json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&self) != nil {
		return fail("daemon unavailable", "invalid World response", "run ailang-worldd doctor", 10)
	}
	if *asJSON {
		_ = json.NewEncoder(stdout).Encode(self)
	} else {
		fmt.Fprintf(stdout, "World · %s\nworktree: %s\nsession: %dh%dm left\ndaemon: %s\nsurface: %s\n", self.Episode, self.Worktree, self.SecondsLeft/3600, self.SecondsLeft/60%60, self.Addr, self.Surface)
		for _, g := range self.Grants {
			fmt.Fprintf(stdout, "  %s %s: %d/%d left\n", g.Effect, g.Scope, g.Left, g.Budget)
		}
	}
	return 0
}
