(() => {
  const frames = document.querySelectorAll("iframe[data-game-src]");
  if (frames.length === 0) return;

  const load = (frame) => {
    const src = frame.dataset.gameSrc;
    if (!src) return;
    frame.addEventListener("load", () => frame.removeAttribute("aria-busy"), { once: true });
    frame.src = src;
    frame.removeAttribute("data-game-src");
  };

  if (!("IntersectionObserver" in window)) {
    frames.forEach(load);
    return;
  }

  const observer = new IntersectionObserver((entries) => {
    for (const entry of entries) {
      if (!entry.isIntersecting) continue;
      observer.unobserve(entry.target);
      load(entry.target);
    }
  }, { rootMargin: "200px 0px" });

  frames.forEach((frame) => {
    frame.setAttribute("aria-busy", "true");
    observer.observe(frame);
  });
})();
