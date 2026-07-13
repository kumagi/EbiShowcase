import { readFileSync, readdirSync, existsSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const home = readFileSync(join(root, "web/ja/index.html"), "utf8");

// The Visual Effects Lab is a distinct group that lives between the 12 core
// lessons and the genre tracks. It is counted separately from the main gate.
const vfxTrack = "visual-effects";

const core = [...home.matchAll(/href="games\/([^/]+)\/"/g)].map((m) => ({
  id: `core/${m[1]}`,
  slug: m[1],
  route: `games/${m[1]}`,
  group: "core",
}));

function trackSteps(track, group) {
  const hub = readFileSync(join(root, `web/ja/tracks/${track}/index.html`), "utf8");
  return [...hub.matchAll(/class="path-step" href="([^/]+)\/"/g)].map((m) => ({
    id: `tracks/${track}/${m[1]}`,
    slug: m[1],
    route: `tracks/${track}/${m[1]}`,
    group,
  }));
}

const trackIDs = [...home.matchAll(/href="tracks\/([^/]+)\/"/g)].map((m) => m[1]);
const vfx = trackIDs.includes(vfxTrack) ? trackSteps(vfxTrack, "vfx") : [];
const tracks = trackIDs
  .filter((track) => track !== vfxTrack)
  .flatMap((track) => trackSteps(track, "track"));

export const curriculum = [...core, ...vfx, ...tracks].map((item, index) => {
  const source = item.id === "core/flappy"
    ? join(root, "game/main.go")
    : join(root, "games", item.id, "main.go");
  return { ...item, order: index + 1, source, playable: existsSync(source) };
});

// Entries that count toward the completion gate (core + genre tracks).
export const gated = curriculum.filter((item) => item.group !== "vfx");

if (process.argv[1] === new URL(import.meta.url).pathname) {
  const command = process.argv[2] || "list";
  if (command === "next") {
    // The main gate ignores the separately counted Visual Effects Lab.
    const next = gated.find((item) => !item.playable);
    if (next) console.log(JSON.stringify(next));
    else console.log(JSON.stringify({ complete: true, total: gated.length }));
  } else if (command === "summary") {
    const playable = gated.filter((item) => item.playable).length;
    const vfxItems = curriculum.filter((item) => item.group === "vfx");
    const vfxPlayable = vfxItems.filter((item) => item.playable).length;
    console.log(JSON.stringify({
      total: gated.length,
      playable,
      remaining: gated.length - playable,
      vfx: { total: vfxItems.length, playable: vfxPlayable },
    }));
  } else {
    for (const item of curriculum) console.log(`${item.playable ? "PLAYABLE" : "PLANNED"}\t[${item.group}]\t${item.id}`);
  }
}
