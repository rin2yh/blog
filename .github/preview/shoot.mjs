// public/ をローカル配信し、変更された各記事ページを Playwright(Chromium) でスクショする。
// デスクトップ/モバイル幅の PNG を OUT_DIR に書き出し、PR コメント本文 Markdown を COMMENT_FILE に書き出す。
//
// 環境変数:
//   URLS        必須  改行区切りの /post/<slug>/ パス一覧 (changed-urls.sh の出力)
//   PUBLIC_DIR  任意  配信するディレクトリ (既定: public)
//   PORT        任意  ローカルサーバのポート (既定: 1313)
//   OUT_DIR     必須  PNG 出力先ディレクトリ
//   COMMENT_FILE必須  生成するコメント Markdown の出力先
//   RAW_PREFIX  必須  画像を参照する raw URL の接頭辞
//                     例: https://raw.githubusercontent.com/rin2yh/blog/pr-preview-assets/pr-79/abc1234
//   PR          任意  PR 番号 (見出し表示用)
//   SHA7        任意  短縮 SHA (脚注表示用)

import { createServer } from 'node:http';
import { mkdir, writeFile } from 'node:fs/promises';
import { join } from 'node:path';

const PUBLIC_DIR = process.env.PUBLIC_DIR || 'public';
const PORT = Number(process.env.PORT || 1313);
const OUT_DIR = process.env.OUT_DIR;
const COMMENT_FILE = process.env.COMMENT_FILE;
const RAW_PREFIX = (process.env.RAW_PREFIX || '').replace(/\/$/, '');
const PR = process.env.PR || '';
const SHA7 = process.env.SHA7 || '';

if (!OUT_DIR || !COMMENT_FILE) {
  console.error('OUT_DIR and COMMENT_FILE are required');
  process.exit(1);
}

const urls = (process.env.URLS || '')
  .split('\n')
  .map((s) => s.trim())
  .filter(Boolean);

// sirv で public/ を配信する (MIME 判定・index フォールバック・トラバーサル対策を内包)。
// sirv は記事変更があるときだけ読み込む (no-changes 経路では未インストールのため)。
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

async function main() {
  await mkdir(OUT_DIR, { recursive: true });

  if (urls.length === 0) {
    await writeFile(COMMENT_FILE, body(['変更された記事ページはありません。']));
    console.log('No article URLs to shoot.');
    return;
  }

  // playwright は記事変更があるときだけ読み込む (no-changes 経路では未インストールのため)
  const { chromium } = await import('playwright');
  const server = await startServer();
  const browser = await chromium.launch();
  const sections = [];

  try {
    for (const u of urls) {
      const slug = slugFromUrl(u);
      const target = `http://127.0.0.1:${PORT}${u}`;
      const shots = [];
      let notFound = false;

      for (const vp of VIEWPORTS) {
        const page = await browser.newPage({ viewport: { width: vp.width, height: vp.height } });
        try {
          const resp = await page.goto(target, { waitUntil: 'load', timeout: 30000 });
          if (!resp || resp.status() >= 400) {
            notFound = true;
            break;
          }
          const file = `${vp.key}--${slug}.png`;
          await page.screenshot({ path: join(OUT_DIR, file), fullPage: true });
          shots.push({ ...vp, file });
        } catch (err) {
          // goto/screenshot の失敗 (タイムアウト等) で全体を落とさず、このURLだけスキップする
          console.error(`Failed to capture ${target} (${vp.label}):`, err);
          notFound = true;
          break;
        } finally {
          await page.close().catch(() => {});
        }
      }

      if (notFound) {
        sections.push(`<details><summary><code>${u}</code></summary>\n\n> ⚠️ ページが見つかりませんでした (\`${target}\`)。ドラフト未ビルドやパス不一致の可能性があります。\n\n</details>`);
        console.warn(`Not found: ${target}`);
        continue;
      }

      const imgs = shots
        .map((s) => `**${s.label}**\n\n![${s.key}](${RAW_PREFIX}/${s.file})`)
        .join('\n\n');
      sections.push(`<details open><summary><code>${u}</code></summary>\n\n${imgs}\n\n</details>`);
      console.log(`Shot: ${target} (${shots.length} viewport(s))`);
    }
  } finally {
    await browser.close().catch(() => {});
    server.close();
  }

  await writeFile(COMMENT_FILE, body(sections));
}

function body(sections) {
  const pr = PR ? ` (PR #${PR})` : '';
  const sha = SHA7 ? ` \`${SHA7}\`` : '';
  return (
    `<!-- article-preview -->\n## 📝 記事プレビュー${pr}\n\n` +
    `${sections.join('\n\n')}\n\n` +
    `<sub>下書き・未来日付の記事も含めてビルドしています。コミット${sha} 時点のプレビューです。</sub>\n`
  );
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
