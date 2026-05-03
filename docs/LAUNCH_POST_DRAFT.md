# ZINTENT 2.1: The Intent-Based CSS Revolution

We are thrilled to announce the official launch of **ZINTENT v2.1.0**, the CSS framework that prioritizes **Intent over Declaration**.

## Why ZINTENT?

Most CSS frameworks force you to think in terms of properties. ZINTENT forces you to think in terms of **Outcomes**. By combining a high-performance Go-based compiler with an intelligent intent registry, we've created a system that is:

- **Zero-Runtime**: No JavaScript required for styling.
- **Ultra-Fast**: Sub-millisecond compilation for local development.
- **Accessibility-First**: Automatic ARIA injection and reduced-motion handling.
- **Container-Aware**: True container-first responsive design.

## Benchmark Performance

We don't just claim speed; we prove it.

| Build Preset | Size (index.html) | Speed (ms) |
| :--- | :--- | :--- |
| `minimal` | ~24.1 KB | < 5ms |
| `core` | ~43.9 KB | < 8ms |
| `full` | ~44.9 KB | < 10ms |

*Results measured on a standard developer workstation using `npm run benchmark`.*

## Migration Proof: Tailwind to ZINTENT

Coming from Tailwind? You'll feel right at home, but with more power. Our mapping system makes the transition seamless.

**Tailwind:**
```html
<div class="flex flex-col items-center justify-center p-8 bg-white shadow-lg rounded-xl">
  ...
</div>
```

**ZINTENT:**
```html
<div class="intent-stack intent-center zi-p-8 surface-elevated">
  ...
</div>
```

ZINTENT reduces the noise and focuses on the **meaning** of your layout.

## Get Started

```bash
npm install zintent
npx zintent --init
```

Check out our [Documentation](docs/BROWSER_COMPATIBILITY.md) and [Plugin API](docs/plugin-api.md) to start extending ZINTENT today.

---
Built with passion by the ZINTENT Core Team.
