/**
 * liquid-glass.js — minimal runtime for the Liquid Glass theme
 *  - GlassTheme : light/dark toggle, persisted in localStorage
 *  - GlassRipple: subtle ripple on glass cards / buttons
 */
(function () {
  "use strict";

  var GlassTheme = {
    KEY: "liquid-glass-theme",
    init: function () {
      var btn = document.querySelector("[data-theme-toggle]");
      var root = document.documentElement;
      if (!btn) return;
      btn.addEventListener("click", function () {
        var next = root.getAttribute("data-theme") === "light" ? "dark" : "light";
        root.setAttribute("data-theme", next);
        try {
          localStorage.setItem(GlassTheme.KEY, next);
        } catch (e) {}
        btn.setAttribute(
          "aria-label",
          next === "dark" ? "Switch to light mode" : "Switch to dark mode"
        );
      });
    }
  };

  var GlassRipple = {
    init: function () {
      document.addEventListener("click", function (evt) {
        var target = evt.target.closest(".glass-btn, .article-card");
        if (!target) return;
        if (target.querySelectorAll("[data-ripple]").length > 2) return;
        var ripple = document.createElement("span");
        ripple.setAttribute("data-ripple", "");
        var rect = target.getBoundingClientRect();
        var size = Math.max(rect.width, rect.height) * 1.2;
        var x = evt.clientX - rect.left - size / 2;
        var y = evt.clientY - rect.top - size / 2;
        Object.assign(ripple.style, {
          position: "absolute",
          left: x + "px",
          top: y + "px",
          width: size + "px",
          height: size + "px",
          borderRadius: "50%",
          background: "rgba(255,255,255,0.16)",
          transform: "scale(0)",
          pointerEvents: "none",
          animation: "lg-ripple 600ms ease-out forwards",
          zIndex: "5"
        });
        target.appendChild(ripple);
        ripple.addEventListener(
          "animationend",
          function () {
            ripple.remove();
          },
          { once: true }
        );
      });
    }
  };

  var ScrollHoverGuard = {
    init: function () {
      var root = document.documentElement;
      var timer = 0;
      window.addEventListener(
        "scroll",
        function () {
          if (!timer) root.classList.add("is-scrolling");
          clearTimeout(timer);
          timer = setTimeout(function () {
            root.classList.remove("is-scrolling");
            timer = 0;
          }, 120);
        },
        { passive: true }
      );
    }
  };

  var ClipboardCopy = {
    init: function () {
      document.addEventListener("click", function (evt) {
        var btn = evt.target.closest("[data-copy]");
        if (!btn) return;
        var text = btn.getAttribute("data-copy");
        var ok = btn.getAttribute("data-copy-success") || "Copied";
        var ng = btn.getAttribute("data-copy-error") || "Copy failed";
        if (!navigator.clipboard) {
          alert(ng);
          return;
        }
        navigator.clipboard.writeText(text).then(
          function () {
            alert(ok);
          },
          function () {
            alert(ng);
          }
        );
      });
    }
  };

  function boot() {
    GlassTheme.init();
    GlassRipple.init();
    ScrollHoverGuard.init();
    ClipboardCopy.init();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot);
  } else {
    boot();
  }
})();
