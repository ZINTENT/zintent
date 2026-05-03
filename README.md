# ZINTENT: The Intent-Based CSS Framework v2.1.0

![CI Status](https://img.shields.io/github/actions/workflow/status/ZINTENT/zintent/ci.yml?branch=main)
![Version](https://img.shields.io/github/v/release/ZINTENT/zintent)
![License](https://img.shields.io/github/license/ZINTENT/zintent)
![Go Version](https://img.shields.io/github/go-mod/go-version/ZINTENT/zintent)

> **"Design by Desire, Not by Declaration."**

ZINTENT is a high-performance, zero-runtime CSS framework powered by a **Go-based Binary Engine**. It translates high-level design "intents" into optimized, tree-shaken CSS, eliminating runtime overhead and massive utility-class bloat.

## ✨ Core Innovations

- 🎯 **Container-First Responsive**: Components adapt to their parent container, not the viewport. Goodbye breakpoint hell.
- 🤖 **AI Design Tokens**: Auto-generate WCAG 2.2 AA compliant design systems from a single brand color.
- ✨ **Intent-Driven Animations**: Physics-based, GPU-accelerated motion with built-in `prefers-reduced-motion` support.
- 🏗️ **Antigravity Layout Engine**: Content-aware intelligent layouts (Bento, Masonry, Auto-Grid) with zero manual math.
- ♿ **Auto-ARIA Engine**: Automatic accessibility attribute injection based on your styling intents.

---

## 🚀 Quick Start

### 1. Installation
Clone the repository and install dependencies:
```bash
git clone https://github.com/ZINTENT/zintent.git
cd zintent
npm install
```

### 2. Build the Framework
Generate the production-ready CSS from the demo files:
```bash
npm run build:demo
```

### 3. Development Mode
Run the compiler in watch mode with hot reload:
```bash
npm run dev
```

---

## 🛠️ Production Workflow

ZINTENT is built for production efficiency. Use **Presets** to keep your final bundle as light as possible.

### Tree-Shaking & Presets
```bash
# Minimal build: Strips all unused core modules (recommended for production)
npm run build:minimal

# Full core build: Complete framework stack minus cross-browser shims
npm run build:core
```

### Performance Benchmarks (v2.1.0)
| Preset | Size (Standard Page) | Compile Speed |
|:---|:---|:---|
| `minimal` | **~24.1 KB** | < 100ms |
| `core` | **~43.9 KB** | < 100ms |
| `full` | **~44.9 KB** | < 120ms |

---

## 🎨 The Intent System

### Layout Intents
| Intent | Description |
| :--- | :--- |
| `container-responsive` | Enables container-queries for the element |
| `intent-auto-grid` | Smart responsive grid (auto-fit columns) |
| `intent-bento` | Irregular "Apple-style" grid layout |
| `intent-masonry` | Pinterest-style masonry distribution |
| `intent-sidebar` | Responsive sidebar + main content pattern |

### Animation Intents
| Intent | Description |
| :--- | :--- |
| `transition-lift` | Physics-based spring lift on hover |
| `transition-glow` | Subtle shadow-glow interaction |
| `animate-fade-in` | Smooth entrance animation |
| `animate-pulse` | Continuous high-performance pulse |

---

## ♿ Professional Accessibility
ZINTENT doesn't just style; it secures your accessibility.
- **Contrast Guard**: Build-time linting ensures your color tokens meet WCAG 2.2 AA standards.
- **Smart ARIA**: `btn-primary` automatically receives `role="button"` and `tabindex="0"`.
- **Motion Safety**: All animations respect OS-level "Reduce Motion" settings by default.

---

## 📦 Framework Integrations
Start building immediately with our production-ready templates:
- [React + Vite](examples/react-vite/)
- [Laravel + Blade](examples/laravel-blade/)
- [Tailwind Migration Guide](docs/MIGRATION_FROM_TAILWIND.md)

---

Built with ❤️ by the **ZINTENT Core Team** for the next generation of web engineers.

[Contributing](CONTRIBUTING.md) | [Roadmap](ROADMAP.md) | [Changelog](CHANGELOG.md)
