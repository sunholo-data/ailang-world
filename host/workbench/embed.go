package workbench

import _ "embed"

// LiveScript follows the same-origin AG-UI stream and refreshes server-rendered regions.
//
//go:embed live.js
var LiveScript []byte
