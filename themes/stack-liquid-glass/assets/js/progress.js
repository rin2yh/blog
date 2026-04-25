/**
 * progress.js — reading progress bar at the top of the page
 *  Activates only when an .article-content is present.
 */
(function () {
  "use strict";

  function init() {
    var article = document.querySelector(".article-content");
    if (!article) return;

    var bar = document.createElement("div");
    bar.className = "reading-progress";
    bar.setAttribute("aria-hidden", "true");
    var fill = document.createElement("div");
    fill.className = "reading-progress__bar";
    bar.appendChild(fill);
    document.body.appendChild(bar);

    var ticking = false;
    function update() {
      var rect = article.getBoundingClientRect();
      var top = rect.top + window.scrollY;
      var bottom = top + article.offsetHeight - window.innerHeight;
      var span = Math.max(1, bottom - top);
      var pct = Math.min(100, Math.max(0, ((window.scrollY - top) / span) * 100));
      fill.style.width = pct + "%";
      ticking = false;
    }

    window.addEventListener(
      "scroll",
      function () {
        if (!ticking) {
          window.requestAnimationFrame(update);
          ticking = true;
        }
      },
      { passive: true }
    );
    update();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
