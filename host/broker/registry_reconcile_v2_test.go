package broker

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

const (
	wcV1 = "sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083"
	wcV2 = "sha256:ifacev2:b25fe03155db0c7bf595cf730295b945d6ac64ec415a1998fae8a693d621e8d8"
)

func served01(t *testing.T) []byte {
	b, err := os.ReadFile("testdata/metadata_world_core_0.1.0.json")
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func cfg01() ReconcileConfig {
	return ReconcileConfig{Vendor: "world", Name: "core", Version: "0.1.0",
		Expected: PublishHashes{
			TarballSHA256: "sha256:44fc9fab7be710f09b84274f744445db41a79df77c71cc43c0e952d75d97c27f",
			ContentHash:   "sha256:0c8c60616e592dc01891e8bbb59350786f242a2f79a9eb2c587ae8b0ca2e00b9",
			InterfaceHash: wcV1},
		ExpectedInterfaceV2: wcV2}
}

func TestReconcilePublished010OverFourDigests(t *testing.T) {
	r, _ := resolvePresent(ReconcileReceipt{}, served01(t), cfg01())
	if r.State != ReconcileSucceededReconciled || !strings.Contains(r.Detail, "four") {
		t.Fatalf("%s %s", r.State, r.Detail)
	}
}

func TestReconcileV2MismatchIsConflict(t *testing.T) {
	c := cfg01()
	c.ExpectedInterfaceV2 = "sha256:ifacev2:0a46f29a089fe1add745412d3e95d5b0093c9bd7e05177b6f3aa8daabf14c30e"
	r, _ := resolvePresent(ReconcileReceipt{}, served01(t), c)
	if r.State != ReconcileConflict || !strings.Contains(r.Detail, "interface-v2") {
		t.Fatalf("%s %s", r.State, r.Detail)
	}
}

func TestReconcileAbsentV2IsNotSuccess(t *testing.T) {
	var doc map[string]any
	if err := json.Unmarshal(served01(t), &doc); err != nil {
		t.Fatal(err)
	}
	// Real-world shape: sunholo/auth@0.4.1 (schema v1) omits interface_hash_v2.
	ctl, err := os.ReadFile("testdata/metadata_sunholo_auth_0.4.1.json")
	if err != nil {
		t.Fatal(err)
	}
	var auth map[string]any
	if err := json.Unmarshal(ctl, &auth); err != nil || len(auth) == 0 {
		t.Fatalf("instrument: auth fixture unreadable: %v", err)
	}
	if v, ok := auth["interface_hash_v2"]; ok {
		t.Fatalf("instrument: auth fixture now serves v2: %v", v)
	}
	for _, mode := range []string{"null", "absent"} {
		d := map[string]any{}
		for k, v := range doc {
			d[k] = v
		}
		if mode == "null" {
			d["interface_hash_v2"] = nil
		} else {
			delete(d, "interface_hash_v2")
		}
		body, _ := json.Marshal(d)
		r, _ := resolvePresent(ReconcileReceipt{}, body, cfg01())
		if r.State != ReconcileConflict || !strings.Contains(r.Detail, "no interface_hash_v2") {
			t.Fatalf("%s: %s %s", mode, r.State, r.Detail)
		}
	}
}

func TestReconcileRefusesEmptyExpectedV2(t *testing.T) {
	c := cfg01()
	c.RegistryOrigin = "http://127.0.0.1:1"
	c.ExpectedInterfaceV2 = ""
	_, err := reconcileLoopback(t.Context(), c)
	if err == nil || !strings.Contains(err.Error(), "four expected digests") {
		t.Fatalf("want R6 refusal, got %v", err)
	}
}
