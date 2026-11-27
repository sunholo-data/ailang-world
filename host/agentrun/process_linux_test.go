package agentrun

import (
	"errors"
	"os"
	"strings"
	"testing"
	"testing/fstest"
)

func TestLinuxProcessIdentity(t *testing.T) {
	stat := func(state, ticks string) string {
		return "123 (a task ) with right parentheses) " + state + strings.Repeat(" 0", 18) + " " + ticks
	}
	for _, tc := range []struct {
		name, stat, boot, want, broken string
		pid                            int
		live                           bool
	}{
		{"live", stat("S", "17"), "boot\n", "boot:17", "", 123, true},
		{"zombie", stat("Z", "17"), "boot", "boot:17", "", 123, false},
		{"partial stat read", stat("S", "17"), "boot", "", "123/stat", 123, false},
		{"partial boot read", stat("S", "17"), "boot", "", "sys/kernel/random/boot_id", 123, false},
		{"missing stat", "", "boot", "", "", 123, false},
		{"missing boot", stat("S", "17"), "", "", "", 123, false},
		{"blank boot", stat("S", "17"), " ", "", "", 123, false},
		{"short stat", "123 (short) S", "boot", "", "", 123, false},
		{"no delimiter", "S" + strings.Repeat(" 0", 18) + " 17", "boot", "", "", 123, false},
		{"invalid ticks", stat("S", "bad"), "boot", "", "", 123, false},
		{"negative pid", stat("S", "17"), "boot", "", "", -1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			proc := fstest.MapFS{}
			if tc.stat != "" {
				proc["123/stat"] = &fstest.MapFile{Data: []byte(tc.stat)}
				proc["-1/stat"] = proc["123/stat"]
			}
			if tc.boot != "" {
				proc["sys/kernel/random/boot_id"] = &fstest.MapFile{Data: []byte(tc.boot)}
			}
			if got, live := linuxProcessIdentity(procFixture{proc, tc.broken}, tc.pid); got != tc.want || live != tc.live {
				t.Fatalf("identity = %q, %t, want %q, %t", got, live, tc.want, tc.live)
			}
		})
	}
	if birth, live := ProcessIdentity(os.Getpid()); !live || birth == "" {
		t.Fatalf("read actual procfs: %q %v", birth, live)
	}
}

type procFixture struct {
	fstest.MapFS
	broken string
}

func (f procFixture) ReadFile(name string) ([]byte, error) {
	b, err := f.MapFS.ReadFile(name)
	if name == f.broken {
		return b, errors.New("fixture: partial filesystem read")
	}
	return b, err
}
