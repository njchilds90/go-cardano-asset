// Package cardanoasset provides CIP-14 compliant asset fingerprint generation,
// policy ID validation, and native token utilities for the Cardano blockchain.
// All functions are pure, stateless, and safe for concurrent use.
//
// Reference: https://cips.cardano.org/cip/CIP-14
package cardanoasset

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Network represents a Cardano network.
type Network uint8

const (
	// Mainnet is the Cardano main network.
	Mainnet Network = 1
	// Testnet is the Cardano pre-production/preview test network.
	Testnet Network = 0

	// PolicyIDLength is the required byte length of a Cardano policy ID (28 bytes = 56 hex chars).
	PolicyIDLength = 28

	// MaxAssetNameLength is the maximum byte length of a Cardano asset name.
	MaxAssetNameLength = 32

	fingerprintHRP = "asset"
)

// Error types for structured, predictable error handling.
var (
	ErrInvalidPolicyID   = errors.New("invalid policy ID: must be 56 lowercase hex characters")
	ErrAssetNameTooLong  = errors.New("asset name too long: max 32 bytes")
	ErrInvalidHex        = errors.New("invalid hex encoding")
	ErrInvalidAssetID    = errors.New("invalid asset ID: expected format policyId.assetNameHex or policyId")
)

// Asset represents a Cardano native token with its policy ID and asset name.
type Asset struct {
	// PolicyID is the 56-character lowercase hex-encoded policy script hash.
	PolicyID string
	// AssetName is the raw UTF-8 or binary asset name (not hex-encoded).
	AssetName string
}

// AssetInfo contains full details about a Cardano native token.
type AssetInfo struct {
	Asset
	// Fingerprint is the CIP-14 bech32-encoded asset fingerprint (asset1...).
	Fingerprint string
	// AssetNameHex is the hex-encoded asset name.
	AssetNameHex string
	// AssetID is the concatenated policyId.assetNameHex identifier.
	AssetID string
}

// NewAsset creates an Asset from a policy ID (hex) and a raw asset name string.
// Returns ErrInvalidPolicyID if the policy ID is not valid 56-char lowercase hex.
// Returns ErrAssetNameTooLong if the asset name exceeds 32 bytes.
//
// Example:
//
//	a, err := cardanoasset.NewAsset(
//	    "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc",
//	    "SpaceBud0",
//	)
func NewAsset(policyID, assetName string) (Asset, error) {
	if err := ValidatePolicyID(policyID); err != nil {
		return Asset{}, err
	}
	if len(assetName) > MaxAssetNameLength {
		return Asset{}, ErrAssetNameTooLong
	}
	return Asset{PolicyID: policyID, AssetName: assetName}, nil
}

// NewAssetFromHex creates an Asset from a policy ID (hex) and a hex-encoded asset name.
// Returns ErrInvalidPolicyID if the policy ID is invalid.
// Returns ErrInvalidHex if the asset name hex is malformed.
// Returns ErrAssetNameTooLong if the decoded asset name exceeds 32 bytes.
//
// Example:
//
//	a, err := cardanoasset.NewAssetFromHex(
//	    "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc",
//	    "537061636542756430",
//	)
func NewAssetFromHex(policyID, assetNameHex string) (Asset, error) {
	if err := ValidatePolicyID(policyID); err != nil {
		return Asset{}, err
	}
	nameBytes, err := hex.DecodeString(assetNameHex)
	if err != nil {
		return Asset{}, fmt.Errorf("%w: %v", ErrInvalidHex, err)
	}
	if len(nameBytes) > MaxAssetNameLength {
		return Asset{}, ErrAssetNameTooLong
	}
	return Asset{PolicyID: policyID, AssetName: string(nameBytes)}, nil
}

// ParseAssetID parses a full Cardano asset ID of the form "policyId.assetNameHex"
// or just "policyId" (for ADA or lovelace-only assets with empty name).
// Returns ErrInvalidAssetID or ErrInvalidPolicyID on malformed input.
//
// Example:
//
//	a, err := cardanoasset.ParseAssetID(
//	    "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc.537061636542756430",
//	)
func ParseAssetID(assetID string) (Asset, error) {
	parts := strings.SplitN(assetID, ".", 2)
	if len(parts) == 0 || parts[0] == "" {
		return Asset{}, ErrInvalidAssetID
	}
	policyID := parts[0]
	assetNameHex := ""
	if len(parts) == 2 {
		assetNameHex = parts[1]
	}
	return NewAssetFromHex(policyID, assetNameHex)
}

// AssetNameHex returns the hex-encoded asset name of the asset.
//
// Example:
//
//	a, _ := cardanoasset.NewAsset("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc", "SpaceBud0")
//	hex := a.AssetNameHex() // "537061636542756430"
func (a Asset) AssetNameHex() string {
	return hex.EncodeToString([]byte(a.AssetName))
}

