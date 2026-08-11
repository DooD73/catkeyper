# CatKeyper

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![macOS](https://img.shields.io/badge/macOS-12%2B-blue)](https://www.apple.com/macos/)
[![Go](https://img.shields.io/badge/Go-1.25%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)

A lightweight macOS keyboard lock for keeping accidental paw input out of calls, editors, and chat windows.

Lock keyboard input system-wide from the menu bar. Unlock with the button or hold **Shift** and type **C**, **A**, **T**.

## Download

Pre-built macOS binaries are available on [GitHub Releases](https://github.com/DooD73/catkeyper/releases/latest):

- [DMG](https://github.com/DooD73/catkeyper/releases/latest/download/catkeyper-macos.dmg) (recommended)
- [ZIP](https://github.com/DooD73/catkeyper/releases/latest/download/catkeyper-macos.zip) (archive alternative)

Both downloads include the same app; only the installation format differs.

Requires macOS 12.0 (Monterey) or later.

## Install

1. Open the recommended DMG and drag `CatKeyper.app` to **Applications**. If you downloaded the ZIP instead, unzip it first, then drag `CatKeyper.app` to **Applications**.
2. Open CatKeyper from **Applications**.
3. Grant **Accessibility** permission when macOS asks.

Copying CatKeyper into Applications does not grant Accessibility permission. macOS keeps that permission separate and requests it when you first open the app.

If macOS does not prompt automatically:

`System Settings > Privacy & Security > Accessibility`

Add and enable `CatKeyper.app`, then restart the app.

> **Gatekeeper note:** Release builds are ad-hoc signed. macOS may show a warning on first launch. Right-click the app and choose **Open**.

## Frequently asked questions

### How do I install CatKeyper?

Open the recommended DMG and drag CatKeyper to **Applications**. If you choose the ZIP instead, unzip it first, then drag `CatKeyper.app` to **Applications**.

### Why does macOS ask for Accessibility permission?

CatKeyper needs Accessibility permission to intercept keyboard input while locked. macOS keeps this access under your control in **System Settings**.

### What if macOS will not open CatKeyper?

On first launch, right-click `CatKeyper.app`, choose **Open**, then confirm. If locking still does not work, enable CatKeyper in **System Settings > Privacy & Security > Accessibility**, then restart CatKeyper.

### Which macOS versions are supported?

CatKeyper requires macOS 12.0 (Monterey) or later.

### What is CatKeyper?

CatKeyper is a free, open-source macOS menu bar app that blocks keyboard input system-wide while locked, so curious cat paws cannot type into your active app.

### How do I unlock the keyboard?

Click **Unlock Keyboard** in the menu bar app, or hold **Shift** and press **C**, **A**, **T** in sequence. Both methods remain available while the keyboard is locked.

### Does CatKeyper block my mouse or trackpad?

No. CatKeyper blocks only keyboard input while locked. Your mouse and trackpad remain usable, so you can still open the menu and click **Unlock Keyboard**.

### Is CatKeyper free?

Yes. CatKeyper is free and open source under the MIT License. Optional donations support development via Buy Me a Coffee.

## Usage

- Click **Lock Keyboard** in the menu bar to block key presses.
- Click **Unlock** or use the Shift+CAT chord to restore input.
- Enable debug logs when troubleshooting:

```bash
CATKEYPER_DEBUG=1 open "build/CatKeyper.app"
```

## Privacy

CatKeyper works locally on your Mac. It checks keyboard events only while locked, so it can block key presses and recognize Shift + CAT. It does not record, store, analyze, or transmit your keystrokes or any other data.

## Build from source

Requirements:

- macOS 12.0+
- Go 1.25+
- Xcode Command Line Tools
- `create-dmg` 1.2.3+ (`brew install create-dmg`)

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
