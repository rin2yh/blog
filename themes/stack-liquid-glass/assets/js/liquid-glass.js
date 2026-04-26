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

  var MobileNav = {
    init: function () {
      var btn = document.querySelector("[data-mobile-nav-toggle]");
      var sidebar = document.querySelector(".left-sidebar");
      if (!btn || !sidebar) return;
      var body = document.body;
      function setOpen(open) {
        sidebar.classList.toggle("is-open", open);
        body.classList.toggle("nav-open", open);
        btn.setAttribute("aria-expanded", open ? "true" : "false");
      }
      btn.addEventListener("click", function (evt) {
        evt.stopPropagation();
        setOpen(!sidebar.classList.contains("is-open"));
      });
      sidebar.addEventListener("click", function (evt) {
        if (evt.target.closest("a")) setOpen(false);
      });
      document.addEventListener("click", function (evt) {
        if (!sidebar.classList.contains("is-open")) return;
        if (evt.target.closest(".left-sidebar")) return;
        if (evt.target.closest("[data-mobile-nav-toggle]")) return;
        setOpen(false);
      });
      document.addEventListener("keydown", function (evt) {
        if (evt.key === "Escape" && sidebar.classList.contains("is-open")) {
          setOpen(false);
          btn.focus();
        }
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

  function showToast(message) {
    var toast = document.createElement("div");
    toast.className = "lg-toast";
    toast.setAttribute("role", "status");
    toast.textContent = message;
    Object.assign(toast.style, {
      position: "fixed",
      left: "50%",
      bottom: "32px",
      transform: "translate(-50%, 8px)",
      padding: "10px 18px",
      borderRadius: "999px",
      background: "rgba(15, 23, 42, 0.85)",
      color: "#fff",
      fontSize: "0.85rem",
      boxShadow: "0 8px 24px rgba(0,0,0,0.25)",
      backdropFilter: "blur(12px)",
      opacity: "0",
      transition: "opacity 200ms ease, transform 200ms ease",
      zIndex: "9999",
      pointerEvents: "none"
    });
    document.body.appendChild(toast);
    requestAnimationFrame(function () {
      toast.style.opacity = "1";
      toast.style.transform = "translate(-50%, 0)";
    });
    setTimeout(function () {
      toast.style.opacity = "0";
      toast.style.transform = "translate(-50%, 8px)";
      setTimeout(function () {
        toast.remove();
      }, 220);
    }, 1800);
  }

  var ClipboardCopy = {
    init: function () {
      document.addEventListener("click", function (evt) {
        var btn = evt.target.closest("[data-copy]");
        if (!btn) return;
        var text = btn.getAttribute("data-copy");
        var ok = btn.getAttribute("data-copy-success") || "Copied";
        var ng = btn.getAttribute("data-copy-error") || "Copy failed";
        if (!navigator.clipboard) {
          showToast(ng);
          return;
        }
        navigator.clipboard.writeText(text).then(
          function () {
            showToast(ok);
          },
          function () {
            showToast(ng);
          }
        );
      });
    }
  };

  function boot() {
    GlassTheme.init();
    GlassRipple.init();
    MobileNav.init();
    ScrollHoverGuard.init();
    ClipboardCopy.init();
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", boot);
  } else {
    boot();
  }
})();
