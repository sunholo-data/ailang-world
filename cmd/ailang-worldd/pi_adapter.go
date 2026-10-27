package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const piAdapterWhy = "pi has no built-in MCP; World's tools come through that extension"
const piAdapterFix = "pi install npm:pi-mcp-adapter, then ailang-worldd hosts init --replace"

func defaultPiMCPAdapter() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".pi", "agent", "npm", "node_modules", "pi-mcp-adapter")
}

func piAdapterWhat(dir string) string {
	return "pi can't reach World: no pi-mcp-adapter at " + dir
}

func piMCPAdapterDigest(dir string) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return "", err
	}
	var pkg struct {
		Pi struct {
			Extensions []string `json:"extensions"`
		} `json:"pi"`
	}
	if err := json.Unmarshal(b, &pkg); err != nil {
		return "", err
	}
	if len(pkg.Pi.Extensions) == 0 {
		return "", fmt.Errorf("package.json must name pi.extensions")
	}
	for _, extension := range pkg.Pi.Extensions {
		if extension == "" {
			return "", fmt.Errorf("pi.extensions contains an empty path")
		}
	}
	return fmt.Sprintf("sha256:%x", sha256.Sum256(b)), nil
}

func doctorPiAdapter(r *doctorReport, dir string) {
	b, err := os.ReadFile(filepath.Join(dir, "hosts.json"))
	var manifest hostsManifest
	if err != nil || json.Unmarshal(b, &manifest) != nil {
		r.add(markFail, "host", "pi adapter manifest unavailable", piAdapterFix)
		return
	}
	digest, err := piMCPAdapterDigest(manifest.PiMCPAdapter)
	if manifest.PiMCPAdapter == "" || err != nil {
		r.add(markFail, "host", piAdapterWhat(manifest.PiMCPAdapter)+"; "+piAdapterWhy, piAdapterFix)
	} else if digest != manifest.PiMCPAdapterSHA256 {
		r.add(markFail, "host", "pi-mcp-adapter package.json changed at "+manifest.PiMCPAdapter, piAdapterFix)
	} else {
		r.add(markOK, "host", "pi-mcp-adapter verified at "+manifest.PiMCPAdapter, "")
	}
}
