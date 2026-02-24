package cardanoasset_test

import (
	"strings"
	"testing"

	cardanoasset "github.com/njchilds90/go-cardano-asset"
)

// Known test vector: SpaceBudz policy and asset name
// Policy ID from the real SpaceBudz collection
const (
	testPolicyID     = "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc"
	testAssetName    = "SpaceBud0"
	testAssetNameHex = "537061636542756430"
)

func TestValidatePolicyID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", testPolicyID, false},
		{"all zeros", strings.Repeat("0", 56), false},
		{"valid all f", strings.Repeat("f", 56), false},
		{"too short", "d5e6bf", true},
		{"too long", testPolicyID + "00", true},
		{"uppercase", strings.ToUpper(testPolicyID), true},
		{"contains g", "g" + testPolicyID[1:], true},
		{"empty", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cardanoasset.ValidatePolicyID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePolicyID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateAssetNameHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", testAssetNameHex, false},
		{"empty", "", false},
		{"max length hex", strings.Repeat("61", 32), false},
		{"too long", strings.Repeat("61", 33), true},
		{"invalid hex", "zz", true},
		{"odd length", "abc", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cardanoasset.ValidateAssetNameHex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAssetNameHex(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestNewAsset(t *testing.T) {
	t.Run("valid asset", func(t *testing.T) {
		a, err := cardanoasset.NewAsset(testPolicyID, testAssetName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.PolicyID != testPolicyID {
			t.Errorf("PolicyID mismatch: got %s", a.PolicyID)
		}
		if a.AssetName != testAssetName {
			t.Errorf("AssetName mismatch: got %s", a.AssetName)
		}
	})

	t.Run("empty asset name allowed", func(t *testing.T) {
		_, err := cardanoasset.NewAsset(testPolicyID, "")
		if err != nil {
			t.Errorf("empty asset name should be valid, got: %v", err)
		}
	})

	t.Run("asset name too long", func(t *testing.T) {
		_, err := cardanoasset.NewAsset(testPolicyID, strings.Repeat("a", 33))
		if err == nil {
			t.Error("expected error for too-long asset name")
		}
	})

	t.Run("invalid policy ID", func(t *testing.T) {
		_, err := cardanoasset.NewAsset("badpolicy", testAssetName)
		if err == nil {
			t.Error("expected error for invalid policy ID")
		}
	})
}

func TestNewAssetFromHex(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		a, err := cardanoasset.NewAssetFromHex(testPolicyID, testAssetNameHex)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if a.AssetName != testAssetName {
			t.Errorf("AssetName mismatch: got %q, want %q", a.AssetName, testAssetName)
		}
	})

	t.Run("empty hex name", func(t *testing.T) {
		_, err := cardanoasset.NewAssetFromHex(testPolicyID, "")
		if err != nil {
			t.Errorf("empty hex should be valid: %v", err)
		}
	})

	t.Run("invalid hex", func(t *testing.T) {
		_, err := cardanoasset.NewAssetFromHex(testPolicyID, "zz")
		if err == nil {
			t.Error("expected error for invalid hex")
		}
	})
}

func TestParseAssetID(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantPolicy  string
		wantName    string
		wantErr     bool
	}{
		{
			name:       "full asset ID",
			input:      testPolicyID + "." + testAssetNameHex,
			wantPolicy: testPolicyID,
			wantName:   testAssetName,
		},
		{
			name:       "policy only",
			input:      testPolicyID,
			wantPolicy: testPolicyID,
			wantName:   "",
		},
		{
			name:    "empty",
			input:   "",
			wantErr: true,
		},
		{
			name:    "bad policy",
			input:   "bad." + testAssetNameHex,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := cardanoasset.ParseAssetID(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseAssetID(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if a.PolicyID != tt.wantPolicy {
					t.Errorf("PolicyID: got %q, want %q", a.PolicyID, tt.wantPolicy)
				}
				if a.AssetName != tt.wantName {
					t.Errorf("AssetName: got %q, want %q", a.AssetName, tt.wantName)
				}
			}
		})
	}
}

func TestAssetNameHex(t *testing.T) {
	a, _ := cardanoasset.NewAsset(testPolicyID, testAssetName)
	got := a.AssetNameHex()
	if got != testAssetNameHex {
		t.Errorf("AssetNameHex() = %q, want %q", got, testAssetNameHex)
	}
}

func TestAssetID(t *testing.T) {
	t.Run("with name", func(t *testing.T) {
		a, _ := cardanoasset.NewAsset(testPolicyID, testAssetName)
		want := testPolicyID + "." + testAssetNameHex
		if a.AssetID() != want {
			t.Errorf("AssetID() = %q, want %q", a.AssetID(), want)
		}
	})
	t.Run("without name", func(t *testing.T) {
		a, _ := cardanoasset.NewAsset(testPolicyID, "")
		if a.AssetID() != testPolicyID {
			t.Errorf("AssetID() = %q, want %q", a.AssetID(), testPolicyID)
		}
	})
}

func TestFingerprint(t *testing.T) {
	t.Run("returns asset1 prefix", func(t *testing.T) {
		fp, err := cardanoasset.Fingerprint(testPolicyID, testAssetName)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(fp, "asset1") {
			t.Errorf("Fingerprint should start with 'asset1', got: %s", fp)
		}
	})

	t.Run("deterministic", func(t *testing.T) {
		fp1, _ := cardanoasset.Fingerprint(testPolicyID, testAssetName)
		fp2, _ := cardanoasset.Fingerprint(testPolicyID, testAssetName)
		if fp1 != fp2 {
			t.Errorf("Fingerprint not deterministic: %s != %s", fp1, fp2)
		}
	})

	t.Run("different names produce different fingerprints", func(t *testing.T) {
		fp1, _ := cardanoasset.Fingerprint(testPolicyID, "SpaceBud0")
		fp2, _ := cardanoasset.Fingerprint(testPolicyID, "SpaceBud1")
		if fp1 == fp2 {
			t.Error("different asset names should produce different fingerprints")
		}
	})

	t.Run("different policies produce different fingerprints", func(t *testing.T) {
		policy2 := strings.Repeat("a", 56)
		fp1, _ := cardanoasset.Fingerprint(testPolicyID, testAssetName)
		fp2, _ := cardanoasset.Fingerprint(policy2, testAssetName)
		if fp1 == fp2 {
			t.Error("different policies should produce different fingerprints")
		}
	})

	t.Run("empty asset name", func(t *testing.T) {
		fp, err := cardanoasset.Fingerprint(testPolicyID, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !strings.HasPrefix(fp, "asset1") {
			t.Errorf("expected asset1 prefix, got: %s", fp)
		}
	})

	t.Run("invalid policy ID", func(t *testing.T) {
		_, err := cardanoasset.Fingerprint("bad", "name")
		if err == nil {
			t.Error("expected error for invalid policy ID")
		}
	})
}

func TestAssetInfo(t *testing.T) {
	a, _ := cardanoasset.NewAsset(testPolicyID, testAssetName)
	info, err := a.Info()
	if err != nil {
		t.Fatalf("Info() error: %v", err)
	}
	if info.PolicyID != testPolicyID {
		t.Errorf("PolicyID mismatch")
	}
	if info.AssetNameHex != testAssetNameHex {
		t.Errorf("AssetNameHex mismatch: %q", info.AssetNameHex)
	}
	if !strings.HasPrefix(info.Fingerprint, "asset1") {
		t.Errorf("Fingerprint prefix wrong: %s", info.Fingerprint)
	}
	if info.AssetID != testPolicyID+"."+testAssetNameHex {
		t.Errorf("AssetID mismatch: %s", info.AssetID)
	}
}

func TestIsValidUTF8Name(t *testing.T) {
	tests := []struct {
		name  string
		asset string
		want  bool
	}{
		{"valid utf8", "SpaceBud0", true},
		{"empty", "", true},
		{"unicode", "🚀NFT", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := cardanoasset.NewAsset(testPolicyID, tt.asset)
			if a.IsValidUTF8Name() != tt.want {
				t.Errorf("IsValidUTF8Name() = %v, want %v", a.IsValidUTF8Name(), tt.want)
			}
		})
	}
}

func BenchmarkFingerprint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = cardanoasset.Fingerprint(testPolicyID, testAssetName)
	}
}
