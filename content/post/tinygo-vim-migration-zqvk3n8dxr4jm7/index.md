+++
date = '2026-07-06T13:00:00+09:00'
draft = true
title = 'tinygo.nvim から tinygo.vim に乗り換えた'
categories = ['tech']
tags = ['neovim', 'tinygo', 'go']
+++

## TL;DR

- Wio Terminal 向けに TinyGo を書き始めたところ、gopls が `machine` パッケージを解決できず補完が効かない状態になりました
- 原因は 2 点ありました。`pcolladosoto/tinygo.nvim` の LSP 再起動が attach 済みクライアントを止めないこと、および mise の `go` shim が GOROOT env を上書きすることです
- `sago35/tinygo.vim` に乗り換え、`env PATH=<mise-go-dir>:...` で gopls コマンドをラップすることで解消しました

## 発端

`tinygobook` というディレクトリを切り、Wio Terminal 向けに TinyGo を書き始めました。`main.go` は以下のように最小構成です。

```go
package main

import "machine"

func main() {
	led := machine.LED
	_ = led
}
```

`.tinygo.json` にターゲットを置いておきます。

```json
{"target": "wioterminal"}
```

`tinygo flash -target=wioterminal` は通るのですが、Neovim で開くと gopls が以下のエラーを返しました。

```
could not import machine (no required module provides package "machine")
```

当然 `machine.LED` の補完も効きません。組み込みを書く上で `machine.` から先が見えないのは大きな痛手です。

## 最初の構成: `pcolladosoto/tinygo.nvim`

もともと dotfiles には `pcolladosoto/tinygo.nvim` を仕込んでおり、`.tinygo.json` を検出して `TinyGoSetTarget <target>` を叩く導線を用意していました。プラグイン本体の `applyConfigFile` は相対パスで `.tinygo.json` を探すため、nvim の cwd がプロジェクト外の場合に検出漏れが発生します。そのため bufname から上向きに探索する版に差し替えて使っていました。

これで target 切替は走るはずなのですが、それでも補完は復活しませんでした。

## 犯人 その 1: `vim.lsp.enable(false)` は attach 済みを止めない

`pcolladosoto/tinygo.nvim` の `setTarget` は、要約すると次のような処理です。

```lua
vim.lsp.config('gopls', { cmd_env = { GOROOT = ..., GOFLAGS = ... } })
vim.lsp.enable('gopls', false)
vim.lsp.enable('gopls', true)
```

Neovim 0.11 の native LSP における `vim.lsp.enable(name, false)` は、**FileType autocmd を解除するだけ**であり、既にバッファに attach 済みのクライアントを stop してはくれません。

流れとしては以下のようになります。

1. Neovim 起動 → `FileType go` → 通常 GOROOT の gopls が起動して attach
2. `.tinygo.json` を検出して `TinyGoSetTarget` → `cmd_env` は書き換わる
3. しかし attach 済みの gopls は古い GOROOT のまま動き続ける
4. 結果として補完は復旧しない

`vim.lsp.get_clients({ name = 'gopls' })` 経由で明示的に `client:stop(true)` を叩き、その後 `:edit` 相当で FileType を再発火させれば直ります。実際、プラグインを monkey-patch して stop と再 attach まで叩き込むことで `machine` は解決できるようになりました。

とはいえ、プラグインの内部挙動に依存した継ぎ接ぎであり、上流のリファクタで簡単に壊れます。もう少し安定した構成にしたいところでした。

## 犯人 その 2: mise の `go` shim が GOROOT env を上書きする

target 切替後の GOROOT が gopls に届いているかを `:LspInfo` と `cmd_env` で確認したところ、gopls プロセスの env には正しい TinyGo overlay GOROOT (`~/Library/Caches/tinygo/goroot-<hash>`) が入っていました。にもかかわらず、gopls 内部から `go env` を叩くと**素の Go の GOROOT が返ってくる**という現象が見られました。

原因は mise の shim です。`mise` は `~/.local/share/mise/shims/go` のようなラッパースクリプトを PATH に置いてバージョン管理を行っていますが、この shim は env の GOROOT を尊重しません。gopls は内部で `go env GOROOT` を実行してモジュール解決するため、いくら親プロセスから GOROOT を渡しても shim の内側でリセットされてしまいます。

