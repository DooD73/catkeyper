## Summary

Adds a foundational automated unit test suite for CatKeyper's pure Go logic. Previously, CI only built the app and quality relied entirely on manual macOS checks.

**What changed:**

- Extracted unlock chord and lock/unlock decision logic from `main.go` into [`pkg/lockmanager/`](pkg/lockmanager/)
- Extracted icon rendering helpers from `cmd/iconrender/main.go` into [`internal/iconrender/`](internal/iconrender/)
- Added 15 unit tests for chord matching (sequential and simultaneous Shift+C+A+T), edge cases, and keyboard decision behavior
- Added 3 unit tests for icon rasterization (`FillRoundedRectInset`, `RenderIcon`, `WritePNG`)
- Added `make test` target (`go test -v ./...`) to the Makefile
- Wired `make test` into CI (`.github/workflows/ci.yml`) after `go mod verify`
- Documented automated vs manual testing in `CONTRIBUTING.md`

**Out of scope (unchanged):** E2E tests, UI tests, Accessibility/CGEventTap integration tests.

**Acceptance criteria:**

- [x] `*_test.go` files cover unlock chord logic (Shift+CAT sequential and simultaneous)
- [x] Tests cover locked vs unlocked keyboard decision behavior
- [x] At least one test for icon rendering (`internal/iconrender`; `cmd/iconrender` delegates to it)
- [x] `make test` runs `go test ./...`
- [x] `CONTRIBUTING.md` documents how to run tests
- [x] `go test ./...` passes locally on macOS

## Related issue

No linked issue.

## Test plan

- [x] `make test` passes (`go test -v ./...`)
- [x] `go test ./...` passes on macOS
- [x] `go build .` succeeds after refactor
- [ ] Built locally with `make release` or `make run`
- [ ] Tested lock/unlock on macOS (manual — confirms CGO hook wiring after refactor)
- [ ] Tested Shift+CAT unlock chord sequential and simultaneous (manual)
- [ ] Updated `CHANGELOG.md` under **Unreleased** (if user-visible) — **N/A**: developer-only test infrastructure, no user-visible behavior change

### Automated test coverage

```bash
make test
# or
go test -v ./pkg/lockmanager/...
go test -v ./internal/iconrender/...
```

| Package | Tests | Covers |
|---------|-------|--------|
| `pkg/lockmanager` | 15 | Sequential/simultaneous chord, wrong order, partial chord, no shift, shift release, timeouts, unlocked passthrough, locked suppress, chord unlock, `SetLocked` reset |
| `internal/iconrender` | 3 | Rounded rect fill, SVG render, PNG write |

## Screenshots / recordings

N/A — no UI changes.

## Checklist

- [x] PR is focused on a single change set
- [x] No build artifacts committed
- [x] Third-party assets documented in `assets/ATTRIBUTION.md` (if added) — N/A, no new assets