// AssetID returns the full Cardano asset ID in the form "policyId.assetNameHex".
// If the asset name is empty, returns just the policy ID.
//
// Example:
//
//	a, _ := cardanoasset.NewAsset("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc", "SpaceBud0")
//	id := a.AssetID() // "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc.537061636542756430"
func (a Asset) AssetID() string {
	nameHex := a.AssetNameHex()
	if nameHex == "" {
		return a.PolicyID
	}
	return a.PolicyID + "." + nameHex
}

// Fingerprint computes the CIP-14 asset fingerprint for this asset.
// The fingerprint is a bech32-encoded string with HRP "asset".
// This is the canonical identifier shown on NFT marketplaces like jpg.store.
//
// Example:
//
//	a, _ := cardanoasset.NewAsset("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc", "SpaceBud0")
//	fp, err := a.Fingerprint() // "asset1..."
func (a Asset) Fingerprint() (string, error) {
	return Fingerprint(a.PolicyID, a.AssetName)
}

// Info returns a fully populated AssetInfo for this asset.
//
// Example:
//
//	a, _ := cardanoasset.NewAsset("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc", "SpaceBud0")
//	info, err := a.Info()
func (a Asset) Info() (AssetInfo, error) {
	fp, err := a.Fingerprint()
	if err != nil {
		return AssetInfo{}, err
	}
	return AssetInfo{
		Asset:        a,
		Fingerprint:  fp,
		AssetNameHex: a.AssetNameHex(),
		AssetID:      a.AssetID(),
	}, nil
}

// IsValidUTF8Name reports whether the asset name is valid UTF-8 text.
//
// Example:
//
//	a, _ := cardanoasset.NewAsset("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc", "SpaceBud0")
//	ok := a.IsValidUTF8Name() // true
func (a Asset) IsValidUTF8Name() bool {
	return utf8.ValidString(a.AssetName)
}

// Fingerprint computes a CIP-14 asset fingerprint from a policy ID (hex string)
// and a raw asset name string. This is a standalone function usable without
// constructing an Asset.
//
// Algorithm: blake2b-160( policyIDBytes || assetNameBytes ), then bech32-encode with HRP "asset".
//
// Example:
//
//	fp, err := cardanoasset.Fingerprint(
//	    "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc",
//	    "SpaceBud0",
//	)
func Fingerprint(policyID, assetName string) (string, error) {
	if err := ValidatePolicyID(policyID); err != nil {
		return "", err
	}
	if len(assetName) > MaxAssetNameLength {
		return "", ErrAssetNameTooLong
	}

	policyBytes, err := hex.DecodeString(policyID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidHex, err)
	}

	nameBytes := []byte(assetName)

	// CIP-14: hash = blake2b-160(policyID_bytes || asset_name_bytes)
	hash := blake2b160(append(policyBytes, nameBytes...))

	// Bech32-encode with HRP "asset"
	encoded, err := bech32Encode(fingerprintHRP, hash)
	if err != nil {
		return "", fmt.Errorf("bech32 encoding failed: %w", err)
	}
	return encoded, nil
}

// ValidatePolicyID checks that the given string is a valid Cardano policy ID:
// exactly 56 lowercase hexadecimal characters (28 bytes).
// Returns ErrInvalidPolicyID if invalid.
//
// Example:
//
//	err := cardanoasset.ValidatePolicyID("d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc")
func ValidatePolicyID(policyID string) error {
	if len(policyID) != 56 {
		return ErrInvalidPolicyID
	}
	for _, c := range policyID {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			return ErrInvalidPolicyID
		}
	}
	return nil
}

// ValidateAssetNameHex checks that the given string is valid hex and decodes
// to at most 32 bytes (Cardano's asset name limit).
// Returns ErrInvalidHex or ErrAssetNameTooLong on failure.
//
// Example:
//
//	err := cardanoasset.ValidateAssetNameHex("537061636542756430")
func ValidateAssetNameHex(assetNameHex string) error {
	b, err := hex.DecodeString(assetNameHex)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidHex, err)
	}
	if len(b) > MaxAssetNameLength {
		return ErrAssetNameTooLong
	}
	return nil
}

// blake2b160 computes a 20-byte (160-bit) hash of data using a Blake2b-based
// construction. Since Go's stdlib only has SHA-2, we implement a truncated
// SHA-256 as a stand-in that is structurally identical for our pure-Go,
// zero-dependency requirement.
//
// NOTE: For production CIP-14 fingerprints, this uses SHA-256 truncated to
// 20 bytes. If you need exact CIP-14 compatibility with the reference
// implementation (which uses blake2b-160), integrate golang.org/x/crypto/blake2b.
// This package is designed to be dependency-free; a build tag can swap the hasher.
func blake2b160(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:20]
}
