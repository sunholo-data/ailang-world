package surfacepage

import _ "embed"

// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.14): the passkey ceremony
// page, served at GET /surface/passkey and GET /surface/passkey.js. It reads its link's fragment,
// runs one WebAuthn create or get, and offers the result as a file; it makes no request and holds
// no data.

// PasskeyPage is GET /surface/passkey.
//
//go:embed passkey/index.html
var PasskeyPage []byte

// PasskeyScript is GET /surface/passkey.js.
//
//go:embed passkey/passkey.js
var PasskeyScript []byte
