# Security Policy

## Supported versions

Security fixes are provided for the latest release on [GitHub Releases](https://github.com/DooD73/catkeyper/releases).

| Version | Supported |
| ------- | --------- |
| 1.0.x   | Yes       |
| < 1.0   | No        |

## Reporting a vulnerability

Please do **not** open a public GitHub issue for security vulnerabilities.

Email **diegofaso00@gmail.com** with:

- A description of the issue
- Steps to reproduce
- Impact assessment, if known
- Your GitHub username, if you would like attribution

You should receive a response within 7 days. If the report is accepted, we will coordinate a fix and disclosure timeline before publishing details.

## Scope

In scope:

- Privilege escalation or sandbox escape in CatKeyper
- Unexpected keyboard input leakage while locked
- Malicious behavior introduced through the app bundle or update flow

Out of scope:

- Social engineering
- Issues that require physical access to an unlocked Mac
- macOS platform bugs outside CatKeyper's control
