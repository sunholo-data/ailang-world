package surfacepage

import _ "embed"

// Row 243 V3 sprint 2 M5 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the confirmation
// card, to be served at GET /surface/confirm/{id} and GET /surface/confirm.js (sprint 3 registers the
// routes; until then only tests serve these bytes). The page reads the card's data from GET
// /surface/confirmations/{id}, checks the request's digest, shows every field of the signed bytes
// through one table (refusing a field it cannot show), and builds the 86-byte message itself: Save
// the message and the ssh-keygen line. Sprint 4 adds the submit half. The bytes are the same for every
// id.

// ConfirmPage is GET /surface/confirm/{id}.
//
//go:embed confirm/index.html
var ConfirmPage []byte

// ConfirmScript is GET /surface/confirm.js.
//
//go:embed confirm/confirm.js
var ConfirmScript []byte
