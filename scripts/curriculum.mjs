import { readFileSync, readdirSync, existsSync } from "node:fs";
import { join } from "node:path";

const root = new URL("..", import.meta.url).pathname;
const home = readFileSync(join(root, "web/ja/index.html"), "utf8");
const core = [...home.matchAll(/href="games\/([^/]+)\/"/g)].map((m) => ({
  id: `core/${m[1]}`,
  slug: m[1],
  route: `games/${m[1]}`,
}));
const trackIDs = [...home.matchAll(/href="tracks\/([^/]+)\/"/g)].map((m) => m[1]);
const tracks = trackIDs.flatMap((track) => {
  const hub = readFileSync(join(root, `web/ja/tracks/${track}/index.html`), "utf8");
  return [...hub.matchAll(/class="path-step" href="([^/]+)\/"/g)].map((m) => ({
    id: `tracks/${track}/${m[1]}`,
    slug: m[1],
    route: `tracks/${track}/${m[1]}`,
  }));
});

export const curriculum = [...core, ...tracks].map((item, index) => {
  const source = item.id === "core/flappy"
    ? join(root, "game/main.go")
    : join(root, "games", item.id, "main.go");
  return { ...item, order: index + 1, source, playable: existsSync(source) };
});

if (process.argv[1] === new URL(import.meta.url).pathname) {
  const command = process.argv[2] || "list";
  if (command === "next") {
    const next = curriculum.find((item) => !item.playable);
    if (next) console.log(JSON.stringify(next));
    else console.log(JSON.stringify({ complete: true, total: curriculum.length }));
  } else if (command === "summary") {
    const playable = curriculum.filter((item) => item.playable).length;
    console.log(JSON.stringify({ total: curriculum.length, playable, remaining: curriculum.length - playable }));
  } else {
    for (const item of curriculum) console.log(`${item.playable ? "PLAYABLE" : "PLANNED"}\t${item.id}`);
  }
}
