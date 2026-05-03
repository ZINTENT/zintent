/**
 * ZINTENT Compiler Prototype (v0.1)
 * The Intelligent Styling Engine
 */

const fs = require('fs');
const path = require('path');

// Core Intent Database
const INTENT_MAP = {
  'center-content': `
    display: grid;
    place-items: center;
    text-align: center;
  `,
  'auto-layout': `
    display: flex;
    flex-direction: column;
    gap: var(--zi-space-4);
  `,
  'stack-h': `
    display: flex;
    flex-direction: row;
    align-items: center;
    gap: var(--zi-space-4);
  `,
  'surface-elevated': `
    background-color: var(--zi-bg-surface);
    border-radius: var(--zi-radius-md);
    box-shadow: var(--zi-shadow-md);
    padding: var(--zi-space-4);
  `,
  'interaction-lift': `
    transition: var(--zi-transition-fast);
  `,
  'interaction-lift:hover': `
    transform: translateY(-2px);
    box-shadow: var(--zi-shadow-lg);
  `,
  // Antigravity Engine (Content-Aware)
  'antigravity-fill': `
    display: flex;
    flex-wrap: wrap;
    gap: var(--zi-space-4);
    align-content: flex-start;
  `,
  'antigravity-center': `
    display: grid;
    place-items: center;
    gap: var(--zi-space-4);
    width: 100%;
  `
};

const ARIA_MAP = {
  'btn-primary': { 'role': 'button', 'tabindex': '0' },
  'btn-dropdown': { 'role': 'button', 'aria-haspopup': 'true', 'aria-expanded': 'false' },
  'surface-elevated': { 'role': 'region', 'aria-label': 'Content Section' }
};

function injectAria(html) {
  console.log(`♿ ZINTENT: Running Accessibility Engine...`);
  let modifiedHtml = html;

  for (const [cls, attrs] of Object.entries(ARIA_MAP)) {
    const regex = new RegExp(`class="([^"]*\\b${cls}\\b[^"]*)"(?![^>]*(${Object.keys(attrs).join('|')}))`, 'g');
    
    modifiedHtml = modifiedHtml.replace(regex, (match, classes) => {
      const attrStr = Object.entries(attrs).map(([k, v]) => `${k}="${v}"`).join(' ');
      return `class="${classes}" ${attrStr}`;
    });
  }

  return modifiedHtml;
}

function compile(inputFile, outputFile) {
  console.log(`🚀 ZINTENT: Compiling ${inputFile}...`);
  
  if (!fs.existsSync(inputFile)) {
    console.error(`❌ Error: Input file ${inputFile} not found.`);
    return;
  }

  let html = fs.readFileSync(inputFile, 'utf-8');
  
  // Phase 1: ARIA Injection
  html = injectAria(html);
  fs.writeFileSync(inputFile, html); // Update HTML with ARIA
  const usedClasses = new Set();
  
  // Simple regex to find classes in HTML
  const classRegex = /class="([^"]+)"/g;
  let match;
  while ((match = classRegex.exec(html)) !== null) {
    match[1].split(/\s+/).forEach(cls => usedClasses.add(cls));
  }

  // Build the CSS output
  let cssOutput = `/* ZINTENT Generated Styles */\n\n`;
  
  // Add tokens first
  const tokensPath = path.join(__dirname, '../core/tokens.css');
  if (fs.existsSync(tokensPath)) {
    cssOutput += fs.readFileSync(tokensPath, 'utf-8') + '\n\n';
  }

  // Add used intents
  usedClasses.forEach(cls => {
    if (INTENT_MAP[cls]) {
      cssOutput += `.${cls} {\n${INTENT_MAP[cls]}\n}\n\n`;
      
      // Check for hover variants
      const hoverKey = `${cls}:hover`;
      if (INTENT_MAP[hoverKey]) {
        cssOutput += `.${cls}:hover {\n${INTENT_MAP[hoverKey]}\n}\n\n`;
      }
    }
  });

  // Ensure output directory exists
  const outputDir = path.dirname(outputFile);
  if (!fs.existsSync(outputDir)) {
    fs.mkdirSync(outputDir, { recursive: true });
  }

  fs.writeFileSync(outputFile, cssOutput);
  console.log(`✅ ZINTENT: Successfully built ${outputFile}`);
}

// Simple CLI handling
const args = process.argv.slice(2);
if (args.length >= 3 && args[1] === '-o') {
  compile(args[0], args[2]);
} else {
  console.log("Usage: node zintent-compiler.js <input.html> -o <output.css>");
}