回避策は shim を経由させないことです。`mise which go` で直パスが取得できるので、gopls 起動時に PATH 先頭にその bin ディレクトリを差し込みます。

```lua
local go_dir = vim.fs.dirname(vim.trim(vim.fn.system('mise which go')))
M.cmd = { 'env', 'PATH=' .. go_dir .. ':' .. vim.env.PATH, 'gopls' }
```

これで gopls 内部の `go env` は TinyGo overlay GOROOT を素直に参照するようになります。

## 乗り換え先: `sago35/tinygo.vim`

犯人 1 の対処を継ぎ接ぎで続けるよりも、target 切替時の LSP 再 attach をきちんと扱うプラグインに移行した方が早いと判断しました。`sago35/tinygo.vim` には以下の利点があります。

- `:TinygoTarget <target>` で GOROOT / GOOS / GOARCH / GOFLAGS をすべて LSP config に流し込む
- attach 済みクライアントを実際に stop する
- ターゲット補完 (`tinygo#TinygoTargets`) も提供している

ただし Neovim 0.11 の native LSP 側で「client:stop 後の再 attach」を叩いてくれないケースがあるため、そこだけラッパーを噛ませます。

```lua
vim.pack.add({ 'https://github.com/sago35/tinygo.vim' })

vim.api.nvim_create_user_command('Tinygo', function(opts)
  vim.cmd.TinygoTarget(opts.args)
  local buf = vim.api.nvim_get_current_buf()
  vim.schedule(function()
    vim.api.nvim_exec_autocmds('FileType', { buffer = buf })
  end)
end, {
  nargs = 1,
  complete = function(arglead) return vim.fn['tinygo#TinygoTargets'](arglead, '', 0) end,
})
```

`FileType` を再発火させれば、`vim.lsp.enable('gopls', true)` 側の autocmd が拾って、新しい設定で gopls を attach し直してくれます。

## `.tinygo.json` の自動検出

`pcolladosoto/tinygo.nvim` の `applyConfigFile` 相当は自前で書きました。`FileType go` のたびに上向きで `.tinygo.json` を探し、見つかれば `:Tinygo <target>` を叩きます。同一ファイル・同一 target であれば 1 セッション 1 回のみ実行するようガードしています。

```lua
local tinygo_applied = {}
vim.api.nvim_create_autocmd('FileType', {
  pattern = 'go',
  group = vim.api.nvim_create_augroup('tinygo-auto', {}),
  callback = function(args)
    local search = vim.fs.dirname(vim.api.nvim_buf_get_name(args.buf))
    local found = vim.fs.find({ '.tinygo.json' }, { upward = true, path = search })
    if vim.tbl_isempty(found) then return end
    local f = io.open(found[1], 'r')
    if not f then return end
    local raw = f:read('a')
    f:close()
    local ok, cfg = pcall(vim.json.decode, raw)
    if not ok or type(cfg) ~= 'table' or cfg.target == nil then return end
    if tinygo_applied[found[1]] == cfg.target then return end
    tinygo_applied[found[1]] = cfg.target
    vim.cmd({ cmd = 'Tinygo', args = { cfg.target } })
  end,
})
```

## ついでに直したところ

`tinygobook` は `.git` を切っていなかったため、gopls の `root_markers` が効かずワークスペースルートを取り違えていました。`go.mod` も root marker に追加しています。

```lua
M.root_markers = { 'go.mod' }
```

## まとめ

- `pcolladosoto/tinygo.nvim` から `sago35/tinygo.vim` へ移行しました
- Neovim 0.11 の native LSP 側は `FileType` 再発火のラッパーで補強しました
- mise の go shim を PATH で迂回することで、gopls 内部の `go env` を TinyGo overlay に寄せました
- `.tinygo.json` を開くだけで target 切替と gopls 再 attach まで自動で走るようになりました

補完が復活してからは、組み込みコードを書く負担が大きく減りました。mise と TinyGo と Neovim を組み合わせて `machine` が解決できない状況に遭遇している方は、おおむね同じ罠にはまっているのではないかと思います。
