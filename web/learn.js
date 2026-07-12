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
  });
  lab.querySelector("[data-lab-step]")?.addEventListener("click", step);
  lab.querySelector("[data-lab-reset]")?.addEventListener("click", () => {
    velocity = 0;
    position = 320;
    render();
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
    dxOut.textContent = dx.toFixed(0);
    dyOut.textContent = dy.toFixed(0);
    distOut.textContent = dist.toFixed(1);
    result.textContent = hit ? hitLabel : missLabel;
    result.dataset.state = hit ? "hit" : "miss";
    board.dataset.state = hit ? "hit" : "miss";
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
    dxOut.textContent = "—";
    dyOut.textContent = "—";
    distOut.textContent = "—";
    result.textContent = lab.dataset.wait || "TAP";
    result.dataset.state = "wait";
    board.dataset.state = "wait";
  });
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

document.querySelectorAll(".motion-lab[data-lab='hitbox']").forEach((lab) => {
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
  const actor = bind(lab, "[data-lab-actor]");
  if (!camOut || !playerOut || !screenOut) return;
  let player = 200, cam = 0;
  const render = () => {
    const target = player - 160;
    cam += (target - cam) * 0.15;
    setText(camOut, cam.toFixed(0));
    setText(playerOut, player.toFixed(0));
    setText(screenOut, (player - cam).toFixed(0));
    if (actor) actor.style.left = `${((player - cam) / 480) * 100}%`;
  };
  bind(lab, "[data-lab-right]")?.addEventListener("click", () => { player += 40; render(); });
  bind(lab, "[data-lab-left]")?.addEventListener("click", () => { player = Math.max(40, player - 40); render(); });
  bind(lab, "[data-lab-step]")?.addEventListener("click", render);
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { player = 200; cam = 0; render(); });
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

document.querySelectorAll(".motion-lab[data-lab='move8']").forEach((lab) => {
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
  if (!targetOut || !cdOut || !shotOut) return;
  const enemies = [{ id: "A", d: 120 }, { id: "B", d: 80 }, { id: "C", d: 200 }];
  let frame = 0, shots = 0, cooldown = 28;
  const nearest = () => enemies.reduce((b, e) => (e.d < b.d ? e : b));
  const render = () => {
    setText(targetOut, nearest().id);
    setText(cdOut, String(frame % cooldown));
    setText(shotOut, String(shots));
  };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    frame++;
    if (frame % cooldown === 0) shots++;
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { frame = 0; shots = 0; render(); });
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
  const stateOut = bind(lab, "[data-lab-state]");
  const states = (lab.dataset.states || "select,player,enemy,win").split(",");
  if (!stateOut) return;
  let i = 0;
  const render = () => setText(stateOut, states[i]);
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { i = (i + 1) % states.length; render(); });
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
  const startup = lab.dataset.startup || "startup";
  const active = lab.dataset.active || "active";
  const recovery = lab.dataset.recovery || "recovery";
  if (!frameOut || !phaseOut) return;
  let frame = 0;
  const phase = () => (frame <= 8 ? startup : frame <= 12 ? active : recovery);
  const render = () => { setText(frameOut, String(frame)); setText(phaseOut, phase()); setState(phaseOut, frame > 8 && frame <= 12); };
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { frame = Math.min(30, frame + 1); render(); });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { frame = 0; render(); });
  render();
});

document.querySelectorAll(".motion-lab[data-lab='react']").forEach((lab) => {
  const stopOut = bind(lab, "[data-lab-stop]");
  const stunOut = bind(lab, "[data-lab-stun]");
  const vOut = bind(lab, "[data-lab-v]");
  if (!stopOut || !stunOut || !vOut) return;
  let hitstop = 0, stun = 0, v = 0;
  const render = () => { setText(stopOut, String(hitstop)); setText(stunOut, String(stun)); setText(vOut, v.toFixed(2)); };
  bind(lab, "[data-lab-hit]")?.addEventListener("click", () => { hitstop = 8; stun = 25; v = 7; render(); });
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => {
    if (hitstop > 0) { hitstop--; render(); return; }
    if (stun > 0) { stun--; v *= 0.86; }
    render();
  });
  bind(lab, "[data-lab-reset]")?.addEventListener("click", () => { hitstop = 0; stun = 0; v = 0; render(); });
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
  bind(lab, "[data-lab-step]")?.addEventListener("click", () => { v += g; y += v; if (y > 400) { y = 400; v = 0; } render(); });
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
