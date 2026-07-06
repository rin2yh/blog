+++
date = '2026-07-06T13:00:00+09:00'
draft = true
title = 'tinygo.nvim から tinygo.vim に乗り換えた'
categories = ['tech']
tags = ['neovim', 'tinygo', 'go']
+++

## 概要

Wio Terminal 向けに TinyGo を書き始めたところ、Neovim で gopls が `machine` パッケージを解決できず、
補完もエラー表示も TinyGo 用の内容にならない状態でした。

`.tinygo.json` を検出して `pcolladosoto/tinygo.nvim` で target を切り替える構成にしていたのですが、
それでも補完が復活しません。原因を追ってみたところ、mise の shim が GOROOT env を上書きしていたこと、
そしてプラグインが LSP に流し込む cmd_env のカバー範囲が狭かったことの 2 点が絡んでいました。

前者を `env PATH=<mise-go-dir>:...` で shim を迂回する形にし、後者を GOROOT / GOOS / GOARCH / GOFLAGS を
全て流し込んでくれる `sago35/tinygo.vim` への乗り換えで解消しました。

### 手短なまとめ

- mise の `go` shim は env の GOROOT を尊重しないので、gopls 内部の `go env` に TinyGo overlay の GOROOT が伝わらない
- `pcolladosoto/tinygo.nvim` の `TinyGoSetTarget` は GOROOT と GOFLAGS しか LSP config に流し込まないため、GOOS / GOARCH が反映されない
- 上記 2 点を `env PATH=<mise-go-dir>:...` gopls ラップと、`sago35/tinygo.vim` への乗り換えで解消しました

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

## mise の `go` shim が GOROOT env を上書きする

### 課題: gopls 内部の `go env` に TinyGo overlay の GOROOT が伝わらない

target 切替後、GOROOT が gopls に届いているかを確認したところ、
gopls プロセスの env には正しい TinyGo overlay の GOROOT (`~/Library/Caches/tinygo/goroot-<hash>`) が入っていました。
にもかかわらず、gopls 内部から `go env` を実行すると **通常の Go の GOROOT** が返ってきます。

追ってみると、mise の shim (`~/.local/share/mise/shims/go`) が env の GOROOT を尊重せず、
mise 側で管理している GOROOT で上書きしていました。
gopls は内部で `go env GOROOT` を呼んでモジュール解決するので、
親プロセスから GOROOT を渡しても shim の内側でリセットされてしまいます。

### 対応: shim を経由せずに gopls を起動する

shim を経由しない go bin を PATH 先頭に置いてから gopls を起動します。
`mise which go` は mise が使用する go のパスを返してくれるので、そこから bin ディレクトリを求めて差し込みました。

```lua
local go_dir = vim.fs.dirname(vim.trim(vim.fn.system('mise which go')))
M.cmd = { 'env', 'PATH=' .. go_dir .. ':' .. vim.env.PATH, 'gopls' }
```

これで gopls 内部の `go env` が TinyGo overlay の GOROOT を素直に返すようになりました。

## `pcolladosoto/tinygo.nvim` は GOOS / GOARCH を流し込まない

### 課題: cmd_env のカバー範囲が狭く、gopls の view と実際のビルド環境がずれる

`pcolladosoto/tinygo.nvim` の `TinyGoSetTarget` は LSP config に GOROOT と GOFLAGS しか設定しません。

```lua
vim.lsp.config('gopls', {
  cmd_env = {
    GOROOT = currentGOROOT,
    GOFLAGS = currentGOFLAGS,
  }
})
```

GOROOT と GOFLAGS のビルドタグだけで `machine` の解決自体は通ります。ただし GOOS と GOARCH が抜けているため、
gopls は宿主環境 (macOS/arm64 など) の値で view を構築し、実際のビルド環境 (Wio Terminal だと linux/arm) と食い違います。
`_arm.go` や `_linux.go` のようなファイル名タグ、`//go:build linux` を含む分岐の解析にずれが出やすい形です。

### 対応: `sago35/tinygo.vim` に乗り換える

`sago35/tinygo.vim` の `:TinygoTarget` は GOROOT / GOOS / GOARCH / GOFLAGS を全て LSP config に流し込みます
(`autoload/tinygo.vim`)。

- `:TinygoTarget <target>` で 4 つの env を LSP config に反映
- `tinygo#TinygoTargets` で target 補完も提供

内部では `vim.lsp.enable('gopls', false)` → `sleep 100m` → `vim.lsp.enable('gopls', true)` の順で gopls を
再起動する実装になっており、Neovim 0.11.2 以降はこの流れで `client:stop()` と
`doautoall('nvim.lsp.enable FileType')` が走って新しい cmd_env で再 attach してくれます
（0.11.2 未満では `enable(false)` が停止しないため別途ラッパーが必要でしたが、現行の 0.12 系ではそのままで動きます）。

## `.tinygo.json` の target が自動反映されない

### 課題: プロジェクトを開くたびに `:TinygoTarget <target>` を呼ぶ必要がある

`sago35/tinygo.vim` は `:TinygoTarget <target>` で明示的に target を切り替える形なので、
プロジェクトを開くたびに実行が必要になります。TinyGo プロジェクトごとに target は `.tinygo.json` で決まっているので、
これを開いた時点で自動的に反映されてほしいところです。

### 対応: `FileType go` の autocmd で自動反映する

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

## 対応後のフロー

`.tinygo.json` を持つプロジェクトを開いた時の挙動はこうなります。

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
shim による GOROOT 上書きに該当することが多いようです。
プラグイン側の cmd_env カバレッジも、`machine` は解決できても細かい挙動がずれる要因になるので、
乗り換え前に一度確認しておくと後戻りが少なくなります。

読んでいただき、ありがとうございました！

## 参考

- `sago35/tinygo.vim`: https://github.com/sago35/tinygo.vim
- `pcolladosoto/tinygo.nvim`: https://github.com/pcolladosoto/tinygo.nvim
