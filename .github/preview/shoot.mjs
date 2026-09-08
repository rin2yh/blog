// public/ をローカル配信し、変更された各記事ページを Playwright(Chromium) でスクショする。
// デスクトップ/モバイル幅の PNG を OUT_DIR に書き出すか、Artifact へのリンクを含む
// PR コメント本文 Markdown を COMMENT_FILE に書き出す。
//
// 環境変数:
//   URLS               任意  改行区切りの /post/<slug>/ パス一覧 (スクリーンショット生成時)
//   URLS_JSON          任意  上記パス一覧の JSON 配列 (コメント生成時)
//   PUBLIC_DIR         任意  配信するディレクトリ (既定: public)
//   PORT               任意  ローカルサーバのポート (既定: 1313)
//   OUT_DIR            任意  PNG 出力先ディレクトリ (指定時はスクリーンショットを生成)
//   COMMENT_FILE       任意  コメント Markdown の出力先 (指定時はコメントを生成)
//   ARTIFACT_URLS_FILE 必須  ファイル名を Artifact URL に対応させた JSON (コメント生成時)
//   PR                 任意  PR 番号 (見出し表示用)
//   SHA7               任意  短縮 SHA (脚注表示用)

import { createServer } from 'node:http';
import { appendFile, mkdir, readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const PUBLIC_DIR = process.env.PUBLIC_DIR || 'public';
const PORT = Number(process.env.PORT || 1313);
const OUT_DIR = process.env.OUT_DIR || '';
const COMMENT_FILE = process.env.COMMENT_FILE || '';
const ARTIFACT_URLS_FILE = process.env.ARTIFACT_URLS_FILE || '';
const PR = process.env.PR || '';
const SHA7 = (process.env.SHA7 || '').slice(0, 7);

if (!OUT_DIR && !COMMENT_FILE) {
  console.error('OUT_DIR or COMMENT_FILE is required');
  process.exit(1);
}
if (COMMENT_FILE && !ARTIFACT_URLS_FILE) {
  console.error('ARTIFACT_URLS_FILE is required when COMMENT_FILE is set');
  process.exit(1);
}

const urls = process.env.URLS_JSON
  ? JSON.parse(process.env.URLS_JSON).filter(Boolean)
  : (process.env.URLS || '').split('\n').map((s) => s.trim()).filter(Boolean);

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

function slugFromUrl(u) {
  // /post/foo/ -> foo
  return u.split('/').filter(Boolean).at(-1) || 'index';
}

const VIEWPORTS = [
  { key: 'desktop', label: 'Desktop', width: 1280, height: 900 },
  { key: 'mobile', label: 'Mobile', width: 390, height: 844 },
];

async function captureScreenshots() {
  await mkdir(OUT_DIR, { recursive: true });
  if (urls.length === 0) {
    console.log('No article URLs to shoot.');
    return;
  }

  const { chromium } = await import('playwright');
  const server = await startServer();
  const browser = await chromium.launch();

  try {
    for (const u of urls) {
      const slug = slugFromUrl(u);
      const target = `http://127.0.0.1:${PORT}${u}`;

      for (const vp of VIEWPORTS) {
        const page = await browser.newPage({ viewport: { width: vp.width, height: vp.height } });
        try {
          const resp = await page.goto(target, { waitUntil: 'load', timeout: 30000 });
          if (!resp || resp.status() >= 400) {
            console.warn(`Not found: ${target} (${vp.label})`);
            continue;
          }

          const file = `${vp.key}--${slug}.png`;
          const path = join(OUT_DIR, file);
          await page.screenshot({ path, fullPage: true });
          if (process.env.GITHUB_OUTPUT) {
            await appendFile(process.env.GITHUB_OUTPUT, `${vp.key}=${path}\n`);
          }
          console.log(`Shot: ${target} (${vp.label})`);
        } catch (err) {
          // goto/screenshot の失敗 (タイムアウト等) で全体を落とさず、この viewport だけスキップする
          console.error(`Failed to capture ${target} (${vp.label}):`, err);
        } finally {
          await page.close().catch(() => {});
        }
      }
    }
  } finally {
    await browser.close().catch(() => {});
    server.close();
  }
}

async function buildComment() {
  const artifactUrls = JSON.parse(await readFile(ARTIFACT_URLS_FILE, 'utf8'));
  const rows = urls.map((u) => {
    const slug = slugFromUrl(u);
    const links = VIEWPORTS.map((vp) => {
      const file = `${vp.key}--${slug}.png`;
      const url = artifactUrls[file];
      return url ? `[開く](${url})` : '⚠️ 生成失敗';
    });
    return `| \`${u}\` | ${links.join(' | ')} |`;
  });

  const content = rows.length === 0
    ? '変更された記事ページはありません。'
    : ['| 記事 | Desktop | Mobile |', '| --- | --- | --- |', ...rows].join('\n');
  await writeFile(COMMENT_FILE, body(content));
}

function body(content) {
  const pr = PR ? ` (PR #${PR})` : '';
  const sha = SHA7 ? ` \`${SHA7}\`` : '';
  return (
    `<!-- article-preview -->\n## 📝 記事プレビュー${pr}\n\n` +
    `${content}\n\n` +
    `<sub>下書き・未来日付の記事も含めてビルドしています。コミット${sha} 時点のプレビューです。</sub>\n`
  );
}

async function main() {
  if (OUT_DIR) {
    await captureScreenshots();
  }
  if (COMMENT_FILE) {
    await buildComment();
  }
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
