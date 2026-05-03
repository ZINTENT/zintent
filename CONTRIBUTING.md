# Contributing to ZINTENT

Thanks for contributing to ZINTENT.

## Development Setup

Prerequisites:
- Go 1.22+
- Node.js 20+

Install and verify:

```bash
npm run build:demo
```

The Go compiler is built as a module from the repo root:

```bash
go run ./compiler --help
```

## Branch and PR Workflow

1. Create a feature branch from `main`.
2. Make focused changes with clear commit messages.
3. Open a PR using the provided template.
4. Ensure CI passes before requesting review.

## Required Local Checks

Run these commands before opening a PR:

```bash
npm run test:go
npm run test:snapshots
npm run test:budget
npm run benchmark:compare
```

If your compiler change intentionally affects generated CSS:

```bash
npm run test:snapshots:update
```

Then include `tests/snapshots.json` in your PR.

## Performance and Size Expectations

- Keep framework output lightweight by default.
- Avoid introducing broad core CSS without intent-based gating.
- If bundle size increases, explain why in the PR template.

## Coding Guidelines

- Preserve backward compatibility for CLI flags when possible.
- Prefer deterministic output (stable ordering of generated CSS).
- Add tests or snapshots for behavior changes.
- Keep documentation updated (`README.md`, workflow docs, CLI options).

## Examples and migration

- App starters: `examples/react-vite/`, `examples/laravel-blade/`
- Tailwind migration guide: `docs/MIGRATION_FROM_TAILWIND.md`

## Issue and Feature Intake

- Use GitHub issue templates for bug reports and feature requests.
- Include repro steps and exact command output for compiler issues.

## Code Ownership

Review ownership is defined in `.github/CODEOWNERS`.
Update owners if team handles or organization names change.
