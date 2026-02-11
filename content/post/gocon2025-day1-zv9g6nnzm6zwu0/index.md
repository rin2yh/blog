+++
date = '2025-09-27T20:40:34+09:00'
draft = false
title = 'Gocon2025 Day1'
categories = ['idea']
tags = ['conference', 'note']
+++

go con day1のメモ。

## Go1.24で進化したmap型について理解する

1.23
段階的に再配置する→マップ自体の負荷が高いから
（結局CPU負荷は高かった）

チェイン法＋オーバーフローバケット（連結リスト的な感じ）
- オーバーフローバケットを辿る過程でキャッシュミスが多い

1.24
オープンアドレス方式
Quadratic probing
- ハッシュの衝突＝クラスタリングを軽減する（オープンアドレスの課題）

Control word
8 slot一括検証してパフォーマンスUP

## Go で WebAssembly を利用した実用的なプラグインシステムの構築方法
WASM
・ブラウザ利用が当初の目的
・Sandbox機構で安全に実行
・いろんな言語から作れる
・言語問わず、WASMランタイムがあれば動く

プラグインの作り方→WASM以外の場合、Shared Library、RPC

WASI：WebAssembly System Interface (規格)
https://speakerdeck.com/goccy/go-de-webassembly-woli-yong-sita-shi-yong-de-napuraguinsisutemunogou-zhu-fang-fa


## GoのinterfaceとGenericsの内部構造と進化

interface→辞書とGcshape Stencilling
モノモーフィゼーション（コンパイル時展開）とランタイムの辞書渡しで実装されている
・1つの関数インスタンスを複数の方で動かせる

Genericsガイドライン：https://go.dev/blog/when-generics

## 全てGoで作るP2P対戦ゲーム入門

https://speakerdeck.com/ponyo877/quan-tegodezuo-rup2pdui-zhan-gemuru-men

## analysis パッケージの仕組みの上でMulti linter with configを実現する
統合型Linterの作り方

なぜパッケージの仕組みの上でMulti linter with configが必要か
・それぞれのルールを設定で調整できるようにしたい

設定ファイルを読み込むだけのAnalyzerを置く

https://github.com/k1LoW/gostyle

vetツールが汎用的になった時にやる。

