package surfacepage

import _ "embed"

// OperatorScript and OperatorStyle are the standalone operator card assets.
//
//go:embed operator/card.js
var OperatorScript []byte

//go:embed operator/card.css
var OperatorStyle []byte
