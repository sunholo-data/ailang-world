package surfacefollow

import (
	"os"
	"testing"
)

// FileInfo.Sys has the same Uid field on the supported Unix platforms. These
// pure cases cover the foreign-owner refusal when chown is unavailable.
type ownerInfo struct {
	os.FileInfo
	uid uint32
}

func (i ownerInfo) Sys() any { return &struct{ Uid uint32 }{i.uid} }

func TestConfigOwnerMetadata(t *testing.T) {
	if !owned(ownerInfo{uid: uint32(os.Getuid())}) {
		t.Fatal("current owner refused")
	}
	if owned(ownerInfo{uid: uint32(os.Getuid() + 1)}) {
		t.Fatal("foreign owner accepted")
	}
}
