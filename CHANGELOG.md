# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Automated unit tests for unlock chord logic, lock/unlock decisions, and icon rendering
- `make test` target and CI test step

## [1.0.0] - 2026-06-24

### Added

- macOS menu bar app that blocks keyboard input system-wide while locked
- One-click lock/unlock from the menu bar
- Secret unlock chord: hold Shift and press C, A, T
- Release packaging via `make release` (`.app`, ZIP, and DMG)
- Open-source release files: license, contributing guide, issue templates, and CI

[Unreleased]: https://github.com/DooD73/catkeyper/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/DooD73/catkeyper/releases/tag/v1.0.0
