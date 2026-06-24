# Releasing CatKeyper

This guide is for maintainers cutting a new public release.

## Prerequisites

- macOS build machine with Go 1.22+ and Xcode Command Line Tools
- Push access to `https://github.com/DooD73/catkeyper`

## Version bump checklist

Update the version everywhere it appears:

1. `Makefile` — `VERSION ?= x.y.z`
2. `main.go` — `appVersion = "x.y.z"` (fallback when not built via Makefile)
3. `CHANGELOG.md` — move **Unreleased** items into a new `x.y.z` section with date
4. Website download URLs in `catkeyper-website` if filenames change

## Build release artifacts

```bash
make clean
make release
```

This produces:

- `build/CatKeyper.app`
- `build/catkeyper-<version>-macos.zip`
- `build/catkeyper-<version>-macos.dmg`

Artifacts are ad-hoc signed.

## Publish to GitHub Releases

1. Commit the version bump and changelog on `main`.
2. Create and push an annotated tag:

    ```bash
    git tag -a v1.0.0 -m "v1.0.0"
    git push origin v1.0.0
    ```

3. Open [GitHub Releases](https://github.com/DooD73/catkeyper/releases/new).
4. Choose the tag you just pushed.
5. Set the release title to `v1.0.0`.
6. Paste the relevant `CHANGELOG.md` section into the release notes.
7. Upload:
    - `build/catkeyper-<version>-macos.zip`
    - `build/catkeyper-<version>-macos.dmg`
8. Publish the release.

## Verify the release

- Download the ZIP from the release page on a clean Mac (or VM).
- Open `CatKeyper.app` and confirm Accessibility permission flow works.
- Lock the keyboard, confirm input is blocked, then unlock with the menu bar button and the Shift+CAT chord.

## Optional automation

CI builds the app on every push and tag via `.github/workflows/ci.yml`. Release uploads are still manual unless you add a dedicated release workflow later.
