document.querySelectorAll(".motion-lab").forEach((lab) => {
  const section = lab.closest(".physics");
  const bird = lab.querySelector("[data-lab-bird]");
  const velocityOutput = lab.querySelector("[data-lab-velocity]");
  const positionOutput = lab.querySelector("[data-lab-position]");
  const directionOutput = lab.querySelector("[data-lab-direction]");
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

  lab.querySelector("[data-lab-step]").addEventListener("click", step);
  lab.querySelector("[data-lab-flap]").addEventListener("click", () => {
    velocity = -7.4;
    step();
  });
  lab.querySelector("[data-lab-reset]").addEventListener("click", () => {
    velocity = 0;
    position = 320;
    render();
  });
  render();
});
