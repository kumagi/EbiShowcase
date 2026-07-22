#!/usr/bin/env node
// SPDX-License-Identifier: Apache-2.0

import assert from "node:assert/strict";
import {
  createFightingState,
  createPlatformerState,
  createTowerDefenseState,
  exportFightingDocument,
  exportPlatformerDocument,
  exportTowerDefenseDocument,
  importFightingDocument,
  importPlatformerDocument,
  importTowerDefenseDocument,
  validateFightingDocument,
  validatePlatformerDocument,
  validateTowerDefenseDocument,
} from "../web/game-data-editor.mjs";

const platformer = exportPlatformerDocument(createPlatformerState());
assert.deepEqual(validatePlatformerDocument(platformer), []);
assert.deepEqual(exportPlatformerDocument(importPlatformerDocument(platformer)), platformer);
assert.deepEqual(validatePlatformerDocument({ ...platformer, entities: platformer.entities.filter((item) => item.kind !== "goal") }), ["platformerGoal"]);

const fighting = exportFightingDocument(createFightingState());
assert.deepEqual(validateFightingDocument(fighting), []);
assert.deepEqual(exportFightingDocument(importFightingDocument(fighting)), fighting);
const impossibleMove = structuredClone(fighting);
impossibleMove.moves[0].startup = 50;
assert.deepEqual(validateFightingDocument(impossibleMove), ["fightingFrames"]);

const defense = exportTowerDefenseDocument(createTowerDefenseState());
assert.deepEqual(validateTowerDefenseDocument(defense), []);
assert.deepEqual(exportTowerDefenseDocument(importTowerDefenseDocument(defense)), defense);
const brokenRoute = structuredClone(defense);
brokenRoute.route.splice(4, 1);
assert.deepEqual(validateTowerDefenseDocument(brokenRoute), ["towerRoute"]);

console.log("OK — game-data editor documents validate and round-trip.");
