package surfacepage

import _ "embed"

// Row 201 A3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.7): the request card, served
// at GET /surface/approve/{id} and GET /surface/approve.js. The page reads the request from GET
// /surface/requests/{id}, shows every signed field of its stored bytes after checking their digest,
// computes the message it signs itself and POSTs a passkey assertion or an SSH signature to POST
// /surface/approve/{id}; the daemon judges it. The bytes are the same for every id.

// ApprovePage is GET /surface/approve/{id}.
//
//go:embed approve/index.html
var ApprovePage []byte

// ApproveScript is GET /surface/approve.js.
//
//go:embed approve/approve.js
var ApproveScript []byte
