# WOL Manager Release Materials

This folder contains the materials prepared for submitting WOL Manager to Synology Package Center review.

## Files

- `synology-submission.md`: package summary, compatibility, support information, and submission checklist.
- `operation-manual.md`: reviewer-facing user manual and core workflows.
- `test-scenarios.md`: installation, usage, update, uninstall, and regression test scenarios.
- `changelog.md`: release notes for version `1.0.0`.
- `email-reply-template.md`: template for replying to Synology when they ask for package and review materials.
- `screenshots/`: screenshots for the Package Center review submission.

## Package

Current package metadata:

- Product name: WOL Manager
- Package ID: WOLManager
- Version: 1.0.0
- Category: Utilities
- DSM compatibility: DSM 7.0 or later
- Architecture: x86_64
- Package path after build: `build/WOLManager-1.0.0-x86_64.spk`

Build command:

```sh
./build.sh
```

Recommended verification command:

```sh
go test ./...
```

