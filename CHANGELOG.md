# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!--
Release section template for the automated release parser:

## [X.Y.Z] - YYYY-MM-DD

### Added

- Change summary
-->

## [Unreleased]

### Changed

- README: promoted DMG as the recommended download format and clarified that both DMG and ZIP contain the same app
- README: expanded install steps to explicitly cover the drag-to-Applications flow for both download formats
- README: added a Frequently Asked Questions section covering installation, Accessibility permission, Gatekeeper, supported macOS versions, unlocking, mouse/trackpad behavior, and licensing

## [1.0.1] - 2026-07-02

### Added

- Automated unit tests for unlock chord logic, lock/unlock decisions, and icon rendering
- `make test` target and CI test step
- "Contact Support" link at the bottom of the app window; clicking it opens the default mail client addressed to `catkeyper.support@gmail.com`

### Changed

- Updated support email to `catkeyper.support@gmail.com` in `SECURITY.md`

## [1.0.0] - 2026-06-24

### Added

- macOS menu bar app that blocks keyboard input system-wide while locked
- One-click lock/unlock from the menu bar
- Secret unlock chord: hold Shift and press C, A, T
- Release packaging via `make release` (`.app`, ZIP, and DMG)
- Open-source release files: license, contributing guide, issue templates, and CI

[Unreleased]: https://github.com/DooD73/catkeyper/compare/v1.0.1...HEAD
[1.0.1]: https://github.com/DooD73/catkeyper/compare/v1.0.0...v1.0.1
[1.0.0]: https://github.com/DooD73/catkeyper/releases/tag/v1.0.0
