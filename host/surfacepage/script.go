package surfacepage

import _ "embed"

// Script is /surface/surface.js (row 188 M4a), the surface's one script: the
// key map, the address bar and the Appearance panel. Both pages load it with
// one <script src defer> element and keep working without it.
//
//go:embed surface.js
var Script []byte

// ScriptElement is the one script element both pages carry.
const ScriptElement = `<script src="/surface/surface.js" defer></script>`

// chromeTemplate is what the two pages share for the script: the address
// bar, the gear (shown by the script), and the CSS for the panels and the
// Appearance attributes. Each page parses it beside its own template.
const chromeTemplate = `
{{- define "chrome-style"}}
@font-face { font-family: "Montserrat"; src: url(/surface/fonts/Montserrat-Variable.ttf) format("truetype"); font-weight: 100 900; font-style: normal; font-display: swap; }
@font-face { font-family: "IBM Plex Mono"; src: url(/surface/fonts/IBMPlexMono-Regular.ttf) format("truetype"); font-weight: 400; font-style: normal; font-display: swap; }
@font-face { font-family: "IBM Plex Mono"; src: url(/surface/fonts/IBMPlexMono-Medium.ttf) format("truetype"); font-weight: 500; font-style: normal; font-display: swap; }
.tick, .read-path, .file, .task-path, .source-chip, .front h1, .card span:nth-child(2), .rail a, code, time { font-family: "IBM Plex Mono", monospace; }

.address { flex:0 1 480px; min-width:0; margin-left:12px; padding:3px 8px; border:1px solid var(--edge); border-radius:3px; background:transparent; color:var(--headink); font:inherit; font-family:"IBM Plex Mono",monospace; font-size:13px; }
.address:focus { outline:2px solid var(--you); outline-offset:1px; }
.address-note { flex:0 1 auto; min-width:0; color:#F9A697; font-size:12px; }
.gear { margin-left:auto; padding:2px 8px; border:1px solid var(--edge); border-radius:3px; background:transparent; color:var(--headink); font:inherit; cursor:pointer; }
.gear[hidden] { display:none; }
.panel { position:fixed; top:52px; right:16px; z-index:10; width:320px; max-width:calc(100vw - 32px); max-height:calc(100vh - 68px); overflow:auto; padding:12px 16px; border:1px solid var(--edge); border-radius:3px; background:var(--bg); color:var(--ink); }
.panel h2 { margin:0 0 8px; padding:0; border:0; background:none; font-size:14px; font-weight:600; text-transform:none; color:var(--ink); }
.panel .close { position:absolute; top:6px; right:8px; border:0; background:none; color:var(--sub); font:inherit; font-size:18px; cursor:pointer; }
.panel fieldset { margin:0 0 10px; padding:6px 10px; border:1px solid var(--line); border-radius:3px; }
.panel legend { padding:0 4px; color:var(--sub); font-size:12px; }
.panel label { display:block; margin:2px 0; cursor:pointer; }
.panel dl { display:grid; grid-template-columns:auto 1fr; gap:4px 12px; margin:0; }
.panel dt { font-weight:600; white-space:nowrap; }
.panel dd { margin:0; color:var(--sub); }
.panel p { margin:4px 0; }
.panel .sunholo { margin-top:12px; color:var(--sub); font-size:12px; font-weight:600; }
:root[data-marks="shape"] { --mark:#526677; }
:root[data-marks="shape"][data-theme="slate"], :root[data-marks="shape"][data-theme="night"] { --mark:#8496A5; }
:root[data-marks="shape"] .lane polygon, :root[data-marks="shape"] .state polygon, :root[data-marks="shape"] .legend polygon, :root[data-marks="shape"] .glyph polygon { fill:var(--mark); }
:root[data-marks="shape"] .ring, :root[data-marks="shape"] .glyph circle, :root[data-marks="shape"] .link { stroke:var(--mark); }
:root[data-motion="reduced"] *, :root[data-motion="reduced"] *::before, :root[data-motion="reduced"] *::after { transition:none !important; animation:none !important; scroll-behavior:auto !important; }
{{- end}}
{{- define "chrome-head"}}<input class="address" type="text" value="{{.}}" aria-label="Address" spellcheck="false" autocomplete="off"><span class="address-note" role="status" aria-live="polite"></span><button class="gear" type="button" aria-label="Appearance" title="Appearance (,)" hidden>⚙</button>{{end}}
`
