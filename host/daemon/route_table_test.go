package daemon

import (
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"testing"
)

// TestFrozenV1RouteTableMatchesMux pins every literal frozen method/path pair.
func TestFrozenV1RouteTableMatchesMux(t *testing.T) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	daemonSource, err := os.ReadFile(filepath.Join(root, "host/daemon/daemon.go"))
	if err != nil {
		t.Fatal(err)
	}
	sketchSource, err := os.ReadFile(filepath.Join(root, "design_docs/sketches/worlddapi.ail"))
	if err != nil {
		t.Fatal(err)
	}
	handler := strings.SplitN(string(daemonSource), "func (d *Daemon) Handler() http.Handler {", 2)
	if len(handler) != 2 {
		t.Fatal("Handler source missing")
	}
	handlerBody := strings.SplitN(handler[1], "\n}\n", 2)[0]
	muxRE := regexp.MustCompile(`mux\.HandleFunc\("(GET|POST) (/v1/[^\"]+)"`)
	var mux []string
	for _, match := range muxRE.FindAllStringSubmatch(handlerBody, -1) {
		mux = append(mux, match[1]+" "+match[2])
	}
	routes := strings.SplitN(string(sketchSource), "export func routes()", 2)
	if len(routes) != 2 {
		t.Fatal("routes() missing")
	}
	routesBody := strings.SplitN(routes[1], "\n}\n", 2)[0]
	sketchRE := regexp.MustCompile(`\{ method: "(GET|POST)", path: "(/v1/[^\"]+)" \}`)
	var sketch []string
	for _, match := range sketchRE.FindAllStringSubmatch(routesBody, -1) {
		sketch = append(sketch, match[1]+" "+match[2])
	}
	sort.Strings(mux)
	sort.Strings(sketch)
	if len(mux) != 10 || len(sketch) != 10 || !reflect.DeepEqual(mux, sketch) {
		t.Fatalf("frozen routes: mux=%v sketch=%v; want identical ten pairs", mux, sketch)
	}
}
