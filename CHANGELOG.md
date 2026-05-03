# Changelog

All notable changes to ZINTENT will be documented in this file.

This project follows [Semantic Versioning](https://semver.org/) and uses
[Keep a Changelog](https://keepachangelog.com/) format.

## [2.1.0] - 2026-05-01

### Added

- **HTML DOM Parser Scanner** (`--scanner parser`): Parser-backed HTML/PHP class
  extraction using `golang.org/x/net/html`. JSX/TSX/Vue continue via regex
  scanner. Select with `--scanner parser` for `.html`/`.htm`/`.php` files.
- **Preset: `--preset minimal`**: Tree-shakes unused core CSS modules
  (`intent-`, `antigravity-`, animation prefixes) for lighter production
  bundles.
- **Preset: `--preset core`**: Full framework stack minus
  `core/cross-browser.css` for modern-browser-only targets.
- **Multi-file content scanning** (`--content <dir>`): Scan an entire source
  directory so classes used in JSX/TSX/Vue/PHP templates are included.
- **CSS snapshot regression tests**: `npm run test:snapshots` /
  `npm run test:snapshots:update` to catch unintended CSS output changes.
- **Performance benchmark suite**: `npm run benchmark`,
  `npm run benchmark:baseline`, `npm run benchmark:compare` with size gate
  against `tests/benchmarks-baseline.json`.
- **Bundle budget gate**: `npm run test:budget` enforces max output size.
- **Go unit tests for compiler**: `npm run test:go` (scanner, parser, edge
  cases).
- **CI pipeline** (`.github/workflows/ci.yml`): Runs Go tests, snapshot
  verification, budget gate, and benchmark comparison on every push/PR.
- **Manual snapshot refresh workflow** (`snapshot-refresh.yml`): Trigger from
  Actions tab to regenerate snapshots and open a PR automatically.
- **PR template** and **issue templates** (bug report, feature request).
- **Repo governance**: `CONTRIBUTING.md`, `.github/CODEOWNERS`.
- **Starter templates**: React + Vite (`examples/react-vite/`) and Laravel +
  Blade (`examples/laravel-blade/`).
- **Tailwind migration guide**: `docs/MIGRATION_FROM_TAILWIND.md`.
- **Improved CLI help**: Clearer `--help` output with common production
  workflow examples.

### Changed

- README performance claims updated to match measured benchmark numbers.
- Watch mode reliability improved (mod-time based change detection).
- CLI hardened: consistent support for `--input`, `--output`, `-o`, and
  positional argument patterns.

### Fixed

- Deterministic CSS output ordering across all presets.
- Scanner edge cases for deeply nested HTML structures.

### Performance

Baseline measurements (from `tests/benchmarks-baseline.json`):

| Build | Size |
|:------|:-----|
| `index-full` | 43.9 KB |
| `index-core` | 42.9 KB |
| `index-minimal` | 23.6 KB |
| `phase1-full` | 49.6 KB |
| `phase1-core` | 48.6 KB |
| `phase1-minimal` | 48.6 KB |

## [2.0.0] - 2026-04-14

### Added

- Initial public release of ZINTENT v2 with Go-based zero-runtime compiler.
- Intent-Based Styling engine with registry-driven CSS generation.
- Antigravity Layout Engine (auto-grid, sidebar, split, bento, masonry).
- Container-First Responsive system (container queries, not viewport).
- AI Design Token Generator with WCAG 2.2 AA compliance.
- Intent-Driven Animations with `prefers-reduced-motion` support.
- Auto-ARIA Accessibility Engine.
- Multi-Theme Runtime Engine (light, dark, high-contrast, midnight, vibrant,
  forest, nordic).
- Phases 1–10 feature demos.
