+++
title = "Welcome to Liquid Glass"
date = 2026-04-01
categories = ["Meta"]
tags = ["welcome", "intro"]
description = "A short tour of the Liquid Glass theme — what it ships with and what you can switch on."
+++

**Liquid Glass** is a Stack-inspired Hugo theme built around three ideas:

- **Real translucency.** Cards use `backdrop-filter` so the background actually shows through. No fake blurs.
- **A scene that moves.** Three drifting color blobs animate behind everything via pure CSS keyframes.
- **Zero JS dependencies.** ~5 KB of vanilla JS handles theme toggle, ripple, mobile nav, and code copy.

## Quick tour

This very page is rendered with the theme's default `single.html` layout. Notice the:

- Glass article container with a soft inner reflection
- Sidebar widget with a TOC that highlights the current heading
- Drifting blobs visible at the edges of the viewport
- Theme toggle in the top-right corner — try it

## What ships in the box

| Area | What's included |
|---|---|
| Layouts | `home`, `single`, `list`, `taxonomy`, `term`, `404`, `archives`, `search` |
| Widgets | search, archives, categories, tag-cloud, TOC, twitter-share |
| Markdown | GitHub-style tables, footnotes, syntax highlighting, mermaid |
| i18n | Japanese (`ja.toml`) included; add your own `i18n/<lang>.toml` |
| OpenGraph | Per-page OG + Twitter Card tags |
| Search | Client-side, indexed via `index.json` |

## Configuration shape

The theme reads parameters under `[params]`. The minimum you need is:

```toml
theme = "stack-liquid-glass"

[params]
description = "Your tagline"
mainSections = ["post"]

  [params.sidebar.avatar]
  enabled = true
  local = true
  src = "image/avatar.svg"

  [[params.widgets.homepage]]
  type = "search"
  [[params.widgets.homepage]]
  type = "archives"
  [[params.widgets.page]]
  type = "toc"
```

Everything else is optional. See `exampleSite/hugo.toml` for the full surface.
