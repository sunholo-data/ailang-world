package daemon

import (
	"net/http"
	"strings"

	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// handleSurfaceFont serves only exact embedded names, including OFL licences.
func (d *Daemon) handleSurfaceFont(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "" || r.URL.ForceQuery {
		http.Error(w, "surface fonts accept no query", http.StatusBadRequest)
		return
	}
	name := strings.TrimPrefix(r.URL.Path, "/surface/fonts/")
	b, contentType, ok := surfacepage.Font(name)
	if !ok || strings.Contains(strings.ToLower(r.URL.EscapedPath()), "%2f") {
		http.Error(w, "surface font not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	_, _ = w.Write(b)
}

// surfaceFontPaths prevents ServeMux's path-cleaning redirects for font GETs.
// Other methods still reach the mux, which supplies Method Not Allowed.
func (d *Daemon) surfaceFontPaths(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == http.MethodGet || r.Method == http.MethodHead) &&
			(r.URL.Path == "/surface/fonts" || strings.HasPrefix(r.URL.Path, "/surface/fonts/")) {
			d.handleSurfaceFont(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
