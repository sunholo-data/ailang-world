package store

import "testing"

func TestOperatorLaunchProcessStatus(t *testing.T) {
	for _, tc := range []struct {
		name, phase, status, birth, want string
		live, noSpawn, reads             bool
	}{
		{"no-spawn", "observed", "", "", "spawn-refused", false, true, false},
		{"reserved", "reserved", "", "", "uncertain", false, false, false},
		{"live", "observed", "launched", "same", "launched", true, false, true},
		{"reused", "observed", "launched", "different", "uncertain", true, false, true},
		{"gone", "observed", "launched", "same", "uncertain", false, false, true},
		{"stopped", "observed", "stopped", "", "stopped", false, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshot := OperatorLaunchSnapshot{Phase: tc.phase, NoSpawn: tc.noSpawn,
				Receipt: OperatorLaunchReceipt{PID: 42, Birth: "same", Status: tc.status}}
			read := false
			got := snapshot.ProcessStatus(func(pid int) (string, bool) {
				if pid != 42 {
					t.Fatalf("looked up pid %d", pid)
				}
				read = true
				return tc.birth, tc.live
			})
			if got != tc.want || read != tc.reads {
				t.Fatalf("status %q, read=%v, want %q, read=%v", got, read, tc.want, tc.reads)
			}
		})
	}
}
