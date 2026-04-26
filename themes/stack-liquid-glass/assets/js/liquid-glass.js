(function () {
  "use strict";

  var GlassTheme = {
    KEY: "liquid-glass-theme",
    init: function () {
      var btn = document.querySelector("[data-theme-toggle]");
      var root = document.documentElement;
      if (!btn) return;
      var i18n = window.lgI18n || {};
      var updateLabel = function (theme) {
        btn.setAttribute(
          "aria-label",
          theme === "dark"
            ? i18n.switchToLight || "Switch to light mode"
            : i18n.switchToDark || "Switch to dark mode"
        );
      };
      updateLabel(root.getAttribute("data-theme"));
      btn.addEventListener("click", function () {
        var next = root.getAttribute("data-theme") === "light" ? "dark" : "light";
        root.setAttribute("data-theme", next);
        try {
          localStorage.setItem(GlassTheme.KEY, next);
        } catch (e) {}
        updateLabel(next);
      });
    }
  };

  var GlassRipple = {
    init: function () {
      document.addEventListener("click", function (evt) {
        var target = evt.target.closest(".glass-btn, .article-card");
        if (!target) return;
        if (target.querySelectorAll("[data-ripple]").length > 5) return;
        var ripple = document.createElement("span");
        ripple.setAttribute("data-ripple", "");
        ripple.setAttribute("aria-hidden", "true");
        var rect = target.getBoundingClientRect();
        var size = Math.max(rect.width, rect.height) * 1.2;
        var x = evt.clientX - rect.left - size / 2;
        var y = evt.clientY - rect.top - size / 2;
        Object.assign(ripple.style, {
          left: x + "px",
          top: y + "px",
          width: size + "px",
          height: size + "px"
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

  var toastTimer = 0;
  var toastFadeTimer = 0;
  function showToast(message) {
    var toast = document.querySelector(".lg-toast");
    if (!toast) {
      toast = document.createElement("div");
      toast.className = "lg-toast";
      toast.setAttribute("role", "status");
      document.body.appendChild(toast);
    }
    toast.textContent = message;
    clearTimeout(toastTimer);
    clearTimeout(toastFadeTimer);
    requestAnimationFrame(function () {
      toast.classList.add("is-visible");
    });
    toastTimer = setTimeout(function () {
      toast.classList.remove("is-visible");
      toastFadeTimer = setTimeout(function () {
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
