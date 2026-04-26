/**
 * toc.js — highlight active TOC item as the reader scrolls
 */
(function () {
  "use strict";

  function init() {
    var toc = document.querySelector(".toc__inner");
    if (!toc) return;

    var links = Array.from(toc.querySelectorAll('a[href^="#"]'));
    if (!links.length) return;

    var headingMap = new Map();
    links.forEach(function (link) {
      var id = decodeURIComponent(link.getAttribute("href").slice(1));
      var heading = document.getElementById(id);
      if (heading) headingMap.set(heading, link);
    });

    if (!headingMap.size) return;

    var setActive = function (entry) {
      links.forEach(function (a) {
        a.classList.remove("is-active");
      });
      var link = headingMap.get(entry);
      if (link) link.classList.add("is-active");
    };

    var observer = new IntersectionObserver(
      function (entries) {
        var visible = entries
          .filter(function (e) {
            return e.isIntersecting;
          })
          .sort(function (a, b) {
            return b.intersectionRatio - a.intersectionRatio;
          });
        if (visible.length) setActive(visible[0].target);
      },
      {
        rootMargin: "-100px 0px -55% 0px",
        threshold: [0, 0.4, 1]
      }
    );

    headingMap.forEach(function (_, heading) {
      observer.observe(heading);
    });
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
