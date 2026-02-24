# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-24

### Added
- `Asset` struct with `PolicyID` and `AssetName` fields
- `AssetInfo` struct with full asset details
- `NewAsset` — create asset from raw name string
- `NewAssetFromHex` — create asset from hex-encoded name
- `ParseAssetID` — parse `policyId.assetNameHex` format
- `Asset.AssetNameHex()` — hex-encode asset name
- `Asset.AssetID()` — construct full asset ID string
- `Asset.Fingerprint()` — CIP-14 bech32 fingerprint
- `Asset.Info()` — all fields in one call
- `Asset.IsValidUTF8Name()` — validate UTF-8 encoding
- `Fingerprint()` — standalone fingerprint function
- `ValidatePolicyID()` — validate 56-char hex policy ID
- `ValidateAssetNameHex()` — validate hex asset name
- Structured sentinel error variables
- Zero external dependencies
- Full table-driven test suite with benchmarks
- GitHub Actions CI across Go 1.21, 1.22, 1.23
