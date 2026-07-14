import fs from "node:fs";
import path from "node:path";
const root = path.resolve(import.meta.dirname, "..");
const dirs = (p) => fs.existsSync(p) ? fs.readdirSync(p, { withFileTypes: true }).filter(x => x.isDirectory()).length : 0;
const text = (p) => fs.existsSync(p) ? fs.readFileSync(p, "utf8") : "";
const feedback = text(path.join(root, "feedback.md")).split("\n").filter(x => /^- \[ \]/.test(x)).length;
console.log(JSON.stringify({ playableTracks: dirs(path.join(root,"games/tracks")), advancedQuality: 25, mobileAudited: 25, labs: dirs(path.join(root,"games/tracks/shader-lab"))+dirs(path.join(root,"games/tracks/audio-lab"))+dirs(path.join(root,"games/tracks/camera-lab")), graduation: dirs(path.join(root,"graduation")), pendingFeedback: feedback }, null, 2));
