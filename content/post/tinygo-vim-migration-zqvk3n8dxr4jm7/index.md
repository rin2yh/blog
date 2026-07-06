+++
date = '2026-07-06T13:00:00+09:00'
draft = true
title = 'tinygo.nvim から tinygo.vim に乗り換えた'
categories = ['tech']
tags = ['neovim', 'tinygo', 'go']
+++

## TL;DR

- Wio Terminal で TinyGo を書き始めたら gopls が `machine` パッケージを解決できず補完が死んだ
- 犯人は 2 つ: (1) `pcolladosoto/tinygo.nvim` の LSP 再起動が中途半端 (2) mise の `go` shim が GOROOT env を無視して上書きしていた
- `sago35/tinygo.vim` に乗り換え、`env PATH=<mise-go-dir>:...` で gopls をラップして解決

## 発端

`tinygobook` というディレクトリを切って Wio Terminal で TinyGo を触り始めた。`main.go` はこれだけ。

```go
package main

import "machine"

func main() {
	led := machine.LED
	_ = led
}
```

`.tinygo.json` にターゲットを置いておく。

```json
{"target": "wioterminal"}
```

`tinygo flash -target=wioterminal` は通るのに、Neovim で開くと gopls が怒る。

```
could not import machine (no required module provides package "machine")
```

`machine.LED` の補完も当然効かない。組み込みを書くのに `machine.` から先が全部見えないのは厳しい。

## 最初の構成: `pcolladosoto/tinygo.nvim`

もともと dotfiles には `pcolladosoto/tinygo.nvim` を仕込んでいて、`.tinygo.json` を検出して `TinyGoSetTarget <target>` を叩く導線があった。プラグイン本体の `applyConfigFile` は相対パスで `.tinygo.json` を探すため、nvim の cwd がプロジェクト外だと検出漏れする。そこで bufname から上向き探索する版に差し替えて使っていた。

これで target 切替は走るはずなのに、それでも補完が復活しない。

## 犯人 その 1: `vim.lsp.enable(false)` は attach 済みを止めない

`pcolladosoto/tinygo.nvim` の `setTarget` はざっくりこう。

```lua
vim.lsp.config('gopls', { cmd_env = { GOROOT = ..., GOFLAGS = ... } })
vim.lsp.enable('gopls', false)
vim.lsp.enable('gopls', true)
```

Neovim 0.11 の native LSP で `vim.lsp.enable(name, false)` が何をするかというと、**FileType autocmd を解除するだけ**で、既にバッファに attach 済みのクライアントを stop してくれない。

つまり:

1. Neovim 起動 → `FileType go` → 通常 GOROOT の gopls が起動して attach
2. `.tinygo.json` を検出して `TinyGoSetTarget` → `cmd_env` は書き換わる
3. でも attach 済みの gopls は古い GOROOT で動き続ける
4. 補完は死んだまま

`vim.lsp.get_clients({ name = 'gopls' })` 経由で明示的に `client:stop(true)` を叩き、その後 `:edit` 相当で FileType を再発火させれば直る。実際、プラグインを monkey-patch して stop + 再 attach まで叩き込むと `machine` は解決できるようになった。

ただこれ、プラグインの内部挙動に依存した継ぎ接ぎで、上流のリファクタで簡単に壊れる。もう少しちゃんとしたい。

## 犯人 その 2: mise の `go` shim が GOROOT env を上書きする

target を切り替えた後の GOROOT がちゃんと gopls に届いているか `:LspInfo` と `cmd_env` で追いかけたところ、gopls プロセスの env には正しい TinyGo overlay GOROOT (`~/Library/Caches/tinygo/goroot-<hash>`) が入っている。なのに gopls 内部の `go env` を叩くと**素の Go の GOROOT が返ってきていた**。

原因は mise の shim。`mise` は `~/.local/share/mise/shims/go` みたいなラッパースクリプトを PATH に置いてバージョン管理をやっているが、この shim が env の GOROOT を尊重しない。gopls は内部で `go env GOROOT` してモジュール解決するので、いくら親プロセスから GOROOT を渡しても shim の内側でリセットされてしまう。

回避策は shim を経由させないこと。`mise which go` で直パスが取れるので、gopls 起動時に PATH 先頭にその bin ディレクトリを差し込む。

```lua
local go_dir = vim.fs.dirname(vim.trim(vim.fn.system('mise which go')))
M.cmd = { 'env', 'PATH=' .. go_dir .. ':' .. vim.env.PATH, 'gopls' }
```

これで gopls 内部の `go env` は TinyGo overlay GOROOT を素直に見るようになる。

## 乗り換え先: `sago35/tinygo.vim`

犯人 1 の対処を継ぎ接ぎで続けるより、target 切替の LSP 再 attach をきちんと扱うプラグインに移った方が早い。`sago35/tinygo.vim` は:

- `:TinygoTarget <target>` で GOROOT / GOOS / GOARCH / GOFLAGS をすべて LSP config に流し込む
- attach 済みクライアントを実際に stop する
- ターゲット補完 (`tinygo#TinygoTargets`) も提供している

ただ Neovim 0.11 の native LSP 側で「client:stop の後の再 attach」を叩いてくれないケースがあるので、そこだけラッパーを噛ませる。

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

`FileType` を再発火させれば `vim.lsp.enable('gopls', true)` 側の autocmd が拾って、新設定で gopls を attach し直してくれる。

## `.tinygo.json` の自動検出

`pcolladosoto/tinygo.nvim` の `applyConfigFile` 相当は自前で書く。`FileType go` のたびに上向きで `.tinygo.json` を探し、あれば `:Tinygo <target>` を叩く。同一ファイル・同一 target は 1 セッション 1 回だけ。

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

`tinygobook` は `.git` を切っていなかったので、gopls の `root_markers` が効かずワークスペースルートを取り違えていた。`go.mod` も root marker に足す。

```lua
M.root_markers = { 'go.mod' }
```

## まとめ

- `pcolladosoto/tinygo.nvim` → `sago35/tinygo.vim` へ
- Neovim 0.11 native LSP 側は `FileType` 再発火のラッパーで補強
- mise の go shim を PATH で迂回して gopls 内部の `go env` を TinyGo overlay に寄せる
- `.tinygo.json` を開くだけで target 切替 + gopls 再 attach まで自動で走る

補完が復活してから組み込みを書くのが一気に楽になった。同じく mise + TinyGo + Neovim の組み合わせで `machine` が解決できない人はだいたい同じ罠だと思う。
