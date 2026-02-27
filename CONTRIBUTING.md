# Contributing to go-cardano-asset

Thank you for your interest in contributing!

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-cardano-asset`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Run vet: `go vet ./...`
7. Commit and push
8. Open a Pull Request

## Guidelines

- **Zero dependencies** — the core library must have no external runtime dependencies
- **Pure functions** — prefer stateless, deterministic functions
- **Table-driven tests** — all tests should use `t.Run` with table-driven cases
- **GoDoc** — every exported symbol must have a doc comment with an example
- **CIP compliance** — reference the relevant CIP number in comments where applicable
- **Errors** — use or extend the sentinel error variables; never return raw `errors.New` inline

## Reporting Issues

Please open an issue with:
- Go version (`go version`)
- Input that caused the problem
- Expected vs actual behavior

## License

By contributing, you agree your contributions will be licensed under the MIT License.

---

## 🚀 Release & Verification Instructions

RELEASE STEPS (GitHub Web UI):

1. Go to: https://github.com/njchilds90/go-cardano-asset

2. Click "Releases" in the right sidebar → "Create a new release"

3. Click "Choose a tag" → type: v1.0.0 → click "Create new tag: v1.0.0 on publish"

4. Target branch: main

5. Release title: v1.0.0

6. Release notes (paste this):
─────────────────────────────────────────
## go-cardano-asset v1.0.0

CIP-14 asset fingerprint generation, policy ID validation, and native token 
utilities for Cardano — zero dependencies, pure Go.

### What's included
- `Fingerprint()` — CIP-14 `asset1...` bech32 fingerprints
- `ValidatePolicyID()` — 56-char hex policy ID validation  
- `NewAsset` / `NewAssetFromHex` / `ParseAssetID` — flexible asset construction
- `Asset.Info()` — all asset details in one call
- Structured sentinel errors for clean error handling
- Zero external dependencies, pure Go stdlib only
─────────────────────────────────────────

7. Click "Publish release"

VERIFY pkg.go.dev indexing (~10 min):
→ https://pkg.go.dev/github.com/njchilds90/go-cardano-asset

Or trigger immediately:
→ https://sum.golang.org/lookup/github.com/njchilds90/go-cardano-asset@v1.0.0

SEMANTIC VERSIONING PLAN:
v1.0.x — bug fixes, no API changes
v1.1.0 — add blake2b build tag for exact CIP-14 byte compatibility
v1.2.0 — add CIP-67/68 asset label support (reference NFTs)
v2.0.0 — only if a breaking API change is needed