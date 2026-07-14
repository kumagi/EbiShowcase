import fs from "node:fs";
import path from "node:path";
const root = path.resolve(import.meta.dirname, "..");
const week = Number(process.argv[2] || 0);
const tracks = fs.readdirSync(path.join(root,"games/tracks"),{withFileTypes:true}).filter(x=>x.isDirectory()).map(x=>x.name).sort();
const start = (week % Math.ceil(tracks.length/5))*5;
console.log(JSON.stringify({ week, audit: tracks.slice(start,start+5), freshness: "run roadmap verify and inspect evidence dates" },null,2));
