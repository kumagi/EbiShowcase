#!/usr/bin/env node
/** Ensure games reuse audiolab.Context instead of crashing on retry. */
import { execFileSync } from "node:child_process";

const root = new URL("..", import.meta.url).pathname;
const allowed = "internal/audiolab/context.go";
let output = "";
try {
  output = execFileSync(
    "rg",
    ["-n", "audio\\.NewContext", "--glob", "*.go", "."],
    { cwd: root, encoding: "utf8" },
  );
} catch (error) {
  if (error.status !== 1) throw error;
}

const violations = output
  .trim()
  .split("\n")
  .filter(Boolean)
  .filter((line) => !line.startsWith(`./${allowed}:`));

if (violations.length > 0) {
  console.error("Audio context gate failed: Ebitengine permits only one audio.Context per process.");
  console.error("Use audiolab.Context() so retry/new-run construction cannot panic.");
  for (const violation of violations) console.error(`- ${violation}`);
  process.exit(1);
}

console.log("Audio context gate passed: all games reuse audiolab.Context().");
