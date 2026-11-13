package surfacepage

import _ "embed"

// Row 201 A4 sprint 4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision card,
// served at GET /surface/decide/{ref} and GET /surface/decide.js. The page reads the subject from GET
// /surface/decisions/{ref}, checks both objects' hashes and the packet's requestRef, shows every
// field of both, builds the 91-byte decision message itself and POSTs a passkey assertion or a pasted
// SSH signature to POST /surface/decide/{ref}; the daemon runs the law again. The bytes are the same
// for every ref.

// DecidePage is GET /surface/decide/{ref}.
//
//go:embed decide/index.html
var DecidePage []byte

// DecideScript is GET /surface/decide.js.
//
//go:embed decide/decide.js
var DecideScript []byte
