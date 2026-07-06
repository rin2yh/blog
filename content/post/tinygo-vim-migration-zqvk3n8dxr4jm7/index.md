+++
date = '2026-07-06T13:00:00+09:00'
draft = true
title = 'Neovim で TinyGo (Wio Terminal) の補完を効かせる'
categories = ['tech']
tags = ['neovim', 'tinygo', 'go']
+++

## 概要

Wio Terminal 向けに TinyGo を書き始めたところ、Neovim で gopls が `machine` パッケージを解決してくれず、
補完もエラー表示も TinyGo 用の内容にならない状態でした。

`.tinygo.json` を検出して target を切り替える構成を組んでも状況は変わらず、
原因を追ってみると 2 つの要因が絡んでいたのでそれぞれ対応した記録です。

### 手短なまとめ

- mise の `go` shim が env の GOROOT を尊重せず、gopls 内部の `go env` に TinyGo overlay の GOROOT が伝わらない → PATH に mise が使用する go 本体の bin を差し込んだ状態で gopls を起動
- プラグインが LSP に流し込む `cmd_env` のカバー範囲が狭く、GOROOT だけしか反映されないと gopls の view が実際のビルド環境とずれる → GOROOT / GOOS / GOARCH / GOFLAGS を全て埋めてくれるプラグインを選ぶ

### 環境

- Neovim 0.12 系 (native LSP)
- Go / TinyGo は mise で管理
- `mise.toml`:

```toml
[tools]
go = "1.25"
tinygo = "0.41.1"
```

- ボード: Wio Terminal
- `main.go`:

```go
package main

import "machine"

func main() {
	led := machine.LED
	_ = led
}
```

- `.tinygo.json`:

```json
{"target": "wioterminal"}
```

## 起きていたこと

Neovim で `main.go` を開くと、gopls が以下のエラーを出していました。

```
could not import machine (no required module provides package "machine")
```

当然 `machine.LED` の補完も効きません。組み込みを書くのに `machine.` の先が見えないと不便なので、
これをなんとかするというのが出発点です。

## 対応したこと

1. shim を経由せずに gopls を起動する
2. `cmd_env` を GOROOT / GOOS / GOARCH / GOFLAGS まで埋めるプラグインを使う
3. `.tinygo.json` を開いた時点で `:TinygoTarget` を自動で走らせる

### shim を経由せずに gopls を起動する

`.tinygo.json` から target を切り替えると、gopls プロセスの env には正しい TinyGo overlay の GOROOT
(`~/Library/Caches/tinygo/goroot-<hash>`) が入っていました。にもかかわらず、
gopls 内部で `go env` を呼ぶと **通常の Go の GOROOT** が返ってきます。

犯人は mise の shim (`~/.local/share/mise/shims/go`) で、env の GOROOT を尊重せず、mise 側で管理している値で
上書きしていました。gopls は `go env GOROOT` を呼んでモジュール解決するので、
親プロセスから渡した GOROOT は shim の内側でリセットされてしまいます。

shim を経由しない go bin を PATH 先頭に置いてから gopls を起動する形にしました。
`mise which go` は mise が使用する go のパスを返してくれるので、そこから bin ディレクトリを求めて差し込みます。

```lua
local go_dir = vim.fs.dirname(vim.trim(vim.fn.system('mise which go')))
M.cmd = { 'env', 'PATH=' .. go_dir .. ':' .. vim.env.PATH, 'gopls' }
```

これで gopls 内部の `go env` が TinyGo overlay の GOROOT を素直に返すようになります。

### `cmd_env` を GOROOT / GOOS / GOARCH / GOFLAGS まで埋めるプラグインを使う

target 切替に使っていた `pcolladosoto/tinygo.nvim` の `TinyGoSetTarget` は、LSP config に
`GOROOT` と `GOFLAGS` しか設定していませんでした。

```lua
vim.lsp.config('gopls', {
  cmd_env = {
    GOROOT = currentGOROOT,
    GOFLAGS = currentGOFLAGS,
  }
})
```

