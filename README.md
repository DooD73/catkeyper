# CatKeyper

A lightweight macOS keyboard lock for keeping accidental paw input out of calls, editors, and chat windows.

## Build

```bash
make release
```

This creates:

- `build/CatKeyper.app`
- `build/catkeyper-1.0.0-macos.zip`
- `build/catkeyper-1.0.0-macos.dmg`

## Install

Open `build/CatKeyper.app`, then grant Accessibility permission when macOS asks.

If macOS does not prompt automatically, go to:

`System Settings > Privacy & Security > Accessibility`

Add and enable `CatKeyper.app`, then restart the app.

## Debug Logs

Production logs are emitted by default. Debug logs are opt-in:

```bash
CATKEYPER_DEBUG=1 open "build/CatKeyper.app"
```

## Signing And Notarization

Developer ID signing requires an Apple Developer certificate installed in Keychain.

```bash
make sign SIGN_IDENTITY="Developer ID Application: Your Name (TEAMID)"
```

Notarization requires a notarytool keychain profile:

```bash
xcrun notarytool store-credentials catkeyper-notary \
  --apple-id "you@example.com" \
  --team-id "TEAMID" \
  --password "app-specific-password"

make notarize \
  SIGN_IDENTITY="Developer ID Application: Your Name (TEAMID)" \
  NOTARY_PROFILE=catkeyper-notary
```

No valid Developer ID code-signing identity is currently installed on this Mac. The default release artifacts are ad-hoc signed, executable, and suitable for local testing, but they are not notarized for frictionless distribution to other Macs.
