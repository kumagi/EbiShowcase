/* --- TRY IT × Go highlight API (see FIX_TRY_IT_RALPH.md) ---
 * Supported data-lab kinds (non-exhaustive; see also scripts/try-it-go-catalog.mjs):
 * loop, hitbox, circles, meter, flappy, bounce, bricks, bullets, camera, gravity,
 * friction, sling-drag, move8, energy-turn, deck-cycle, input-buffer, turn, aim,
 * translate, geom, colorscale, opacity, blend, sheet, spray, spellbook,
 * fx-split, and various track-specific labs.
 */

function focusGo(lab, ids, caption) {
  if (!lab) return;
  const list = Array.isArray(ids) ? ids : String(ids || "").split(/\s+/).filter(Boolean);
  lab.querySelectorAll(".lab-go [data-go-line]").forEach((el) => {
    el.classList.toggle("is-active", list.includes(el.dataset.goLine));
  });
  const cap = lab.querySelector("[data-lab-go-caption]");
  if (cap) {
    let text = caption;
    if (text && typeof text === "object") {
      const lang = document.documentElement.lang === "en" ? "en" : "ja";
      text = text[lang] || text.ja || text.en || "—";
    }
    if (text == null || text === "") text = list.length ? list.join(" · ") : "—";
    cap.textContent = text;
  }
}

function clearGo(lab) {
  focusGo(lab, [], "—");
}

function initFlappyLab(lab) {
  const section = lab.closest(".physics");
  const bird = lab.querySelector("[data-lab-bird]");
  const velocityOutput = lab.querySelector("[data-lab-velocity]");
  const positionOutput = lab.querySelector("[data-lab-position]");
  const directionOutput = lab.querySelector("[data-lab-direction]");
  if (!bird || !velocityOutput || !positionOutput || !directionOutput) return;

  let velocity = 0;
  let position = 320;

  const render = () => {
    const percent = 8 + (Math.max(0, Math.min(720, position)) / 720) * 84;
    bird.style.top = percent + "%";
    velocityOutput.textContent = velocity.toFixed(2);
    positionOutput.textContent = position.toFixed(2);
    directionOutput.textContent = velocity < -0.01 ? section.dataset.up : velocity > 0.01 ? section.dataset.down : section.dataset.still;
  };

  const step = () => {
    velocity += 0.42;
    position = Math.max(0, Math.min(720, position + velocity));
    render();
  };

  lab.querySelector("[data-lab-flap]")?.addEventListener("click", () => {
    velocity = -7.4;
    step();
    focusGo(lab, ["flap", "grav", "pos"]);
  });
  lab.querySelector("[data-lab-step]")?.addEventListener("click", () => {
    step();
    focusGo(lab, ["grav", "pos"]);
  });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    velocity = 0;
    position = 320;
    render();
    clearGo(lab);
  });
  render();
}

document.querySelectorAll(".motion-lab").forEach((lab) => {
  const kind = lab.dataset.lab || (lab.querySelector("[data-lab-bird]") ? "flappy" : "");
  if (kind === "flappy") initFlappyLab(lab);
});

document.querySelectorAll(".motion-lab[data-lab='hit-test']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const tap = lab.querySelector("[data-lab-tap]");
  const result = lab.querySelector("[data-lab-result]");
  const pointOut = lab.querySelector("[data-lab-point]");
  const dxOut = lab.querySelector("[data-lab-dx]");
  const dyOut = lab.querySelector("[data-lab-dy]");
  const distOut = lab.querySelector("[data-lab-dist]");
  if (!board || !tap || !result || !dxOut || !dyOut || !distOut) return;

  const hitLabel = lab.dataset.hit || "HIT";
  const missLabel = lab.dataset.miss || "MISS";
  const circleX = 240;
  const circleY = 260;
  const radius = 70;

  const evaluate = (x, y) => {
    const dx = x - circleX;
    const dy = y - circleY;
    const dist = Math.hypot(dx, dy);
    const hit = dist <= radius;
    tap.style.left = `${(x / 480) * 100}%`;
    tap.style.top = `${(y / 520) * 100}%`;
    tap.hidden = false;
    if (pointOut) pointOut.textContent = `(${x.toFixed(0)}, ${y.toFixed(0)})`;
    dxOut.textContent = dx.toFixed(0);
    dyOut.textContent = dy.toFixed(0);
    distOut.textContent = dist.toFixed(1);
    result.textContent = hit ? hitLabel : missLabel;
    result.dataset.state = hit ? "hit" : "miss";
    board.dataset.state = hit ? "hit" : "miss";
    focusGo(lab, hit ? ["dxdy", "hypot", "hit"] : ["dxdy", "hypot", "miss"]);
  };

  const placeFromEvent = (event) => {
    const rect = board.getBoundingClientRect();
    const clientX = event.touches ? event.touches[0].clientX : event.clientX;
    const clientY = event.touches ? event.touches[0].clientY : event.clientY;
    const x = ((clientX - rect.left) / rect.width) * 480;
    const y = ((clientY - rect.top) / rect.height) * 520;
    evaluate(Math.max(0, Math.min(480, x)), Math.max(0, Math.min(520, y)));
  };

  board.addEventListener("click", placeFromEvent);
  board.addEventListener("touchstart", (event) => {
    event.preventDefault();
    placeFromEvent(event);
  }, { passive: false });

  lab.querySelector("[data-lab-sample-hit]")?.addEventListener("click", () => evaluate(250, 240));
  lab.querySelector("[data-lab-sample-miss]")?.addEventListener("click", () => evaluate(390, 120));
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    tap.hidden = true;
    if (pointOut) pointOut.textContent = "—";
    dxOut.textContent = "—";
    dyOut.textContent = "—";
    distOut.textContent = "—";
    result.textContent = lab.dataset.wait || "TAP";
    result.dataset.state = "wait";
    board.dataset.state = "wait";
    clearGo(lab);
  });
});

/* Optical-trick labs used by VFX BASIC 14–17. */
document.querySelectorAll(".motion-lab[data-lab='squash']").forEach((lab) => {
  const actor = lab.querySelector("[data-lab-actor]");
  const xOut = lab.querySelector("[data-lab-xscale]");
  const yOut = lab.querySelector("[data-lab-yscale]");
  const setShape = (name) => {
    const values = name === "squash" ? [1.28, .72] : name === "stretch" ? [.72, 1.28] : [1, 1];
    actor.style.transform = `translateX(-50%) scale(${values[0]}, ${values[1]})`;
    setText(xOut, values[0].toFixed(2)); setText(yOut, values[1].toFixed(2));
  };
  lab.querySelectorAll("[data-lab-shape]").forEach((b) => b.addEventListener("click", () => setShape(b.dataset.labShape)));
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => setShape("neutral"));
  setShape("neutral");
});

document.querySelectorAll(".motion-lab[data-lab='outline']").forEach((lab) => {
  const actor = lab.querySelector("[data-lab-outline]");
  const widthOut = lab.querySelector("[data-lab-width]");
  const modeOut = lab.querySelector("[data-lab-mode]");
  let width = 4, light = false;
  const render = () => {
    const c = light ? "#fff" : "#091126";
    actor.style.webkitTextStroke = `${width}px ${c}`;
    setText(widthOut, `${width}px`); setText(modeOut, light ? "light" : "dark");
  };
  lab.querySelector("[data-lab-outline-down]")?.addEventListener("click", () => { width = Math.max(1, width - 1); render(); });
  lab.querySelector("[data-lab-outline-up]")?.addEventListener("click", () => { width = Math.min(10, width + 1); render(); });
  lab.querySelector("[data-lab-outline-color]")?.addEventListener("click", () => { light = !light; render(); });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => { width = 4; light = false; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='impact']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const raysNode = lab.querySelector("[data-lab-rays]");
  const out = lab.querySelector("[data-lab-rays-value]");
  let rays = 16;
  const render = () => {
    raysNode.innerHTML = Array.from({length: rays}, (_, i) => `<i style="transform:rotate(${i * 360 / rays}deg)"></i>`).join("");
    setText(out, String(rays));
  };
  const burst = () => { board.classList.remove("is-burst"); void board.offsetWidth; board.classList.add("is-burst"); };
  lab.querySelector("[data-lab-impact]")?.addEventListener("click", burst);
  lab.querySelector("[data-lab-rays-down]")?.addEventListener("click", () => { rays = Math.max(6, rays - 2); render(); });
  lab.querySelector("[data-lab-rays-up]")?.addEventListener("click", () => { rays = Math.min(32, rays + 2); render(); });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => { rays = 16; render(); burst(); });
  render(); burst();
});

document.querySelectorAll(".motion-lab[data-lab='bloom']").forEach((lab) => {
  const orb = lab.querySelector("[data-lab-bloom]");
  const copiesOut = lab.querySelector("[data-lab-copies]");
  const modeOut = lab.querySelector("[data-lab-mode]");
  let copies = 5, on = true;
  const render = () => {
    const spread = 18 + copies * 8;
    orb.style.boxShadow = on ? `0 0 ${spread}px ${Math.round(spread*.45)}px rgba(73,210,255,.48),0 0 ${spread*2}px ${Math.round(spread*.2)}px rgba(126,92,255,.28)` : "none";
    setText(copiesOut, String(copies)); setText(modeOut, on ? "ON" : "OFF");
  };
  lab.querySelector("[data-lab-bloom-down]")?.addEventListener("click", () => { copies = Math.max(1, copies - 1); render(); });
  lab.querySelector("[data-lab-bloom-up]")?.addEventListener("click", () => { copies = Math.min(9, copies + 1); render(); });
  lab.querySelector("[data-lab-bloom-toggle]")?.addEventListener("click", () => { on = !on; render(); });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => { copies = 5; on = true; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='meter']").forEach((lab) => {
  const marker = lab.querySelector("[data-lab-marker]");
  const xOut = lab.querySelector("[data-lab-x]");
  const speedOut = lab.querySelector("[data-lab-speed]");
  const stateOut = lab.querySelector("[data-lab-state]");
  const scoreOut = lab.querySelector("[data-lab-score]");
  if (!marker || !xOut || !speedOut || !stateOut || !scoreOut) return;

  const movingLabel = lab.dataset.moving || "moving";
  const stoppedLabel = lab.dataset.stopped || "stopped";
  let x = 45;
  let speed = 8;
  let stopped = false;
  let score = 0;
  const minX = 45;
  const maxX = 435;
  const center = 240;

  const render = () => {
    marker.style.left = `${((x - minX) / (maxX - minX)) * 100}%`;
    xOut.textContent = x.toFixed(0);
    speedOut.textContent = speed.toFixed(1);
    stateOut.textContent = stopped ? stoppedLabel : movingLabel;
    scoreOut.textContent = String(score);
  };

  lab.querySelector("[data-lab-step]")?.addEventListener("click", () => {
    if (stopped) return;
    x += speed;
    if (x < minX || x > maxX) {
      speed = -speed;
      x = Math.max(minX, Math.min(maxX, x));
    }
    render();
  });

  lab.querySelector("[data-lab-stop]")?.addEventListener("click", () => {
    if (stopped) {
      stopped = false;
      render();
      return;
    }
    stopped = true;
    const distance = Math.abs(x - center);
    if (distance <= 8) score += 100;
    else if (distance <= 28) score += 50;
    else if (distance <= 55) score += 10;
    render();
  });

  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    x = 45;
    speed = 8;
    stopped = false;
    score = 0;
    render();
  });

  render();
});

document.querySelectorAll(".motion-lab[data-lab='aabb']").forEach((lab) => {
  const star = lab.querySelector("[data-lab-star]");
  const basket = lab.querySelector("[data-lab-basket]");
  const yOut = lab.querySelector("[data-lab-y]");
  const overlapOut = lab.querySelector("[data-lab-overlap]");
  if (!star || !basket || !yOut || !overlapOut) return;

  const yes = lab.dataset.yes || "YES";
  const no = lab.dataset.no || "NO";
  let y = 80;
  const starSize = 30;
  const basketY = 360;
  const basketH = 38;
  const basketX = 180;
  const basketW = 116;
  const starX = 210;

  const overlaps = () => {
    const left = starX;
    const right = starX + starSize;
    const top = y;
    const bottom = y + starSize;
    return left < basketX + basketW && right > basketX && top < basketY + basketH && bottom > basketY;
  };

  const render = () => {
    star.style.top = `${(y / 480) * 100}%`;
    star.style.left = `${(starX / 480) * 100}%`;
    basket.style.left = `${(basketX / 480) * 100}%`;
    basket.style.top = `${(basketY / 480) * 100}%`;
    yOut.textContent = y.toFixed(0);
    const hit = overlaps();
    overlapOut.textContent = hit ? yes : no;
    overlapOut.dataset.state = hit ? "hit" : "miss";
  };

  lab.querySelector("[data-lab-step]")?.addEventListener("click", () => {
    y = Math.min(450, y + 28);
    render();
  });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    y = 80;
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='jump']").forEach((lab) => {
  const actor = lab.querySelector("[data-lab-actor]");
  const yOut = lab.querySelector("[data-lab-y]");
  const vOut = lab.querySelector("[data-lab-v]");
  const gOut = lab.querySelector("[data-lab-grounded]");
  if (!actor || !yOut || !vOut || !gOut) return;

  const yes = lab.dataset.yes || "YES";
  const no = lab.dataset.no || "NO";
  const floor = 420;
  const gravity = 0.55;
  const jumpPower = -9.5;
  let y = floor;
  let vy = 0;
  let grounded = true;

  const render = () => {
    actor.style.top = `${(y / 520) * 100}%`;
    yOut.textContent = y.toFixed(0);
    vOut.textContent = vy.toFixed(2);
    gOut.textContent = grounded ? yes : no;
    gOut.dataset.state = grounded ? "hit" : "miss";
  };

  lab.querySelector("[data-lab-jump]")?.addEventListener("click", () => {
    if (!grounded) return;
    vy = jumpPower;
    grounded = false;
    render();
  });
  lab.querySelector("[data-lab-step]")?.addEventListener("click", () => {
    vy += gravity;
    y += vy;
    if (y >= floor) {
      y = floor;
      vy = 0;
      grounded = true;
    } else {
      grounded = false;
    }
    render();
  });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    y = floor;
    vy = 0;
    grounded = true;
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='cost']").forEach((lab) => {
  const countOut = lab.querySelector("[data-lab-count]");
  const costOut = lab.querySelector("[data-lab-cost]");
  const rateOut = lab.querySelector("[data-lab-rate]");
  if (!countOut || !costOut || !rateOut) return;

  let count = 0;
  const base = Number(lab.dataset.base || 10);
  const growth = Number(lab.dataset.growth || 1.18);

  const format = (n) => {
    if (n >= 1e6) return (n / 1e6).toFixed(2) + "M";
    if (n >= 1e3) return (n / 1e3).toFixed(2) + "K";
    return n.toFixed(1);
  };

  const render = () => {
    const cost = base * Math.pow(growth, count);
    const rate = count * count * 8;
    countOut.textContent = String(count);
    costOut.textContent = format(cost);
    rateOut.textContent = format(rate) + "/s";
  };

  lab.querySelector("[data-lab-buy]")?.addEventListener("click", () => {
    count += 1;
    render();
  });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    count = 0;
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='hitbox'], .motion-lab[data-lab='circles']").forEach((lab) => {
  const atk = lab.querySelector("[data-lab-attack]");
  const hurt = lab.querySelector("[data-lab-hurt]");
  const result = lab.querySelector("[data-lab-result]");
  const distOut = lab.querySelector("[data-lab-sep]");
  if (!atk || !hurt || !result || !distOut) return;

  const hitLabel = lab.dataset.hit || "HIT";
  const missLabel = lab.dataset.miss || "MISS";
  let sep = 90;

  const render = () => {
    atk.style.left = `${((180) / 480) * 100}%`;
    hurt.style.left = `${((180 + sep) / 480) * 100}%`;
    const overlap = sep < 70;
    result.textContent = overlap ? hitLabel : missLabel;
    result.dataset.state = overlap ? "hit" : "miss";
    distOut.textContent = String(sep);
  };

  lab.querySelector("[data-lab-closer]")?.addEventListener("click", () => {
    sep = Math.max(20, sep - 15);
    render();
  });
  lab.querySelector("[data-lab-farther]")?.addEventListener("click", () => {
    sep = Math.min(160, sep + 15);
    render();
  });
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    sep = 90;
    render();
  });
  render();
});


function bind(lab, sel) {
  return lab.querySelector(sel);
}
function setText(el, v) { if (el) el.textContent = v; }
function setState(el, hit) {
  if (!el) return;
  el.dataset.state = hit ? "hit" : "miss";
}

