# ZINTENT Plugin API

The ZINTENT Plugin API allows developers to extend the core framework with custom intent macros, semantic tokens, and layout engines. Because ZINTENT is a zero-runtime, build-time compiler, plugins operate by injecting definitions directly into the compiler's intent registry or by passing structured JavaScript/JSON configurations that ZINTENT merges during the build phase.

## Overview

A ZINTENT Plugin provides a standardized way to define:
1. **Custom Intents** (e.g., `.intent-glitch`, `.transition-snap`)
2. **Custom Macros** (e.g., `<z-alert>`, `<z-hero>`)
3. **Custom Design Tokens** (e.g., `theme: cyber`)
4. **Auto-ARIA Mappings** (Custom accessibility logic)

## 1. Plugin Configuration format

Plugins are typically defined as JSON files (or exported from JS modules) that match the structure of the core `intent-registry-v2.json`.

```json
{
  "name": "zintent-plugin-cyberpunk",
  "version": "1.0.0",
  "intents": {
    "intent-glitch": {
      "styles": "animation: glitch 1s linear infinite; position: relative;",
      "hover": "animation: glitch-anim 0.3s cubic-bezier(.25, .46, .45, .94) both infinite;",
      "nested": {
        "&::before, &::after": {
          "styles": "content: attr(data-text); position: absolute; top: 0; left: 0;"
        }
      }
    }
  },
  "macros": {
    "z-alert": "<div class=\"intent-stack zi-bg-red-500 zi-text-white zi-p-4\">{{content}}</div>"
  },
  "aria_mappings": {
    "intent-glitch": {
      "role": "presentation"
    }
  }
}
```

## 2. Consuming Plugins via CLI

We are extending the ZINTENT CLI to accept a `--plugin` flag, which can point to local JSON files. Multiple plugins can be chained; they merge from right to left (later plugins override earlier ones).

```bash
# Using a single plugin
go run ./compiler -i src/index.html -o dist/styles.css --plugin ./plugins/cyberpunk.json

# Using multiple plugins
go run ./compiler -i src/index.html -o dist/styles.css --plugin ./plugins/typography.json --plugin ./plugins/forms.json
```

## 3. Developing a Plugin

### Intent Requirements
Intents must follow the schema:
- `styles` (required): The base CSS for the class.
- `hover`, `focus`, `active`, `disabled` (optional): Interactive states.
- `nested` (optional): Complex media queries or pseudo-elements (`@media (min-width: 768px)`, `&::before`).

### Macro Requirements
Macros must utilize the `{{placeholder}}` syntax:
- `{{content}}`: Replaced by the inner HTML of the macro tag.
- `{{attribute_name}}`: Replaced by the specific attribute provided to the custom HTML tag.
- `{{slot_name}}`: Targets specific `<z-slot name="name">` declarations.

## 4. Native Go Plugins (Roadmap)

In addition to JSON configuration passing, future ZINTENT versions may support native Go plugins (`buildmode=plugin`) that can inspect the AST/DOM during the compile phase to rewrite nodes intelligently.

```go
// Example Native Plugin Interface (Proposed)
type ZintentPlugin interface {
    Name() string
    OnParseClassList(classes []string) []string
    OnGenerateAST(tag string, attrs map[string]string) string
}
```

## 5. Publishing

To ensure discovery, we recommend naming your GitHub repositories and npm packages using the prefix `zintent-plugin-*`.
