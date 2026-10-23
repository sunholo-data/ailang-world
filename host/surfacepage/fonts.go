package surfacepage

import (
	"embed"
	"strings"
)

// fontFiles contains the pinned OFL fonts and their accompanying licences.
//
//go:embed fonts/*.ttf fonts/*.txt
var fontFiles embed.FS

// Font returns an embedded font or licence by exact file name.
func Font(name string) ([]byte, string, bool) {
	if name == "" || strings.ContainsAny(name, "/\\") {
		return nil, "", false
	}
	contentType := ""
	switch {
	case strings.HasSuffix(name, ".ttf"):
		contentType = "font/ttf"
	case strings.HasSuffix(name, ".txt"):
		contentType = "text/plain; charset=utf-8"
	default:
		return nil, "", false
	}
	b, err := fontFiles.ReadFile("fonts/" + name)
	if err != nil {
		return nil, "", false
	}
	return b, contentType, true
}
