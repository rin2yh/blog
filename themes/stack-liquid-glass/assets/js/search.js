/**
 * search.js — incremental client-side search against /index.json
 */
(function () {
  "use strict";

  function escapeHtml(str) {
    return (str || "")
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#39;");
  }

  var MAX_QUERY_LEN = 100;

  function highlight(text, query) {
    if (!query) return escapeHtml(text);
    if (query.length > MAX_QUERY_LEN) query = query.slice(0, MAX_QUERY_LEN);
    var terms = query.split(/\s+/).filter(Boolean);
    if (!terms.length) return escapeHtml(text);
    var pattern = terms
      .map(function (w) {
        return w.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
      })
      .join("|");
    var re = new RegExp("(" + pattern + ")", "gi");
    return text
      .split(re)
      .map(function (part, i) {
        return i % 2 === 1
          ? "<mark>" + escapeHtml(part) + "</mark>"
          : escapeHtml(part);
      })
      .join("");
  }

  function snippet(content, query, len) {
    if (!content) return "";
    if (!query) return content.slice(0, len) + (content.length > len ? "…" : "");
    var lc = content.toLowerCase();
    var first = -1;
    query
      .split(/\s+/)
      .filter(Boolean)
      .forEach(function (w) {
        var idx = lc.indexOf(w.toLowerCase());
        if (idx !== -1 && (first === -1 || idx < first)) first = idx;
      });
    if (first === -1) return content.slice(0, len) + (content.length > len ? "…" : "");
    var start = Math.max(0, first - 40);
    var end = Math.min(content.length, start + len);
    return (start > 0 ? "…" : "") + content.slice(start, end) + (end < content.length ? "…" : "");
  }

  function score(entry, terms) {
    var s = 0;
    var t = (entry.title || "").toLowerCase();
    var c = (entry.content || "").toLowerCase();
    var tags = (entry.tags || []).join(" ").toLowerCase();
    var cats = (entry.categories || []).join(" ").toLowerCase();
    terms.forEach(function (w) {
      if (t.indexOf(w) !== -1) s += 10;
      if (tags.indexOf(w) !== -1) s += 4;
      if (cats.indexOf(w) !== -1) s += 3;
      if (c.indexOf(w) !== -1) s += 1;
    });
    return s;
  }

  function init() {
    var input = document.querySelector("[data-search-input]");
    var resultsEl = document.querySelector("[data-search-results]");
    if (!input || !resultsEl) return;

    var endpoint = input.getAttribute("data-search-endpoint");
    if (!endpoint) return;

    var data = null;
    var dataPromise = null;

    function ensureData() {
      if (data) return Promise.resolve(data);
      if (dataPromise) return dataPromise;
      dataPromise = fetch(endpoint)
        .then(function (r) {
          if (!r.ok) throw new Error("Search index fetch failed");
          return r.json();
        })
        .then(function (json) {
          data = json && Array.isArray(json) ? json : (json && json.entries) || [];
          dataPromise = null;
          return data;
        })
        .catch(function () {
          data = [];
          dataPromise = null;
          return data;
        });
      return dataPromise;
    }

    function render(query) {
      var q = query.trim();
      resultsEl.innerHTML = "";
      if (!q) {
        var emptyMsg = (window.lgI18n && window.lgI18n.searchEmpty) || "キーワードを入力してください。";
        resultsEl.innerHTML =
          '<p class="search-page__empty">' + escapeHtml(emptyMsg) + '</p>';
        return;
      }
      ensureData().then(function (entries) {
        if (input.value.trim() !== q) return;
        var terms = q.toLowerCase().split(/\s+/).filter(Boolean);
        var matches = entries
          .map(function (e) {
            return { entry: e, score: score(e, terms) };
          })
          .filter(function (m) {
            return m.score > 0;
          })
          .sort(function (a, b) {
            return b.score - a.score;
          })
          .slice(0, 30);

        if (!matches.length) {
          var noResultsMsg = (window.lgI18n && window.lgI18n.searchNoResults) || "該当する記事がありません。";
          resultsEl.innerHTML =
            '<p class="search-page__empty">' + escapeHtml(noResultsMsg) + '</p>';
          return;
        }

        matches.forEach(function (m) {
          var e = m.entry;
          var el = document.createElement("a");
          el.className = "glass search-result";
          el.href = e.url;
          var meta = [];
          if (e.date) meta.push(escapeHtml(e.date));
          if (e.categories && e.categories.length) {
            meta.push(escapeHtml(e.categories.join(" / ")));
          }
          el.innerHTML =
            '<div class="search-result__title">' +
            highlight(e.title || "", q) +
            "</div>" +
            (meta.length
              ? '<div class="search-result__meta">' + meta.join(" · ") + "</div>"
              : "") +
            '<div class="search-result__excerpt">' +
            highlight(snippet(e.content || e.summary || "", q, 200), q) +
            "</div>";
          resultsEl.appendChild(el);
        });
      });
    }

    var t = null;
    input.addEventListener("input", function () {
      clearTimeout(t);
      t = setTimeout(function () {
        render(input.value);
      }, 120);
    });

    var form = input.closest("[data-search-form]");
    if (form) {
      form.addEventListener("submit", function (e) {
        e.preventDefault();
      });
    }

    var q = new URLSearchParams(window.location.search).get("q");
    if (q && !input.value) input.value = q;

    if (input.value) render(input.value);
  }

  if (document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", init);
  } else {
    init();
  }
})();
