#!/usr/bin/env node
/* eslint-disable no-console */
const crypto = require("crypto");
const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");

const ROOT = path.resolve(__dirname, "..");
const SNAPSHOT_PATH = path.join(ROOT, "tests", "snapshots.json");

const CASES = [
  {
    name: "index-full",
    args: ["--input", "src/index.html", "-o", "dist/test-index-full.css", "--preset", "full"],
    output: "dist/test-index-full.css",
  },
  {
    name: "index-minimal",
    args: ["--input", "src/index.html", "-o", "dist/test-index-minimal.css", "--preset", "minimal"],
    output: "dist/test-index-minimal.css",
  },
  {
    name: "phase1-full",
    args: ["--input", "src/phase1-demo.html", "-o", "dist/test-phase1-full.css", "--preset", "full"],
    output: "dist/test-phase1-full.css",
  },
];

function runCompiler(args) {
  const result = spawnSync("go", ["run", "./compiler", ...args], {
    cwd: ROOT,
    encoding: "utf8",
  });
  if (result.status !== 0) {
    throw new Error(`Compiler failed:\n${result.stdout}\n${result.stderr}`);
  }
}

function fileHash(filePath) {
  const bytes = fs.readFileSync(filePath);
  const hash = crypto.createHash("sha256").update(bytes).digest("hex");
  return { hash, sizeBytes: bytes.length };
}

function ensureSnapshotDir() {
  fs.mkdirSync(path.dirname(SNAPSHOT_PATH), { recursive: true });
}

function loadSnapshots() {
  if (!fs.existsSync(SNAPSHOT_PATH)) {
    return { version: 1, generatedAt: "", snapshots: {} };
  }
  return JSON.parse(fs.readFileSync(SNAPSHOT_PATH, "utf8"));
}

function saveSnapshots(data) {
  ensureSnapshotDir();
  fs.writeFileSync(SNAPSHOT_PATH, `${JSON.stringify(data, null, 2)}\n`, "utf8");
}

function main() {
  const updateMode = process.argv.includes("--update");
  const snapshots = loadSnapshots();
  snapshots.version = 1;
  snapshots.generatedAt = new Date().toISOString();
  snapshots.snapshots = snapshots.snapshots || {};

  let hasFailure = false;

  for (const c of CASES) {
    runCompiler(c.args);
    const absOutput = path.join(ROOT, c.output);
    const current = fileHash(absOutput);

    if (updateMode) {
      snapshots.snapshots[c.name] = current;
      console.log(`updated: ${c.name} (${current.sizeBytes} bytes)`);
      continue;
    }

    const expected = snapshots.snapshots[c.name];
    if (!expected) {
      console.error(`missing snapshot: ${c.name}`);
      hasFailure = true;
      continue;
    }

    if (expected.hash !== current.hash) {
      console.error(
        `snapshot mismatch: ${c.name}\n  expected: ${expected.hash}\n  current:  ${current.hash}\n  expected size: ${expected.sizeBytes} bytes\n  current size:  ${current.sizeBytes} bytes`
      );
      hasFailure = true;
    } else {
      console.log(`ok: ${c.name} (${current.sizeBytes} bytes)`);
    }
  }

  if (updateMode) {
    saveSnapshots(snapshots);
    console.log(`snapshots saved: ${path.relative(ROOT, SNAPSHOT_PATH)}`);
    return;
  }

  if (hasFailure) {
    process.exit(1);
  }

  console.log("all CSS snapshots passed");
}

main();
