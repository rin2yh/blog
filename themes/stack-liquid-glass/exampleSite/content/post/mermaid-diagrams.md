+++
title = "Mermaid diagrams"
date = 2026-04-12
categories = ["Showcase"]
tags = ["mermaid", "diagrams"]
description = "Mermaid is loaded lazily — only when a page actually contains a `mermaid` fenced block."
+++

The theme detects mermaid code fences during render and injects the mermaid runtime *only on those pages*, so unrelated posts pay zero cost. The runtime auto-picks `default` or `dark` based on the active theme attribute.

## Sequence

```mermaid
sequenceDiagram
    participant Reader
    participant Browser
    participant Hugo
    Reader->>Browser: requests /post/welcome/
    Browser->>Hugo: serves HTML
    Hugo-->>Browser: HTML + bundled CSS
    Browser-->>Reader: paints
    Note over Browser: data-theme="dark"
```

## Flowchart

```mermaid
flowchart LR
    A[Markdown] -->|goldmark| B(HTML)
    B --> C{Has mermaid?}
    C -->|yes| D[Inject runtime]
    C -->|no| E[Skip]
    D --> F((Render diagram))
```

## Notes

- The runtime is loaded from a CDN ESM URL. Self-host it by replacing the `import` URL in `layouts/_default/baseof.html`.
- Diagram colors follow mermaid's built-in `default` / `dark` themes; for finer control wrap your blocks in `%%{init: ...}%%` directives.
