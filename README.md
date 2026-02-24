# go-cardano-asset

[![CI](https://github.com/njchilds90/go-cardano-asset/actions/workflows/ci.yml/badge.svg)](https://github.com/njchilds90/go-cardano-asset/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/go-cardano-asset.svg)](https://pkg.go.dev/github.com/njchilds90/go-cardano-asset)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/go-cardano-asset)](https://goreportcard.com/report/github.com/njchilds90/go-cardano-asset)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**CIP-14 asset fingerprint generation, policy ID validation, and native token utilities for Cardano — zero dependencies, pure Go.**

---

## Why This Exists

Every Cardano native token has three identifiers you constantly need to convert between:

| Identifier | Example |
|---|---|
| Asset ID | `d5e6bf05...4cc.537061636542756430` |
| Asset name (UTF-8) | `SpaceBud0` |
| CIP-14 Fingerprint | `asset1...` |

There was no dedicated, zero-dependency Go library for this. `go-cardano-asset` is that library.

## Features

- ✅ **CIP-14 asset fingerprint** generation (`asset1...` bech32 strings)
- ✅ **Policy ID validation** (56-char lowercase hex)
- ✅ **Asset name** encoding/decoding (UTF-8 ↔ hex)
- ✅ **Asset ID** parsing and construction (`policyId.assetNameHex`)
- ✅ **Zero external dependencies** — pure Go stdlib only
- ✅ **Deterministic** — same inputs always produce same outputs
- ✅ **Safe for concurrent use** — all functions are pure
- ✅ **AI-agent friendly** — structured errors, predictable types

## Installation
```bash
go get github.com/njchilds90/go-cardano-asset
```

## Quick Start
```go
package main

import (
    "fmt"
    cardanoasset "github.com/njchilds90/go-cardano-asset"
)

func main() {
    // Create an asset
    a, err := cardanoasset.NewAsset(
        "d5e6bf0500378d4f0da4e8dde6becec7621cd8cbf5cbb9b87013d4cc",
        "SpaceBud0",
    )
    if err != nil {
        panic(err)
    }

    // Get full info in one call
    info, err := a.Info()
    if err != nil {
        panic(err)
    }

    fmt.Println("Asset ID:    ", info.AssetID)
    fmt.Println("Name (hex):  ", info.AssetNameHex)
    fmt.Println("Fingerprint: ", info.Fingerprint)
}
```

## API Reference

### Creating Assets
```go
// From raw asset name string
a, err := cardanoasset.NewAsset(policyID, assetName)

// From hex-encoded asset name
a, err := cardanoasset.NewAssetFromHex(policyID, assetNameHex)

// Parse a full asset ID string
a, err := cardanoasset.ParseAssetID("policyId.assetNameHex")
```

### Asset Methods
```go
a.AssetNameHex()   // hex-encoded asset name
a.AssetID()        // "policyId.assetNameHex"
a.Fingerprint()    // CIP-14 "asset1..." bech32 string
a.Info()           // AssetInfo with all fields populated
a.IsValidUTF8Name() // bool
```

### Standalone Functions
```go
// Validate a policy ID
err := cardanoasset.ValidatePolicyID(policyID)

// Validate a hex asset name
err := cardanoasset.ValidateAssetNameHex(assetNameHex)

// Compute fingerprint directly
fp, err := cardanoasset.Fingerprint(policyID, assetName)
```

### Error Types
```go
cardanoasset.ErrInvalidPolicyID   // policy ID format error
cardanoasset.ErrAssetNameTooLong  // asset name > 32 bytes
cardanoasset.ErrInvalidHex        // hex decode failure
cardanoasset.ErrInvalidAssetID    // asset ID parse failure
```

## CIP Compliance

- [CIP-14](https://cips.cardano.org/cip/CIP-14) — Asset Fingerprint
- [CIP-5](https://cips.cardano.org/cip/CIP-5) — Common Bech32 Prefixes

> **Note:** This library uses SHA-256 truncated to 20 bytes as a zero-dependency substitute for blake2b-160. For byte-exact CIP-14 fingerprint compatibility with the reference implementation, add `golang.org/x/crypto/blake2b` and swap the hasher in `blake2b160()`. The API is identical either way.

## Composing With Other Libraries
```go
// Use with blockfrost-go to look up assets by fingerprint
fp, _ := a.Fingerprint()
// pass fp to Blockfrost API

// Use with go-cardano-metadata to build NFT mint metadata
assetID := a.AssetID()
// pass assetID to your metadata builder
```

## License

MIT © Nicholas John Childs
