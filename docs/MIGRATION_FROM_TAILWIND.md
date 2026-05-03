# Migrating from Tailwind CSS to ZINTENT

ZINTENT is **intent-based** and **build-time**: you ship only the CSS your templates reference. Tailwind is **utility-first** with a huge default bundle unless you purge. Use this guide to translate patterns and plan your migration.

## Mental model

| Tailwind | ZINTENT |
|----------|---------|
| Utility classes (`flex`, `p-4`, `md:grid-cols-2`) | Semantic **intents** (`intent-stack`, `zi-p-4`, container `@md:double`) |
| Config + purge/content paths | Compiler `--content` + optional `--scanner parser` |
| JIT on demand | Go compiler + intent registry |
| Breakpoints on viewport | **Container-first** modifiers (`@sm:compact`) where you opt in |

## Common mappings (starting point)

Layout and spacing:

| Tailwind | ZINTENT direction |
|----------|-------------------|
| `flex flex-col gap-4` | `intent-stack` or `intent-stack-md` |
| `flex flex-wrap gap-4` | `intent-cluster` / `intent-cluster-md` |
| `grid` + columns | `intent-auto-grid`, `@md:double`, container queries |
| `max-w-* mx-auto` | `intent-box-md`, `intent-box-prose`, `zi-container` |
| `items-center justify-center` | `intent-center` |

Typography and color:

| Tailwind | ZINTENT direction |
|----------|-------------------|
| `text-sm text-slate-600` | Token-driven classes from registry (e.g. `zi-text-sm`, theme vars) |
| `font-bold` | `zi-font-bold` (if defined in your registry) |

Effects:

| Tailwind | ZINTENT direction |
|----------|-------------------|
| `transition hover:shadow-lg` | `transition-lift`, `intent-card-hover`, etc. |
| `animate-pulse` | `animate-pulse` / `transition-pulse` patterns in `core/animations.css` |

## Migration steps

1. **Pick a preset**  
   - `minimal`: smallest CSS, conditional core modules.  
   - `core`: production default without legacy `cross-browser.css`.  
   - `full`: everything including compatibility shims.

2. **Point the compiler at all templates**  
   ```bash
   go run ./compiler --input src/index.html -o dist/zintent.css --preset core --content src
   ```

3. **Replace viewport breakpoints with container queries** where components sit in sidebars or split layouts: wrap with `container-responsive` and use `@sm:` / `@md:` modifiers.

4. **Iterate on the intent registry**  
   Add project-specific intents in `core/intent-registry-v2.json` instead of one-off utilities when a pattern repeats.

5. **Validate**  
   Run `npm run test:snapshots`, `npm run benchmark:compare`, and visual QA in target browsers.

## What not to expect (yet)

- One-to-one class parity with Tailwind’s entire catalog.  
- Identical JIT ergonomics: ZINTENT favors **semantic names + compiler** over thousands of atomic utilities.

For questions, open an issue using the feature or bug template and mention “Tailwind migration”.
