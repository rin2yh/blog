+++
title = "Typography & Markdown"
date = 2026-04-05
categories = ["Showcase"]
tags = ["markdown", "typography"]
description = "Every Markdown construct rendered in the Liquid Glass theme — for visual review."
image = "cover.svg"
+++

This page exists so you can eyeball every Markdown construct in one scroll. If something looks broken, it's a theme bug.

## Headings

# H1 — never used in body, but defined
## H2 — section
### H3 — subsection
#### H4 — minor
##### H5
###### H6

## Inline elements

This paragraph contains **bold**, *italic*, ***bold-italic***, ~~strikethrough~~, `inline code`, and an [external link](https://gohugo.io/). Footnotes work too[^1].

[^1]: Footnotes render at the bottom of the article.

## Lists

Unordered:

- First item
- Second item
  - Nested item
  - Another nested item
- Third item

Ordered:

1. Plan
2. Build
3. Ship

Task list:

- [x] Decide on design tokens
- [x] Implement glass utilities
- [ ] Write the docs

## Blockquote

> Glass is not transparent. It is what you see *through* — and what you see *of* — at the same time.
>
> — someone, probably

## Table

| Token | Light | Dark |
|---|---|---|
| `--glass-white` | `rgba(255,255,255,0.62)` | `rgba(255,255,255,0.12)` |
| `--blur-md` | `blur(18px)` | `blur(18px)` |
| `--shadow-glass` | soft white inset | warm dark inset |

## Horizontal rule

---

## Image

Images render full-width inside the article container, with rounded corners inherited from the glass card.

![placeholder](https://picsum.photos/seed/liquid-glass/1200/600)

## Shortcode — `{{</* qr */>}}`

The theme ships a single shortcode for embedding a QR code generated from arbitrary text. The image is fetched from `api.qrserver.com`, so it's a network dependency, not a bundled asset.

{{< qr text="https://gohugo.io" title="hugo.io" />}}

Self-closing form is accepted; `text` is required, `alt` and `title` default to it.

## Definition list (raw HTML)

<dl>
  <dt>Glassmorphism</dt>
  <dd>A UI style using translucency, blur, and subtle borders to simulate frosted glass.</dd>
  <dt>Liquid Glass</dt>
  <dd>This theme's take: glass + drifting color, with motion that responds to scroll.</dd>
</dl>
