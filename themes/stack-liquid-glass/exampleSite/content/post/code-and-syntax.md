+++
title = "Code blocks & syntax highlighting"
date = 2026-04-08
categories = ["Showcase"]
tags = ["code", "syntax"]
description = "How code looks in Liquid Glass — fenced blocks, inline code, copy button, and a few languages side by side."
+++

Fenced code blocks are highlighted via Hugo's Chroma. Each block has a copy button in the top-right that reads from the rendered DOM, so it survives any post-processing.

## Go

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func fetchWithTimeout(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return "", errors.New("timed out")
	case <-time.After(500 * time.Millisecond):
		return fmt.Sprintf("response from %s", url), nil
	}
}
```

## TypeScript

```ts
type Token = `--${string}`;

interface GlassConfig<T extends Token = Token> {
  blur: T;
  background: T;
  border: T;
}

const defaults: GlassConfig = {
  blur: "--blur-md",
  background: "--glass-white",
  border: "--glass-border",
};
```

## CSS

```css
.glass {
  backdrop-filter: var(--blur-md);
  -webkit-backdrop-filter: var(--blur-md);
  background: var(--glass-white);
  border: 1px solid var(--glass-border);
  box-shadow: var(--shadow-glass);
}

.glass::before {
  content: "";
  position: absolute;
  inset: 0;
  background: var(--reflection-top);
  pointer-events: none;
  border-radius: inherit;
}
```

## Shell

```bash
hugo server --buildDrafts --bind 0.0.0.0
```

Inline code like `--blur-md` keeps the same monospace stack but skips the box and copy button.
