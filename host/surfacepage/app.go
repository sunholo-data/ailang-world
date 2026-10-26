package surfacepage

import (
	"bytes"
	_ "embed"
)

// Row 188 N2b: the surface app, the approved v3 working copy
// (design_docs/mockups/surface-app-v3/dist) embedded and served at /surface.
// It renders GET /surface/model.json and GET /surface/stack.json; it carries
// no data of its own.

// AppScript is /surface/app.js: the app, with its pure part exported for node.
//
//go:embed app/app.js
var AppScript []byte

// AppStyle is /surface/app.css: the working copy's stylesheet, its fonts
// self-hosted under /surface/fonts/.
//
//go:embed app/app.css
var AppStyle []byte

// AppLogo is /surface/ailang-logo.svg, the AILANG mark the shell inlines.
//
//go:embed app/ailang-logo.svg
var AppLogo []byte

//go:embed app/index.html
var appIndex []byte

// AppStylesheetLink is index.html's stylesheet link, which AppShell replaces
// with the stylesheet itself: /surface's CSP allows inline styles and no
// stylesheet from 'self', so a <link> would be refused.
const AppStylesheetLink = `<link rel="stylesheet" href="/surface/app.css">`

// AppShell is GET /surface: index.html with app.css inlined in one <style>.
var AppShell = inlineStyle(appIndex, AppStyle)

func inlineStyle(index, css []byte) []byte {
	if bytes.Count(index, []byte(AppStylesheetLink)) != 1 || bytes.Contains(bytes.ToLower(css), []byte("</style")) {
		panic("surfacepage: app/index.html must link app.css exactly once, and app.css must not close a style element")
	}
	style := append(append([]byte("<style>\n"), css...), "</style>"...)
	return bytes.Replace(index, []byte(AppStylesheetLink), style, 1)
}
