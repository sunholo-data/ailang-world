package surfacepage

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSurfaceFontEmbed(t *testing.T) {
	want := []string{"IBMPlexMono-Medium.ttf", "IBMPlexMono-Regular.ttf", "Montserrat-Variable.ttf", "OFL-IBMPlexMono.txt", "OFL-Montserrat.txt"}
	entries, err := fontFiles.ReadDir("fonts")
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	if !reflect.DeepEqual(names, want) {
		t.Fatalf("embedded names = %v, want %v", names, want)
	}
	for _, name := range want {
		b, ct, ok := Font(name)
		disk, err := os.ReadFile("fonts/" + name)
		if err != nil {
			t.Fatal(err)
		}
		wantType := "font/ttf"
		if strings.HasSuffix(name, ".txt") {
			wantType = "text/plain; charset=utf-8"
		}
		if !ok || !bytes.Equal(b, disk) || ct != wantType {
			t.Errorf("lookup %s differs from pinned file or content type", name)
		}
	}
	for _, name := range []string{"", "..", "../fonts/Montserrat-Variable.ttf", "fonts/Montserrat-Variable.ttf", "montserrat-variable.ttf", "Montserrat-Variable.ttf/extra", "missing.ttf", "OFL-Montserrat.TXT", "fonts.go", "Montserrat-Variable.ttf?x=1", "foo\\bar.ttf"} {
		if b, ct, ok := Font(name); ok || b != nil || ct != "" {
			t.Errorf("accepted %q", name)
		}
	}
}

func TestSurfaceFontFaces(t *testing.T) {
	timeline, err := RenderTimeline(Timeline{}, Now{})
	if err != nil {
		t.Fatal(err)
	}
	stack, err := RenderStack(Stack{})
	if err != nil {
		t.Fatal(err)
	}
	for _, body := range [][]byte{timeline, stack} {
		s := string(body)
		if strings.Count(s, "@font-face") != 3 || strings.Count(s, "font-display: swap") != 3 {
			t.Error("page must define three swapping font faces")
		}
		for _, rule := range []string{
			`src: url(/surface/fonts/Montserrat-Variable.ttf) format("truetype"); font-weight: 100 900; font-style: normal;`,
			`src: url(/surface/fonts/IBMPlexMono-Regular.ttf) format("truetype"); font-weight: 400; font-style: normal;`,
			`src: url(/surface/fonts/IBMPlexMono-Medium.ttf) format("truetype"); font-weight: 500; font-style: normal;`,
			`font-family:"IBM Plex Mono",monospace;`,
		} {
			if !strings.Contains(s, rule) {
				t.Errorf("page lacks %s", rule)
			}
		}
	}
}