document.querySelectorAll(".motion-lab[data-lab='bounce']").forEach((lab) => {
  const ball = bind(lab, "[data-lab-ball]");
  const vxOut = bind(lab, "[data-lab-vx]");
  const vyOut = bind(lab, "[data-lab-vy]");
  const note = bind(lab, "[data-lab-note]");
  if (!ball || !vxOut || !vyOut) return;
  let x = 120, y = 200, vx = 6, vy = 4;
  const wall = lab.dataset.wall || "wall";
  const paddle = lab.dataset.paddle || "paddle";
  const render = () => {
    ball.style.left = `${(x / 480) * 100}%`;
    ball.style.top = `${(y / 480) * 100}%`;
    setText(vxOut, vx.toFixed(1));
    setText(vyOut, vy.toFixed(1));
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    x += vx; y += vy;
    let msg = "—";
    if (x < 20 || x > 460) { vx = -vx; x = Math.max(20, Math.min(460, x)); msg = wall + " → vx = -vx"; }
    if (y < 20) { vy = -vy; y = 20; msg = wall + " → vy = -vy"; }
    if (y > 400) { vy = -Math.abs(vy); y = 400; msg = paddle + " → vy = -vy"; }
    setText(note, msg);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { x = 120; y = 200; vx = 6; vy = 4; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='bricks']").forEach((lab) => {
  const aliveOut = bind(lab, "[data-lab-alive]");
  const scoreOut = bind(lab, "[data-lab-score]");
  const grid = bind(lab, "[data-lab-grid]");
  if (!aliveOut || !scoreOut || !grid) return;
  let alive = Array(12).fill(true);
  let score = 0;
  const render = () => {
    grid.innerHTML = alive.map((a, i) => `<span class="${a ? "on" : "off"}" data-i="${i}"></span>`).join("");
    setText(aliveOut, String(alive.filter(Boolean).length));
    setText(scoreOut, String(score));
  };
  grid.addEventListener("click", (e) => {
    const cell = e.target.closest("[data-i]");
    if (!cell) return;
    const i = Number(cell.dataset.i);
    if (!alive[i]) return;
    alive[i] = false;
    score += 10;
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { alive = Array(12).fill(true); score = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='snake']").forEach((lab) => {
  const lenOut = bind(lab, "[data-lab-len]");
  const bodyOut = bind(lab, "[data-lab-body]");
  const ateLabel = lab.dataset.ate || "ate";
  const moveLabel = lab.dataset.move || "move";
  if (!lenOut || !bodyOut) return;
  let body = [{ x: 2, y: 2 }, { x: 1, y: 2 }, { x: 0, y: 2 }];
  let grow = false;
  const render = () => {
    setText(lenOut, String(body.length));
    setText(bodyOut, body.map((p) => `(${p.x},${p.y})`).join(" "));
  };
  const step = (eat) => {
    const h = body[0];
    const head = { x: h.x + 1, y: h.y };
    body = [head, ...body];
    if (!eat && !grow) body.pop();
    grow = false;
    render();
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => step(false));
  bind(lab, "[data-lab-eat]")?.addEventListener("click", () => { grow = true; step(true); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { body = [{ x: 2, y: 2 }, { x: 1, y: 2 }, { x: 0, y: 2 }]; grow = false; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='entities']").forEach((lab) => {
  const countOut = bind(lab, "[data-lab-count]");
  const listOut = bind(lab, "[data-lab-list]");
  if (!countOut || !listOut) return;
  let shots = [];
  let id = 1;
  const render = () => {
    setText(countOut, String(shots.length));
    setText(listOut, shots.length ? shots.map((s) => `#${s.id}@${s.y}`).join(" · ") : "—");
  };
  bind(lab, "[data-lab-fire]")?.addEventListener("click", () => {
    shots.push({ id: id++, y: 400 });
    render();
  });
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    shots = shots.map((s) => ({ ...s, y: s.y - 40 })).filter((s) => s.y > 0);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { shots = []; id = 1; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='push']").forEach((lab) => {
  const mapOut = bind(lab, "[data-lab-map]");
  const note = bind(lab, "[data-lab-note]");
  if (!mapOut) return;
  // 5 cells: . P B . #
  let cells = [".", "P", ".", "B", "."];
  const wallMsg = lab.dataset.blocked || "blocked";
  const pushedMsg = lab.dataset.pushed || "pushed";
  const render = () => setText(mapOut, cells.join(" "));
  bind(lab, "[data-lab-right]")?.addEventListener("click", () => {
    const p = cells.indexOf("P");
    const next = p + 1;
    if (next >= cells.length || cells[next] === "#") { setText(note, wallMsg); return; }
    if (cells[next] === "B") {
      const n2 = next + 1;
      if (n2 >= cells.length || cells[n2] !== ".") { setText(note, wallMsg); return; }
      cells[n2] = "B";
      cells[next] = "P";
      cells[p] = ".";
      setText(note, pushedMsg);
    } else {
      cells[next] = "P";
      cells[p] = ".";
      setText(note, "→");
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { cells = [".", "P", ".", "B", "."]; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='camera']").forEach((lab) => {
  const camOut = bind(lab, "[data-lab-cam]");
  const playerOut = bind(lab, "[data-lab-player]");
  const screenOut = bind(lab, "[data-lab-screen]");
  const factorOut = bind(lab, "[data-lab-factor]");
  const actor = bind(lab, "[data-lab-actor]");
  const ruler = bind(lab, "[data-camera-ruler]");
  if (!camOut || !playerOut || !screenOut) return;
  let player = 200, cam = 0, factor = Number(lab.dataset.factor || 0.09);
  const render = () => {
    const target = player - 160;
    cam += (target - cam) * factor;
    setText(camOut, cam.toFixed(0));
    setText(playerOut, player.toFixed(0));
    setText(screenOut, (player - cam).toFixed(0));
    setText(factorOut, factor.toFixed(2));
    if (actor) actor.style.left = `${((player - cam) / 480) * 100}%`;
    if (ruler) ruler.style.transform = `translateX(-${(cam / 960) * 100}%)`;
  };
  bind(lab, "[data-lab-right]")?.addEventListener("click", () => { player += 40; render(); });
  bind(lab, "[data-lab-left]")?.addEventListener("click", () => { player = Math.max(40, player - 40); render(); });
  bind(lab, "[data-lab-step]")?.addEventListener("click", render);
  lab.querySelectorAll("[data-camera-factor]").forEach((button) => button.addEventListener("click", () => { factor = Number(button.dataset.cameraFactor); render(); }));
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { player = 200; cam = 0; factor = Number(lab.dataset.factor || 0.09); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='ai']").forEach((lab) => {
  const stateOut = bind(lab, "[data-lab-state]");
  const distOut = bind(lab, "[data-lab-dist]");
  const patrol = lab.dataset.patrol || "patrol";
  const chase = lab.dataset.chase || "chase";
  if (!stateOut || !distOut) return;
  let dist = 180;
  const render = () => {
    const chasing = dist < 120;
    setText(stateOut, chasing ? chase : patrol);
    setState(stateOut, chasing);
    setText(distOut, dist.toFixed(0));
  };
  bind(lab, "[data-lab-closer]")?.addEventListener("click", () => { dist = Math.max(20, dist - 30); render(); });
  bind(lab, "[data-lab-farther]")?.addEventListener("click", () => { dist = Math.min(260, dist + 30); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { dist = 180; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='burst']").forEach((lab) => {
  const countOut = bind(lab, "[data-lab-count]");
  const board = bind(lab, "[data-lab-board]");
  if (!countOut || !board) return;
  let count = 8;
  const render = () => {
    setText(countOut, String(count));
    const dots = [];
    for (let i = 0; i < count; i++) {
      const a = (i / count) * Math.PI * 2 - Math.PI / 2;
      const x = 50 + Math.cos(a) * 32;
      const y = 50 + Math.sin(a) * 32;
      dots.push(`<span style="left:${x}%;top:${y}%"></span>`);
    }
    board.innerHTML = `<div class="lab-burst">${dots.join("")}</div>`;
  };
  bind(lab, "[data-lab-more]")?.addEventListener("click", () => { count = Math.min(24, count + 2); render(); });
  bind(lab, "[data-lab-less]")?.addEventListener("click", () => { count = Math.max(4, count - 2); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { count = 8; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='carry']").forEach((lab) => {
  const platOut = bind(lab, "[data-lab-plat]");
  const playerOut = bind(lab, "[data-lab-player]");
  const deltaOut = bind(lab, "[data-lab-delta]");
  if (!platOut || !playerOut || !deltaOut) return;
  let plat = 100, prev = 100, player = 120, riding = true;
  const render = () => {
    setText(platOut, plat.toFixed(0));
    setText(playerOut, player.toFixed(0));
    setText(deltaOut, (plat - prev).toFixed(0));
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    prev = plat;
    plat = 100 + Math.sin(Date.now() / 400) * 0;
    // deterministic step: move +12 then -12 alternating via dataset
    const dir = Number(lab.dataset.dir || 1);
    plat = prev + 12 * dir;
    lab.dataset.dir = String(-dir);
    if (riding) player += plat - prev;
    render();
  });
  bind(lab, "[data-lab-toggle]")?.addEventListener("click", () => { riding = !riding; });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { plat = 100; prev = 100; player = 120; lab.dataset.dir = "1"; riding = true; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='stomp']").forEach((lab) => {
  const result = bind(lab, "[data-lab-result]");
  const vyOut = bind(lab, "[data-lab-vy]");
  const stomp = lab.dataset.stomp || "STOMP";
  const hurt = lab.dataset.hurt || "HURT";
  if (!result || !vyOut) return;
  let vy = 4;
  const render = () => {
    const ok = vy > 0;
    setText(result, ok ? stomp : hurt);
    setState(result, ok);
    setText(vyOut, vy.toFixed(1));
  };
  bind(lab, "[data-lab-fall]")?.addEventListener("click", () => { vy = 4; render(); });
  bind(lab, "[data-lab-rise]")?.addEventListener("click", () => { vy = -4; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { vy = 4; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='power']").forEach((lab) => {
  const poweredOut = bind(lab, "[data-lab-powered]");
  const result = bind(lab, "[data-lab-result]");
  const yes = lab.dataset.yes || "ON";
  const no = lab.dataset.no || "OFF";
  const stomp = lab.dataset.stomp || "stomp";
  const damage = lab.dataset.damage || "damage";
  if (!poweredOut || !result) return;
  let powered = false;
  const render = () => {
    setText(poweredOut, powered ? yes : no);
    setState(poweredOut, powered);
    setText(result, powered ? stomp : damage);
    setState(result, powered);
  };
  bind(lab, "[data-lab-toggle]")?.addEventListener("click", () => { powered = !powered; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { powered = false; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='move8'], .motion-lab[data-lab='sling-drag']").forEach((lab) => {
  const rawOut = bind(lab, "[data-lab-raw]");
  const normOut = bind(lab, "[data-lab-norm]");
  if (!rawOut || !normOut) return;
  let dx = 1, dy = 1, speed = 4;
  const render = () => {
    const raw = Math.hypot(dx * speed, dy * speed);
    const len = Math.hypot(dx, dy) || 1;
    const nx = (dx / len) * speed;
    const ny = (dy / len) * speed;
    setText(rawOut, raw.toFixed(2));
    setText(normOut, Math.hypot(nx, ny).toFixed(2));
  };
  bind(lab, "[data-lab-cardinal]")?.addEventListener("click", () => { dx = 1; dy = 0; render(); });
  bind(lab, "[data-lab-diag]")?.addEventListener("click", () => { dx = 1; dy = 1; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { dx = 1; dy = 1; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='aim']").forEach((lab) => {
  const targetOut = bind(lab, "[data-lab-target]");
  const cdOut = bind(lab, "[data-lab-cd]");
  const shotOut = bind(lab, "[data-lab-shots]");
  const note = bind(lab, "[data-lab-note]");
  if (!targetOut || !cdOut || !shotOut) return;
  const enemies = [{ id: "A", d: 120 }, { id: "B", d: 80 }, { id: "C", d: 200 }];
  let frame = 0, shots = 0, cooldown = 8;
  const nearest = () => enemies.reduce((b, e) => (e.d < b.d ? e : b));
  const render = () => {
    setText(targetOut, nearest().id);
    setText(cdOut, String(frame % cooldown));
    setText(shotOut, String(shots));
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    // Advance until the next shot so one click always feels like an action.
    let fired = false;
    for (let i = 0; i < cooldown; i++) {
      frame++;
      if (frame % cooldown === 0) {
        shots++;
        fired = true;
        break;
      }
    }
    setText(note, fired ? (lab.dataset.fire || `PEW → ${nearest().id}`) : "…");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    frame = 0; shots = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='pool']").forEach((lab) => {
  const slotsOut = bind(lab, "[data-lab-slots]");
  const killsOut = bind(lab, "[data-lab-kills]");
  if (!slotsOut || !killsOut) return;
  let slots = Array.from({ length: 8 }, (_, i) => `E${i}`);
  let kills = 0;
  const render = () => {
    setText(slotsOut, slots.join(" "));
    setText(killsOut, String(kills));
  };
  bind(lab, "[data-lab-kill]")?.addEventListener("click", () => {
    const i = kills % slots.length;
    kills++;
    slots[i] = `N${kills}`;
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { slots = Array.from({ length: 8 }, (_, i) => `E${i}`); kills = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='draft']").forEach((lab) => {
  const modeOut = bind(lab, "[data-lab-mode]");
  const pickOut = bind(lab, "[data-lab-pick]");
  const combat = lab.dataset.combat || "combat";
  const draft = lab.dataset.draft || "draft";
  if (!modeOut || !pickOut) return;
  let drafting = false, pick = "—";
  const render = () => {
    setText(modeOut, drafting ? draft : combat);
    setState(modeOut, drafting);
    setText(pickOut, pick);
  };
  bind(lab, "[data-lab-level]")?.addEventListener("click", () => { drafting = true; pick = "—"; render(); });
  lab.querySelectorAll("[data-lab-card]").forEach((btn) => {
    btn.addEventListener("click", () => {
      if (!drafting) return;
      pick = btn.dataset.card || "A";
      drafting = false;
      render();
    });
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { drafting = false; pick = "—"; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='evolve']").forEach((lab) => {
  const nameOut = bind(lab, "[data-lab-name]");
  const countOut = bind(lab, "[data-lab-wcount]");
  const cdOut = bind(lab, "[data-lab-wcd]");
  if (!nameOut || !countOut || !cdOut) return;
  let evolved = false;
  const render = () => {
    setText(nameOut, evolved ? "Ebi Storm" : "Ebi Needle");
    setText(countOut, evolved ? "3" : "1");
    setText(cdOut, evolved ? "16" : "32");
  };
  bind(lab, "[data-lab-evolve]")?.addEventListener("click", () => { evolved = true; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { evolved = false; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='curve']").forEach((lab) => {
  const secOut = bind(lab, "[data-lab-sec]");
  const intervalOut = bind(lab, "[data-lab-interval]");
  const speedOut = bind(lab, "[data-lab-speed]");
  if (!secOut || !intervalOut || !speedOut) return;
  let sec = 0;
  const render = () => {
    const interval = Math.max(14, 42 - Math.floor(sec / 2));
    const speed = 0.85 + sec * 0.018;
    setText(secOut, String(sec));
    setText(intervalOut, String(interval));
    setText(speedOut, speed.toFixed(2));
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { sec = Math.min(45, sec + 5); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { sec = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='click']").forEach((lab) => {
  const countOut = bind(lab, "[data-lab-count]");
  if (!countOut) return;
  let n = 0;
  const render = () => setText(countOut, String(n));
  bind(lab, "[data-lab-tap]")?.addEventListener("click", () => { n++; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { n = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='shop']").forEach((lab) => {
  const goldOut = bind(lab, "[data-lab-gold]");
  const costOut = bind(lab, "[data-lab-cost]");
  const ownedOut = bind(lab, "[data-lab-owned]");
  const note = bind(lab, "[data-lab-note]");
  if (!goldOut || !costOut || !ownedOut) return;
  let gold = 0, owned = 0, cost = 10;
  const ok = lab.dataset.bought || "bought";
  const no = lab.dataset.cant || "not enough";
  const render = () => {
    setText(goldOut, String(gold));
    setText(costOut, String(cost));
    setText(ownedOut, String(owned));
  };
  bind(lab, "[data-lab-earn]")?.addEventListener("click", () => { gold += 5; render(); });
  bind(lab, "[data-lab-buy]")?.addEventListener("click", () => {
    if (gold < cost) { setText(note, no); return; }
    gold -= cost; owned++; cost += 5; setText(note, ok); render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { gold = 0; owned = 0; cost = 10; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='idle']").forEach((lab) => {
  const pointsOut = bind(lab, "[data-lab-points]");
  const rateOut = bind(lab, "[data-lab-rate]");
  if (!pointsOut || !rateOut) return;
  let points = 0, machines = 2, dt = 0.25;
  const render = () => {
    setText(pointsOut, points.toFixed(1));
    setText(rateOut, (machines).toFixed(0) + "/s");
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { points += machines * dt; render(); });
  bind(lab, "[data-lab-buy]")?.addEventListener("click", () => { machines++; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { points = 0; machines = 2; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='save']").forEach((lab) => {
  const awayOut = bind(lab, "[data-lab-away]");
  const gainedOut = bind(lab, "[data-lab-gained]");
  if (!awayOut || !gainedOut) return;
  let away = 0, rate = 8;
  const render = () => {
    setText(awayOut, away + "s");
    setText(gainedOut, (away * rate).toFixed(0));
  };
  bind(lab, "[data-lab-away]")?.addEventListener("click", () => { away += 30; render(); });
  // button uses data-lab-add to avoid conflict
  const addBtn = lab.querySelector("[data-lab-add]");
  addBtn?.addEventListener("click", () => { away += 30; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { away = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='tile']").forEach((lab) => {
  const posOut = bind(lab, "[data-lab-pos]");
  const faceOut = bind(lab, "[data-lab-face]");
  const note = bind(lab, "[data-lab-note]");
  if (!posOut || !faceOut) return;
  let x = 1, y = 1, face = "S";
  const walls = new Set(["2,1"]);
  const blocked = lab.dataset.blocked || "wall";
  const render = () => { setText(posOut, `${x},${y}`); setText(faceOut, face); };
  const tryMove = (dx, dy, f) => {
    face = f;
    const nx = x + dx, ny = y + dy;
    if (walls.has(`${nx},${ny}`)) { setText(note, blocked); render(); return; }
    x = nx; y = ny; setText(note, "→"); render();
  };
  bind(lab, "[data-lab-up]")?.addEventListener("click", () => tryMove(0, -1, "N"));
  bind(lab, "[data-lab-down]")?.addEventListener("click", () => tryMove(0, 1, "S"));
  bind(lab, "[data-lab-left]")?.addEventListener("click", () => tryMove(-1, 0, "W"));
  bind(lab, "[data-lab-right]")?.addEventListener("click", () => tryMove(1, 0, "E"));
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { x = 1; y = 1; face = "S"; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='flag']").forEach((lab) => {
  const flagOut = bind(lab, "[data-lab-flag]");
  const textOut = bind(lab, "[data-lab-text]");
  const yes = lab.dataset.yes || "true";
  const no = lab.dataset.no || "false";
  const before = lab.dataset.before || "Hello.";
  const after = lab.dataset.after || "You found it!";
  if (!flagOut || !textOut) return;
  let flag = false;
  const render = () => {
    setText(flagOut, flag ? yes : no);
    setState(flagOut, flag);
    setText(textOut, flag ? after : before);
  };
  bind(lab, "[data-lab-toggle]")?.addEventListener("click", () => { flag = !flag; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { flag = false; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='turn']").forEach((lab) => {
  const stateOut = bind(lab, "[data-lab-state]") || lab.querySelector(".lab-readout");
  const board = lab.querySelector("[data-lab-board], .lab-board, .lab-entities");
  const states = (lab.dataset.states || "select,player,enemy,win").split(",");
  if (!stateOut) return;
  let i = 0;
  const render = () => {
    setText(stateOut, states[i]);
    if (board) {
      board.innerHTML = `<div class="lab-turn-row">${states.map((state, n) =>
        `<span class="${n === i ? "now" : n < i ? "done" : ""}">${state}</span>`).join("<b>→</b>")}</div>`;
    }
  };
  (bind(lab, "[data-lab-step]") || lab.querySelector(".lab-action"))?.addEventListener("click", () => { i = (i + 1) % states.length; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { i = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='damage']").forEach((lab) => {
  const dmgOut = bind(lab, "[data-lab-dmg]");
  const atkOut = bind(lab, "[data-lab-atk]");
  const defOut = bind(lab, "[data-lab-def]");
  if (!dmgOut || !atkOut || !defOut) return;
  let atk = 10, def = 3, buff = 0;
  const render = () => {
    const dmg = Math.max(1, atk + buff - def);
    setText(atkOut, String(atk + buff));
    setText(defOut, String(def));
    setText(dmgOut, String(dmg));
  };
  bind(lab, "[data-lab-buff]")?.addEventListener("click", () => { buff = 5; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { buff = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='inv']").forEach((lab) => {
  const goldOut = bind(lab, "[data-lab-gold]");
  const itemOut = bind(lab, "[data-lab-item]");
  const note = bind(lab, "[data-lab-note]");
  if (!goldOut || !itemOut) return;
  let gold = 100, item = "—";
  const render = () => { setText(goldOut, String(gold)); setText(itemOut, item); };
  bind(lab, "[data-lab-buy]")?.addEventListener("click", () => {
    if (gold < 50) { setText(note, lab.dataset.cant || "need gold"); return; }
    gold -= 50; item = "Sword"; setText(note, lab.dataset.ok || "equipped"); render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { gold = 100; item = "—"; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='scene']").forEach((lab) => {
  const sceneOut = bind(lab, "[data-lab-scene]");
  const enemyOut = bind(lab, "[data-lab-enemy]");
  const field = lab.dataset.field || "field";
  const battle = lab.dataset.battle || "battle";
  if (!sceneOut || !enemyOut) return;
  let scene = field, region = "grass", enemy = "—";
  const tables = { grass: "Slime", desert: "Cactus", snow: "Wolf" };
  const render = () => { setText(sceneOut, scene); setText(enemyOut, enemy); };
  bind(lab, "[data-lab-walk]")?.addEventListener("click", () => {
    scene = battle;
    enemy = tables[region] || "Slime";
    render();
  });
  lab.querySelectorAll("[data-lab-region]").forEach((btn) => {
    btn.addEventListener("click", () => {
      region = btn.dataset.region || "grass";
      if (scene === battle) enemy = tables[region];
      render();
    });
  });
  bind(lab, "[data-lab-back]")?.addEventListener("click", () => { scene = field; enemy = "—"; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { scene = field; region = "grass"; enemy = "—"; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='quest']").forEach((lab) => {
  const questOut = bind(lab, "[data-lab-quest]");
  const savedOut = bind(lab, "[data-lab-saved]");
  if (!questOut || !savedOut) return;
  let quest = 0, saved = "—";
  const render = () => { setText(questOut, String(quest)); setText(savedOut, saved); };
  bind(lab, "[data-lab-advance]")?.addEventListener("click", () => { quest++; render(); });
  bind(lab, "[data-lab-save]")?.addEventListener("click", () => { saved = `quest=${quest}`; render(); });
  bind(lab, "[data-lab-load]")?.addEventListener("click", () => {
    const m = /quest=(\d+)/.exec(saved);
    if (m) quest = Number(m[1]);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { quest = 0; saved = "—"; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='frames']").forEach((lab) => {
  const frameOut = bind(lab, "[data-lab-frame]");
  const phaseOut = bind(lab, "[data-lab-phase]");
  const board = lab.querySelector(".lab-board");
  const startup = lab.dataset.startup || "startup";
  const active = lab.dataset.active || "active";
  const recovery = lab.dataset.recovery || "recovery";
  if (!frameOut || !phaseOut) return;
  let frame = 0;
  const phase = () => {
    if (frame === 0) return "ready";
    if (frame <= 8) return startup;
    if (frame <= 12) return active;
    return recovery;
  };
  const render = () => {
    const p = phase();
    setText(frameOut, String(frame));
    setText(phaseOut, p);
    setState(phaseOut, frame > 8 && frame <= 12);
    if (!board) return;
    let reach = 0;
    if (frame >= 1 && frame <= 8) reach = Math.round((frame / 8) * 56);
    else if (frame >= 9 && frame <= 12) reach = 100;
    else if (frame >= 13) reach = Math.max(0, 100 - (frame - 12) * 6);
    const hitbox = frame >= 9 && frame <= 12
      ? `<div class="lab-frame-hitbox" style="width:${reach}%"></div>`
      : `<div class="lab-frame-arm" style="width:${reach}%"></div>`;
    board.className = "lab-board lab-frame-stage";
    board.dataset.phase = p;
    board.innerHTML = `<div class="lab-frame-fighter you"></div>${hitbox}<div class="lab-frame-fighter foe"></div><strong class="lab-frame-label">${p.toUpperCase()} · F${frame}</strong>`;
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { frame = Math.min(30, frame + 1); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { frame = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='react']").forEach((lab) => {
  const stopOut = bind(lab, "[data-lab-stop]");
  const stunOut = bind(lab, "[data-lab-stun]");
  const vOut = bind(lab, "[data-lab-v]");
  const board = lab.querySelector(".lab-board");
  if (!stopOut || !stunOut || !vOut) return;
  let hitstop = 0, stun = 0, v = 0, x = 58;
  const render = () => {
    setText(stopOut, String(hitstop));
    setText(stunOut, String(stun));
    setText(vOut, v.toFixed(2));
    if (!board) return;
    let mode = "idle";
    if (hitstop > 0) mode = "hitstop";
    else if (stun > 0) mode = "stun";
    board.className = "lab-board lab-react-stage";
    board.dataset.mode = mode;
    board.innerHTML = `<div class="lab-react-you"></div><div class="lab-react-foe" style="left:${x}%"></div><strong class="lab-react-label">${mode.toUpperCase()}</strong>`;
  };
  bind(lab, "[data-lab-hit]")?.addEventListener("click", () => {
    hitstop = 8;
    stun = 25;
    v = 7;
    x = 58;
    render();
  });
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    if (hitstop > 0) {
      hitstop--;
      render();
      return;
    }
    if (stun > 0) {
      stun--;
      x = Math.min(88, x + v * 0.55);
      v *= 0.86;
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    hitstop = 0;
    stun = 0;
    v = 0;
    x = 58;
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='rps']").forEach((lab) => {
  const result = bind(lab, "[data-lab-result]");
  if (!result) return;
  const win = lab.dataset.win || "WIN";
  const lose = lab.dataset.lose || "LOSE";
  const clash = lab.dataset.clash || "CLASH";
  const beats = { strike: "throw", throw: "guard", guard: "strike" };
  let you = "guard";
  const renderPick = () => {};
  lab.querySelectorAll("[data-lab-pick]").forEach((btn) => {
    btn.addEventListener("click", () => {
      you = btn.dataset.labPick;
      const enemy = lab.dataset.enemy || "strike";
      let msg = clash;
      if (you === enemy) msg = clash;
      else if (beats[you] === enemy) msg = win;
      else msg = lose;
      setText(result, `${you} vs ${enemy} → ${msg}`);
      setState(result, msg === win);
    });
  });
  lab.querySelectorAll("[data-lab-enemy]").forEach((btn) => {
    btn.addEventListener("click", () => { lab.dataset.enemy = btn.dataset.enemy; });
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { setText(result, "—"); });
});

document.querySelectorAll(".motion-lab[data-lab='buffer']").forEach((lab) => {
  const bufOut = bind(lab, "[data-lab-buf]");
  const lifeOut = bind(lab, "[data-lab-life]");
  const note = bind(lab, "[data-lab-note]");
  if (!bufOut || !lifeOut) return;
  let buffer = "none", life = 0, frame = 0, move = false;
  const render = () => { setText(bufOut, buffer); setText(lifeOut, String(life)); };
  bind(lab, "[data-lab-press]")?.addEventListener("click", () => { buffer = "light"; life = 8; setText(note, "buffered"); render(); });
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    frame++;
    if (life > 0) life--; else buffer = "none";
    const cancel = move && frame >= 8 && frame <= 13;
    if (cancel && buffer !== "none") { setText(note, "cancel → " + buffer); move = true; frame = 0; buffer = "none"; life = 0; }
    else if (!move && buffer !== "none") { move = true; frame = 1; setText(note, "start " + buffer); buffer = "none"; life = 0; }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { buffer = "none"; life = 0; frame = 0; move = false; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='command']").forEach((lab) => {
  const histOut = bind(lab, "[data-lab-hist]");
  const matchOut = bind(lab, "[data-lab-match]");
  if (!histOut || !matchOut) return;
  let hist = [];
  const render = () => {
    setText(histOut, hist.join("") || "—");
    setText(matchOut, hist.join("") === "↓↘→" ? "HADOKEN" : "—");
    setState(matchOut, hist.join("") === "↓↘→");
  };
  lab.querySelectorAll("[data-lab-dir]").forEach((btn) => {
    btn.addEventListener("click", () => {
      hist.push(btn.dataset.labDir);
      if (hist.length > 3) hist = hist.slice(-3);
      render();
    });
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { hist = []; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='rounds']").forEach((lab) => {
  const p1Out = bind(lab, "[data-lab-p1]");
  const p2Out = bind(lab, "[data-lab-p2]");
  const roundOut = bind(lab, "[data-lab-round]");
  if (!p1Out || !p2Out || !roundOut) return;
  let p1 = 0, p2 = 0, round = 1;
  const render = () => { setText(p1Out, String(p1)); setText(p2Out, String(p2)); setText(roundOut, String(round)); };
  bind(lab, "[data-lab-p1hit]")?.addEventListener("click", () => {
    if (p1 >= 2 || p2 >= 2) return;
    p1++;
    if (p1 < 2 && p2 < 2) round++;
    render();
  });
  bind(lab, "[data-lab-p2hit]")?.addEventListener("click", () => {
    if (p1 >= 2 || p2 >= 2) return;
    p2++;
    if (p1 < 2 && p2 < 2) round++;
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { p1 = 0; p2 = 0; round = 1; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='gravity']").forEach((lab) => {
  const yOut = bind(lab, "[data-lab-y]");
  const vOut = bind(lab, "[data-lab-v]");
  const actor = bind(lab, "[data-lab-actor]");
  if (!yOut || !vOut) return;
  let y = 40, v = 0, g = 0.5;
  const render = () => {
    setText(yOut, y.toFixed(0));
    setText(vOut, v.toFixed(2));
    if (actor) actor.style.top = `${Math.min(85, (y / 480) * 100)}%`;
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    v += g;
    y += v;
    if (y > 400) {
      y = 400;
      v = -v * 0.45;
      if (Math.abs(v) < 0.8) v = 0;
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { y = 40; v = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='loop']").forEach((lab) => {
  const phaseOuts = lab.querySelectorAll("[data-lab-phase]");
  const frameOut = bind(lab, "[data-lab-frame]");
  const noteOuts = lab.querySelectorAll("[data-lab-note]");
  const board = bind(lab, "[data-lab-board]");
  if (!phaseOuts.length || !frameOut) return;
  const updateLabel = lab.dataset.update || "UPDATE";
  const drawLabel = lab.dataset.draw || "DRAW";
  const updateNote = lab.dataset.updateNote || "Change numbers (score, positions).";
  const drawNote = lab.dataset.drawNote || "Paint the current numbers on screen.";
  let frame = 0;
  let phase = "update";
  const render = () => {
    const isUpdate = phase === "update";
    const label = isUpdate ? updateLabel : drawLabel;
    const note = isUpdate ? updateNote : drawNote;
    phaseOuts.forEach((el) => setText(el, label));
    noteOuts.forEach((el) => setText(el, note));
    setText(frameOut, String(frame));
    if (board) board.dataset.phase = phase;
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    if (phase === "update") {
      phase = "draw";
    } else {
      phase = "update";
      frame += 1;
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    frame = 0;
    phase = "update";
    render();
  });
  render();
});

// ---- Visual Effects Lab handlers -------------------------------------------

function fxBoard(lab) {
  return lab.querySelector("[data-lab-board]");
}

document.querySelectorAll(".motion-lab[data-lab='translate']").forEach((lab) => {
  const board = fxBoard(lab);
  const xo = bind(lab, "[data-lab-x]");
  const yo = bind(lab, "[data-lab-y]");
  if (!board || !xo || !yo) return;
  const dot = document.createElement("div");
  dot.className = "lab-fx";
  dot.hidden = true;
  board.appendChild(dot);
  const place = (event) => {
    const rect = board.getBoundingClientRect();
    const cx = event.touches ? event.touches[0].clientX : event.clientX;
    const cy = event.touches ? event.touches[0].clientY : event.clientY;
    const x = Math.max(0, Math.min(480, ((cx - rect.left) / rect.width) * 480));
    const y = Math.max(0, Math.min(520, ((cy - rect.top) / rect.height) * 520));
    dot.style.left = `${(x / 480) * 100}%`;
    dot.style.top = `${(y / 520) * 100}%`;
    dot.hidden = false;
    setText(xo, x.toFixed(0));
    setText(yo, y.toFixed(0));
  };
  board.addEventListener("click", place);
  board.addEventListener("touchstart", (event) => { event.preventDefault(); place(event); }, { passive: false });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { dot.hidden = true; setText(xo, "—"); setText(yo, "—"); });
});

document.querySelectorAll(".motion-lab[data-lab='geom']").forEach((lab) => {
  const board = fxBoard(lab);
  const angleOut = bind(lab, "[data-lab-angle]");
  const scaleOut = bind(lab, "[data-lab-scale]");
  const pivotOut = bind(lab, "[data-lab-pivot]");
  if (!board || !angleOut || !scaleOut) return;
  const dot = document.createElement("div");
  dot.className = "lab-fx";
  board.appendChild(dot);
  let angle = 0, scale = 1, center = true;
  const render = () => {
    dot.style.transformOrigin = center ? "center" : "top left";
    dot.style.transform = `rotate(${angle}deg) scale(${scale})`;
    setText(angleOut, `${angle}°`);
    setText(scaleOut, `x${scale.toFixed(2)}`);
    if (pivotOut) { setText(pivotOut, center ? "center" : "corner"); setState(pivotOut, center); }
  };
  bind(lab, "[data-lab-rotl]")?.addEventListener("click", () => { angle -= 15; render(); });
  bind(lab, "[data-lab-rotr]")?.addEventListener("click", () => { angle += 15; render(); });
  bind(lab, "[data-lab-sdown]")?.addEventListener("click", () => { scale = Math.max(0.4, scale - 0.15); render(); });
  bind(lab, "[data-lab-sup]")?.addEventListener("click", () => { scale = Math.min(2.4, scale + 0.15); render(); });
  bind(lab, "[data-lab-pivot]")?.addEventListener("click", () => { center = !center; render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { angle = 0; scale = 1; center = true; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='colorscale']").forEach((lab) => {
  const board = fxBoard(lab);
  const modeOut = bind(lab, "[data-lab-mode]");
  const codeOut = bind(lab, "[data-lab-code]");
  if (!board || !modeOut) return;
  const dot = document.createElement("div");
  dot.className = "lab-fx";
  board.appendChild(dot);
  const codes = {
    normal: "ColorScale (none)",
    tint: "Scale(1, 0.4, 0.4, 1)",
    flash: "Scale(6, 6, 6, 1)",
    shadow: "Scale(0, 0, 0, 0.5)",
  };
  const set = (mode) => {
    dot.className = "lab-fx" + (mode === "tint" ? " is-tint" : mode === "flash" ? " is-flash" : mode === "shadow" ? " is-shadow" : "");
    setText(modeOut, mode);
    if (codeOut) setText(codeOut, codes[mode]);
  };
  lab.querySelectorAll("[data-lab-mode-set]").forEach((btn) => btn.addEventListener("click", () => set(btn.dataset.labModeSet)));
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => set("normal"));
  set("normal");
});

document.querySelectorAll(".motion-lab[data-lab='opacity']").forEach((lab) => {
  const board = fxBoard(lab);
  const valOut = bind(lab, "[data-lab-alpha]");
  if (!board || !valOut) return;
  const ghosts = [];
  for (let i = 0; i < 5; i++) {
    const g = document.createElement("div");
    g.className = "lab-fx lab-ghost";
    g.style.left = `${20 + i * 15}%`;
    board.appendChild(g);
    ghosts.push(g);
  }
  let alpha = 1, trail = 1;
  const render = () => {
    ghosts.forEach((g, i) => {
      const active = i >= ghosts.length - trail;
      g.style.opacity = active ? String(alpha * ((i + 1) / ghosts.length)) : "0";
    });
    setText(valOut, alpha.toFixed(2));
  };
  bind(lab, "[data-lab-adown]")?.addEventListener("click", () => { alpha = Math.max(0.1, alpha - 0.15); render(); });
  bind(lab, "[data-lab-aup]")?.addEventListener("click", () => { alpha = Math.min(1, alpha + 0.15); render(); });
  bind(lab, "[data-lab-tup]")?.addEventListener("click", () => { trail = Math.min(5, trail + 1); render(); });
  bind(lab, "[data-lab-tdown]")?.addEventListener("click", () => { trail = Math.max(1, trail - 1); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { alpha = 1; trail = 1; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='blend']").forEach((lab) => {
  const board = fxBoard(lab);
  const modeOut = bind(lab, "[data-lab-mode]");
  if (!board) return;
  const a = document.createElement("div");
  a.className = "lab-orb";
  a.style.left = "38%"; a.style.top = "50%"; a.style.background = "radial-gradient(circle,#ff4661,transparent 70%)";
  const b = document.createElement("div");
  b.className = "lab-orb";
  b.style.left = "62%"; b.style.top = "50%"; b.style.background = "radial-gradient(circle,#46e6c8,transparent 70%)";
  board.appendChild(a); board.appendChild(b);
  const set = (mode) => { board.dataset.blend = mode; if (modeOut) setText(modeOut, mode === "add" ? "additive" : "normal"); };
  bind(lab, "[data-lab-toggle]")?.addEventListener("click", () => set(board.dataset.blend === "add" ? "normal" : "add"));
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => set("add"));
  set("add");
});

document.querySelectorAll(".motion-lab[data-lab='sheet']").forEach((lab) => {
  const frameOut = bind(lab, "[data-lab-frame]");
  const cells = lab.querySelectorAll("[data-lab-cell]");
  if (!frameOut || !cells.length) return;
  let frame = 0, timer = null;
  const render = () => {
    setText(frameOut, String(frame));
    cells.forEach((c, i) => c.classList.toggle("on", i === frame));
  };
  const stepFrame = () => { frame = (frame + 1) % cells.length; render(); };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { if (timer) { clearInterval(timer); timer = null; } stepFrame(); });
  bind(lab, "[data-lab-play]")?.addEventListener("click", () => {
    if (timer) { clearInterval(timer); timer = null; return; }
    timer = setInterval(stepFrame, 180);
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { if (timer) { clearInterval(timer); timer = null; } frame = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='spray']").forEach((lab) => {
  const board = fxBoard(lab);
  const countOut = bind(lab, "[data-lab-count]");
  if (!board) return;
  let total = 0;
  const burst = () => {
    for (let i = 0; i < 16; i++) {
      const s = document.createElement("div");
      s.className = "lab-spark";
      s.style.left = "50%";
      s.style.top = "60%";
      board.appendChild(s);
      const angle = Math.random() * Math.PI * 2;
      const dist = 40 + Math.random() * 120;
      requestAnimationFrame(() => {
        s.style.transform = `translate(${Math.cos(angle) * dist}px, ${Math.sin(angle) * dist}px)`;
        s.style.opacity = "0";
      });
      setTimeout(() => s.remove(), 700);
    }
    total += 16;
    if (countOut) setText(countOut, String(total));
  };
  bind(lab, "[data-lab-burst]")?.addEventListener("click", burst);
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { total = 0; if (countOut) setText(countOut, "0"); });
});

document.querySelectorAll(".motion-lab[data-lab='spellbook']").forEach((lab) => {
  const doneOut = bind(lab, "[data-lab-done]");
  const board = bind(lab, "[data-lab-board]");
  const cast = new Set();
  // Lesson pages live at web/{ja,en}/tracks/visual-effects/<slug>/
  const src = {
    fire: "../../../../assets/vfx-fire.png",
    water: "../../../../assets/vfx-water.png",
    spark: "../../../../assets/vfx-spark.png",
    bolt: "../../../../assets/vfx-bolt.png",
  };
  const recipes = {
    fire: { ja: "炎 = 炎スプライト + 加算 + 上昇", en: "FIRE = flame PNG + additive + rise" },
    water: { ja: "水 = 水滴スプライト + 半透明 + 重力", en: "WATER = drop PNG + alpha + gravity" },
    thunder: { ja: "雷 = 稲妻スプライト + 閃光 + 粒", en: "THUNDER = bolt PNG + flash + sparks" },
  };
  const lang = (document.documentElement.lang || "ja").startsWith("en") ? "en" : "ja";
  const idleHTML = board ? board.innerHTML : "";
  let animTimer = null;

  const render = () => {
    if (doneOut) {
      setText(doneOut, `${cast.size}/3`);
      setState(doneOut, cast.size === 3);
    }
  };

  const clearAnim = () => {
    if (animTimer) {
      clearInterval(animTimer);
      animTimer = null;
    }
  };

  const spawnFire = (stage) => {
    for (let i = 0; i < 16; i++) {
      const el = document.createElement("img");
      el.src = src.fire;
      el.alt = "";
      el.className = "lab-vfx-sprite lab-vfx-flame";
      el.style.left = `${42 + Math.random() * 16}%`;
      el.style.setProperty("--drift", `${(Math.random() - 0.5) * 60}px`);
      el.style.setProperty("--rise", `${140 + Math.random() * 120}px`);
      el.style.setProperty("--size", `${48 + Math.random() * 56}px`);
      el.style.animationDelay = `${Math.random() * 0.5}s`;
      stage.appendChild(el);
    }
    for (let i = 0; i < 18; i++) {
      const el = document.createElement("img");
      el.src = src.spark;
      el.alt = "";
      el.className = "lab-vfx-sprite lab-vfx-ember";
      el.style.left = `${40 + Math.random() * 20}%`;
      el.style.setProperty("--drift", `${(Math.random() - 0.5) * 90}px`);
      el.style.setProperty("--rise", `${100 + Math.random() * 140}px`);
      el.style.animationDelay = `${Math.random() * 0.6}s`;
      stage.appendChild(el);
    }
  };

  const spawnWater = (stage) => {
    for (let i = 0; i < 20; i++) {
      const el = document.createElement("img");
      el.src = src.water;
      el.alt = "";
      el.className = "lab-vfx-sprite lab-vfx-drop";
      el.style.left = `${35 + Math.random() * 30}%`;
      el.style.setProperty("--arc-x", `${(Math.random() - 0.5) * 120}px`);
      el.style.setProperty("--fall", `${110 + Math.random() * 90}px`);
      el.style.setProperty("--size", `${28 + Math.random() * 34}px`);
      el.style.animationDelay = `${Math.random() * 0.45}s`;
      stage.appendChild(el);
    }
  };

  const spawnThunder = (stage) => {
    stage.classList.add("is-flash");
    setTimeout(() => stage.classList.remove("is-flash"), 180);
    for (let i = 0; i < 5; i++) {
      const el = document.createElement("img");
      el.src = src.bolt;
      el.alt = "";
      el.className = "lab-vfx-sprite lab-vfx-bolt";
      el.style.left = `${30 + i * 10 + Math.random() * 8}%`;
      el.style.setProperty("--rot", `${(Math.random() - 0.5) * 28}deg`);
      el.style.animationDelay = `${i * 0.04}s`;
      stage.appendChild(el);
    }
    for (let i = 0; i < 22; i++) {
      const el = document.createElement("img");
      el.src = src.spark;
      el.alt = "";
      el.className = "lab-vfx-sprite lab-vfx-zap";
      el.style.left = `${45 + (Math.random() - 0.5) * 30}%`;
      el.style.top = `${30 + Math.random() * 40}%`;
      el.style.setProperty("--dx", `${(Math.random() - 0.5) * 140}px`);
      el.style.setProperty("--dy", `${(Math.random() - 0.5) * 140}px`);
      stage.appendChild(el);
    }
  };

  const showSpell = (kind) => {
    if (!board) return;
    clearAnim();
    board.className = `lab-board lab-spell-stage is-${kind}`;
    board.innerHTML = "";
    const label = document.createElement("p");
    label.className = "lab-spell-label";
    label.textContent = recipes[kind][lang];
    board.appendChild(label);
    if (kind === "fire") spawnFire(board);
    else if (kind === "water") spawnWater(board);
    else spawnThunder(board);
    if (kind === "fire") {
      animTimer = setInterval(() => {
        if (!board.isConnected) return clearAnim();
        const el = document.createElement("img");
        el.src = src.fire;
        el.alt = "";
        el.className = "lab-vfx-sprite lab-vfx-flame";
        el.style.left = `${44 + Math.random() * 12}%`;
        el.style.setProperty("--drift", `${(Math.random() - 0.5) * 50}px`);
        el.style.setProperty("--rise", `${150 + Math.random() * 100}px`);
        el.style.setProperty("--size", `${40 + Math.random() * 50}px`);
        board.appendChild(el);
        setTimeout(() => el.remove(), 1400);
      }, 220);
    }
  };

  lab.querySelectorAll("[data-lab-spell]").forEach((btn) => {
    btn.addEventListener("click", () => {
      const kind = btn.dataset.labSpell;
      cast.add(kind);
      btn.dataset.cast = "1";
      showSpell(kind);
      render();
    });
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    cast.clear();
    clearAnim();
    lab.querySelectorAll("[data-lab-spell]").forEach((btn) => (btn.dataset.cast = "0"));
    if (board) {
      board.className = "lab-board lab-spell-stage";
      board.innerHTML = idleHTML;
    }
    render();
  });
  render();
});

// Advanced magic showcase previews (STEPS 09–13).
document.querySelectorAll(".motion-lab[data-lab^='magic-']").forEach((lab) => {
  const kind = (lab.dataset.lab || "").replace("magic-", "");
  const board = bind(lab, "[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const idleHTML = board ? board.innerHTML : "";
  const asset = {
    fire: "../../../../assets/vfx-fire.png",
    ice: "../../../../assets/vfx-ice.png",
    thunder: "../../../../assets/vfx-bolt.png",
    light: "../../../../assets/vfx-light.png",
    dark: "../../../../assets/vfx-dark.png",
    spark: "../../../../assets/vfx-spark.png",
    ring: "../../../../assets/vfx-ring.png",
  };
  let total = 0;

  const sprite = (src, cls, style = {}) => {
    const el = document.createElement("img");
    el.src = src;
    el.alt = "";
    el.className = `lab-vfx-sprite ${cls}`;
    Object.entries(style).forEach(([k, v]) => {
      if (k.startsWith("--") || k === "left" || k === "top") el.style.setProperty(k, v);
      else el.style[k] = v;
    });
    return el;
  };

  const burst = () => {
    if (!board) return;
    board.className = `lab-board lab-spell-stage lab-magic-stage is-${kind}`;
    board.innerHTML = "";
    const n = 18 + Math.floor(Math.random() * 12);
    total += n;
    if (countOut) setText(countOut, String(total));
    if (kind === "fire") {
      for (let i = 0; i < n; i++) {
        board.appendChild(sprite(asset.fire, "lab-vfx-flame", {
          left: `${40 + Math.random() * 20}%`,
          "--drift": `${(Math.random() - 0.5) * 70}px`,
          "--rise": `${130 + Math.random() * 130}px`,
          "--size": `${40 + Math.random() * 60}px`,
        }));
      }
      board.appendChild(sprite(asset.ring, "lab-vfx-ring", { left: "50%", top: "70%" }));
    } else if (kind === "ice") {
      for (let i = 0; i < n; i++) {
        board.appendChild(sprite(asset.ice, "lab-vfx-ice", {
          left: "50%",
          top: "45%",
          "--dx": `${(Math.random() - 0.5) * 220}px`,
          "--dy": `${(Math.random() - 0.5) * 180}px`,
          "--rot": `${(Math.random() - 0.5) * 120}deg`,
          "--size": `${36 + Math.random() * 48}px`,
        }));
      }
    } else if (kind === "thunder") {
      board.classList.add("is-flash");
      setTimeout(() => board.classList.remove("is-flash"), 180);
      for (let i = 0; i < 6; i++) {
        board.appendChild(sprite(asset.thunder, "lab-vfx-bolt", {
          left: `${28 + i * 9}%`,
          "--rot": `${(Math.random() - 0.5) * 30}deg`,
        }));
      }
      for (let i = 0; i < 16; i++) {
        board.appendChild(sprite(asset.spark, "lab-vfx-zap", {
          left: `${45 + (Math.random() - 0.5) * 30}%`,
          top: `${30 + Math.random() * 40}%`,
          "--dx": `${(Math.random() - 0.5) * 140}px`,
          "--dy": `${(Math.random() - 0.5) * 140}px`,
        }));
      }
    } else if (kind === "light") {
      board.appendChild(sprite(asset.light, "lab-vfx-flare", { left: "50%", top: "45%" }));
      board.appendChild(sprite(asset.ring, "lab-vfx-ring", { left: "50%", top: "45%" }));
      for (let i = 0; i < n; i++) {
        board.appendChild(sprite(asset.spark, "lab-vfx-zap", {
          left: "50%",
          top: "45%",
          "--dx": `${(Math.random() - 0.5) * 200}px`,
          "--dy": `${(Math.random() - 0.5) * 200}px`,
        }));
      }
    } else {
      for (let i = 0; i < n; i++) {
        board.appendChild(sprite(asset.dark, "lab-vfx-dark", {
          left: "50%",
          top: "48%",
          "--ang": `${(i / n) * 360}deg`,
          "--rad": `${40 + Math.random() * 90}px`,
          "--size": `${32 + Math.random() * 40}px`,
        }));
      }
      board.appendChild(sprite(asset.ring, "lab-vfx-ring", { left: "50%", top: "48%" }));
    }
  };

  bind(lab, "[data-lab-burst]")?.addEventListener("click", burst);
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    total = 0;
    if (countOut) setText(countOut, "0");
    if (board) {
      board.className = "lab-board lab-spell-stage lab-magic-stage";
      board.innerHTML = idleHTML;
    }
  });
});

document.querySelectorAll(".motion-lab[data-lab='fx-split']").forEach((lab) => {
  const playList = bind(lab, "[data-lab-play-list]");
  const fxList = bind(lab, "[data-lab-fx-list]");
  const playN = bind(lab, "[data-lab-play-n]");
  const fxN = bind(lab, "[data-lab-fx-n]");
  const scoreOut = bind(lab, "[data-lab-score]");
  const stage = bind(lab, "[data-lab-fx-stage]");
  let score = 0;
  let particles = [];
  let pid = 0;
  const render = () => {
    if (playList) {
      playList.innerHTML = [
        `<li>score = ${score}</li>`,
        "<li>player</li>",
        "<li>target</li>",
      ].join("");
    }
    if (fxList) {
      fxList.innerHTML = particles.length
        ? particles.map((p) => `<li>${p.name} life=${p.life}</li>`).join("")
        : `<li class="lab-fx-empty">${lab.closest("[lang]")?.lang === "en" ? "(empty)" : "（空）"}</li>`;
    }
    if (playN) setText(playN, String(3));
    if (fxN) setText(fxN, String(particles.length));
    if (scoreOut) setText(scoreOut, String(score));
    if (stage) {
      stage.innerHTML = particles.map((p) => {
        const left = 20 + (p.id * 17) % 60;
        const top = 25 + (p.id * 13) % 50;
        return `<span class="lab-fx-dot" style="left:${left}%;top:${top}%;opacity:${Math.max(0.25, p.life / 8)}"></span>`;
      }).join("");
    }
  };
  const ja = document.documentElement.lang === "ja";
  bind(lab, "[data-lab-fx-ping]")?.addEventListener("click", () => {
    score += 1;
    for (let i = 0; i < 4; i++) {
      pid += 1;
      particles.push({ id: pid, name: `spark#${pid}`, life: 8 });
    }
    if (particles.length > 16) particles = particles.slice(-16);
    render();
  });
  bind(lab, "[data-lab-fx-tick]")?.addEventListener("click", () => {
    particles = particles
      .map((p) => ({ ...p, life: p.life - 1 }))
      .filter((p) => p.life > 0);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    score = 0;
    particles = [];
    pid = 0;
    render();
  });
  const hint = bind(lab, "[data-lab-fx-hint]");
  if (hint) {
    setText(hint, ja
      ? "「命中！」→ score と spark が増える。「fx 1F」→ spark だけ減る（score はそのまま）。"
      : "Hit! grows score + sparks. Tick FX ages sparks only — score stays.");
  }
  render();
});

document.querySelectorAll(".motion-lab[data-lab='fx-meter-grade']").forEach((lab) => {
  const position = bind(lab, "[data-lab-meter-position]");
  const marker = bind(lab, "[data-lab-grade-marker]");
  const resultOut = bind(lab, "[data-lab-grade-result]");
  const ruleOut = bind(lab, "[data-lab-grade-rule]");
  const strengthOut = bind(lab, "[data-lab-grade-strength]");
  const distanceOut = bind(lab, "[data-lab-grade-distance]");
  const particlesOut = bind(lab, "[data-lab-grade-particles]");
  const burst = bind(lab, "[data-lab-grade-burst]");
  const explain = bind(lab, "[data-lab-grade-explain]");
  const ja = document.documentElement.lang === "ja";

  const moveMarker = () => {
    if (marker && position) marker.style.left = `${Number(position.value)}%`;
  };
  const judge = () => {
    const stoppedAt = Number(position?.value || 50);
    const distance = Math.abs(stoppedAt - 50);
    const grade = distance <= 5 ? "PERFECT" : distance <= 18 ? "OK" : "MISS";
    const score = grade === "PERFECT" ? 100 : grade === "OK" ? 40 : 0;
    const particles = grade === "PERFECT" ? 18 : grade === "OK" ? 8 : 2;
    const strength = grade === "PERFECT" ? (ja ? "大花火" : "BIG BURST") : grade === "OK" ? (ja ? "中花火" : "MEDIUM") : (ja ? "小さな煙" : "TINY PUFF");

    setText(resultOut, `${grade} / +${score}`);
    setText(ruleOut, ja ? `|${stoppedAt} − 50| = ${distance} → ${grade}` : `|${stoppedAt} − 50| = ${distance} → ${grade}`);
    setText(strengthOut, strength);
    setText(distanceOut, String(distance));
    setText(particlesOut, String(particles));
    setText(explain, ja
      ? `play が ${grade} を決定 → fx は「${particles}粒」として見せる。fx は位置を判定し直しません。`
      : `play decides ${grade} → fx presents it as ${particles} particles. fx does not judge the position again.`);
    if (burst) {
      burst.className = `lab-grade-burst is-${grade.toLowerCase()}`;
      burst.innerHTML = Array.from({ length: particles }, (_, i) =>
        `<span style="--r:${18 + i * 2}px;--a:${(i * 137.5).toFixed(1)}deg"></span>`).join("");
    }
  };
  const reset = () => {
    if (position) position.value = "50";
    moveMarker();
    setText(resultOut, "—");
    setText(ruleOut, ja ? "中心との差を計算" : "measure center distance");
    setText(strengthOut, "—");
    setText(distanceOut, "—");
    setText(particlesOut, "—");
    setText(explain, ja ? "マーカーを動かして「判定！」を押そう。" : "Move the marker, then press Judge.");
    if (burst) { burst.className = "lab-grade-burst"; burst.innerHTML = ""; }
  };

  position?.addEventListener("input", moveMarker);
  bind(lab, "[data-lab-meter-judge]")?.addEventListener("click", judge);
  bind(lab, "[data-lab-meter-miss]")?.addEventListener("click", () => {
    if (position) position.value = "15";
    moveMarker();
    judge();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", reset);
  reset();
});

const escapeCodeHTML = (value) => value.replace(/[&<>"']/g, (character) => ({
  "&": "&amp;",
  "<": "&lt;",
  ">": "&gt;",
  '"': "&quot;",
  "'": "&#39;",
}[character]));

const highlightCode = (code) => {
  const source = code.textContent || "";
  const tokenPattern = /(\/\/[^\n]*|\/\*[\s\S]*?\*\/|"(?:\\.|[^"\\])*"|`[\s\S]*?`|'(?:\\.|[^'\\])*'|\b\d+(?:\.\d+)?\b|\b(?:package|import|func|type|struct|interface|const|var|if|else|for|range|switch|case|default|return|go|defer|select|break|continue|fallthrough)\b|\b(?:true|false|nil)\b|\b(?:bool|string|byte|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|complex64|complex128|error)\b|\b[A-Z][A-Za-z0-9_]*\b)/g;
  let html = "";
  let cursor = 0;

  for (const match of source.matchAll(tokenPattern)) {
    const token = match[0];
    const start = match.index ?? 0;
    html += escapeCodeHTML(source.slice(cursor, start));
    let kind = "code-number";
    if (token.startsWith("//") || token.startsWith("/*")) kind = "code-comment";
    else if (/^["'`]/.test(token)) kind = "code-string";
    else if (/^(package|import|func|type|struct|interface|const|var|if|else|for|range|switch|case|default|return|go|defer|select|break|continue|fallthrough)$/.test(token)) kind = "code-keyword";
    else if (/^(true|false|nil)$/.test(token)) kind = "code-constant";
    else if (/^(bool|string|byte|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|complex64|complex128|error)$/.test(token)) kind = "code-type";
    else if (/^[A-Z]/.test(token)) kind = "code-type";
    else if (/^\s*\(/.test(source.slice(start + token.length))) kind = "code-function";
    html += `<span class="${kind}">${escapeCodeHTML(token)}</span>`;
    cursor = start + token.length;
  }
  code.innerHTML = html + escapeCodeHTML(source.slice(cursor));
};

document.querySelectorAll("pre code").forEach(highlightCode);

document.querySelectorAll(".feedback-form").forEach((form) => {
  const button = form.querySelector(".feedback-submit");
  const status = form.querySelector(".feedback-status");
  const message = form.querySelector(".feedback-message");
  if (!button || !status || !message) return;
  const syncFeedbackButton = () => {
    const empty = !message.value.trim();
    button.disabled = empty;
    button.setAttribute("aria-disabled", String(empty));
  };
  message.addEventListener("input", syncFeedbackButton);
  syncFeedbackButton();
  form.addEventListener("submit", async (event) => {
    event.preventDefault();
    button.disabled = true;
    status.textContent = message.dataset.sending || "送信中…";
    status.classList.remove("is-sent");
    try {
      await fetch(form.action, {
        method: "POST",
        body: new FormData(form),
        mode: "no-cors",
        credentials: "omit",
      });
      message.value = "";
      status.textContent = message.dataset.sent || "送信しました。ありがとうございます！";
      status.classList.add("is-sent");
    } catch {
      status.textContent = message.dataset.failed || "送信できませんでした。時間をおいて再試行してください。";
    } finally {
      syncFeedbackButton();
    }
  });
});

document.querySelectorAll(".full-code").forEach((block) => {
  const button = block.querySelector("[data-copy]");
  const code = block.querySelector("[data-embed-slot], pre code");
  if (!button || !code) return;
  const idle = button.textContent.trim() || "Copy";
  const done = button.dataset.copiedLabel || "Copied!";
  button.addEventListener("click", async () => {
    const text = code.textContent || "";
    try {
      await navigator.clipboard.writeText(text);
    } catch {
      const area = document.createElement("textarea");
      area.value = text;
      area.setAttribute("readonly", "");
      area.style.position = "fixed";
      area.style.left = "-9999px";
      document.body.appendChild(area);
      area.select();
      document.execCommand("copy");
      area.remove();
    }
    button.dataset.copied = "1";
    button.textContent = done;
    window.setTimeout(() => {
      button.dataset.copied = "0";
      button.textContent = idle;
    }, 1600);
  });
});

/* --- Concept labs for thin track articles (enrichment) --- */

document.querySelectorAll(".motion-lab[data-lab='drop-timer']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const timerOut = bind(lab, "[data-lab-timer]");
  const yOut = bind(lab, "[data-lab-y]");
  const note = bind(lab, "[data-lab-note]");
  if (!board || !timerOut) return;
  const need = Number(lab.dataset.need || 8);
  const rows = 8;
  let timer = 0;
  let y = 0;
  let locked = Array.from({ length: rows }, () => false);
  const render = () => {
    setText(timerOut, `${timer}/${need}`);
    if (yOut) setText(yOut, String(y));
    const fill = Math.min(100, (timer / need) * 100);
    board.innerHTML = `<div class="lab-drop-bar"><i style="width:${fill}%"></i></div><div class="lab-drop-grid">${locked.map((on, i) => {
      const cls = on ? "locked" : i === y ? "active" : "";
      return `<span class="${cls}"></span>`;
    }).join("")}</div>`;
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    timer += 1;
    if (timer >= need) {
      timer = 0;
      if (y + 1 >= rows || locked[y + 1]) {
        locked[y] = true;
        setText(note, lab.dataset.lock || "LOCK");
        y = 0;
        if (locked[0]) setText(note, lab.dataset.top || "TOP BLOCKED");
      } else {
        y += 1;
        setText(note, lab.dataset.drop || "DROP");
      }
    } else {
      setText(note, lab.dataset.wait || "wait…");
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    timer = 0; y = 0; locked = Array.from({ length: rows }, () => false);
    setText(note, "—");
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='card-play']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const energyOut = bind(lab, "[data-lab-energy]");
  const hpOut = bind(lab, "[data-lab-hp]");
  const enemyOut = bind(lab, "[data-lab-enemy]");
  const note = bind(lab, "[data-lab-note]");
  if (!energyOut || !enemyOut) return;
  let energy = 3, hp = 40, enemy = 60, block = 0;
  const cards = [
    { id: "damage", cost: 2, value: 12 },
    { id: "block", cost: 1, value: 8 },
    { id: "heal", cost: 2, value: 10 },
  ];
  const render = () => {
    setText(energyOut, String(energy));
    if (hpOut) setText(hpOut, String(hp));
    setText(enemyOut, String(enemy));
    if (board) {
      board.innerHTML = `<div class="lab-card-stage"><div class="lab-card-you"><b>YOU</b><span>HP ${hp}</span><span>BLK ${block}</span></div><div class="lab-card-foe"><b>FOE</b><span>HP ${enemy}</span></div></div><div class="lab-card-energy">⚡ ${energy}</div>`;
    }
  };
  const play = (kind) => {
    const c = cards.find((x) => x.id === kind);
    if (!c) return;
    if (energy < c.cost) {
      setText(note, lab.dataset.poor || "not enough energy");
      return;
    }
    energy -= c.cost;
    if (kind === "damage") { enemy = Math.max(0, enemy - c.value); setText(note, `-${c.value} HP`); }
    if (kind === "block") { block += c.value; setText(note, `+${c.value} block`); }
    if (kind === "heal") { hp = Math.min(40, hp + c.value); setText(note, `+heal`); }
    render();
  };
  lab.querySelectorAll("[data-lab-card]").forEach((btn) => {
    btn.addEventListener("click", () => play(btn.dataset.labCard));
  });
  bind(lab, "[data-lab-end]")?.addEventListener("click", () => {
    const hit = Math.max(0, 7 - block);
    hp = Math.max(0, hp - hit);
    block = 0;
    energy = 3;
    setText(note, lab.dataset.enemy || `enemy hits ${hit}`);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    energy = 3; hp = 40; enemy = 60; block = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='friction']").forEach((lab) => {
  const xOut = bind(lab, "[data-lab-x]");
  const speedOut = bind(lab, "[data-lab-speed]");
  const stateOut = bind(lab, "[data-lab-state]");
  if (!xOut || !speedOut) return;
  const friction = Number(lab.dataset.friction || 0.92);
  const stopAt = Number(lab.dataset.stop || 0.35);
  let x = 40, vx = 8;
  const moving = lab.dataset.moving || "moving";
  const stopped = lab.dataset.stopped || "stopped";
  const render = () => {
    const speed = Math.abs(vx);
    setText(xOut, x.toFixed(0));
    setText(speedOut, speed.toFixed(2));
    if (stateOut) setText(stateOut, speed < stopAt ? stopped : moving);
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    x += vx;
    vx *= friction;
    if (Math.abs(vx) < stopAt) vx = 0;
    if (x > 220) x = 40;
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    x = 40;
    vx = 8;
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='energy-turn'], .motion-lab[data-lab='deck-cycle']").forEach((lab) => {
  const energyOut = bind(lab, "[data-lab-energy]");
  const turnOut = bind(lab, "[data-lab-turn]");
  const spentOut = bind(lab, "[data-lab-spent]");
  const board = lab.querySelector("[data-lab-board]");
  const note = bind(lab, "[data-lab-note]");
  if (!energyOut || !turnOut) return;
  const max = Number(lab.dataset.max || 3);
  let energy = max, turn = 1, spent = 0;
  const render = () => {
    setText(energyOut, `${energy}/${max}`);
    setText(turnOut, String(turn));
    if (spentOut) setText(spentOut, String(spent));
    if (board) {
      board.innerHTML = `<div class="lab-energy-pips">${Array.from({ length: max }, (_, i) =>
        `<span class="${i < energy ? "on" : ""}"></span>`).join("")}</div><p class="lab-energy-label">turn ${turn}</p>`;
    }
  };
  bind(lab, "[data-lab-spend]")?.addEventListener("click", () => {
    if (energy <= 0) { setText(note, lab.dataset.empty || "0 energy"); return; }
    energy -= 1; spent += 1;
    setText(note, lab.dataset.spend || "played 1 cost");
    render();
  });
  bind(lab, "[data-lab-end]")?.addEventListener("click", () => {
    turn += 1; energy = max; spent = 0;
    setText(note, lab.dataset.refill || "energy refilled");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    energy = max; turn = 1; spent = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='merge-same']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const tierOut = bind(lab, "[data-lab-tier]");
  const note = bind(lab, "[data-lab-note]");
  if (!board) return;
  let left = 1, right = 1;
  const render = () => {
    if (tierOut) setText(tierOut, `${left} + ${right}`);
    board.innerHTML = `<div class="lab-merge-stage"><span class="lab-fruit t${left}" style="--s:${16 + left * 8}px">T${left}</span><span class="lab-merge-plus">+</span><span class="lab-fruit t${right}" style="--s:${16 + right * 8}px">T${right}</span></div>`;
  };
  bind(lab, "[data-lab-bump-left]")?.addEventListener("click", () => { left = Math.min(5, left + 1); render(); });
  bind(lab, "[data-lab-bump-right]")?.addEventListener("click", () => { right = Math.min(5, right + 1); render(); });
  bind(lab, "[data-lab-merge]")?.addEventListener("click", () => {
    if (left !== right) { setText(note, lab.dataset.mismatch || "tiers differ — no merge"); return; }
    const next = Math.min(6, left + 1);
    setText(note, (lab.dataset.merge || "merged → T{n}").replace("{n}", String(next)));
    left = next; right = Math.max(1, next - 1);
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { left = 1; right = 1; setText(note, "—"); render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='match-scan']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  if (!board) return;
  // 5x5 simple grid of colors 0-3
  let grid = [
    [0, 0, 0, 1, 2],
    [1, 2, 1, 1, 1],
    [2, 2, 3, 0, 0],
    [3, 1, 1, 1, 2],
    [0, 0, 2, 2, 2],
  ];
  let marks = new Set();
  const key = (x, y) => `${x},${y}`;
  const scan = () => {
    marks = new Set();
    const h = grid.length, w = grid[0].length;
    for (let y = 0; y < h; y++) {
      let x = 0;
      while (x < w) {
        let n = 1;
        while (x + n < w && grid[y][x + n] === grid[y][x]) n++;
        if (n >= 3) for (let i = 0; i < n; i++) marks.add(key(x + i, y));
        x += n;
      }
    }
    for (let x = 0; x < w; x++) {
      let y = 0;
      while (y < h) {
        let n = 1;
        while (y + n < h && grid[y + n][x] === grid[y][x]) n++;
        if (n >= 3) for (let i = 0; i < n; i++) marks.add(key(x, y + i));
        y += n;
      }
    }
  };
  const render = () => {
    if (countOut) setText(countOut, String(marks.size));
    board.innerHTML = `<div class="lab-match-grid">${grid.map((row, y) => row.map((c, x) =>
      `<span class="c${c}${marks.has(key(x, y)) ? " hit" : ""}"></span>`).join("")).join("")}</div>`;
  };
  bind(lab, "[data-lab-scan]")?.addEventListener("click", () => {
    scan();
    setText(note, marks.size ? (lab.dataset.found || "match found") : (lab.dataset.none || "no match"));
    render();
  });
  bind(lab, "[data-lab-clear]")?.addEventListener("click", () => {
    if (!marks.size) scan();
    if (!marks.size) { setText(note, lab.dataset.none || "no match"); render(); return; }
    for (const k of marks) {
      const [x, y] = k.split(",").map(Number);
      grid[y][x] = -1;
    }
    // gravity
    for (let x = 0; x < 5; x++) {
      const col = [];
      for (let y = 4; y >= 0; y--) if (grid[y][x] >= 0) col.push(grid[y][x]);
      for (let y = 4; y >= 0; y--) grid[y][x] = col[4 - y] ?? Math.floor(Math.random() * 4);
    }
    marks = new Set();
    setText(note, lab.dataset.cleared || "cleared + gravity");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    grid = [
      [0, 0, 0, 1, 2],
      [1, 2, 1, 1, 1],
      [2, 2, 3, 0, 0],
      [3, 1, 1, 1, 2],
      [0, 0, 2, 2, 2],
    ];
    marks = new Set();
    setText(note, "—");
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='line-clear']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const clearedOut = bind(lab, "[data-lab-cleared]");
  const note = bind(lab, "[data-lab-note]");
  if (!board) return;
  let rows = [
    [1, 1, 1, 1, 1, 1],
    [0, 1, 1, 0, 1, 0],
    [1, 1, 1, 1, 1, 1],
    [1, 0, 0, 1, 0, 1],
  ];
  let cleared = 0;
  const render = () => {
    if (clearedOut) setText(clearedOut, String(cleared));
    board.innerHTML = `<div class="lab-line-grid">${rows.map((row) => {
      const full = row.every(Boolean);
      return `<div class="lab-line-row${full ? " full" : ""}">${row.map((c) => `<span class="${c ? "on" : ""}"></span>`).join("")}</div>`;
    }).join("")}</div>`;
  };
  bind(lab, "[data-lab-clear]")?.addEventListener("click", () => {
    const keep = [];
    let n = 0;
    for (const row of rows) {
      if (row.every(Boolean)) n += 1;
      else keep.push(row);
    }
    while (keep.length < rows.length) keep.unshift(Array(6).fill(0));
    rows = keep;
    cleared += n;
    setText(note, n ? (lab.dataset.cleared || `cleared ${n}`) : (lab.dataset.none || "no full rows"));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    rows = [
      [1, 1, 1, 1, 1, 1],
      [0, 1, 1, 0, 1, 0],
      [1, 1, 1, 1, 1, 1],
      [1, 0, 0, 1, 0, 1],
    ];
    cleared = 0; setText(note, "—"); render();
  });
  render();
});

// Enhance plain tile labs with a mini board when [data-lab-board] is present.
document.querySelectorAll(".motion-lab[data-lab='tile']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  if (!board || board.dataset.enhanced) return;
  board.dataset.enhanced = "1";
  const posOut = bind(lab, "[data-lab-pos]");
  const paint = () => {
    const raw = posOut?.textContent || "1,1";
    const [px, py] = raw.split(",").map(Number);
    const cells = [];
    for (let y = 0; y < 3; y++) for (let x = 0; x < 3; x++) {
      const wall = x === 2 && y === 1;
      const here = x === px && y === py;
      cells.push(`<span class="${wall ? "wall" : here ? "you" : "path"}"></span>`);
    }
    board.innerHTML = `<div class="lab-tile-grid">${cells.join("")}</div>`;
  };
  // 値の変更はボタン操作から明示的に再描画する。MutationObserverに
  // 頼らないので、iframeやモバイルブラウザでも同じ挙動になる。
  lab.querySelector(".lab-controls")?.addEventListener("click", () => requestAnimationFrame(paint));
  paint();
});

/* --- Wipe entities: concept-specific replacements --- */

document.querySelectorAll(".motion-lab[data-lab='bullets']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  if (!countOut) return;
  let shots = [];
  let id = 1;
  const render = () => {
    setText(countOut, String(shots.length));
    if (board) {
      board.innerHTML = `<div class="lab-bullets-stage"><span class="lab-bullets-ship"></span>${shots.map((s) =>
        `<span class="lab-bullet" style="bottom:${12 + (400 - s.y) * 0.55}%"></span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-fire]")?.addEventListener("click", () => {
    shots.push({ id: id++, y: 400 });
    setText(note, lab.dataset.fire || "pew");
    render();
  });
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    shots = shots.map((s) => ({ ...s, y: s.y - 50 })).filter((s) => s.y > 0);
    setText(note, lab.dataset.step || "fly");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    shots = []; id = 1; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='preview-next']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const nextOut = bind(lab, "[data-lab-next]");
  const note = bind(lab, "[data-lab-note]");
  let queue = [1, 2, 1];
  const render = () => {
    if (nextOut) setText(nextOut, queue[0] ? `T${queue[0]}` : "—");
    if (board) {
      board.innerHTML = `<div class="lab-preview-row">${queue.map((t, i) =>
        `<span class="lab-fruit t${t}${i === 0 ? " next" : ""}" style="--s:${20 + t * 6}px">T${t}</span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-enqueue]")?.addEventListener("click", () => {
    queue.push(1 + Math.floor(Math.random() * 3));
    if (queue.length > 5) queue = queue.slice(-5);
    setText(note, lab.dataset.enqueued || "queued");
    render();
  });
  bind(lab, "[data-lab-drop]")?.addEventListener("click", () => {
    if (!queue.length) { setText(note, lab.dataset.empty || "empty"); return; }
    const t = queue.shift();
    setText(note, (lab.dataset.dropped || "drop T{n}").replace("{n}", String(t)));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    queue = [1, 2, 1]; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='status-ticks']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  let effects = [];
  const render = () => {
    if (countOut) setText(countOut, String(effects.length));
    if (board) {
      board.innerHTML = effects.length
        ? `<div class="lab-status-row">${effects.map((e) =>
          `<span class="lab-status-chip">${e.name}<b>${e.left}</b></span>`).join("")}</div>`
        : `<p class="lab-empty">${lab.dataset.none || "no status"}</p>`;
    }
  };
  bind(lab, "[data-lab-add]")?.addEventListener("click", () => {
    effects.push({ name: lab.dataset.status || "POISON", left: 3 });
    setText(note, lab.dataset.added || "applied");
    render();
  });
  bind(lab, "[data-lab-tick]")?.addEventListener("click", () => {
    effects = effects.map((e) => ({ ...e, left: e.left - 1 })).filter((e) => e.left > 0);
    setText(note, lab.dataset.ticked || "turn passed");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    effects = []; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='deck-pick']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const deckOut = bind(lab, "[data-lab-deck]");
  const note = bind(lab, "[data-lab-note]");
  const pool = (lab.dataset.cards || "Strike,Guard,Heal").split(",");
  let deck = ["Strike", "Strike", "Guard"];
  const render = () => {
    if (deckOut) setText(deckOut, String(deck.length));
    if (board) {
      board.innerHTML = `<div class="lab-deck-picks">${pool.map((c, i) =>
        `<button type="button" class="lab-deck-card" data-pick="${i}">${c}</button>`).join("")}</div>
        <p class="lab-deck-list">${deck.join(" · ")}</p>`;
      board.querySelectorAll("[data-pick]").forEach((btn) => {
        btn.addEventListener("click", () => {
          const card = pool[Number(btn.dataset.pick)];
          deck.push(card);
          setText(note, (lab.dataset.picked || "added {c}").replace("{c}", card));
          render();
        });
      });
    }
  };
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    deck = ["Strike", "Strike", "Guard"]; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='map-nodes']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const pathOut = bind(lab, "[data-lab-path]");
  const note = bind(lab, "[data-lab-note]");
  let path = ["start"];
  let node = "start";
  const edges = { start: ["fight", "shop"], fight: ["elite", "rest"], shop: ["rest"], elite: ["boss"], rest: ["boss"], boss: [] };
  const render = () => {
    if (pathOut) setText(pathOut, path.join(" → "));
    const next = edges[node] || [];
    if (board) {
      board.innerHTML = `<div class="lab-map-now">${node}</div><div class="lab-map-choices">${next.map((n) =>
        `<button type="button" class="lab-button" data-goto="${n}">→ ${n}</button>`).join("") || `<span class="lab-empty">END</span>`}</div>`;
      board.querySelectorAll("[data-goto]").forEach((btn) => {
        btn.addEventListener("click", () => {
          node = btn.dataset.goto;
          path.push(node);
          setText(note, node);
          render();
        });
      });
    }
  };
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    path = ["start"]; node = "start"; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='stage-goals']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const movesOut = bind(lab, "[data-lab-moves]");
  const scoreOut = bind(lab, "[data-lab-score]");
  const note = bind(lab, "[data-lab-note]");
  const goal = Number(lab.dataset.goal || 30);
  const startMoves = Number(lab.dataset.moves || 10);
  let moves = startMoves, score = 0;
  const render = () => {
    if (movesOut) setText(movesOut, String(moves));
    if (scoreOut) setText(scoreOut, `${score}/${goal}`);
    if (board) {
      const pct = Math.min(100, (score / goal) * 100);
      board.innerHTML = `<div class="lab-goal-bar"><i style="width:${pct}%"></i></div><p class="lab-goal-label">${score >= goal ? "CLEAR" : `${moves} moves left`}</p>`;
    }
  };
  bind(lab, "[data-lab-match]")?.addEventListener("click", () => {
    if (moves <= 0) { setText(note, lab.dataset.out || "no moves"); return; }
    moves -= 1;
    score += 8;
    setText(note, score >= goal ? (lab.dataset.clear || "clear!") : (lab.dataset.hit || "+8"));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    moves = startMoves; score = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='shape-cells']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const nameOut = bind(lab, "[data-lab-name]");
  const note = bind(lab, "[data-lab-note]");
  const shapes = {
    I: [[0, 1], [1, 1], [2, 1], [3, 1]],
    O: [[1, 1], [2, 1], [1, 2], [2, 2]],
    T: [[0, 1], [1, 1], [2, 1], [1, 2]],
    L: [[1, 0], [1, 1], [1, 2], [2, 2]],
  };
  const names = Object.keys(shapes);
  let idx = 0;
  let rot = 0;
  const cells = () => {
    const base = shapes[names[idx]];
    return base.map(([x, y]) => {
      let nx = x - 1.5, ny = y - 1.5;
      for (let r = 0; r < rot; r++) { const t = nx; nx = -ny; ny = t; }
      return [Math.round(nx + 1.5), Math.round(ny + 1.5)];
    });
  };
  const render = () => {
    if (nameOut) setText(nameOut, `${names[idx]} r${rot}`);
    const set = new Set(cells().map(([x, y]) => `${x},${y}`));
    if (board) {
      let html = '<div class="lab-shape-grid">';
      for (let y = 0; y < 4; y++) for (let x = 0; x < 4; x++) {
        html += `<span class="${set.has(`${x},${y}`) ? "on" : ""}"></span>`;
      }
      html += "</div>";
      board.innerHTML = html;
    }
  };
  bind(lab, "[data-lab-next]")?.addEventListener("click", () => {
    idx = (idx + 1) % names.length; rot = 0;
    setText(note, names[idx]);
    render();
  });
  bind(lab, "[data-lab-rot]")?.addEventListener("click", () => {
    rot = (rot + 1) % 4;
    setText(note, lab.dataset.rotated || "rotated");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    idx = 0; rot = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='kick-try']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const tryOut = bind(lab, "[data-lab-try]");
  const note = bind(lab, "[data-lab-note]");
  const kicks = (lab.dataset.kicks || "0,0;1,0;-1,0;0,-1").split(";").map((s) => s.split(",").map(Number));
  let i = 0;
  let pos = [2, 2];
  const wall = new Set(["3,2", "3,1"]);
  const render = () => {
    if (tryOut) setText(tryOut, `${i}/${kicks.length}`);
    if (board) {
      let html = '<div class="lab-kick-grid">';
      for (let y = 0; y < 4; y++) for (let x = 0; x < 4; x++) {
        const k = `${x},${y}`;
        const cls = wall.has(k) ? "wall" : (pos[0] === x && pos[1] === y) ? "you" : "";
        html += `<span class="${cls}"></span>`;
      }
      html += "</div>";
      board.innerHTML = html;
    }
  };
  bind(lab, "[data-lab-kick]")?.addEventListener("click", () => {
    if (i >= kicks.length) { setText(note, lab.dataset.fail || "all kicks failed"); return; }
    const [dx, dy] = kicks[i];
    const nx = 2 + dx, ny = 2 + dy;
    i += 1;
    if (wall.has(`${nx},${ny}`) || nx < 0 || ny < 0 || nx > 3 || ny > 3) {
      setText(note, (lab.dataset.blocked || "kick ({dx},{dy}) blocked").replace("{dx}", dx).replace("{dy}", dy));
    } else {
      pos = [nx, ny];
      setText(note, (lab.dataset.ok || "kick ({dx},{dy}) OK").replace("{dx}", dx).replace("{dy}", dy));
      i = kicks.length;
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    i = 0; pos = [2, 2]; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='bag-draw']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const leftOut = bind(lab, "[data-lab-left]");
  const note = bind(lab, "[data-lab-note]");
  const refill = () => ["I", "O", "T", "L", "J", "S", "Z"];
  let bag = refill();
  let hold = "—";
  const render = () => {
    if (leftOut) setText(leftOut, String(bag.length));
    if (board) {
      board.innerHTML = `<div class="lab-bag-row">${bag.map((s) => `<span>${s}</span>`).join("")}</div>
        <p class="lab-bag-hold">HOLD ${hold}</p>`;
    }
  };
  bind(lab, "[data-lab-draw]")?.addEventListener("click", () => {
    if (!bag.length) bag = refill();
    const i = Math.floor(Math.random() * bag.length);
    const piece = bag.splice(i, 1)[0];
    setText(note, (lab.dataset.drew || "drew {p}").replace("{p}", piece));
    render();
  });
  bind(lab, "[data-lab-hold]")?.addEventListener("click", () => {
    if (!bag.length) bag = refill();
    const i = Math.floor(Math.random() * bag.length);
    const piece = bag.splice(i, 1)[0];
    const prev = hold;
    hold = piece;
    if (prev !== "—") bag.push(prev);
    setText(note, lab.dataset.held || "swapped hold");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    bag = refill(); hold = "—"; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='pipeline']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const stepOut = bind(lab, "[data-lab-stepn]");
  const note = bind(lab, "[data-lab-note]");
  const steps = (lab.dataset.steps || "input,update,draw").split(",");
  let i = 0;
  const render = () => {
    if (stepOut) setText(stepOut, `${i}/${steps.length}`);
    if (board) {
      board.innerHTML = `<div class="lab-pipe-row">${steps.map((s, n) =>
        `<span class="${n < i ? "done" : n === i ? "now" : ""}">${s}</span>`).join("<b>→</b>")}</div>`;
    }
  };
  bind(lab, "[data-lab-next]")?.addEventListener("click", () => {
    if (i >= steps.length) { i = 0; setText(note, lab.dataset.loop || "next frame"); }
    else {
      setText(note, steps[i]);
      i += 1;
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    i = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='event-queue']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  let q = [];
  let n = 1;
  const render = () => {
    if (countOut) setText(countOut, String(q.length));
    if (board) {
      board.innerHTML = q.length
        ? `<div class="lab-event-row">${q.map((e) => `<span>${e}</span>`).join("")}</div>`
        : `<p class="lab-empty">${lab.dataset.empty || "queue empty"}</p>`;
    }
  };
  bind(lab, "[data-lab-push]")?.addEventListener("click", () => {
    q.push(`${lab.dataset.event || "hit"}#${n++}`);
    setText(note, lab.dataset.pushed || "queued");
    render();
  });
  bind(lab, "[data-lab-pop]")?.addEventListener("click", () => {
    if (!q.length) { setText(note, lab.dataset.empty || "empty"); return; }
    const e = q.shift();
    setText(note, (lab.dataset.popped || "resolve {e}").replace("{e}", e));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    q = []; n = 1; setText(note, "—"); render();
  });
  render();
});

// Ebi Strike: connect the visible physics path to the event queue it creates.
document.querySelectorAll(".motion-lab[data-lab='shot-event-queue']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const hitsOut = bind(lab, "[data-lab-hits]");
  const note = bind(lab, "[data-lab-note]");
  const launchButton = bind(lab, "[data-lab-launch]");
  const labels = (lab.dataset.events || "HIT A,HIT B,HIT C,TURN END").split(",");
  const points = [[8, 78], [31, 25], [57, 70], [83, 30], [94, 58]];
  let queue = [];
  let hits = 0;
  let hp = [2, 2, 2];
  let launched = false;
  let progress = 0;
  let run = 0;

  const message = (key, replacements = {}) => {
    let value = lab.dataset[key] || "";
    for (const [name, replacement] of Object.entries(replacements)) {
      value = value.replace(`{${name}}`, replacement);
    }
    return value;
  };
  const positionAt = (t) => {
    const segment = Math.min(points.length - 2, Math.floor(t * (points.length - 1)));
    const local = t * (points.length - 1) - segment;
    return [
      points[segment][0] + (points[segment + 1][0] - points[segment][0]) * local,
      points[segment][1] + (points[segment + 1][1] - points[segment][1]) * local,
    ];
  };
  const render = () => {
    if (countOut) setText(countOut, String(queue.length));
    if (hitsOut) setText(hitsOut, `${hits}/3`);
    if (!board) return;
    const [x, y] = positionAt(progress);
    const trail = points.map((point) => point.join(",")).join(" ");
    board.innerHTML = `<div class="lab-shot-queue-demo">
      <div class="lab-shot-arena" aria-label="shot path through three enemies">
        <svg viewBox="0 0 100 100" aria-hidden="true"><polyline points="${trail}"/></svg>
        <i class="lab-shot-ball" style="left:${x}%;top:${y}%"></i>
        ${hp.map((value, i) => `<span class="lab-shot-enemy ${hits > i ? "contact" : ""} ${value < 2 ? "resolved" : ""}" style="left:${points[i + 1][0]}%;top:${points[i + 1][1]}%"><b>${String.fromCharCode(65 + i)}</b><small>HP ${value}</small></span>`).join("")}
      </div>
      <div class="lab-shot-queue" aria-label="event queue">
        <strong>EVENT QUEUE</strong>
        <div>${queue.length ? queue.map((event, i) => `<span class="${i === 0 ? "front" : ""}"><b>${i + 1}</b>${event.label}</span>`).join("") : `<em>${lab.dataset.empty || "queue empty"}</em>`}</div>
      </div>
    </div>`;
  };
  const enqueue = (index) => {
    const label = labels[index] || `HIT ${String.fromCharCode(65 + index)}`;
    queue.push({ label, target: index < 3 ? String.fromCharCode(65 + index) : "" });
    if (index < 3) {
      hits = index + 1;
      setText(note, message("contact", { e: label }));
    } else {
      setText(note, message("finished"));
    }
    render();
  };
  const launch = () => {
    run += 1;
    const thisRun = run;
    queue = [];
    hits = 0;
    hp = [2, 2, 2];
    progress = 0;
    launched = true;
    if (launchButton) launchButton.disabled = true;
    setText(note, message("launched"));
    render();
    let previous = 0;
    let lastTime = performance.now();
    const frame = (now) => {
      if (thisRun !== run) return;
      const elapsed = Math.min(34, now - lastTime);
      lastTime = now;
      progress = Math.min(1, progress + elapsed / 2100);
      [0.25, 0.5, 0.75].forEach((threshold, index) => {
        if (previous < threshold && progress >= threshold) enqueue(index);
      });
      previous = progress;
      render();
      if (progress < 1) requestAnimationFrame(frame);
      else {
        enqueue(3);
        launched = false;
        if (launchButton) launchButton.disabled = false;
      }
    };
    requestAnimationFrame(frame);
  };
  launchButton?.addEventListener("click", launch);
  bind(lab, "[data-lab-pop]")?.addEventListener("click", () => {
    if (!queue.length) {
      setText(note, lab.dataset.empty || "queue empty");
      return;
    }
    const event = queue.shift();
    if (event.target) {
      const index = event.target.charCodeAt(0) - 65;
      hp[index] = Math.max(0, hp[index] - 1);
      setText(note, message("resolved", { e: event.label, target: event.target }));
    } else {
      setText(note, message("turn"));
    }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    run += 1;
    queue = [];
    hits = 0;
    hp = [2, 2, 2];
    progress = 0;
    launched = false;
    if (launchButton) launchButton.disabled = false;
    setText(note, "—");
    render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='height-layers']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const hOut = bind(lab, "[data-lab-h]");
  const note = bind(lab, "[data-lab-note]");
  let heights = [2, 3, 4, 3, 5, 4, 2];
  const render = () => {
    if (hOut) setText(hOut, heights.join(","));
    if (board) {
      board.innerHTML = `<div class="lab-height-row">${heights.map((h) =>
        `<span style="height:${18 + h * 18}px"><i>${h}</i></span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-noise]")?.addEventListener("click", () => {
    heights = heights.map(() => 1 + Math.floor(Math.random() * 5));
    setText(note, lab.dataset.sampled || "sampled");
    render();
  });
  bind(lab, "[data-lab-carve]")?.addEventListener("click", () => {
    heights = heights.map((h) => Math.max(1, h - 1));
    setText(note, lab.dataset.carved || "carved");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    heights = [2, 3, 4, 3, 5, 4, 2]; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='craft-recipe']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const invOut = bind(lab, "[data-lab-inv]");
  const note = bind(lab, "[data-lab-note]");
  let wood = 0, string = 0, bow = 0;
  const render = () => {
    if (invOut) setText(invOut, `W${wood} S${string} B${bow}`);
    if (board) {
      board.innerHTML = `<div class="lab-craft-inv"><span>🪵 ${wood}</span><span>🧵 ${string}</span><span>🏹 ${bow}</span></div>
        <p class="lab-craft-recipe">2 wood + 1 string → bow</p>`;
    }
  };
  bind(lab, "[data-lab-wood]")?.addEventListener("click", () => { wood += 1; setText(note, "+wood"); render(); });
  bind(lab, "[data-lab-string]")?.addEventListener("click", () => { string += 1; setText(note, "+string"); render(); });
  bind(lab, "[data-lab-craft]")?.addEventListener("click", () => {
    if (wood < 2 || string < 1) { setText(note, lab.dataset.need || "need materials"); return; }
    wood -= 2; string -= 1; bow += 1;
    setText(note, lab.dataset.made || "crafted bow");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    wood = 0; string = 0; bow = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='light-flood']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const litOut = bind(lab, "[data-lab-lit]");
  const note = bind(lab, "[data-lab-note]");
  const size = 5;
  let light = Array.from({ length: size }, () => Array(size).fill(0));
  const render = () => {
    let lit = 0;
    light.flat().forEach((v) => { if (v > 0) lit += 1; });
    if (litOut) setText(litOut, String(lit));
    if (board) {
      board.innerHTML = `<div class="lab-light-grid">${light.map((row) => row.map((v) =>
        `<span style="opacity:${0.15 + v * 0.2}" class="${v ? "on" : ""}"></span>`).join("")).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-torch]")?.addEventListener("click", () => {
    light = Array.from({ length: size }, () => Array(size).fill(0));
    light[2][2] = 4;
    setText(note, lab.dataset.torch || "torch placed");
    render();
  });
  bind(lab, "[data-lab-flood]")?.addEventListener("click", () => {
    const next = light.map((row) => row.slice());
    for (let y = 0; y < size; y++) for (let x = 0; x < size; x++) {
      if (!light[y][x]) continue;
      for (const [dx, dy] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
        const nx = x + dx, ny = y + dy;
        if (nx < 0 || ny < 0 || nx >= size || ny >= size) continue;
        next[ny][nx] = Math.max(next[ny][nx], light[y][x] - 1);
      }
    }
    light = next;
    setText(note, lab.dataset.flooded || "flooded");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    light = Array.from({ length: size }, () => Array(size).fill(0));
    setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='species-inst']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  const maxHP = Number(lab.dataset.maxhp || 8);
  let units = [];
  let id = 1;
  const render = () => {
    if (countOut) setText(countOut, String(units.length));
    if (board) {
      board.innerHTML = `<p class="lab-species-def">DEF slime maxHP=${maxHP}</p>
        <div class="lab-species-row">${units.map((u) =>
          `<span>sl#${u.id} HP${u.hp}</span>`).join("") || `<i class="lab-empty">no instances</i>`}</div>`;
    }
  };
  bind(lab, "[data-lab-spawn]")?.addEventListener("click", () => {
    units.push({ id: id++, hp: maxHP });
    setText(note, lab.dataset.spawned || "spawned");
    render();
  });
  bind(lab, "[data-lab-hit]")?.addEventListener("click", () => {
    if (!units.length) { setText(note, lab.dataset.none || "none"); return; }
    units[0].hp -= 3;
    if (units[0].hp <= 0) units.shift();
    setText(note, lab.dataset.hit || "hit -3");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    units = []; id = 1; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='party-swap']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const activeOut = bind(lab, "[data-lab-active]");
  const note = bind(lab, "[data-lab-note]");
  const party = (lab.dataset.party || "Ebi,Shell,Coral").split(",");
  let active = 0;
  const render = () => {
    if (activeOut) setText(activeOut, party[active]);
    if (board) {
      board.innerHTML = `<div class="lab-party-row">${party.map((p, i) =>
        `<span class="${i === active ? "on" : ""}">${p}</span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-swap]")?.addEventListener("click", () => {
    active = (active + 1) % party.length;
    setText(note, (lab.dataset.swapped || "active → {p}").replace("{p}", party[active]));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    active = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='capture-roll']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const rateOut = bind(lab, "[data-lab-rate]");
  const note = bind(lab, "[data-lab-note]");
  const base = 15;
  let hp = 100;
  let sleep = false;
  let orb = false;
  const chance = () => {
    const hpBonus = Math.floor(((100 - hp) * 55) / 100);
    const status = sleep ? 20 : 0;
    const item = orb ? 15 : 0;
    return Math.min(95, base + hpBonus + status + item);
  };
  const render = () => {
    const rate = chance();
    if (rateOut) setText(rateOut, `${rate}%`);
    if (board) {
      const ja = document.documentElement.lang === "ja";
      board.innerHTML = `<div class="lab-capture-stage">
        <div class="lab-capture-monster" style="opacity:${0.35 + (hp / 100) * 0.65}"></div>
        <div class="lab-capture-formula">
          <span>15</span><span>+ HP↓ ${Math.floor(((100 - hp) * 55) / 100)}</span>
          <span>+ ${sleep ? "SLEEP 20" : "awake 0"}</span>
          <span>+ ${orb ? "ORB 15" : "ball 0"}</span>
          <strong>= ${rate}%</strong>
        </div>
        <div class="lab-capture-bar"><i style="width:${rate}%"></i><b class="lab-capture-needle" data-lab-needle></b></div>
        <p class="lab-capture-label">${ja ? `あと HP ${hp} · 捕獲率 ${rate}%` : `HP ${hp} left · ${rate}% catch`}</p>
      </div>`;
    }
  };
  const bumpHP = () => {
    hp = Math.max(5, hp - 20);
    setText(note, (lab.dataset.weakened || "HP down → rate up").replace("{hp}", String(hp)));
    render();
  };
  bind(lab, "[data-lab-bait]")?.addEventListener("click", bumpHP);
  bind(lab, "[data-lab-weaken]")?.addEventListener("click", bumpHP);
  bind(lab, "[data-lab-sleep]")?.addEventListener("click", () => {
    sleep = !sleep;
    setText(note, sleep ? (lab.dataset.slept || "SLEEP +20%") : (lab.dataset.woke || "awake"));
    render();
  });
  bind(lab, "[data-lab-orb]")?.addEventListener("click", () => {
    orb = !orb;
    setText(note, orb ? (lab.dataset.orbed || "ORB +15%") : (lab.dataset.basic || "basic ball"));
    render();
  });
  bind(lab, "[data-lab-roll]")?.addEventListener("click", () => {
    const rate = chance();
    const roll = Math.floor(Math.random() * 100);
    const ok = roll < rate;
    const needle = lab.querySelector("[data-lab-needle]");
    if (needle) needle.style.left = `${roll}%`;
    setText(note, ok
      ? (lab.dataset.caught || "caught! roll {r} < {p}%").replace("{r}", String(roll)).replace("{p}", String(rate))
      : (lab.dataset.missed || "broke free — roll {r} ≥ {p}%").replace("{r}", String(roll)).replace("{p}", String(rate)));
    render();
    const n2 = lab.querySelector("[data-lab-needle]");
    if (n2) n2.style.left = `${roll}%`;
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    hp = 100; sleep = false; orb = false; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='xp-level']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const levelOut = bind(lab, "[data-lab-level]");
  const xpOut = bind(lab, "[data-lab-xp]");
  const note = bind(lab, "[data-lab-note]");
  let level = 1, xp = 0;
  const need = () => level * 10;
  const render = () => {
    if (levelOut) setText(levelOut, String(level));
    if (xpOut) setText(xpOut, `${xp}/${need()}`);
    if (board) {
      const pct = Math.min(100, (xp / need()) * 100);
      board.innerHTML = `<div class="lab-xp-bar"><i style="width:${pct}%"></i></div>
        <p class="lab-xp-label">Lv ${level}${level >= 5 ? " → evolve?" : ""}</p>`;
    }
  };
  bind(lab, "[data-lab-gain]")?.addEventListener("click", () => {
    xp += 4;
    while (xp >= need()) { xp -= need(); level += 1; }
    setText(note, level >= 5 ? (lab.dataset.evolve || "ready to evolve") : (lab.dataset.gained || "+4 XP"));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    level = 1; xp = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='bfs-flood']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  let grid = [
    [1, 1, 0, 2],
    [1, 0, 0, 2],
    [1, 1, 2, 2],
    [0, 3, 3, 3],
  ];
  let marks = new Set();
  const key = (x, y) => `${x},${y}`;
  const render = () => {
    if (countOut) setText(countOut, String(marks.size));
    if (board) {
      board.innerHTML = `<div class="lab-bfs-grid">${grid.map((row, y) => row.map((c, x) =>
        `<span class="c${c}${marks.has(key(x, y)) ? " hit" : ""}" data-x="${x}" data-y="${y}"></span>`).join("")).join("")}</div>`;
      board.querySelectorAll("span").forEach((cell) => {
        cell.addEventListener("click", () => {
          const x = Number(cell.dataset.x), y = Number(cell.dataset.y);
          const color = grid[y][x];
          marks = new Set();
          const q = [[x, y]];
          marks.add(key(x, y));
          while (q.length) {
            const [cx, cy] = q.shift();
            for (const [dx, dy] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
              const nx = cx + dx, ny = cy + dy;
              if (ny < 0 || nx < 0 || ny >= 4 || nx >= 4) continue;
              if (grid[ny][nx] !== color) continue;
              const k = key(nx, ny);
              if (marks.has(k)) continue;
              marks.add(k);
              q.push([nx, ny]);
            }
          }
          setText(note, (lab.dataset.group || "group size {n}").replace("{n}", String(marks.size)));
          render();
        });
      });
    }
  };
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    marks = new Set(); setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='pellet-count']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const leftOut = bind(lab, "[data-lab-left]");
  const note = bind(lab, "[data-lab-note]");
  let dots = Array.from({ length: 12 }, (_, i) => i % 5 !== 2);
  const render = () => {
    const left = dots.filter(Boolean).length;
    if (leftOut) setText(leftOut, String(left));
    if (board) {
      board.innerHTML = `<div class="lab-pellet-grid">${dots.map((on) =>
        `<span class="${on ? "on" : ""}"></span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-eat]")?.addEventListener("click", () => {
    const i = dots.findIndex(Boolean);
    if (i < 0) { setText(note, lab.dataset.clear || "all clear!"); return; }
    dots[i] = false;
    const left = dots.filter(Boolean).length;
    setText(note, left ? (lab.dataset.ate || "ate one") : (lab.dataset.clear || "all clear!"));
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    dots = Array.from({ length: 12 }, (_, i) => i % 5 !== 2);
    setText(note, "—"); render();
  });
  render();
});

/* --- Late-track polish labs (maze-chase / bomb-maze) --- */

document.querySelectorAll(".motion-lab[data-lab='input-buffer']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const curOut = bind(lab, "[data-lab-current]");
  const qOut = bind(lab, "[data-lab-queued]");
  const note = bind(lab, "[data-lab-note]");
  let current = "E", queued = "E", atCenter = true, slide = 0, anim = null;
  const delta = { N: [0, -1], S: [0, 1], W: [-1, 0], E: [1, 0] };
  let px = 2, py = 2;
  const render = () => {
    if (curOut) setText(curOut, current);
    if (qOut) setText(qOut, queued);
    if (!board) return;
    const [dx, dy] = delta[current] || [1, 0];
    const drawX = atCenter ? px : px + dx * slide;
    const drawY = atCenter ? py : py + dy * slide;
    let cells = "";
    for (let y = 0; y < 5; y++) {
      for (let x = 0; x < 5; x++) {
        const wall = (x === 0 || y === 0 || x === 4 || y === 4) && !(x === 2 || y === 2);
        cells += `<span class="${wall ? "wall" : "path"}"></span>`;
      }
    }
    board.innerHTML = `<div class="lab-buffer-maze">${cells}
      <i class="lab-buffer-runner" style="left:${10 + drawX * 18}%;top:${10 + drawY * 18}%"></i>
      <div class="lab-buffer-arrows"><span class="cur">CUR ${current}</span><span class="q">Q ${queued}</span></div>
      <p class="lab-buffer-pos">${atCenter ? "AT CENTER — ready" : `SLIDING ${(slide * 100) | 0}%`}</p>
    </div>`;
  };
  const stopAnim = () => { if (anim) { cancelAnimationFrame(anim); anim = null; } };
  lab.querySelectorAll("[data-lab-dir]").forEach((btn) => {
    btn.addEventListener("click", () => {
      queued = btn.dataset.labDir || queued;
      setText(note, lab.dataset.queued || "queued");
      render();
    });
  });
  bind(lab, "[data-lab-center]")?.addEventListener("click", () => {
    stopAnim();
    atCenter = true;
    slide = 0;
    if (queued !== current) {
      current = queued;
      setText(note, lab.dataset.turned || "turned at center");
    } else setText(note, lab.dataset.same || "kept going");
    render();
  });
  bind(lab, "[data-lab-move]")?.addEventListener("click", () => {
    stopAnim();
    atCenter = false;
    slide = 0;
    setText(note, lab.dataset.moving || "moving");
    const tick = () => {
      slide = Math.min(1, slide + 0.08);
      render();
      if (slide < 1) anim = requestAnimationFrame(tick);
      else {
        const [dx, dy] = delta[current] || [0, 0];
        px = Math.max(1, Math.min(3, px + dx));
        py = Math.max(1, Math.min(3, py + dy));
        atCenter = true;
        slide = 0;
        setText(note, lab.dataset.arrived || "arrived — apply queue next");
        render();
      }
    };
    anim = requestAnimationFrame(tick);
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    stopAnim();
    current = "E"; queued = "E"; atCenter = true; slide = 0; px = 2; py = 2;
    setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='bomb-timer']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const timerOut = bind(lab, "[data-lab-timer]");
  const stateOut = bind(lab, "[data-lab-state]");
  const note = bind(lab, "[data-lab-note]");
  const fuse = Number(lab.dataset.fuse || 8);
  let timer = fuse, blasting = false, alive = false;
  const render = () => {
    if (timerOut) setText(timerOut, alive ? String(timer) : "—");
    if (stateOut) setText(stateOut, !alive ? "none" : blasting ? "BLAST" : "armed");
    if (board) {
      const pct = alive && !blasting ? (timer / fuse) * 100 : blasting ? 100 : 0;
      board.innerHTML = `<div class="lab-bomb-stage ${blasting ? "blast" : alive ? "armed" : ""}">
        <div class="lab-bomb-fuse"><i style="width:${pct}%"></i></div>
        <strong>${!alive ? "NO BOMB" : blasting ? "BOOM" : "TICK"}</strong>
      </div>`;
    }
  };
  bind(lab, "[data-lab-place]")?.addEventListener("click", () => {
    alive = true; blasting = false; timer = fuse;
    setText(note, lab.dataset.placed || "placed");
    render();
  });
  bind(lab, "[data-lab-tick]")?.addEventListener("click", () => {
    if (!alive) { setText(note, lab.dataset.none || "place first"); return; }
    if (blasting) {
      alive = false; blasting = false;
      setText(note, lab.dataset.gone || "removed");
      render();
      return;
    }
    timer -= 1;
    if (timer <= 0) { blasting = true; timer = 0; setText(note, lab.dataset.boom || "explode!"); }
    else setText(note, lab.dataset.tick || "tick");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    alive = false; blasting = false; timer = fuse; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='cross-blast']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const reachOut = bind(lab, "[data-lab-reach]");
  const note = bind(lab, "[data-lab-note]");
  const size = 7;
  const mid = 3;
  let power = 1;
  const walls = new Set(["1,3", "5,1"]);
  const render = () => {
    if (reachOut) setText(reachOut, String(power));
    const lit = new Set([`${mid},${mid}`]);
    for (const [dx, dy] of [[1, 0], [-1, 0], [0, 1], [0, -1]]) {
      for (let i = 1; i <= power; i++) {
        const x = mid + dx * i, y = mid + dy * i;
        if (x < 0 || y < 0 || x >= size || y >= size) break;
        if (walls.has(`${x},${y}`)) { lit.add(`${x},${y}`); break; }
        lit.add(`${x},${y}`);
      }
    }
    if (board) {
      let html = '<div class="lab-blast-grid">';
      for (let y = 0; y < size; y++) for (let x = 0; x < size; x++) {
        const k = `${x},${y}`;
        const cls = walls.has(k) ? "wall" : lit.has(k) ? (x === mid && y === mid ? "core" : "flame") : "";
        html += `<span class="${cls}"></span>`;
      }
      html += "</div>";
      board.innerHTML = html;
    }
  };
  bind(lab, "[data-lab-power]")?.addEventListener("click", () => {
    power = Math.min(3, power + 1);
    setText(note, lab.dataset.powered || "power up");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    power = 1; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='chain-bomb']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const countOut = bind(lab, "[data-lab-count]");
  const note = bind(lab, "[data-lab-note]");
  // cells: 0 empty, 1 bomb, 2 blasting
  let cells = [0, 1, 0, 1, 0, 1, 0];
  let blasts = 0;
  const render = () => {
    if (countOut) setText(countOut, String(blasts));
    if (board) {
      board.innerHTML = `<div class="lab-chain-row">${cells.map((c) =>
        `<span class="${c === 1 ? "bomb" : c === 2 ? "blast" : ""}"></span>`).join("")}</div>`;
    }
  };
  bind(lab, "[data-lab-ignite]")?.addEventListener("click", () => {
    // ignite first bomb
    const i = cells.indexOf(1);
    if (i < 0 && !cells.includes(2)) { setText(note, lab.dataset.done || "chain done"); return; }
    if (i >= 0) cells[i] = 2;
    // propagate: blasting lights neighbor bombs
    let changed = true;
    while (changed) {
      changed = false;
      for (let x = 0; x < cells.length; x++) {
        if (cells[x] !== 2) continue;
        for (const n of [x - 1, x + 1]) {
          if (n >= 0 && n < cells.length && cells[n] === 1) { cells[n] = 2; changed = true; }
        }
      }
    }
    blasts = cells.filter((c) => c === 2).length;
    setText(note, lab.dataset.chained || "chain reaction");
    render();
  });
  bind(lab, "[data-lab-clear]")?.addEventListener("click", () => {
    cells = cells.map((c) => (c === 2 ? 0 : c));
    blasts = 0;
    setText(note, lab.dataset.cleared || "blasts cleared");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    cells = [0, 1, 0, 1, 0, 1, 0]; blasts = 0; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='escape-timing']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const fuseOut = bind(lab, "[data-lab-fuse]");
  const etaOut = bind(lab, "[data-lab-eta]");
  const note = bind(lab, "[data-lab-note]");
  let fuse = 90, eta = 22, margin = 8;
  const render = () => {
    if (fuseOut) setText(fuseOut, String(fuse));
    if (etaOut) setText(etaOut, String(eta));
    const ok = eta + margin < fuse;
    if (board) {
      board.innerHTML = `<div class="lab-escape-compare ${ok ? "ok" : "bad"}">
        <div><span>ETA+margin</span><strong>${eta + margin}</strong></div>
        <div><span>vs fuse</span><strong>${fuse}</strong></div>
        <p>${ok ? "SAFE ROUTE" : "TOO SLOW — reroute"}</p>
      </div>`;
    }
  };
  bind(lab, "[data-lab-tick]")?.addEventListener("click", () => {
    fuse = Math.max(0, fuse - 10);
    eta = Math.max(5, eta - 2);
    const ok = eta + margin < fuse;
    setText(note, ok ? (lab.dataset.ok || "still safe") : (lab.dataset.bad || "reroute!"));
    render();
  });
  bind(lab, "[data-lab-far]")?.addEventListener("click", () => {
    eta += 15;
    setText(note, lab.dataset.far || "longer path");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    fuse = 90; eta = 22; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='junction-pick']").forEach((lab) => {
  const board = lab.querySelector("[data-lab-board]");
  const choiceOut = bind(lab, "[data-lab-choice]");
  const note = bind(lab, "[data-lab-note]");
  let target = "right", mode = "chase";
  const render = () => {
    if (choiceOut) setText(choiceOut, `${mode}:${target}`);
    if (board) {
      board.innerHTML = `<div class="lab-junction">
        <span class="up ${target === "up" ? "on" : ""}">↑</span>
        <div class="mid">
          <span class="left ${target === "left" ? "on" : ""}">←</span>
          <span class="core">AI</span>
          <span class="right ${target === "right" ? "on" : ""}">→</span>
        </div>
        <span class="down ${target === "down" ? "on" : ""}">↓</span>
        <p class="lab-junction-mode">${mode.toUpperCase()}</p>
      </div>`;
    }
  };
  bind(lab, "[data-lab-chase]")?.addEventListener("click", () => {
    mode = "chase"; target = "right";
    setText(note, lab.dataset.chase || "pick toward player");
    render();
  });
  bind(lab, "[data-lab-scatter]")?.addEventListener("click", () => {
    mode = "scatter"; target = "up";
    setText(note, lab.dataset.scatter || "pick corner");
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => {
    mode = "chase"; target = "right"; setText(note, "—"); render();
  });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='fx-breakout']").forEach((lab) => {
  const brickList = bind(lab, "[data-lab-brick-list]"), shardList = bind(lab, "[data-lab-shard-list]");
  const brickOut = bind(lab, "[data-lab-bricks]"), shardOut = bind(lab, "[data-lab-shards]");
  const hint = bind(lab, "[data-lab-breakout-hint]"), ja = document.documentElement.lang === "ja";
  let bricks = 3, shards = [], sid = 0;
  const render = () => {
    if (brickList) brickList.innerHTML = Array.from({length: bricks}, (_, i) => `<li>brick#${i+1}</li>`).join("") || `<li>${ja ? "（空）" : "(empty)"}</li>`;
    if (shardList) shardList.innerHTML = shards.length ? shards.map((s) => `<li>shard#${s.id} life=${s.life}</li>`).join("") : `<li>${ja ? "（空）" : "(empty)"}</li>`;
    setText(brickOut, String(bricks)); setText(shardOut, String(shards.length));
  };
  bind(lab, "[data-lab-breakout-hit]")?.addEventListener("click", () => {
    if (bricks > 0) { bricks--; for (let i=0;i<4;i++) shards.push({id:++sid,life:8}); }
    setText(hint, ja ? "ボールは残り、ブロックだけ削除。破片はfxへ追加。" : "Ball stays; brick removed; shards enter FX."); render();
  });
  bind(lab, "[data-lab-breakout-tick]")?.addEventListener("click", () => { shards=shards.map((s)=>({...s,life:s.life-1})).filter((s)=>s.life>0);setText(hint,ja?"破片だけ進み、ブロックは変わりません。":"Only shards advance; bricks stay.");render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click",()=>{bricks=3;shards=[];sid=0;render();}); render();
});

document.querySelectorAll(".motion-lab[data-lab='fx-snake']").forEach((lab) => {
  const grid=bind(lab,"[data-lab-snake-grid]"), bodyList=bind(lab,"[data-lab-snake-body]"), fxList=bind(lab,"[data-lab-snake-fx]");
  const lengthOut=bind(lab,"[data-lab-snake-length]"), partsOut=bind(lab,"[data-lab-snake-parts]"), hint=bind(lab,"[data-lab-snake-hint]");
  const ja=document.documentElement.lang==="ja";let body=[2,1,0],food=3,parts=[],id=0;
  const render=()=>{setText(grid,`[${Array.from({length:7},(_,x)=>body.includes(x)?"●":x===food?"◎":"·").join(" ")}]`);if(bodyList)bodyList.innerHTML=body.map((x,i)=>`<li>${i?"body":"head"}[${i}] x=${x}</li>`).join("");if(fxList)fxList.innerHTML=parts.length?parts.map((p)=>`<li>spark#${p.id} life=${p.life}</li>`).join(""):`<li>${ja?"（空）":"(empty)"}</li>`;setText(lengthOut,String(body.length));setText(partsOut,String(parts.length));};
  bind(lab,"[data-lab-snake-step]")?.addEventListener("click",()=>{const head=(body[0]+1)%7,ate=head===food;body.unshift(head);if(ate){food=(food+3)%7;for(let n=0;n<4;n++)parts.push({id:++id,life:8});setText(hint,ja?"エサを消して体セルを追加。捕食位置にfxを生成。":"Food removed, body grows, FX spawns.");}else body.pop();render();});
  bind(lab,"[data-lab-snake-tick]")?.addEventListener("click",()=>{parts=parts.map((p)=>({...p,life:p.life-1})).filter((p)=>p.life>0);setText(hint,ja?"キラキラだけ進み、体セルは変わりません。":"Only sparkles advance; body stays.");render();});
  bind(lab,"[data-lab-reset]")?.addEventListener("click",()=>{body=[2,1,0];food=3;parts=[];id=0;render();});render();
});

/* --- Mono-board visual upgrade (paint empty data-lab-board from live values) --- */

function labBoard(lab) {
  return lab.querySelector("[data-lab-board]") || lab.querySelector(".lab-board") || lab.querySelector(".lab-entities");
}

function watchLabPaint(lab, paint) {
  const board = labBoard(lab);
  if (!board || board.dataset.visualized === "1") return;
  board.dataset.visualized = "1";
  const run = () => paint(board, lab);
  // MutationObserverは使わず、操作イベントで再描画する（iframe互換）。
  lab.querySelector(".lab-controls")?.addEventListener("click", () => requestAnimationFrame(run));
  run();
}

document.querySelectorAll(".motion-lab[data-lab='pool']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const slots = (lab.querySelector("[data-lab-slots]")?.textContent || "").trim().split(/\s+/).filter(Boolean);
    board.innerHTML = `<div class="lab-pool-grid">${slots.map((s) =>
      `<span class="${String(s).startsWith("N") ? "new" : ""}">${s}</span>`).join("")}</div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='ai']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const dist = Number(lab.querySelector("[data-lab-dist]")?.textContent || 180);
    const state = lab.querySelector("[data-lab-state]")?.textContent || "";
    const chasing = dist < 120;
    const pct = Math.min(100, (dist / 260) * 100);
    board.innerHTML = `<div class="lab-ai-stage ${chasing ? "chase" : "patrol"}"><span class="foe"></span><i style="width:${pct}%"></i><span class="you" style="left:${pct}%"></span><p>${state}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='carry']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const plat = Number(lab.querySelector("[data-lab-plat]")?.textContent || 100);
    const player = Number(lab.querySelector("[data-lab-player]")?.textContent || 120);
    const px = Math.max(0, Math.min(80, ((plat - 40) / 160) * 100));
    const py = Math.max(0, Math.min(80, ((player - 40) / 160) * 100));
    board.innerHTML = `<div class="lab-carry-stage"><span class="plat" style="left:${px}%"></span><span class="you" style="left:${py}%"></span></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='stomp']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const result = lab.querySelector("[data-lab-result]")?.textContent || "";
    const vy = Number(lab.querySelector("[data-lab-vy]")?.textContent || 0);
    const ok = vy > 0;
    board.innerHTML = `<div class="lab-stomp-stage ${ok ? "ok" : "bad"}"><span class="you ${ok ? "down" : "up"}"></span><span class="foe"></span><p>${result}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='power']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const powered = (lab.querySelector("[data-lab-powered]")?.textContent || "").toLowerCase();
    const on = powered === "on" || powered === "あり" || powered.includes("on");
    const result = lab.querySelector("[data-lab-result]")?.textContent || "";
    board.innerHTML = `<div class="lab-power-stage ${on ? "on" : ""}"><span class="you"></span><span class="item"></span><p>${result}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='move8'], .motion-lab[data-lab='sling-drag']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const raw = lab.querySelector("[data-lab-raw]")?.textContent || "—";
    const norm = lab.querySelector("[data-lab-norm]")?.textContent || "—";
    board.innerHTML = `<div class="lab-move8-stage"><span class="raw"></span><span class="norm"></span><p>raw ${raw} → unit ${norm}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='aim']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const t = lab.querySelector("[data-lab-target]")?.textContent || "?";
    const cd = lab.querySelector("[data-lab-cd]")?.textContent || "0";
    const shots = lab.querySelector("[data-lab-shots]")?.textContent || "0";
    const enemies = [{ id: "A", d: 120 }, { id: "B", d: 80 }, { id: "C", d: 200 }];
    board.innerHTML = `<div class="lab-aim-stage"><span class="turret"></span>${enemies.map((e) =>
      `<span class="mob ${e.id === t ? "lock" : ""}" style="--d:${e.d}">${e.id}</span>`).join("")}<p>cd ${cd} · shots ${shots}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='draft']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const mode = lab.querySelector("[data-lab-mode]")?.textContent || "";
    const pick = lab.querySelector("[data-lab-pick]")?.textContent || "—";
    const drafting = /draft|下書き|選択/i.test(mode);
    board.innerHTML = drafting
      ? `<div class="lab-draft-stage open"><span>A</span><span>B</span><span>C</span><p>${mode}</p></div>`
      : `<div class="lab-draft-stage"><p>${mode}${pick !== "—" ? " · " + pick : ""}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='evolve']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const name = lab.querySelector("[data-lab-name]")?.textContent || "";
    const count = Number(lab.querySelector("[data-lab-wcount]")?.textContent || 1);
    const storm = /storm/i.test(name);
    board.innerHTML = `<div class="lab-evolve-stage ${storm ? "storm" : ""}">${Array.from({ length: Math.max(1, count) }, () => "<span></span>").join("")}<p>${name}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='curve']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const sec = Number(lab.querySelector("[data-lab-sec]")?.textContent || 0);
    const interval = lab.querySelector("[data-lab-interval]")?.textContent || "";
    const bars = Array.from({ length: 9 }, (_, i) => {
      const s = i * 5;
      const iv = Math.max(14, 42 - Math.floor(s / 2));
      const h = Math.round(((42 - iv) / 28) * 100);
      return `<span class="${s === sec ? "now" : ""}" style="height:${20 + h * 0.7}%"></span>`;
    }).join("");
    board.innerHTML = `<div class="lab-curve-stage">${bars}<p>t=${sec}s · iv=${interval}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='click']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const n = lab.querySelector("[data-lab-count]")?.textContent || "0";
    board.innerHTML = `<div class="lab-click-stage"><strong>${n}</strong><p>TAP</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='shop']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const gold = lab.querySelector("[data-lab-gold]")?.textContent || "0";
    const cost = lab.querySelector("[data-lab-cost]")?.textContent || "0";
    const owned = Number(lab.querySelector("[data-lab-owned]")?.textContent || 0);
    board.innerHTML = `<div class="lab-shop-stage"><p class="g">★ ${gold}</p><p class="c">cost ${cost}</p><div class="owned">${Array.from({ length: owned }, () => "<span></span>").join("")}</div></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='idle']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const points = lab.querySelector("[data-lab-points]")?.textContent || "0";
    const rate = Number.parseInt(lab.querySelector("[data-lab-rate]")?.textContent || "2", 10) || 2;
    board.innerHTML = `<div class="lab-idle-stage"><div class="machines">${Array.from({ length: Math.min(12, rate) }, () => "<span></span>").join("")}</div><strong>${points}</strong></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='save']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const away = lab.querySelector("[data-lab-away]")?.textContent || "0s";
    const gained = lab.querySelector("[data-lab-gained]")?.textContent || "0";
    board.innerHTML = `<div class="lab-save-stage"><p>${away}</p><strong>+${gained}</strong></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='turn']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const states = (lab.dataset.states || "select,player,enemy,win").split(",");
    const cur = lab.querySelector("[data-lab-state]")?.textContent || states[0];
    const i = Math.max(0, states.indexOf(cur));
    board.innerHTML = `<div class="lab-turn-row">${states.map((s, n) =>
      `<span class="${n === i ? "now" : n < i ? "done" : ""}">${s}</span>`).join("<b>→</b>")}</div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='flag']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const flag = lab.querySelector("[data-lab-flag]")?.textContent || "";
    const text = lab.querySelector("[data-lab-text]")?.textContent || "";
    const on = /true|on|あり|はい/i.test(flag);
    board.innerHTML = `<div class="lab-flag-stage ${on ? "on" : ""}"><code>flag=${flag}</code><p>${text}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='damage']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const atk = lab.querySelector("[data-lab-atk]")?.textContent || "?";
    const def = lab.querySelector("[data-lab-def]")?.textContent || "?";
    const dmg = lab.querySelector("[data-lab-dmg]")?.textContent || "?";
    board.innerHTML = `<div class="lab-dmg-stage"><p>${atk} − ${def}</p><strong>${dmg}</strong></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='inv']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const gold = lab.querySelector("[data-lab-gold]")?.textContent || "0";
    const item = lab.querySelector("[data-lab-item]")?.textContent || "—";
    board.innerHTML = `<div class="lab-inv-stage"><p>★ ${gold}</p><strong>${item}</strong></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='scene']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const scene = lab.querySelector("[data-lab-scene]")?.textContent || "";
    const enemy = lab.querySelector("[data-lab-enemy]")?.textContent || "—";
    board.innerHTML = `<div class="lab-scene-stage"><strong>${scene}</strong><p>${enemy}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='quest']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const quest = lab.querySelector("[data-lab-quest]")?.textContent || "0";
    const saved = lab.querySelector("[data-lab-saved]")?.textContent || "—";
    board.innerHTML = `<div class="lab-quest-stage"><strong>Q${quest}</strong><p>${saved}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='buffer']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const buffer = lab.querySelector("[data-lab-buf]")?.textContent || "none";
    const life = Number(lab.querySelector("[data-lab-life]")?.textContent || 0);
    board.innerHTML = `<div class="lab-buf-stage"><strong>${buffer}</strong><div class="lab-bomb-fuse"><i style="width:${(life / 8) * 100}%"></i></div><p>life ${life}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='command']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const hist = lab.querySelector("[data-lab-hist]")?.textContent || "—";
    const match = lab.querySelector("[data-lab-match]")?.textContent || "—";
    const ok = match && match !== "—";
    board.innerHTML = `<div class="lab-cmd-stage ${ok ? "ok" : ""}"><div class="hist">${hist}</div><p>${ok ? match : "↓↘→"}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='rounds']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const p1 = Number(lab.querySelector("[data-lab-p1]")?.textContent || 0);
    const p2 = Number(lab.querySelector("[data-lab-p2]")?.textContent || 0);
    const round = lab.querySelector("[data-lab-round]")?.textContent || "1";
    const dots = (n) => Array.from({ length: 2 }, (_, i) => `<span class="${i < n ? "on" : ""}"></span>`).join("");
    board.innerHTML = `<div class="lab-rounds-stage"><div>${dots(p1)}</div><p>R${round}</p><div>${dots(p2)}</div></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='rps']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const result = lab.querySelector("[data-lab-result]")?.textContent || "—";
    board.innerHTML = `<div class="lab-rps-stage"><p>${result}</p></div>`;
  });
});

document.querySelectorAll(".motion-lab[data-lab='snake']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const bodyText = lab.querySelector("[data-lab-body]")?.textContent || "";
    const body = [...bodyText.matchAll(/\((\d+),(\d+)\)/g)].map((m) => ({ x: Number(m[1]), y: Number(m[2]) }));
    let html = '<div class="lab-snake-grid">';
    for (let y = 0; y < 5; y++) for (let x = 0; x < 6; x++) {
      const i = body.findIndex((p) => p.x === x && p.y === y);
      html += `<span class="${i === 0 ? "head" : i > 0 ? "body" : ""}"></span>`;
    }
    html += "</div>";
    board.innerHTML = html;
  });
});

document.querySelectorAll(".motion-lab[data-lab='push']").forEach((lab) => {
  watchLabPaint(lab, (board) => {
    const cells = (lab.querySelector("[data-lab-map]")?.textContent || ". P . B .").trim().split(/\s+/);
    board.innerHTML = `<div class="lab-push-row">${cells.map((c) =>
      `<span class="${c === "P" ? "you" : c === "B" ? "box" : c === "#" ? "wall" : ""}">${c === "." ? "" : c}</span>`).join("")}</div>`;
  });
});

// Match-three swap lab: the board only resolves after the moving pieces meet.
document.querySelectorAll(".motion-lab[data-lab='match-swap']").forEach((lab) => {
  const grid = lab.querySelector("[data-match-grid]");
  const phaseOut = lab.querySelector("[data-match-phase]");
  const readout = lab.querySelector("[data-match-readout]");
  const ja = lab.dataset.lang === "ja";
  const colors = ["coral", "sun", "leaf", "sky", "violet"];
  let run = 0;
  let progress = 0;
  let mode = "idle";

  const render = () => {
    const phase = mode === "idle" ? (ja ? "待機" : "READY") : mode === "valid" ? (progress < .42 ? "SWAP" : progress < .64 ? "CHECK MATCH" : progress < .85 ? "CLEAR" : "FALL") : (progress < .42 ? "SWAP" : progress < .7 ? "CHECK MATCH" : "RETURN");
    setText(phaseOut, phase);
    if (grid) {
      grid.innerHTML = Array.from({ length: 25 }, (_, i) => {
        const special = i === 11 ? "swap-a" : i === 12 ? "swap-b" : "";
        const x = special === "swap-a" ? progress * 100 : special === "swap-b" ? -progress * 100 : 0;
        const value = mode === "valid" && progress > .64 && i === 11 ? "" : colors[i % colors.length];
        return `<span class="match-tile ${value} ${special} ${mode === "valid" && progress > .64 && i === 11 ? "is-cleared" : ""}" style="--swap-x:${x}%"></span>`;
      }).join("");
    }
  };
  const play = (nextMode) => {
    if (mode !== "idle") return;
    run++;
    const current = run;
    mode = nextMode;
    progress = 0;
    setText(readout, nextMode === "valid" ? (ja ? "2個が動き始めた" : "Two pieces started moving") : (ja ? "交換を確認中" : "Checking the swap"));
    const start = performance.now();
    const frame = (now) => {
      if (current !== run) return;
      progress = Math.min(1, (now - start) / 1500);
      render();
      if (progress < 1) requestAnimationFrame(frame);
      else {
        setText(readout, nextMode === "valid" ? (ja ? "3個が揃った → 消去して落下" : "Three match → clear, then fall") : (ja ? "揃わない → 元の場所へ戻る" : "No match → return to the original cells"));
        mode = "idle";
        progress = 0;
        render();
      }
    };
    requestAnimationFrame(frame);
  };
  lab.querySelector("[data-match-action='valid']")?.addEventListener("click", () => play("valid"));
  lab.querySelector("[data-match-action='invalid']")?.addEventListener("click", () => play("invalid"));
  lab.querySelector("[data-match-reset]")?.addEventListener("click", () => { run++; mode = "idle"; progress = 0; setText(readout, ja ? "ボタンを押して、状態の順番を見よう。" : "Choose a swap and follow each phase."); render(); });
  render();
});

// Reversi evaluation lab: make the weighted sum visible instead of leaving
// the CPU's choice as a mysterious number.
document.querySelectorAll(".motion-lab[data-lab='reversi-eval']").forEach((lab) => {
  const map = lab.querySelector("[data-reversi-map]");
  const score = lab.querySelector("[data-reversi-score]");
  const ja = lab.dataset.lang === "ja";
  const weights = [
    [120, -20, 20, 5, 5, 20, -20, 120], [-20, -40, -5, -5, -5, -5, -40, -20],
    [20, -5, 15, 3, 3, 15, -5, 20], [5, -5, 3, 3, 3, 3, -5, 5],
    [5, -5, 3, 3, 3, 3, -5, 5], [20, -5, 15, 3, 3, 15, -5, 20],
    [-20, -40, -5, -5, -5, -5, -40, -20], [120, -20, 20, 5, 5, 20, -20, 120],
  ];
  let stones = new Map();
  const key = (x, y) => `${x},${y}`;
  const render = () => {
    let total = 0;
    if (map) {
      map.innerHTML = Array.from({ length: 64 }, (_, i) => {
        const x = i % 8, y = Math.floor(i / 8), value = weights[y][x], stone = stones.get(key(x, y)) || "";
        if (stone === "black") total += value;
        if (stone === "white") total -= value;
        const corner = (x === 0 || x === 7) && (y === 0 || y === 7) ? " is-corner" : "";
        return `<span class="reversi-eval-cell ${stone ? `is-${stone}` : ""}${corner}">${stone === "black" ? "●" : stone === "white" ? "○" : value}</span>`;
      }).join("");
    }
    if (score) score.textContent = ja ? `評価値 = ${total}（黒の石×マス点 − 白の石×マス点）` : `EVALUATION = ${total} (black stone×map − white stone×map)`;
  };
  lab.querySelectorAll("[data-reversi-place]").forEach((button) => button.addEventListener("click", () => {
    const mode = button.dataset.reversiPlace;
    if (mode === "corner") stones = new Map([[key(0, 0), "black"], [key(1, 0), "white"]]);
    if (mode === "risky") stones = new Map([[key(1, 0), "black"], [key(0, 0), "white"]]);
    if (mode === "center") stones = new Map([[key(3, 3), "black"], [key(4, 3), "white"]]);
    if (mode === "reset") stones = new Map();
    render();
  }));
  render();
});

// Metroidvania ability route: make the old ledge a small, readable challenge.
document.querySelectorAll(".motion-lab[data-lab='ability-gate']").forEach((lab) => {
  const player = lab.querySelector("[data-lab-player]");
  const ledge = lab.querySelector("[data-lab-ledge]");
  const wings = lab.querySelector("[data-lab-wings]");
  const status = lab.querySelector("[data-lab-status]");
  const ability = lab.querySelector("[data-lab-ability]");
  const ja = lab.dataset.lang === "ja";
  let x = 82;
  let hasWings = false;
  let jumpCount = 0;
  let won = false;

  const render = () => {
    if (player) player.style.left = `${x}%`;
    player?.classList.toggle("is-air", jumpCount > 0 && !won);
    wings?.classList.toggle("is-collected", hasWings);
    ledge?.classList.toggle("is-reachable", won);
    setText(ability, hasWings ? (ja ? "能力: 二段ジャンプ" : "ABILITY: DOUBLE JUMP") : (ja ? "能力: なし" : "ABILITY: none"));
  };
  const move = (dir) => {
    if (won) return;
    x = Math.max(8, Math.min(92, x + dir * 13));
    if (!hasWings && x >= 72) {
      hasWings = true;
      setText(status, ja ? "羽を取得！左へ戻って高台を目指そう。" : "WINGS FOUND! Return left to the ledge.");
    } else if (hasWings && x <= 31 && jumpCount >= 2) {
      won = true;
      setText(status, ja ? "高台に到着！同じ場所が道へ変わった。" : "LEDGE REACHED! The same place became a route.");
    } else {
      setText(status, x <= 35 && !hasWings ? (ja ? "まだ高さが足りない。右の羽を探そう。" : "Too low for this ledge. Find the wings on the right.") : (ja ? `位置 ${Math.round(x)}%` : `position ${Math.round(x)}%`));
    }
    render();
  };
  lab.querySelectorAll("[data-ability-move]").forEach((button) => button.addEventListener("click", () => move(button.dataset.abilityMove === "left" ? -1 : 1)));
  lab.querySelector("[data-ability-jump]")?.addEventListener("click", () => {
    if (won) return;
    if (!hasWings) {
      setText(status, ja ? "一段では届かない。まず右の羽を取得しよう。" : "One jump is not enough. Find the wings first.");
    } else if (jumpCount < 2) {
      jumpCount++;
      setText(status, jumpCount === 1 ? (ja ? "空中！もう一度押せば二段ジャンプ。" : "In the air! Press again for the second jump.") : (ja ? "二段目！左へ戻ろう。" : "Second jump! Move back to the left."));
    }
    render();
  });
  lab.querySelector("[data-ability-reset]")?.addEventListener("click", () => {
    x = 82; hasWings = false; jumpCount = 0; won = false;
    setText(status, ja ? "右の羽を拾おう" : "Find the wings on the right");
    render();
  });
  render();
});

// Metroidvania map lab: movement and the visited-room set are one visible action.
document.querySelectorAll(".motion-lab[data-lab='metroid-map']").forEach((lab) => {
  const path = lab.querySelector("[data-map-path]");
  const readout = lab.querySelector("[data-map-readout]");
  const ja = lab.dataset.lang === "ja";
  let room = 0;
  let visited = new Set([0]);
  const render = () => {
    if (path) {
      path.innerHTML = Array.from({ length: 6 }, (_, i) => `<span class="map-room ${visited.has(i) ? "is-visited" : "is-hidden"} ${room === i ? "is-current" : ""}">${visited.has(i) ? `R${i}` : "?"}</span>`).join("");
    }
    setText(readout, ja ? `room ID ${room} / 訪問 ${visited.size}/6 — room := int(playerX/roomWidth)` : `room ID ${room} / visited ${visited.size}/6 — room := int(playerX/roomWidth)`);
  };
  lab.querySelectorAll("[data-map-move]").forEach((button) => button.addEventListener("click", () => {
    room = Math.max(0, Math.min(5, room + (button.dataset.mapMove === "next" ? 1 : -1)));
    visited.add(room);
    render();
  }));
  lab.querySelector("[data-map-reset]")?.addEventListener("click", () => { room = 0; visited = new Set([0]); render(); });
  render();
});

/* Auto-wire TRY IT · GO highlights from data-go-focus / data-go-clear */
document.querySelectorAll(".motion-lab").forEach((lab) => {
  if (!lab.querySelector(".lab-go")) return;
  lab.addEventListener("click", (event) => {
    const btn = event.target.closest("button");
    if (!btn || !lab.contains(btn)) return;
    if (btn.hasAttribute("data-go-clear")) {
      clearGo(lab);
      return;
    }
    if (!btn.hasAttribute("data-go-focus")) return;
    const ids = btn.getAttribute("data-go-focus").trim().split(/\s+/).filter(Boolean);
    focusGo(lab, ids, btn.getAttribute("data-go-caption") || undefined);
  });
});
