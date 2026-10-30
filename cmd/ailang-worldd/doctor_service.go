package main

import (
	"fmt"
	"os"
	"regexp"
)

var doctorServiceEnv = realServiceEnv

// doctorService queries the supervisor only for the service being inspected.
func doctorService(report *doctorReport, addr string) {
	r, ok := homeServiceRecord()
	if !ok {
		report.add(markWarn, "service", "no World service is installed", "ailang-worldd service install")
		return
	}
	e := doctorServiceEnv()
	detail := fmt.Sprintf("World service at %s; plist %s; installed binary %s", r.Address, e.paths.Plist, e.paths.Binary)
	if !serviceAddressMatches(addr, r) {
		report.add(markOK, "service", detail, "")
	} else {
		state, err := e.controller.Print(e.domain, e.label)
		switch {
		case err != nil:
			report.add(markFail, "service", detail+"; the supervisor reports "+state, "ailang-worldd service start")
		case report.serviceBuild.ExecutableSHA256 == serviceFileHash(r.Binary) && report.serviceBuild.ExecutableSHA256 != "":
			pid := regexp.MustCompile("(?m)pid = ([0-9]+)").FindStringSubmatch(state)
			if len(pid) > 1 {
				detail += "; pid " + pid[1]
			}
			report.add(markOK, "service", detail, "")
		default:
			report.add(markWarn, "service", detail+"; restart pending", "ailang-worldd service upgrade")
		}
	}
	if serviceFileHash(e.paths.Plist) != r.PlistSHA256 {
		report.add(markWarn, "service", "the plist differs from the one installed", "ailang-worldd service start shows the diff, or service install --replace")
	}
	path, _ := os.Executable()
	cli, installed := serviceFileHash(path), serviceFileHash(r.Binary)
	if cli != installed {
		report.add(markWarn, "service", fmt.Sprintf("this ailang-worldd (%s, sha %s) is not the installed one (%s)", path, pointerSHA7(cli), pointerSHA7(installed)),
			fmt.Sprintf("run %s, or install this one: %s service upgrade", r.Binary, path))
	}
	if info, err := os.Stat(e.paths.Log); err == nil && info.Size() > 100<<20 {
		report.add(markWarn, "service", "the service log is over 100 MiB: "+e.paths.Log, "rotate the log while the service is stopped")
	}
}
