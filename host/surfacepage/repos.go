package surfacepage

import _ "embed"

// Row 269 Phase 1 (design_docs/planned/w-ui-profile-preset-onboarding.md):
// standalone read-only Repos, Presets & Exec Profiles page served at
// GET /surface/repos and GET /surface/repos.js.

// ReposPage is GET /surface/repos.
//
//go:embed repos/index.html
var ReposPage []byte

// ReposScript is GET /surface/repos.js.
//
//go:embed repos/repos.js
var ReposScript []byte
