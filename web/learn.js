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
