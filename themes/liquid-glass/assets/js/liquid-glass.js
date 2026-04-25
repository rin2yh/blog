/**
 * liquid-glass.js — minimal runtime for the Liquid Glass theme
 *  - GlassTheme : light/dark toggle, persisted in localStorage
 *  - GlassRipple: subtle ripple on glass cards / buttons
 *  - GlassParallax: mouse-driven tilt on .glass-card (skipped on touch / reduced-motion)
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

  var GlassParallax = {
    init: function () {
      if (window.matchMedia("(hover: none)").matches) return;
      if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;

      document.querySelectorAll(".article-card").forEach(function (card) {
        var rect = null;
        var rafId = 0;
        var pendingX = 0;
        var pendingY = 0;
        var lastX = 0;
        var lastY = 0;

        function flush() {
          rafId = 0;
          if (
            Math.abs(pendingX - lastX) < 0.005 &&
            Math.abs(pendingY - lastY) < 0.005
          )
            return;
          lastX = pendingX;
          lastY = pendingY;
          card.style.transform =
            "translateY(-6px) scale(1.01) rotateX(" +
            -pendingY * 3 +
            "deg) rotateY(" +
            pendingX * 3 +
            "deg)";
        }

        card.addEventListener("mouseenter", function () {
          rect = card.getBoundingClientRect();
          card.style.transition = "none";
        });
        card.addEventListener(
          "mousemove",
          function (evt) {
            if (!rect) return;
            pendingX =
              (evt.clientX - (rect.left + rect.width / 2)) / (rect.width / 2);
            pendingY =
              (evt.clientY - (rect.top + rect.height / 2)) / (rect.height / 2);
            if (!rafId) rafId = requestAnimationFrame(flush);
          },
          { passive: true }
        );
        card.addEventListener("mouseleave", function () {
          if (rafId) {
            cancelAnimationFrame(rafId);
            rafId = 0;
          }
          rect = null;
          lastX = 0;
          lastY = 0;
          card.style.transform = "";
          card.style.transition = "";
        });
      });
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
    GlassParallax.init();
    ClipboardCopy.init();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot);
  } else {
    boot();
  }
})();
