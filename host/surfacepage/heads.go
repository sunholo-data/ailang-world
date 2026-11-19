package surfacepage

import (
	"bytes"
	_ "embed"
)

//go:embed app/heads.html
var headsIndex []byte

//go:embed app/heads.css
var headsStyle []byte

//go:embed app/heads-data.js
var headsData []byte

//go:embed app/heads-activity.js
var headsActivity []byte

//go:embed app/heads-links.js
var headsLinks []byte

//go:embed app/heads-links.css
var headsLinksStyle []byte

//go:embed app/heads.js
var headsView []byte

//go:embed app/heads-page.js
var headsPage []byte

// HeadsShell shares the timeline's fonts, appearance tokens and embedded World mark.
var HeadsShell = inlineStyle(bytes.Replace(headsIndex, []byte("<!-- World mark -->"), AppLogo, 1),
	bytes.Join([][]byte{AppStyle, headsStyle, headsLinksStyle}, []byte("\n")))

// HeadsScript renders the same indexed states the timeline and brane stack expose.
var HeadsScript = bytes.Join([][]byte{braneScript, projectionMotion, headsData, headsActivity, headsLinks, headsView, headsPage}, []byte("\n"))
