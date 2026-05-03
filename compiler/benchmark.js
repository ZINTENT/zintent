#!/usr/bin/env node
/* eslint-disable no-console */
const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");

const ROOT = path.resolve(__dirname, "..");
const BASELINE_PATH = path.join(ROOT, "tests", "benchmarks-baseline.json");

const RUNS = [
  {
    name: "index-full",
    args: ["--input", "src/index.html", "-o", "dist/bench-index-full.css", "--preset", "full"],
    output: "dist/bench-index-full.css",
  },
  {
    name: "index-core",
    args: ["--input", "src/index.html", "-o", "dist/bench-index-core.css", "--preset", "core"],
    output: "dist/bench-index-core.css",
  },
  {
    name: "index-minimal",
    args: ["--input", "src/index.html", "-o", "dist/bench-index-min.css", "--preset", "minimal"],
    output: "dist/bench-index-min.css",
  },
  {
    name: "phase1-full",
    args: ["--input", "src/phase1-demo.html", "-o", "dist/bench-phase1-full.css", "--preset", "full"],
    output: "dist/bench-phase1-full.css",
  },
  {
    name: "phase1-core",
    args: ["--input", "src/phase1-demo.html", "-o", "dist/bench-phase1-core.css", "--preset", "core"],
    output: "dist/bench-phase1-core.css",
  },
  {
    name: "phase1-minimal",
    args: ["--input", "src/phase1-demo.html", "-o", "dist/bench-phase1-min.css", "--preset", "minimal"],
    output: "dist/bench-phase1-min.css",
  },
];

function runCase(item) {
  const started = process.hrtime.bigint();
  const result = spawnSync("go", ["run", "./compiler", ...item.args], {
    cwd: ROOT,
    encoding: "utf8",
  });
  const elapsedMs = Number(process.hrtime.bigint() - started) / 1_000_000;

  if (result.status !== 0) {
    throw new Error(`benchmark failed (${item.name}):\n${result.stdout}\n${result.stderr}`);
  }

  const outputPath = path.join(ROOT, item.output);
  const size = fs.statSync(outputPath).size;
  return {
    name: item.name,
    elapsedMs,
    sizeBytes: size,
    sizeKB: size / 1024,
  };
}

function pad(value, length) {
  return String(value).padEnd(length, " ");
}

function printTable(rows) {
  console.log("ZINTENT Build Benchmark");
  console.log(pad("case", 18), pad("time(ms)", 12), pad("size(kb)", 10), "size(bytes)");
  console.log("-".repeat(56));

  for (const r of rows) {
    console.log(
      pad(r.name, 18),
      pad(r.elapsedMs.toFixed(2), 12),
      pad(r.sizeKB.toFixed(2), 10),
      r.sizeBytes
    );
  }

  const indexFull = rows.find((r) => r.name === "index-full");
  const indexMin = rows.find((r) => r.name === "index-minimal");
  if (indexFull && indexMin) {
    const saved = indexFull.sizeBytes - indexMin.sizeBytes;
    const pct = indexFull.sizeBytes > 0 ? (saved / indexFull.sizeBytes) * 100 : 0;
    console.log(`\nindex minimal saved vs full: ${saved} bytes (${pct.toFixed(2)}%)`);
  }

  const phaseFull = rows.find((r) => r.name === "phase1-full");
  const phaseMin = rows.find((r) => r.name === "phase1-minimal");
  if (phaseFull && phaseMin) {
    const saved = phaseFull.sizeBytes - phaseMin.sizeBytes;
    const pct = phaseFull.sizeBytes > 0 ? (saved / phaseFull.sizeBytes) * 100 : 0;
    console.log(`phase1 minimal saved vs full: ${saved} bytes (${pct.toFixed(2)}%)`);
  }
}

function main() {
  if (process.argv.includes("--update-baseline")) {
    const rows = RUNS.map(runCase);
    const runs = {};
    for (const r of rows) {
      runs[r.name] = { sizeBytes: r.sizeBytes };
    }
    const payload = {
      version: 1,
      generatedAt: new Date().toISOString(),
      runs,
    };
    fs.mkdirSync(path.dirname(BASELINE_PATH), { recursive: true });
    fs.writeFileSync(BASELINE_PATH, `${JSON.stringify(payload, null, 2)}\n`, "utf8");
    console.log(`benchmark baseline saved: ${path.relative(ROOT, BASELINE_PATH)}`);
    printTable(rows);
    return;
  }

  if (process.argv.includes("--compare")) {
    if (!fs.existsSync(BASELINE_PATH)) {
      console.error(`missing baseline file: ${BASELINE_PATH}`);
      console.error("run: node compiler/benchmark.js --update-baseline");
      process.exit(1);
    }
    const baseline = JSON.parse(fs.readFileSync(BASELINE_PATH, "utf8"));
    const rows = RUNS.map(runCase);
    let failed = false;
    for (const r of rows) {
      const exp = baseline.runs && baseline.runs[r.name];
      if (!exp) {
        console.error(`baseline missing case: ${r.name}`);
        failed = true;
        continue;
      }
      if (exp.sizeBytes !== r.sizeBytes) {
        console.error(
          `size drift: ${r.name} expected ${exp.sizeBytes} bytes, got ${r.sizeBytes} bytes`
        );
        failed = true;
      } else {
        console.log(`ok: ${r.name} (${r.sizeBytes} bytes)`);
      }
    }
    if (failed) {
      process.exit(1);
    }
    console.log("benchmark compare passed");
    return;
  }

  const rows = RUNS.map(runCase);
  printTable(rows);
}

main();
