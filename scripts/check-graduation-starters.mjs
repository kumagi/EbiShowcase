#!/usr/bin/env node
import { spawnSync } from "node:child_process";
const starters = ["arcade-60", "exploration-3rooms", "puzzle-3stages"];
for (const name of starters) {
  const result = spawnSync("go", ["test", `./graduation/${name}/starter`], { encoding: "utf8" });
  if (result.status === 0) throw new Error(`${name} starter unexpectedly passed: red TODO tests are required.`);
  process.stdout.write(`Expected red starter: ${name}\n`);
  const reference = spawnSync("go", ["test", `./graduation/${name}/reference`], { encoding: "utf8" });
  if (reference.status !== 0) throw new Error(`${name} reference failed:\n${reference.stdout}${reference.stderr}`);
  process.stdout.write(`Green reference: ${name}\n`);
}
console.log("Graduation starter contract passed.");
