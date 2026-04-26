/**
 * copy-code.js — adds a copy button to every code fence inside .article-content
 */
(function () {
  "use strict";

  var COPY_SVG =
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>';
  var DONE_SVG =
    '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>';

  function init() {
    var blocks = document.querySelectorAll(".article-content pre");
    blocks.forEach(function (pre) {
      if (pre.parentElement.classList.contains("code-block")) return;
      var wrap = document.createElement("div");
      wrap.className = "code-block";
      pre.parentNode.insertBefore(wrap, pre);
      wrap.appendChild(pre);

      var btn = document.createElement("button");
      btn.className = "code-block__copy";
      btn.type = "button";
      btn.setAttribute("aria-label", (window.lgI18n && window.lgI18n.copyCode) || "Copy code");
      btn.innerHTML = COPY_SVG;
      wrap.appendChild(btn);

      btn.addEventListener("click", function () {
        var code = pre.querySelector("code");
        var text = (code || pre).innerText;
        if (navigator.clipboard) {
          navigator.clipboard
            .writeText(text)
            .then(function () {
              btn.classList.add("is-copied");
              btn.innerHTML = DONE_SVG;
              setTimeout(function () {
                btn.classList.remove("is-copied");
                btn.innerHTML = COPY_SVG;
              }, 1600);
            })
            .catch(function () {});
        }
      });
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
