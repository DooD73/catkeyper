# Contributing to CatKeyper

Thanks for your interest in contributing. CatKeyper is a small macOS menu bar app written in Go.

## Before you start

- Read the [README](README.md) for build and install instructions.
- Search [existing issues](https://github.com/DooD73/catkeyper/issues) to avoid duplicate work.
- For large changes, open an issue first so we can agree on the approach.

## Development setup

Requirements:

- macOS 12.0 (Monterey) or later
- Go 1.22 or later
- Xcode Command Line Tools (`xcode-select --install`)

```bash
git clone https://github.com/DooD73/catkeyper.git
cd catkeyper
make run
```

For a release-style build:

```bash
make release
```

## Testing

### Automated unit tests

Run the full suite:

```bash
make test
```

Or directly:

```bash
go test -v ./...
```

Run a single package:

```bash
go test -v ./pkg/lockmanager/...
go test -v ./internal/iconrender/...
```

Unit tests cover pure Go logic only:

- Unlock chord state machine (sequential and simultaneous Shift+C+A+T)
- Lock suppress / unlock decision engine
- Icon rasterization helpers

They do **not** exercise macOS hooks, Accessibility permissions, Fyne UI, or packaging.

### Manual macOS integration testing

Test on a real Mac when your change touches:

- CGEventTap keyboard hook behavior
- Accessibility or Input Monitoring permissions
- Fyne UI and menu bar presentation
- Packaging (`make release`, codesign, `.app` bundle)

For hook debugging, run with `CATKEYPER_DEBUG=1` (not required for unit tests).

## Making changes

1. Fork the repository and create a branch from `main`.
2. Make focused changes with clear commit messages.
3. Run `make test` for logic changes. Test on a real Mac when your change touches keyboard hooks, permissions, UI, or packaging.
4. Open a pull request against `main` and fill out the PR template.

## Pull request guidelines

- Keep PRs small and scoped to one fix or feature.
- Update `CHANGELOG.md` under **Unreleased** for user-visible changes.
- Do not commit build artifacts from `build/`.
- If you add third-party assets, document them in `assets/ATTRIBUTION.md`.

## Reporting bugs

Use the [bug report template](https://github.com/DooD73/catkeyper/issues/new?template=bug_report.yml) and include:

- macOS version
- CatKeyper version
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs (`CATKEYPER_DEBUG=1` if helpful)

## Feature requests

Use the [feature request template](https://github.com/DooD73/catkeyper/issues/new?template=feature_request.yml) and explain the problem you are trying to solve.

## Code style

- Match the existing Go style in the repository.
- Prefer simple, readable changes over large refactors.
- Keep macOS-specific behavior documented when it is not obvious.

## Questions

Open a GitHub issue with the **Question** label or reach out through the project website.
