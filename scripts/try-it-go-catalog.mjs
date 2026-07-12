/**
 * Shared Go-snippet catalog for TRY IT × Go highlighting.
 * Keyed by data-lab kind. Each entry: { lines: [{ id, code }], buttons: { attr: { ids, clear?, caption? } } }
 * caption may be { ja, en } or omitted (wire script fills generic).
 */
export const catalog = {
  loop: {
    lines: [
      { id: "update", code: "func (g *game) Update() error { /* numbers */ }" },
      { id: "draw", code: "func (g *game) Draw(screen *ebiten.Image) { /* paint */ }" },
      { id: "run", code: "ebiten.RunGame(g) // Update → Draw → repeat" },
    ],
    buttons: {
      "data-lab-step": { ids: ["update", "draw"], cycle: true },
      "data-lab-reset": { clear: true },
    },
  },
  "hit-test": {
    lines: [
      { id: "dxdy", code: "dx := touchX - circleX; dy := touchY - circleY" },
      { id: "hypot", code: "dist := math.Hypot(dx, dy)" },
      { id: "hit", code: "if dist <= radius { score++; moveTarget() }" },
      { id: "miss", code: "// dist > radius → no score" },
    ],
    buttons: {
      "data-lab-sample-hit": { ids: ["dxdy", "hypot", "hit"] },
      "data-lab-sample-miss": { ids: ["dxdy", "hypot", "miss"] },
      "data-lab-reset": { clear: true },
    },
    board: { ids: ["dxdy", "hypot"], thenHit: ["hit"], thenMiss: ["miss"] },
  },
  meter: {
    lines: [
      { id: "move", code: "x += speed; if x > max || x < min { speed = -speed }" },
      { id: "stop", code: "if justPressed { stopped = true }" },
      { id: "score", code: "if abs(x-center) < zone { score++ }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["move"] },
      "data-lab-stop": { ids: ["stop", "score"] },
      "data-lab-reset": { clear: true },
    },
  },
  aabb: {
    lines: [
      { id: "fall", code: "star.y += star.speed" },
      { id: "move", code: "basket.x = clamp(pointerX)" },
      { id: "hit", code: "if overlaps(starBox, basketBox) { score++ }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["fall"] },
      "data-lab-catch": { ids: ["hit"] },
      "data-lab-reset": { clear: true },
    },
  },
  flappy: {
    lines: [
      { id: "grav", code: "vy += 0.42 // gravity each frame" },
      { id: "pos", code: "y += vy" },
      { id: "flap", code: "if pressed { vy = -7.4 // jump impulse }" },
    ],
    buttons: {
      "data-lab-flap": { ids: ["flap", "grav", "pos"] },
      "data-lab-step": { ids: ["grav", "pos"] },
      "data-lab-reset": { clear: true },
    },
  },
  bounce: {
    lines: [
      { id: "move", code: "x += vx; y += vy" },
      { id: "wall", code: "if x < 0 || x > w { vx = -vx }" },
      { id: "paddle", code: "if hitPaddle { vy = -vy }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["move"] },
      "data-lab-bounce": { ids: ["wall"] },
      "data-lab-reset": { clear: true },
    },
  },
  bricks: {
    lines: [
      { id: "grid", code: "bricks[row][col] = alive" },
      { id: "hit", code: "if hitBrick { bricks[r][c] = false; score++ }" },
      { id: "clear", code: "if allClear() { win() }" },
    ],
    buttons: {
      "data-lab-hit": { ids: ["hit"] },
      "data-lab-step": { ids: ["grid"] },
      "data-lab-reset": { clear: true },
    },
  },
  snake: {
    lines: [
      { id: "head", code: "head := body[0]; next := Point{head.X+dx, head.Y+dy}" },
      { id: "grow", code: "body = append([]Point{next}, body...)" },
      { id: "trim", code: "if !ate { body = body[:len(body)-1] }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["head", "grow", "trim"] },
      "data-lab-eat": { ids: ["head", "grow"] },
      "data-lab-reset": { clear: true },
    },
  },
  bullets: {
    lines: [
      { id: "fire", code: "shots = append(shots, Shot{x: px, y: py})" },
      { id: "fly", code: "for i := range shots { shots[i].y -= 8 }" },
      { id: "cull", code: "shots = filterOnScreen(shots)" },
    ],
    buttons: {
      "data-lab-fire": { ids: ["fire"] },
      "data-lab-step": { ids: ["fly", "cull"] },
      "data-lab-reset": { clear: true },
    },
  },
  push: {
    lines: [
      { id: "try", code: "next := ahead(player, dir)" },
      { id: "box", code: "if next == Box && ahead(box, dir) == Empty { push }" },
      { id: "block", code: "if next == Wall { /* stay */ }" },
    ],
    buttons: {
      "data-lab-right": { ids: ["try", "box"] },
      "data-lab-reset": { clear: true },
    },
  },
  camera: {
    lines: [
      { id: "target", code: "target := playerX - screenW/2" },
      { id: "lerp", code: "camX += (target - camX) * 0.15" },
      { id: "draw", code: "drawX := worldX - camX" },
    ],
    buttons: {
      "data-lab-right": { ids: ["target"] },
      "data-lab-left": { ids: ["target"] },
      "data-lab-step": { ids: ["lerp", "draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  ai: {
    lines: [
      { id: "dist", code: "d := math.Hypot(ex-px, ey-py)" },
      { id: "chase", code: "if d < limit { mode = chase } else { mode = patrol }" },
      { id: "step", code: "moveToward(target)" },
    ],
    buttons: {
      "data-lab-closer": { ids: ["dist", "chase"] },
      "data-lab-farther": { ids: ["dist", "chase"] },
      "data-lab-reset": { clear: true },
    },
  },
  burst: {
    lines: [
      { id: "spawn", code: "for i := 0; i < n; i++ { angle := tau * i / n }" },
      { id: "vel", code: "vx, vy := cos(a)*speed, sin(a)*speed" },
      { id: "count", code: "n = clamp(n, 4, 24)" },
    ],
    buttons: {
      "data-lab-more": { ids: ["count", "spawn"] },
      "data-lab-less": { ids: ["count", "spawn"] },
      "data-lab-reset": { clear: true },
    },
  },
  translate: {
    lines: [
      { id: "geom", code: "var op ebiten.DrawImageOptions" },
      { id: "move", code: "op.GeoM.Translate(x, y)" },
      { id: "draw", code: "screen.DrawImage(img, &op)" },
    ],
    buttons: { "data-lab-reset": { clear: true } },
    board: { ids: ["move", "draw"] },
  },
  geom: {
    lines: [
      { id: "pivot", code: "op.GeoM.Translate(-cx, -cy)" },
      { id: "rot", code: "op.GeoM.Rotate(angle)" },
      { id: "scale", code: "op.GeoM.Scale(sx, sy)" },
      { id: "place", code: "op.GeoM.Translate(x, y)" },
    ],
    buttons: {
      "data-lab-rot": { ids: ["pivot", "rot", "place"] },
      "data-lab-scale": { ids: ["pivot", "scale", "place"] },
      "data-lab-reset": { clear: true },
    },
  },
  colorscale: {
    lines: [
      { id: "cs", code: "var op ebiten.DrawImageOptions" },
      { id: "tint", code: "op.ColorScale.Scale(r, g, b, 1)" },
      { id: "draw", code: "screen.DrawImage(img, &op)" },
    ],
    buttons: {
      "data-lab-tint": { ids: ["tint", "draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  opacity: {
    lines: [
      { id: "alpha", code: "op.ColorScale.ScaleAlpha(a) // 0..1" },
      { id: "ghost", code: "for _, p := range trail { drawFaded(p, a) }" },
    ],
    buttons: {
      "data-lab-alpha": { ids: ["alpha"] },
      "data-lab-reset": { clear: true },
    },
  },
  blend: {
    lines: [
      { id: "blend", code: "op.Blend = ebiten.BlendLighter" },
      { id: "draw", code: "screen.DrawImage(glow, &op)" },
    ],
    buttons: {
      "data-lab-blend": { ids: ["blend", "draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  sheet: {
    lines: [
      { id: "frame", code: "frame := frames[i % len(frames)]" },
      { id: "sub", code: "// SubImage cuts one cell from the atlas" },
      { id: "draw", code: "screen.DrawImage(frame, &op)" },
    ],
    buttons: {
      "data-lab-step": { ids: ["frame", "draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  spray: {
    lines: [
      { id: "spawn", code: "p := Particle{x, y, vx, vy, life}" },
      { id: "update", code: "p.x += p.vx; p.life--" },
      { id: "draw", code: "for _, p := range parts { draw(p) }" },
    ],
    buttons: {
      "data-lab-spray": { ids: ["spawn"] },
      "data-lab-step": { ids: ["update", "draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  spellbook: {
    lines: [
      { id: "pick", code: "spell := spells[selected]" },
      { id: "cast", code: "fx = append(fx, spell.Spawn(x, y)...)" },
      { id: "tick", code: "updateFX(fx)" },
    ],
    buttons: {
      "data-lab-cast": { ids: ["pick", "cast"] },
      "data-lab-reset": { clear: true },
    },
  },
  "magic-fire": {
    lines: [
      { id: "core", code: "drawFlame(x, y, t) // warm ColorScale" },
      { id: "add", code: "op.Blend = ebiten.BlendLighter" },
    ],
    buttons: { "data-lab-cast": { ids: ["core", "add"] }, "data-lab-reset": { clear: true } },
  },
  "magic-ice": {
    lines: [
      { id: "core", code: "drawCrystal(x, y, t) // cool tint" },
      { id: "fade", code: "op.ColorScale.ScaleAlpha(a)" },
    ],
    buttons: { "data-lab-cast": { ids: ["core", "fade"] }, "data-lab-reset": { clear: true } },
  },
  "magic-thunder": {
    lines: [
      { id: "bolt", code: "drawBolt(path, flash)" },
      { id: "flash", code: "if frame%3==0 { BlendLighter }" },
    ],
    buttons: { "data-lab-cast": { ids: ["bolt", "flash"] }, "data-lab-reset": { clear: true } },
  },
  "magic-light": {
    lines: [
      { id: "ring", code: "drawRing(x, y, r)" },
      { id: "glow", code: "op.Blend = ebiten.BlendLighter" },
    ],
    buttons: { "data-lab-cast": { ids: ["ring", "glow"] }, "data-lab-reset": { clear: true } },
  },
  "magic-dark": {
    lines: [
      { id: "veil", code: "drawVeil(x, y, a)" },
      { id: "tint", code: "op.ColorScale.Scale(0.4, 0.2, 0.6, a)" },
    ],
    buttons: { "data-lab-cast": { ids: ["veil", "tint"] }, "data-lab-reset": { clear: true } },
  },
  "fx-split": {
    lines: [
      { id: "logic", code: "// Update: change game numbers only" },
      { id: "fx", code: "fx = append(fx, Burst{x, y, life})" },
      { id: "draw", code: "// Draw: paint world, then paint fx" },
    ],
    buttons: {
      "data-lab-fx-ping": { ids: ["fx", "draw"] },
      "data-lab-fx-tick": { ids: ["logic", "fx"] },
      "data-lab-logic": { ids: ["logic"] },
      "data-lab-fx": { ids: ["fx", "draw"] },
      "data-lab-step": { ids: ["logic", "fx"] },
      "data-lab-reset": { clear: true },
    },
  },
  "fx-breakout": {
    lines: [
      { id: "hit", code: "if ballHitsBrick { brick.alive = false }" },
      { id: "fx", code: "fx = append(fx, Shatter{at: brick})" },
      { id: "draw", code: "drawBricks(); drawFX()" },
    ],
    buttons: {
      "data-lab-hit": { ids: ["hit", "fx"] },
      "data-lab-step": { ids: ["draw"] },
      "data-lab-reset": { clear: true },
    },
  },
  "fx-snake": {
    lines: [
      { id: "eat", code: "if head == food { grow(); spawnFood() }" },
      { id: "fx", code: "fx = append(fx, Munch{at: head})" },
    ],
    buttons: {
      "data-lab-eat": { ids: ["eat", "fx"] },
      "data-lab-step": { ids: ["eat"] },
      "data-lab-reset": { clear: true },
    },
  },
  jump: {
    lines: [
      { id: "grav", code: "vy += gravity" },
      { id: "jump", code: "if grounded && pressed { vy = jumpV }" },
      { id: "land", code: "if y >= floor { y, vy, grounded = floor, 0, true }" },
    ],
    buttons: {
      "data-lab-jump": { ids: ["jump"] },
      "data-lab-step": { ids: ["grav", "land"] },
      "data-lab-reset": { clear: true },
    },
  },
  carry: {
    lines: [
      { id: "delta", code: "d := plat.x - plat.prevX" },
      { id: "ride", code: "if riding { player.x += d }" },
      { id: "step", code: "plat.prevX = plat.x; plat.x += vx" },
    ],
    buttons: {
      "data-lab-step": { ids: ["step", "delta", "ride"] },
      "data-lab-toggle": { ids: ["ride"] },
      "data-lab-reset": { clear: true },
    },
  },
  stomp: {
    lines: [
      { id: "check", code: "if vy > 0 && fromAbove(player, enemy)" },
      { id: "stomp", code: "{ enemy.dead = true; vy = bounce }" },
      { id: "hurt", code: "else { player.hurt() }" },
    ],
    buttons: {
      "data-lab-fall": { ids: ["check", "stomp"] },
      "data-lab-rise": { ids: ["check", "hurt"] },
      "data-lab-reset": { clear: true },
    },
  },
  power: {
    lines: [
      { id: "flag", code: "powered := g.powered" },
      { id: "branch", code: "if powered { stomp } else { takeDamage }" },
    ],
    buttons: {
      "data-lab-toggle": { ids: ["flag", "branch"] },
      "data-lab-reset": { clear: true },
    },
  },
  move8: {
    lines: [
      { id: "raw", code: "dx, dy := inputX, inputY" },
      { id: "norm", code: "if dx!=0 && dy!=0 { dx, dy = dx*0.707, dy*0.707 }" },
      { id: "apply", code: "x += dx * speed; y += dy * speed" },
    ],
    buttons: {
      "data-lab-cardinal": { ids: ["raw", "apply"] },
      "data-lab-diag": { ids: ["raw", "norm", "apply"] },
      "data-lab-reset": { clear: true },
    },
  },
  aim: {
    lines: [
      { id: "scan", code: "best, dist := nearestEnemy(enemies, px, py)" },
      { id: "cd", code: "if frame%cooldown == 0 { fire(best) }" },
      { id: "shot", code: "shots = append(shots, toward(best))" },
    ],
    buttons: {
      "data-lab-step": { ids: ["scan", "cd", "shot"] },
      "data-lab-reset": { clear: true },
    },
  },
  pool: {
    lines: [
      { id: "slot", code: "mobs [100]mob // fixed seats" },
      { id: "kill", code: "mobs[i] = spawn() // overwrite, do not delete" },
    ],
    buttons: {
      "data-lab-kill": { ids: ["kill"] },
      "data-lab-reset": { clear: true },
    },
  },
  draft: {
    lines: [
      { id: "pause", code: "if xp >= need { drafting = true }" },
      { id: "pick", code: "applyCard(choice); drafting = false" },
    ],
    buttons: {
      "data-lab-level": { ids: ["pause"] },
      "data-lab-card": { ids: ["pick"] },
      "data-lab-reset": { clear: true },
    },
  },
  evolve: {
    lines: [
      { id: "rule", code: "if weapon == needle && copies >= 3" },
      { id: "evo", code: "{ weapon = storm; cooldown /= 2 }" },
    ],
    buttons: {
      "data-lab-evolve": { ids: ["rule", "evo"] },
      "data-lab-reset": { clear: true },
    },
  },
  curve: {
    lines: [
      { id: "sec", code: "sec := frame / 60" },
      { id: "iv", code: "interval := max(10, 34-sec/2)" },
      { id: "spd", code: "speed := 1.15 + float64(sec)*0.028" },
    ],
    buttons: {
      "data-lab-step": { ids: ["sec", "iv", "spd"] },
      "data-lab-reset": { clear: true },
    },
  },
  click: {
    lines: [
      { id: "tap", code: "if justPressed { count++ }" },
      { id: "draw", code: "drawNumber(count)" },
    ],
    buttons: {
      "data-lab-tap": { ids: ["tap"] },
      "data-lab-reset": { clear: true },
    },
  },
  shop: {
    lines: [
      { id: "earn", code: "gold += 5" },
      { id: "buy", code: "if gold >= cost { gold -= cost; owned++; cost += 5 }" },
    ],
    buttons: {
      "data-lab-earn": { ids: ["earn"] },
      "data-lab-buy": { ids: ["buy"] },
      "data-lab-reset": { clear: true },
    },
  },
  idle: {
    lines: [
      { id: "rate", code: "rate := float64(machines)" },
      { id: "tick", code: "points += rate * dt" },
    ],
    buttons: {
      "data-lab-step": { ids: ["tick"] },
      "data-lab-buy": { ids: ["rate"] },
      "data-lab-reset": { clear: true },
    },
  },
  cost: {
    lines: [
      { id: "curve", code: "cost = base * pow(growth, owned)" },
      { id: "buy", code: "if gold >= cost { purchase() }" },
    ],
    buttons: {
      "data-lab-buy": { ids: ["curve", "buy"] },
      "data-lab-step": { ids: ["curve"] },
      "data-lab-reset": { clear: true },
    },
  },
  save: {
    lines: [
      { id: "away", code: "away := now.Sub(savedAt).Seconds()" },
      { id: "gain", code: "gold += away * rate" },
    ],
    buttons: {
      "data-lab-add": { ids: ["away", "gain"] },
      "data-lab-away": { ids: ["away", "gain"] },
      "data-lab-reset": { clear: true },
    },
  },
  tile: {
    lines: [
      { id: "try", code: "nx, ny := x+dx, y+dy" },
      { id: "wall", code: "if walls[ny][nx] { return /* blocked */ }" },
      { id: "go", code: "x, y = nx, ny" },
    ],
    buttons: {
      "data-lab-up": { ids: ["try", "go"] },
      "data-lab-down": { ids: ["try", "go"] },
      "data-lab-left": { ids: ["try", "go"] },
      "data-lab-right": { ids: ["try", "go"] },
      "data-lab-reset": { clear: true },
    },
  },
  flag: {
    lines: [
      { id: "set", code: "flags[\"found\"] = true" },
      { id: "say", code: "if flags[\"found\"] { text = after } else { text = before }" },
    ],
    buttons: {
      "data-lab-toggle": { ids: ["set", "say"] },
      "data-lab-reset": { clear: true },
    },
  },
  turn: {
    lines: [
      { id: "state", code: "switch phase {" },
      { id: "next", code: "phase = nextPhase(phase)" },
      { id: "act", code: "case player: playCard(); case enemy: enemyAct()" },
    ],
    buttons: {
      "data-lab-step": { ids: ["next", "act"] },
      "data-lab-reset": { clear: true },
    },
  },
  damage: {
    lines: [
      { id: "raw", code: "dmg := atk + buff - def" },
      { id: "min", code: "if dmg < 1 { dmg = 1 }" },
    ],
    buttons: {
      "data-lab-buff": { ids: ["raw", "min"] },
      "data-lab-reset": { clear: true },
    },
  },
  inv: {
    lines: [
      { id: "pay", code: "if gold >= price { gold -= price }" },
      { id: "equip", code: "slot = item" },
    ],
    buttons: {
      "data-lab-buy": { ids: ["pay", "equip"] },
      "data-lab-reset": { clear: true },
    },
  },
  scene: {
    lines: [
      { id: "table", code: "enemy := tables[region]" },
      { id: "swap", code: "scene = battle" },
      { id: "back", code: "scene = field" },
    ],
    buttons: {
      "data-lab-walk": { ids: ["table", "swap"] },
      "data-lab-region": { ids: ["table"] },
      "data-lab-back": { ids: ["back"] },
      "data-lab-reset": { clear: true },
    },
  },
  quest: {
    lines: [
      { id: "adv", code: "quest++" },
      { id: "save", code: "save.quest = quest" },
      { id: "load", code: "quest = save.quest" },
    ],
    buttons: {
      "data-lab-advance": { ids: ["adv"] },
      "data-lab-save": { ids: ["save"] },
      "data-lab-load": { ids: ["load"] },
      "data-lab-reset": { clear: true },
    },
  },
  hitbox: {
    lines: [
      { id: "atk", code: "atk := Rect{x, y, w, h} // active only" },
      { id: "hurt", code: "hurt := Rect{ox, oy, ow, oh}" },
      { id: "hit", code: "if overlaps(atk, hurt) { applyHit() }" },
    ],
    buttons: {
      "data-lab-attack": { ids: ["atk", "hit"] },
      "data-lab-step": { ids: ["atk", "hurt", "hit"] },
      "data-lab-reset": { clear: true },
    },
  },
  frames: {
    lines: [
      { id: "start", code: "if pressed { frame = 1 }" },
      { id: "active", code: "if frame >= 9 && frame <= 12 { hitboxOn = true }" },
      { id: "rec", code: "if frame > 12 { /* recovery */ }" },
      { id: "tick", code: "frame++" },
    ],
    buttons: {
      "data-lab-step": { ids: ["tick", "active", "rec"] },
      "data-lab-reset": { clear: true },
    },
  },
  react: {
    lines: [
      { id: "stop", code: "hitstop = 8" },
      { id: "stun", code: "stun = 25; vx = knockback" },
      { id: "slide", code: "x += vx; vx *= 0.86" },
    ],
    buttons: {
      "data-lab-hit": { ids: ["stop", "stun"] },
      "data-lab-step": { ids: ["slide"] },
      "data-lab-reset": { clear: true },
    },
  },
  rps: {
    lines: [
      { id: "pick", code: "you, enemy := choice, foeChoice" },
      { id: "rule", code: "strike > throw > guard > strike" },
      { id: "res", code: "result := resolve(you, enemy)" },
    ],
    buttons: {
      "data-lab-pick": { ids: ["pick", "rule", "res"] },
      "data-lab-enemy": { ids: ["pick"] },
      "data-lab-reset": { clear: true },
    },
  },
  buffer: {
    lines: [
      { id: "press", code: "buffer, life = \"light\", 8" },
      { id: "tick", code: "if life > 0 { life-- } else { buffer = \"none\" }" },
      { id: "cancel", code: "if cancelWindow && buffer != \"none\" { start(buffer) }" },
    ],
    buttons: {
      "data-lab-press": { ids: ["press"] },
      "data-lab-step": { ids: ["tick", "cancel"] },
      "data-lab-reset": { clear: true },
    },
  },
  command: {
    lines: [
      { id: "hist", code: "hist = append(hist, dir); hist = hist[max(0,len-3):]" },
      { id: "match", code: "if hist == ↓↘→ { hadoken() }" },
    ],
    buttons: {
      "data-lab-dir": { ids: ["hist", "match"] },
      "data-lab-reset": { clear: true },
    },
  },
  rounds: {
    lines: [
      { id: "point", code: "if roundWin { p1++ }" },
      { id: "match", code: "if p1 >= 2 { matchOver = true }" },
    ],
    buttons: {
      "data-lab-p1hit": { ids: ["point", "match"] },
      "data-lab-p2hit": { ids: ["point", "match"] },
      "data-lab-reset": { clear: true },
    },
  },
  gravity: {
    lines: [
      { id: "v", code: "vy += g" },
      { id: "y", code: "y += vy" },
      { id: "floor", code: "if y > floor { y, vy = floor, 0 }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["v", "y", "floor"] },
      "data-lab-reset": { clear: true },
    },
  },
  "merge-same": {
    lines: [
      { id: "touch", code: "if sameTier(a, b) && overlapping(a, b)" },
      { id: "merge", code: "{ spawn(tier+1); remove(a, b) }" },
    ],
    buttons: {
      "data-lab-merge": { ids: ["touch", "merge"] },
      "data-lab-step": { ids: ["touch"] },
      "data-lab-reset": { clear: true },
    },
  },
  "preview-next": {
    lines: [
      { id: "q", code: "next := queue[0]" },
      { id: "drop", code: "queue = queue[1:]; spawn(next)" },
      { id: "enq", code: "queue = append(queue, randomTier())" },
    ],
    buttons: {
      "data-lab-drop": { ids: ["q", "drop"] },
      "data-lab-enqueue": { ids: ["enq"] },
      "data-lab-reset": { clear: true },
    },
  },
  "card-play": {
    lines: [
      { id: "cost", code: "if energy < card.cost { return }" },
      { id: "pay", code: "energy -= card.cost" },
      { id: "fx", code: "apply(card)" },
    ],
    buttons: {
      "data-lab-card": { ids: ["cost", "pay", "fx"] },
      "data-lab-damage": { ids: ["cost", "pay", "fx"] },
      "data-lab-block": { ids: ["cost", "pay", "fx"] },
      "data-lab-heal": { ids: ["cost", "pay", "fx"] },
      "data-lab-reset": { clear: true },
    },
  },
  "energy-turn": {
    lines: [
      { id: "refill", code: "energy = maxEnergy // start of turn" },
      { id: "end", code: "drawHand(); enemyTurn()" },
    ],
    buttons: {
      "data-lab-step": { ids: ["refill", "end"] },
      "data-lab-end": { ids: ["end"] },
      "data-lab-reset": { clear: true },
    },
  },
  "status-ticks": {
    lines: [
      { id: "add", code: "status = append(status, Effect{name, left: 3})" },
      { id: "tick", code: "for i := range status { status[i].left-- }" },
      { id: "drop", code: "status = filterAlive(status)" },
    ],
    buttons: {
      "data-lab-add": { ids: ["add"] },
      "data-lab-tick": { ids: ["tick", "drop"] },
      "data-lab-reset": { clear: true },
    },
  },
  "deck-pick": {
    lines: [
      { id: "offer", code: "choices := randomThree(cardPool)" },
      { id: "take", code: "deck = append(deck, chosen)" },
    ],
    buttons: {
      "data-lab-reset": { clear: true },
    },
    board: { ids: ["take"] },
  },
  "map-nodes": {
    lines: [
      { id: "edge", code: "next := edges[node]" },
      { id: "go", code: "node = choice; path = append(path, node)" },
    ],
    buttons: { "data-lab-reset": { clear: true } },
    board: { ids: ["go"] },
  },
  "stage-goals": {
    lines: [
      { id: "move", code: "moves--; score += cleared" },
      { id: "clear", code: "if score >= goal { stageClear = true }" },
    ],
    buttons: {
      "data-lab-match": { ids: ["move", "clear"] },
      "data-lab-reset": { clear: true },
    },
  },
  "match-scan": {
    lines: [
      { id: "scan", code: "for each cell { if run >= 3 { mark } }" },
      { id: "clear", code: "removeMarked()" },
      { id: "fall", code: "compactColumns()" },
    ],
    buttons: {
      "data-lab-scan": { ids: ["scan"] },
      "data-lab-clear": { ids: ["clear"] },
      "data-lab-fall": { ids: ["fall"] },
      "data-lab-step": { ids: ["scan", "clear", "fall"] },
      "data-lab-reset": { clear: true },
    },
  },
  "drop-timer": {
    lines: [
      { id: "tick", code: "timer++" },
      { id: "drop", code: "if timer >= 38 { y++; timer = 0 }" },
      { id: "lock", code: "if !canFall { lockPiece() }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["tick", "drop"] },
      "data-lab-reset": { clear: true },
    },
  },
  "shape-cells": {
    lines: [
      { id: "shape", code: "cells := shapes[kind] // relative offsets" },
      { id: "rot", code: "cells = rotate(cells)" },
    ],
    buttons: {
      "data-lab-next": { ids: ["shape"] },
      "data-lab-rot": { ids: ["rot"] },
      "data-lab-reset": { clear: true },
    },
  },
  "kick-try": {
    lines: [
      { id: "try", code: "for _, k := range kicks { if canPlace(x+k.dx, y+k.dy) }" },
      { id: "ok", code: "{ x, y = x+k.dx, y+k.dy; return }" },
      { id: "fail", code: "/* all kicks failed → keep old rot */" },
    ],
    buttons: {
      "data-lab-kick": { ids: ["try", "ok"] },
      "data-lab-reset": { clear: true },
    },
  },
  "line-clear": {
    lines: [
      { id: "full", code: "if rowAllFilled { cleared++ }" },
      { id: "pack", code: "write surviving rows from bottom" },
    ],
    buttons: {
      "data-lab-clear": { ids: ["full", "pack"] },
      "data-lab-reset": { clear: true },
    },
  },
  "bag-draw": {
    lines: [
      { id: "bag", code: "if len(bag)==0 { bag = shuffle7() }" },
      { id: "draw", code: "piece := popRandom(&bag)" },
      { id: "hold", code: "held, current = current, held" },
    ],
    buttons: {
      "data-lab-draw": { ids: ["bag", "draw"] },
      "data-lab-hold": { ids: ["hold"] },
      "data-lab-reset": { clear: true },
    },
  },
  pipeline: {
    lines: [
      { id: "order", code: "// fixed order each frame:" },
      { id: "step", code: "input → sim → resolve → draw" },
      { id: "loop", code: "// then next frame" },
    ],
    buttons: {
      "data-lab-next": { ids: ["order", "step"] },
      "data-lab-reset": { clear: true },
    },
  },
  "event-queue": {
    lines: [
      { id: "push", code: "q = append(q, Event{kind, at})" },
      { id: "pop", code: "e, q := q[0], q[1:]; resolve(e)" },
    ],
    buttons: {
      "data-lab-push": { ids: ["push"] },
      "data-lab-pop": { ids: ["pop"] },
      "data-lab-reset": { clear: true },
    },
  },
  "height-layers": {
    lines: [
      { id: "noise", code: "h = sampleNoise(x, z)" },
      { id: "carve", code: "h = max(1, h-1)" },
    ],
    buttons: {
      "data-lab-noise": { ids: ["noise"] },
      "data-lab-carve": { ids: ["carve"] },
      "data-lab-reset": { clear: true },
    },
  },
  "craft-recipe": {
    lines: [
      { id: "need", code: "if has(ingredients) { craft(recipe) }" },
      { id: "out", code: "inventory = append(inventory, result)" },
    ],
    buttons: {
      "data-lab-craft": { ids: ["need", "out"] },
      "data-lab-reset": { clear: true },
    },
  },
  "light-flood": {
    lines: [
      { id: "seed", code: "light[y][x] = 15 // torch" },
      { id: "spread", code: "n = max(n, c-1) // 4-neighborhood" },
    ],
    buttons: {
      "data-lab-step": { ids: ["spread"] },
      "data-lab-seed": { ids: ["seed"] },
      "data-lab-reset": { clear: true },
    },
  },
  "species-inst": {
    lines: [
      { id: "table", code: "speciesTable[id] = Species{hp, atk}" },
      { id: "inst", code: "c := Creature{speciesID: id, hp: table.hp}" },
    ],
    buttons: {
      "data-lab-spawn": { ids: ["inst"] },
      "data-lab-step": { ids: ["table", "inst"] },
      "data-lab-reset": { clear: true },
    },
  },
  "party-swap": {
    lines: [
      { id: "idx", code: "active = slot // do not delete party" },
      { id: "turn", code: "if fainted { forceSwap() }" },
    ],
    buttons: {
      "data-lab-swap": { ids: ["idx"] },
      "data-lab-reset": { clear: true },
    },
  },
  "capture-roll": {
    lines: [
      { id: "rate", code: "p := catchRate(hp, ball)" },
      { id: "roll", code: "if rng.Float64() < p { caught = true }" },
    ],
    buttons: {
      "data-lab-roll": { ids: ["rate", "roll"] },
      "data-lab-reset": { clear: true },
    },
  },
  "xp-level": {
    lines: [
      { id: "xp", code: "xp += gained" },
      { id: "lv", code: "for xp >= need { level++; xp -= need }" },
    ],
    buttons: {
      "data-lab-xp": { ids: ["xp", "lv"] },
      "data-lab-step": { ids: ["xp", "lv"] },
      "data-lab-reset": { clear: true },
    },
  },
  "bfs-flood": {
    lines: [
      { id: "seed", code: "q = []Pos{start}; seen[start]=true" },
      { id: "step", code: "for len(q)>0 { pop; push same-color neighbors }" },
    ],
    buttons: {
      "data-lab-step": { ids: ["step"] },
      "data-lab-flood": { ids: ["seed", "step"] },
      "data-lab-reset": { clear: true },
    },
  },
  "input-buffer": {
    lines: [
      { id: "queue", code: "queued = inputDir // do not write current yet" },
      { id: "center", code: "if atCenter && open(queued) { current = queued }" },
    ],
    buttons: {
      "data-lab-dir": { ids: ["queue"] },
      "data-lab-move": { ids: ["queue"] },
      "data-lab-center": { ids: ["center"] },
      "data-lab-reset": { clear: true },
    },
  },
  "pellet-count": {
    lines: [
      { id: "eat", code: "if tile == pellet { pellets--; score++ }" },
      { id: "win", code: "if pellets == 0 { clear = true }" },
    ],
    buttons: {
      "data-lab-eat": { ids: ["eat", "win"] },
      "data-lab-step": { ids: ["eat"] },
      "data-lab-reset": { clear: true },
    },
  },
  "junction-pick": {
    lines: [
      { id: "mode", code: "target := player // or corner if scatter" },
      { id: "pick", code: "dir = argmin(dist(nextTile, target))" },
    ],
    buttons: {
      "data-lab-chase": { ids: ["mode", "pick"] },
      "data-lab-scatter": { ids: ["mode", "pick"] },
      "data-lab-reset": { clear: true },
    },
  },
  "bomb-timer": {
    lines: [
      { id: "place", code: "bombs = append(bombs, Bomb{x, y, timer: fuse})" },
      { id: "tick", code: "b.timer--; if b.timer <= 0 { b.state = blast }" },
      { id: "gone", code: "remove(b)" },
    ],
    buttons: {
      "data-lab-place": { ids: ["place"] },
      "data-lab-tick": { ids: ["tick"] },
      "data-lab-reset": { clear: true },
    },
  },
  "cross-blast": {
    lines: [
      { id: "arms", code: "for dir in ±x,±y { for i:=1; i<=power; i++ }" },
      { id: "stop", code: "if wall { break } ; light(cell)" },
    ],
    buttons: {
      "data-lab-power": { ids: ["arms", "stop"] },
      "data-lab-reset": { clear: true },
    },
  },
  "chain-bomb": {
    lines: [
      { id: "ignite", code: "blast(i)" },
      { id: "chain", code: "if neighbor is bomb { blast(neighbor) }" },
    ],
    buttons: {
      "data-lab-ignite": { ids: ["ignite", "chain"] },
      "data-lab-clear": { ids: ["ignite"] },
      "data-lab-reset": { clear: true },
    },
  },
  "escape-timing": {
    lines: [
      { id: "eta", code: "eta := pathCost(start, safe)" },
      { id: "ok", code: "if eta+margin < fuse { takeRoute } else { reroute }" },
    ],
    buttons: {
      "data-lab-tick": { ids: ["eta", "ok"] },
      "data-lab-far": { ids: ["eta", "ok"] },
      "data-lab-reset": { clear: true },
    },
  },
};

/** Default captions when a line id lights up */
export const lineCaptions = {
  ja: {
    update: "Update で数字を変える",
    draw: "Draw で絵を描く",
    run: "RunGame がこの順を繰り返す",
    dxdy: "押した点と中心の差",
    hypot: "距離を求める",
    hit: "半径以下ならヒット",
    miss: "遠いのではずれる",
    move: "位置を進める",
    stop: "入力で止める",
    score: "中央判定で得点",
    fall: "毎フレーム落ちる",
    grav: "加速度を足す",
    pos: "速度を位置へ足す",
    flap: "上昇の初速を入れる",
    wall: "壁で速度を反転",
    fire: "弾をスライスへ追加",
    fly: "弾を動かす",
    cull: "画面外を消す",
    kill: "席を上書きリスポーン",
    overwrite: "席を上書きリスポーン",
    place: "爆弾を置く",
    tick: "タイマーを減らす",
    queue: "向きを予約する",
    center: "マス中心で確定する",
  },
  en: {
    update: "Update changes numbers",
    draw: "Draw paints the frame",
    run: "RunGame repeats this order",
    dxdy: "Pointer minus center",
    hypot: "Measure distance",
    hit: "Inside radius → hit",
    miss: "Too far → miss",
    move: "Advance position",
    stop: "Input freezes motion",
    score: "Score near the center",
    fall: "Fall each frame",
    grav: "Add acceleration",
    pos: "Add velocity to position",
    flap: "Apply upward impulse",
    wall: "Flip velocity on walls",
    fire: "Append a shot",
    fly: "Move shots",
    cull: "Drop off-screen shots",
    kill: "Overwrite the seat",
    overwrite: "Overwrite the seat",
    place: "Place a bomb",
    tick: "Tick the fuse",
    queue: "Queue facing",
    center: "Commit at tile center",
  },
};
