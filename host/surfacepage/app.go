package surfacepage

import (
	"bytes"
	_ "embed"
	"html"
)

// Row 188 N2b: the surface app, the approved v3 working copy
// (design_docs/mockups/surface-app-v3/dist) embedded and served at /surface.
// It renders GET /surface/model.json and GET /surface/stack.json; it carries
// no data of its own.

// AppScript is /surface/app.js: the app, with its pure part exported for node.
//
//go:embed app/app.js
var timelineScript []byte

//go:embed app/head-drafts.js
var headDrafts []byte

//go:embed app/head-camera.js
var headCamera []byte

//go:embed app/head-travel-input.js
var headTravelInput []byte

//go:embed app/head-ground-input.js
var headGroundInput []byte

//go:embed app/head-lanes.js
var headLanes []byte

//go:embed app/head-node-preview.js
var headNodePreview []byte

//go:embed app/head-lineage.js
var headLineage []byte

//go:embed app/head-scene.js
var headScene []byte

//go:embed app/camera-controller.js
var cameraController []byte

//go:embed app/composer.js
var composerScript []byte

var AppScript = bytes.Join([][]byte{headDrafts, headCamera, headTravelInput, headGroundInput, headLanes, headNodePreview, headLineage, headScene, cameraController, timelineScript, composerScript}, []byte("\n"))

// AppStyle is /surface/app.css: the working copy's stylesheet, its fonts
// self-hosted under /surface/fonts/.
//
//go:embed app/app.css
var timelineStyle []byte

//go:embed app/head-scene.css
var headStyle []byte

//go:embed app/head-inspect.css
var headInspectStyle []byte

//go:embed app/head-drafts.css
var headDraftStyle []byte

//go:embed app/head-travel.css
var headTravelStyle []byte

//go:embed app/head-lineage.css
var headLineageStyle []byte

//go:embed app/head-lanes.css
var headLaneStyle []byte

//go:embed app/head-status.css
var headStatusStyle []byte

//go:embed app/time-ribbon.css
var timeRibbonStyle []byte

var AppStyle = bytes.Join([][]byte{timelineStyle, headStyle, headInspectStyle, headDraftStyle, headTravelStyle, headLineageStyle, headLaneStyle, headStatusStyle, timeRibbonStyle}, []byte("\n"))

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

// ShellWithDraftCatalog embeds only the presets explicitly exposed by the operator.
func ShellWithDraftCatalog(catalog []byte) []byte {
	if len(catalog) == 0 {
		return AppShell
	}
	return bytes.Replace(AppShell, []byte(`id="episode-draft" data-catalog=""`),
		[]byte(`id="episode-draft" data-catalog="`+html.EscapeString(string(catalog))+`"`), 1)
}
