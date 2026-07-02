# CatKeyper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![macOS](https://img.shields.io/badge/macOS-12%2B-blue)](https://www.apple.com/macos/)
[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)

A lightweight macOS keyboard lock for keeping accidental paw input out of calls, editors, and chat windows.

Lock keyboard input system-wide from the menu bar. Unlock with the button or hold **Shift** and type **C**, **A**, **T**.

## Download

Pre-built macOS binaries are available on [GitHub Releases](https://github.com/DooD73/catkeyper/releases/latest):

- [ZIP](https://github.com/DooD73/catkeyper/releases/latest/download/catkeyper-macos.zip) (recommended)
- [DMG](https://github.com/DooD73/catkeyper/releases/latest/download/catkeyper-macos.dmg)

Requires macOS 12.0 (Monterey) or later.

## Install

1. Download and open `CatKeyper.app`.
2. Grant **Accessibility** permission when macOS asks.

If macOS does not prompt automatically:

`System Settings > Privacy & Security > Accessibility`

Add and enable `CatKeyper.app`, then restart the app.

> **Gatekeeper note:** Release builds are ad-hoc signed. macOS may show a warning on first launch. Right-click the app and choose **Open**.

## Usage

- Click **Lock Keyboard** in the menu bar to block key presses.
- Click **Unlock** or use the Shift+CAT chord to restore input.
- Enable debug logs when troubleshooting:

```bash
CATKEYPER_DEBUG=1 open "build/CatKeyper.app"
```

## Build from source

Requirements:

- macOS 12.0+
- Go 1.22+
- Xcode Command Line Tools

```bash
git clone https://github.com/DooD73/catkeyper.git
cd catkeyper
make release
```

This creates:

- `build/CatKeyper.app`
- `build/catkeyper-<version>-macos.zip`
- `build/catkeyper-<version>-macos.dmg`

Run locally without packaging:

```bash
make run
```

Maintainers: see [RELEASING.md](RELEASING.md) for the full release checklist.

## Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](CONTRIBUTING.md) before opening a pull request.

1. Fork the repo
2. Create a feature branch
3. Submit a pull request with a clear description

## Security

To report a security issue, follow [SECURITY.md](SECURITY.md). Please do not file public issues for vulnerabilities.

## License

MIT License. See [LICENSE](LICENSE).

Copyright (c) 2026 Diego Fasolo

## Assets

The menu bar icon is derived from [OpenMoji](https://openmoji.org/). See [assets/ATTRIBUTION.md](assets/ATTRIBUTION.md).

## Links

- Website: https://catkeyper.com
- Repository: https://github.com/DooD73/catkeyper
- Issues: https://github.com/DooD73/catkeyper/issues