GOROOT が届いていれば `machine` の解決自体は通ります。ただし GOOS と GOARCH が抜けているため、
gopls はホスト環境 (macOS / arm64) の値で view を構築し、実際のビルド環境 (Wio Terminal は linux / arm) と食い違います。
`_arm.go` や `_linux.go` のファイル名タグ、`//go:build linux` を含む分岐の解析でずれが出やすい形です。

そこで `cmd_env` に GOROOT / GOOS / GOARCH / GOFLAGS を全部流し込んでくれる `sago35/tinygo.vim` に差し替えました
(`autoload/tinygo.vim`)。`:TinygoTarget <target>` で 4 つの env を LSP config に反映し、
`tinygo#TinygoTargets` で target 補完も提供してくれます。

内部では `vim.lsp.enable('gopls', false)` → `sleep 100m` → `vim.lsp.enable('gopls', true)` の順で
gopls を再起動する実装になっており、Neovim 0.11.2 以降であればこの流れで `client:stop()` と
`doautoall('nvim.lsp.enable FileType')` が走って新しい `cmd_env` で再 attach してくれます
（0.11.2 未満では `enable(false)` が停止しないため別途ラッパーが必要でしたが、現行の 0.12 系ではそのままで動きます）。

### `.tinygo.json` を開いた時点で `:TinygoTarget` を自動で走らせる

`:TinygoTarget <target>` は明示的に呼ぶ必要があるので、プロジェクトを開くたびに実行するのは煩わしいです。
target は `.tinygo.json` に書いてあるので、これを開いた時点で反映されて欲しいです。

`FileType go` のタイミングで上向きに `.tinygo.json` を探し、あれば `:TinygoTarget` を呼ぶ autocmd を足しました。
`:TinygoTarget` は毎回呼ぶと gopls を再起動して重いので、同じディレクトリ・同じ target なら
1 セッションで 1 回だけ実行されるようガードを入れています。

```lua
local applied = {}
vim.api.nvim_create_autocmd('FileType', {
  pattern = 'go',
  group = vim.api.nvim_create_augroup('tinygo-auto', {}),
  callback = function(args)
    local root = vim.fs.root(args.buf, '.tinygo.json')
    if not root then return end
    local config_path = vim.fs.joinpath(root, '.tinygo.json')
    local ok, cfg = pcall(function()
      return vim.json.decode(table.concat(vim.fn.readfile(config_path), '\n'))
    end)
    if not ok or type(cfg) ~= 'table' or not cfg.target then return end
    if applied[root] == cfg.target then return end
    applied[root] = cfg.target
    vim.cmd({ cmd = 'TinygoTarget', args = { cfg.target } })
  end,
})
```

## 全体のフロー

`.tinygo.json` を持つプロジェクトを開いた時の挙動は最終的にこうなります。

```mermaid
sequenceDiagram
    participant Nvim as Neovim
    participant Plugin as tinygo.vim
    participant Gopls as gopls
    Nvim->>Plugin: .tinygo.json 検出 → :TinygoTarget &lt;target&gt;
    Plugin->>Plugin: cmd_env に GOROOT/GOOS/GOARCH/GOFLAGS を書換
    Plugin->>Nvim: vim.lsp.enable(false → true)
    Note over Nvim,Gopls: enable(false) が client:stop()<br/>enable(true) が FileType 再発火
    Nvim->>Gopls: 新 cmd_env で起動<br/>(PATH に mise の go bin)
    Note over Gopls: go env GOROOT が<br/>TinyGo overlay を返す
    Gopls-->>Nvim: TinyGo 用の補完・エラー表示が動く
```

## 終わりに

補完が戻ってきてからは、Wio Terminal 向けの TinyGo コードを快適に書けるようになりました。
mise + TinyGo + Neovim の組み合わせで `machine` が解決できない場合、
shim による GOROOT の上書きが典型的な引っかかりどころでした。
`machine` の解決だけならプラグイン側の env カバレッジは狭くても通りますが、
`_linux.go` や `//go:build` の解析まで気にするなら GOOS / GOARCH も含めて流し込めるプラグインを
選んだほうが後で困らないと思います。

読んでいただき、ありがとうございました！

## 参考

- `sago35/tinygo.vim`: https://github.com/sago35/tinygo.vim
- `pcolladosoto/tinygo.nvim`: https://github.com/pcolladosoto/tinygo.nvim
