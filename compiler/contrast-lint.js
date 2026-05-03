#!/usr/bin/env node
/**
 * ZINTENT Contrast Lint (WCAG signal)
 * Checks a small set of token pairs across themes and reports contrast ratios.
 *
 * Usage:
 *   node compiler/contrast-lint.js
 */

const fs = require('fs');
const path = require('path');

const TOKENS_PATH = path.join(__dirname, '..', 'core', 'tokens.css');
const THEMES_PATH = path.join(__dirname, '..', 'core', 'themes.css');

function parseVarBlock(cssText) {
  const vars = new Map();
  const re = /--([a-zA-Z0-9_-]+)\s*:\s*([^;]+);/g;
  for (const m of cssText.matchAll(re)) {
    vars.set(`--${m[1]}`, m[2].trim());
  }
  return vars;
}

function parseThemes(cssText) {
  // Very small parser for blocks like: [data-theme="midnight"] { ... }
  const themes = new Map();
  const blockRe = /\[data-theme="([^"]+)"\]\s*\{([\s\S]*?)\}/g;
  for (const m of cssText.matchAll(blockRe)) {
    const name = m[1];
    const body = m[2];
    themes.set(name, parseVarBlock(body));
  }
  return themes;
}

function parseColor(value) {
  const v = value.trim().toLowerCase();
  if (v.startsWith('#')) {
    const hex = v.slice(1);
    if (hex.length === 3) {
      const r = parseInt(hex[0] + hex[0], 16);
      const g = parseInt(hex[1] + hex[1], 16);
      const b = parseInt(hex[2] + hex[2], 16);
      return { r, g, b, a: 1 };
    }
    if (hex.length === 6) {
      const r = parseInt(hex.slice(0, 2), 16);
      const g = parseInt(hex.slice(2, 4), 16);
      const b = parseInt(hex.slice(4, 6), 16);
      return { r, g, b, a: 1 };
    }
    return null;
  }

  const rgba = v.match(/^rgba?\(\s*([0-9.]+)\s*,\s*([0-9.]+)\s*,\s*([0-9.]+)(?:\s*,\s*([0-9.]+)\s*)?\)$/);
  if (rgba) {
    const r = Number(rgba[1]);
    const g = Number(rgba[2]);
    const b = Number(rgba[3]);
    const a = rgba[4] == null ? 1 : Number(rgba[4]);
    if ([r, g, b, a].some(Number.isNaN)) return null;
    return { r, g, b, a };
  }

  return null;
}

function srgbToLinear(c) {
  const x = c / 255;
  return x <= 0.03928 ? x / 12.92 : Math.pow((x + 0.055) / 1.055, 2.4);
}

function relativeLuminance({ r, g, b }) {
  const rl = srgbToLinear(r);
  const gl = srgbToLinear(g);
  const bl = srgbToLinear(b);
  return 0.2126 * rl + 0.7152 * gl + 0.0722 * bl;
}

function blendOver(bg, fg) {
  // alpha blend fg over bg
  const a = fg.a;
  const r = Math.round(fg.r * a + bg.r * (1 - a));
  const g = Math.round(fg.g * a + bg.g * (1 - a));
  const b = Math.round(fg.b * a + bg.b * (1 - a));
  return { r, g, b, a: 1 };
}

function contrastRatio(fg, bg) {
  const L1 = relativeLuminance(fg);
  const L2 = relativeLuminance(bg);
  const brightest = Math.max(L1, L2);
  const darkest = Math.min(L1, L2);
  return (brightest + 0.05) / (darkest + 0.05);
}

function getVar(vars, name) {
  return vars.get(name);
}

function buildThemeVars(baseVars, themeVars) {
  const merged = new Map(baseVars);
  for (const [k, v] of themeVars.entries()) merged.set(k, v);
  return merged;
}

function main() {
  const tokensCss = fs.readFileSync(TOKENS_PATH, 'utf8');
  const themesCss = fs.readFileSync(THEMES_PATH, 'utf8');

  const rootRe = /:root\s*\{([\s\S]*?)\}/;
  const rootMatch = tokensCss.match(rootRe);
  if (!rootMatch) {
    console.error('Could not find :root block in core/tokens.css');
    process.exit(1);
  }

  const baseVars = parseVarBlock(rootMatch[1]);
  const themes = parseThemes(themesCss);

  const checks = [
    { name: 'text-base on bg-base', fg: '--zi-text-base', bg: '--zi-bg-base', min: 4.5 },
    { name: 'text-muted on bg-base', fg: '--zi-text-muted', bg: '--zi-bg-base', min: 4.5 },
    { name: 'text-base on bg-surface', fg: '--zi-text-base', bg: '--zi-bg-surface', min: 4.5 },
    { name: 'text-muted on bg-surface', fg: '--zi-text-muted', bg: '--zi-bg-surface', min: 4.5 },
  ];

  const themeNames = [...themes.keys()];
  if (themeNames.length === 0) {
    console.error('No themes found in core/themes.css');
    process.exit(1);
  }

  let hasFailures = false;

  for (const themeName of themeNames) {
    const merged = buildThemeVars(baseVars, themes.get(themeName));

    console.log(`\nTheme: ${themeName}`);
    for (const check of checks) {
      const fgRaw = getVar(merged, check.fg);
      const bgRaw = getVar(merged, check.bg);
      if (!fgRaw || !bgRaw) {
        console.log(`  - ${check.name}: SKIP (missing var)`);
        continue;
      }

      const bg = parseColor(bgRaw);
      const fg = parseColor(fgRaw);
      if (!bg || !fg) {
        console.log(`  - ${check.name}: SKIP (unparsed color) fg=${fgRaw} bg=${bgRaw}`);
        continue;
      }

      const fgResolved = fg.a < 1 ? blendOver(bg, fg) : fg;
      const ratio = contrastRatio(fgResolved, bg);
      const ok = ratio >= check.min;
      if (!ok) hasFailures = true;
      console.log(`  - ${check.name}: ${ratio.toFixed(2)}:1 ${ok ? 'OK' : 'FAIL'}`);
    }
  }

  console.log('');
  if (hasFailures) {
    process.exitCode = 2;
  }
}

main();

