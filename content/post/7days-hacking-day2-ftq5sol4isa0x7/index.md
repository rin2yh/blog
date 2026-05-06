+++
date = '2026-05-06T16:56:57+09:00'
draft = false
title = '7days Hacking Day2'
categories = ["tech"]
tags = ["note", "security"]
+++

使ったツールについてメモ。

kali linuxの標準ツール

- dirb: 辞書を使って指定のwebサイトのアクセスできるパスを特定するツール
- hydra: sshのパスワードを辞書攻撃するツール
- ssh2john: ssh秘密鍵ファイルのパスフレーズをクラックするために使用
    - sshの鍵ファイルをjohnが解析可能なハッシュ形式に変換する
- john
    - ハッシュ値を総当たりで解析するツール
