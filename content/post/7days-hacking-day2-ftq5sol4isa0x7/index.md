+++
date = '2026-05-04T16:56:57+09:00'
draft = false
title = '7days Hacking Day2'
categories = ["tech"]
tags = ["memo", "security"]
+++

https://www.seshop.com/product/detail/26456?srsltid=AfmBOopr_-WJwEPXdy0dGw5oVG1yDJe9FnfaULTQgR2VN4MG_DGFl2XR

使ったツールについてメモ。

kali linuxの標準ツール

- dirb: 辞書を使って指定のwebサイトのアクセスできるパスを特定するツール
- hydra: sshのパスワードを辞書攻撃するツール
- ssh2john: ssh秘密鍵ファイルのパスフレーズをクラックするために使用
    - sshの鍵ファイルをjohnが解析可能なハッシュ形式に変換する
- john
    - ハッシュ値を総当たりで解析するツール
