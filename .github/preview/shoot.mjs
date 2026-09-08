// public/ をローカル配信し、変更された記事ページを Playwright (Chromium) でスクショする。
// デスクトップ/モバイル幅の PNG を OUT_DIR に書き出す。

import { createServer } from 'node:http';
import { appendFile, mkdir } from 'node:fs/promises';
import { join } from 'node:path';

const URL = process.env.URL;
const PUBLIC_DIR = process.env.PUBLIC_DIR || 'public';
const PORT = Number(process.env.PORT || 1313);
const OUT_DIR = process.env.OUT_DIR;
const GITHUB_OUTPUT = process.env.GITHUB_OUTPUT;

// sirv で public/ を配信する (MIME 判定・index フォールバック・トラバーサル対策を内包)。
async function startServer() {
  const { default: sirv } = await import('sirv');
  const serve = sirv(PUBLIC_DIR, { dev: true });
  const server = createServer((req, res) =>
    serve(req, res, () => {
      res.statusCode = 404;
      res.end('Not Found');
    })
  );
  return new Promise((ready) => server.listen(PORT, '127.0.0.1', () => ready(server)));
}

function getSlugFromUrl(u) {
  // /post/foo/ -> foo
  return u.split('/').filter(Boolean).at(-1) || 'index';
}

const VIEWPORTS = [
  { key: 'desktop', label: 'Desktop', width: 1280, height: 900 },
  { key: 'mobile', label: 'Mobile', width: 390, height: 844 },
];

async function captureScreenshot(browser, target, slug, viewport) {
  const page = await browser.newPage({ viewport: { width: viewport.width, height: viewport.height } });
  try {
    const resp = await page.goto(target, { waitUntil: 'load', timeout: 30000 });
    if (!resp?.ok()) {
      throw new Error(`Failed to load ${target}: ${resp?.status()}`);
    }

    const file = `${viewport.key}--${slug}.png`;
    const path = join(OUT_DIR, file);
    await page.screenshot({ path, fullPage: true });
    await appendFile(GITHUB_OUTPUT, `${viewport.key}=${path}\n`);
    console.log(`Shot: ${target} (${viewport.label})`);
  } finally {
    await page.close();
  }
}

await mkdir(OUT_DIR, { recursive: true });
const { chromium } = await import('playwright');
const server = await startServer();
const browser = await chromium.launch();
const slug = getSlugFromUrl(URL);
const target = `http://127.0.0.1:${PORT}${URL}`;

try {
  for (const viewport of VIEWPORTS) {
    await captureScreenshot(browser, target, slug, viewport);
  }
} finally {
  await browser.close();
  server.close();
}
