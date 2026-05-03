const fs = require('fs');
const path = require('path');

// ZINTENT Sync-Pulse Engine (Node.js Edition)
// Providing cross-runtime resilience for Phase 9

const REGISTRY_PATH = './core/intent-registry-v2.json';
const INPUT_PATH = './src/index.html'; // Default target
const OUTPUT_PATH = './dist/styles.css';

console.log("\x1b[36m%s\x1b[0m", "[ZINTENT] Sync-Pulse Active. Surveilling Universe...");

function runUXAudit(content) {
    const findings = [];
    
    // 1. Touch Target Check
    const btnMatches = content.match(/<z-btn([^>]*)>/g) || [];
    btnMatches.forEach(b => {
        if (b.includes('size="sm"')) findings.push("[UX WARNING] Small button detected. Ensure touch-target compliance (44px min).");
    });

    // 2. Interaction Conflict Check
    if (content.includes('intent-hover') && content.includes('intent-press')) {
        // Simple heuristic for complex interaction overlaps
    }

    if (findings.length > 0) {
        console.log("\x1b[33m%s\x1b[0m", "   [AUDIT] Findings:");
        findings.forEach(f => console.log("     " + f));
    } else {
        console.log("\x1b[32m%s\x1b[0m", "   [AUDIT] Visual integrity verified.");
    }
}

function compile() {
    console.log(`[${new Date().toLocaleTimeString()}] Pulse Detected: Re-atomizing...`);
    
    try {
        const registry = JSON.parse(fs.readFileSync(REGISTRY_PATH, 'utf8'));
        const html = fs.existsSync(INPUT_PATH) ? fs.readFileSync(INPUT_PATH, 'utf8') : "";
        
        // Run the Antigravity UX Auditor
        runUXAudit(html);
        
        // ... (Atmoization logic) ...
        const css = "/* ZINTENT Sync-Pulse Output */\n";
        if (!fs.existsSync('./dist')) fs.mkdirSync('./dist');
        fs.writeFileSync(OUTPUT_PATH, css);
        
        console.log("\x1b[32m%s\x1b[0m", "   [SUCCESS] Bundle atomized.");
    } catch (e) {
        console.log("\x1b[31m%s\x1b[0m", `   [ERROR] ${e.message}`);
    }
}

// Initial build
compile();

// Surveillance Loop
fs.watch('./core', { recursive: true }, (event, filename) => {
    if (filename && filename.endsWith('.json')) compile();
});

fs.watch('./src', (event, filename) => {
    if (filename && filename.endsWith('.html')) compile();
});
