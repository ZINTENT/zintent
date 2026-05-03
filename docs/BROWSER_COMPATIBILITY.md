# Browser Compatibility Matrix

ZINTENT is designed to provide cutting-edge CSS features (like Container Queries, Cascade Layers, and CSS Logical Properties) while maintaining a graceful degradation path for older browsers.

## Core Framework Support

| Browser | Minimum Version | Notes |
|---------|-----------------|-------|
| **Chrome / Edge** | 105+ | Full support (Container Queries, `:has()`, Cascade Layers) |
| **Firefox** | 121+ | Full support |
| **Safari (macOS / iOS)** | 16.0+ | Full support |
| **Opera** | 91+ | Full support |

## Progressive Enhancement Architecture

We use progressive enhancement to ensure functionality without breaking the layout on older browser versions.

### 1. Fallback Layouts (Pre-Container Queries)
Browsers that do not support `@container` queries will fall back to standard `min-width` `@media` queries where appropriate. The framework compiler injects `.container-responsive > *` fallbacks using flexbox wrapping when container units (`cqi`, `cqw`) are not understood.

### 2. Antigravity Animation Engine
Animations rely on spring physics translated into discrete `cubic-bezier` keyframes.
* All animations respect `prefers-reduced-motion: reduce`.
* If `@keyframes` or 3D transforms (`translate3d`) fail, elements will snap to their final state instantly rather than remaining invisible.

### 3. Glassmorphism and Backdrop Filters
The `backdrop-filter: blur(...)` property is used extensively in our Phase 6 `intent-glass` components.
* Browsers lacking this support will fall back to a subtly opaque background color (`rgba()`) which ensures text remains highly legible.

## Polyfills

Because ZINTENT prides itself on "Zero Runtime", we do **not** ship global polyfills out of the box. 

If you must support IE11 or legacy Safari (<= 14.x), do not use ZINTENT. The structural features required for Intention-Driven CSS are incompatible with deep legacy browser engines.
