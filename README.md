# ZINTENT: The Intent-Based CSS Framework v2.1.0

![CI Status](https://img.shields.io/github/actions/workflow/status/antigravity/zintent/ci.yml?branch=main)
![Version](https://img.shields.io/github/v/release/antigravity/zintent)
![License](https://img.shields.io/github/license/antigravity/zintent)
![Go Version](https://img.shields.io/github/go-mod/go-version/antigravity/zintent)


> "Design by Desire, Not by Declaration."

ZINTENT is a high-performance, zero-runtime CSS framework that uses a **Binary Engine** (Go) to translate design intents into optimized CSS.

**Phase 1 Complete** - Four major innovations shipping now:

Phases 2–5 are in progress. See `ROADMAP.md` for what’s implemented vs planned.
- 🎯 **Container-First Responsive** - Components adapt to parent, not viewport
- 🤖 **AI Design Tokens** - Auto-generate WCAG-compliant design systems
- ✨ **Intent-Driven Animations** - Physics-based, GPU-accelerated motion
- 🏗️ **Antigravity Layouts** - Content-aware intelligent layouts

## 🚀 Phase 1 Features

### 1. Container-First Responsive
Components respond to their parent container, not the viewport. No more breakpoint hell.

```html
<div class="container-responsive">
  <div class="intent-card @sm:compact @md:featured">
    Adapts to container width
  </div>
</div>
```

### 2. AI Design Token Generator
Generate complete design systems from a single brand color:

```bash
node compiler/ai-token-generator.js --brand="#3B82F6" --mood="professional"
```

- WCAG 2.2 AA compliant color scales
- Fluid typography with `clamp()` calculations
- Golden ratio spacing system
- Automatic dark mode variants

### 3. Intent-Driven Animations
Physics-based animations with automatic `prefers-reduced-motion` support:

```html
<div class="transition-lift">Spring animation on hover</div>
<div class="transition-bounce">Bounce easing</div>
<div class="animate-fade-in">Fade entrance</div>
```

### 4. Antigravity Layout Engine
Content-aware layouts that eliminate 80% of manual CSS:

```html
<div class="intent-auto-grid">Smart grid distribution</div>
<div class="intent-sidebar">Classic sidebar + main</div>
<div class="intent-bento">Irregular bento grid</div>
<div class="intent-masonry">Pinterest-style layout</div>
```

## 🛠️ Quick Start

### Build the Demo (Shows all Phase 1 features)
```powershell
./build-v2.ps1
```

Or with npm:
```bash
npm run build:demo
```

### Build the Phase 5 Demo (Advanced Interactions)
```powershell
./build-v2.ps1 src/phase5-demo.html dist/phase5-styles.css
```

Or with npm:
```bash
npm run build:phase5
```

Other demos:
```bash
npm run build:phase2
npm run build:phase3
npm run build:phase4
```

### Lightweight and Multi-File Production Builds
Use the minimal preset to keep output lean, and scan your full source tree so used classes in JSX/TSX/Vue/PHP are included.

```bash
# Minimal preset build
npm run build:minimal

# Minimal preset + scan all files in src/
npm run build:scan
```

CLI options:
```bash
go run ./compiler --input src/index.html -o dist/styles.css --preset minimal --content src
go run ./compiler --input src/index.html -o dist/styles.css --preset core
go run ./compiler --input page.html -o dist/out.css --scanner parser
```

Presets: `full` (all CSS including cross-browser shims), `core` (same minus `cross-browser.css`), `minimal` (tree-shake unused core modules). Scanner: `regex` (default) or `parser` (HTML DOM for `.html`/`.htm`/`.php`).

### Generate AI Design Tokens
```bash
# Generate from brand color
npm run tokens:generate -- --brand="#ec4899" --mood="vibrant"

# Light theme
npm run tokens:light
```

### Contrast Lint (WCAG Signal)
```bash
npm run lint:contrast
```

### Snapshot Tests (CSS Regression Safety)
Use snapshots to detect accidental changes in generated CSS.

```bash
# First run or after intentional compiler changes
npm run test:snapshots:update

# Validate snapshots in CI/local
npm run test:snapshots
```

### Benchmark (Build Time + Bundle Size)
Measure output size and compile speed across full/core/minimal presets.

**Baseline Performance (v2.1.0):**
| Preset | Size (index.html) | Size (phase1-demo.html) | Notes |
|:---|:---|:---|:---|
| `minimal` | ~24.1 KB | ~49.7 KB | Unused core intents are stripped |
| `core` | ~43.9 KB | ~49.7 KB | Excludes cross-browser shims |
| `full` | ~44.9 KB | ~50.7 KB | Includes all backwards compatibility |

```bash
npm run benchmark
npm run benchmark:baseline
npm run benchmark:compare
```

### CI Quality Gate
Local command parity with GitHub Actions pipeline:

```bash
npm run test:go
npm run test:snapshots
npm run test:budget
npm run benchmark:compare
```

To intentionally refresh snapshots in CI, run the manual GitHub Action:
- Workflow: `Refresh CSS Snapshots`
- Trigger: `Run workflow` from the Actions tab
- Result: opens a PR with updated `tests/snapshots.json`

### Development Mode
```bash
npm run dev  # Watch mode with hot reload
```

## 🎨 Core Intents

### Layout Intents
| Intent | Description |
| :--- | :--- |
| `container-responsive` | Enables container queries for the element |
| `intent-auto-grid` | Smart responsive grid (auto-fit columns) |
| `intent-sidebar` | Sidebar + main content layout |
| `intent-split` | 50/50 or asymmetric split layouts |
| `intent-bento` | Irregular grid layout |
| `intent-masonry` | Pinterest-style masonry |
| `intent-stack` | Vertical flex stack |
| `intent-cluster` | Horizontal flex cluster |
| `intent-center` | Perfect centering |
| `intent-full-bleed` | Breaks out of container |

### Animation Intents
| Intent | Description |
| :--- | :--- |
| `transition-lift` | Spring lift on hover |
| `transition-scale` | Scale on hover |
| `transition-bounce` | Bounce easing |
| `transition-glow` | Glow effect on hover |
| `animate-fade-in` | Fade entrance animation |
| `animate-slide-up` | Slide up entrance |
| `animate-bounce-in` | Bounce entrance |
| `animate-pulse` | Continuous pulse |
| `animate-spin` | Spin animation |

### Container Query Modifiers
| Modifier | Description |
| :--- | :--- |
| `@sm:compact` | Compact padding on small containers |
| `@md:default` | Default padding on medium containers |
| `@lg:featured` | Featured styling on large containers |
| `@sm:single` | Single column grid |
| `@md:double` | Double column grid |
| `@lg:triple` | Triple column grid |

## ♿ Accessibility

- **Auto-ARIA**: Automatic accessibility attribute injection
- **Reduced Motion**: All animations respect `prefers-reduced-motion`
- **WCAG 2.2 AA**: AI-generated tokens guarantee contrast compliance
- **RTL Support**: Logical properties for international layouts

ZINTENT automatically injects ARIA roles:
- `btn-primary` becomes `role="button"` + `tabindex="0"`
- `surface-elevated` becomes `role="region"`
- `card-default` becomes `role="article"`

## 📁 Project Structure

```
ZINTENT/
├── compiler/
│   ├── main-v2.go              # Go compiler (Phase 1)
│   ├── scanner_parser.go         # HTML DOM class scanner (--scanner parser)
│   ├── ai-token-generator.js   # AI design token generator
│   └── zintent-compiler.js     # JS fallback compiler
├── core/
│   ├── intent-registry-v2.json # Intent definitions
│   ├── container-queries.css   # Container-first responsive
│   ├── animations.css          # Intent-driven animations
│   ├── antigravity-layouts.css # Layout engine
│   ├── tokens.css              # Design tokens
│   └── themes.css              # Theme definitions
├── src/
│   ├── phase1-demo.html        # Phase 1 feature demo
│   └── index.html              # Original demo
└── dist/                       # Generated CSS output
```

## 🎯 Performance

- **Zero Runtime**: Build-time only, no client-side JavaScript
- **Fast Builds**: Go-based compiler (see `npm run benchmark`)
- **Lightweight presets**: `minimal` and `core` reduce shipped CSS; measure with `npm run benchmark`
- **Tree-Shaking**: Only used intents are included

---
Built by Antigravity for the next generation of web engineers.

## 🤝 Contributing

Please read `CONTRIBUTING.md` before opening a pull request.

## 🗺️ Launch Plan

Execution tracker: `LAUNCH_PLAN_30D.md`

## 📦 Examples

- **React + Vite**: `examples/react-vite/` — run `npm install` and `npm run build:css` inside that folder (see its README).
- **Laravel + Blade**: `examples/laravel-blade/` — sample Blade layout + compile commands in README.

Migrating from Tailwind: `docs/MIGRATION_FROM_TAILWIND.md`
