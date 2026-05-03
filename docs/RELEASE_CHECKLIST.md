# ZINTENT Release Checklist

This document outlines the standard operating procedure for releasing a new version of the ZINTENT Framework.

## Pre-Release Phase

1. [ ] **Verify Quality Gates (CI)**
   - All tests pass (`npm run test:go`).
   - Snapshot tests maintain integrity (`npm run test:snapshots`).
   - CSS Budget constraints are satisfied (`npm run test:budget`).

2. [ ] **Update Documentation**
   - Confirm all new intents are documented.
   - If there are breaking changes, add an explicit migration path to `MIGRATION_FROM_TAILWIND.md` or a new migration doc.

3. [ ] **Benchmark Validation**
   - Run `npm run benchmark:compare` to verify no performance regressions exist.
   - If performance dropped intentionally due to structural changes, update `tests/benchmarks-baseline.json`.

## Version Update

4. [ ] **Bump Version Strings**
   - Update `package.json` version.
   - Update header/labels in `README.md`.
   - Update embedded CLI version string in `compiler/main-v2.go`.

5. [ ] **Prepare Changelog**
   - Review commit history and finalize `CHANGELOG.md` following [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
   - Ensure the new version section is changed from `[Unreleased]` to the explicit version number and date.

## Promotion Phase

6. [ ] **Git Tagging**
   - Commit the version bumps and changelog: `git commit -m "chore: prepare vX.Y.Z release"`
   - Create a corresponding git tag: `git tag -a vX.Y.Z -m "Release vX.Y.Z"`
   - Push tags to main repo: `git push origin main --tags`

7. [ ] **Publish**
   - Wait for CI artifact build triggered by git tag.
   - Create GitHub Release, copying the exact markdown block from `CHANGELOG.md` for that version.
   - Run `npm publish --access public` (if distributed via npm registry as an executable wrapper).

## Post-Release

8. [ ] **Create new `[Unreleased]` block** in `CHANGELOG.md` for the next cycle.
9. [ ] Announce release in community channels, highlighting major benchmarks or new intents.
