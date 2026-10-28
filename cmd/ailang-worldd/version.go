package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/sunholo-data/ailang-world/host/daemon"
)

const versionHelp = "usage: ailang-worldd version [--json]\n"

func runVersion(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("version", flag.ContinueOnError)
	fs.SetOutput(stderr)
	asJSON := fs.Bool("json", false, "emit JSON")
	if code, done := parseVerbFlags(fs, args, versionHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() != 0 {
		fmt.Fprint(stderr, versionHelp)
		return exitUsage
	}
	v := struct {
		Build              daemon.BuildIdentity `json:"build"`
		Platform           string               `json:"platform"`
		ToolRelease        string               `json:"tool_release"`
		InterpreterRelease string               `json:"interpreter_release"`
	}{daemon.CurrentBuild(), runtime.GOOS + "/" + runtime.GOARCH, toolRelease(), interpreterRelease}
	if *asJSON {
		if err := json.NewEncoder(stdout).Encode(v); err != nil {
			return exitUsage
		}
	} else {
		fmt.Fprintf(stdout, "ailang-worldd %s\nplatform: %s\nbuild: %s\ngo: %s\ntool: %s\ninterpreter: %s\n", v.Build.String(), v.Platform, v.Build.ExecutableSHA256, v.Build.GoVersion, v.ToolRelease, v.InterpreterRelease)
	}
	return exitOK
}
