+++
date = '2026-09-09T01:21:46+09:00'
draft = false
title = 'Claude Maxプラン(5x)を半年使って見直した話'
categories = ['tech']
tags = ['Claude Code', 'ChatGPT', 'Codex']
+++

## はじめに

生成AIの台頭により、開発スタイルが大きく変わり、[Claude Code](https://claude.com/code)の利用量が一気に増えました。
私は元々Claude Proプランを契約していたのですが、足りなくなったため、8月まで約半年間Claude Maxプラン(5x)を契約しました。（以降、本稿ではClaude Proプランを「Proプラン」と表記します。Maxプラン(5x)を「Maxプラン」と表記します）

Maxプランに変更してから、様々なことにチャレンジでき、大満足でした。しかし、半年使った後に料金と使い方を改めて見直し、Maxプランを解約しました。

本稿では、Maxプランを半年程使った経験と解約に至った判断を振り返ります。

## Maxプランを契約した理由

Maxプランを契約した理由は、単純にProプランの利用量では足りなかったからです。当時は5時間の利用枠でも1〜1.5時間程使うと上限に達することが大半で、週次の上限も概ね使い切りました。開発中にセッションや週次の利用制限を気にすること自体がかなりストレスになってきたため、より利用量の多いMaxプランへ切り替えました。

少なくとも契約した時点では、Maxプランにする理由は明確でした。

## Maxプランを使ってどうだったか

Maxプランを契約していた半年間は、Claude Codeを使って様々なことに取り組みました。

取り組んだことは次の通りです。

1. [pumlv](https://github.com/rin2yh/pumlv)の開発・改善
2. [gostty](https://github.com/rin2yh/gostty)でguiguiとlibghosttyを繋ぐ興味本位の自作ターミナル
3. [dinosaur-game](https://github.com/rin2yh/dinosaur-game)や[games](https://github.com/rin2yh/games)でEbitengineを使ったゲーム開発
4. [dotfiles](https://github.com/rin2yh/dotfiles)にNix導入・改善
5. [slides](https://github.com/rin2yh/slides)でスライドテーマの自作とテーマ含めたMarp移行
6. [hugo-theme-stack-liquid-glass](https://github.com/rin2yh/hugo-theme-stack-liquid-glass)でブログのテーマ自作と改善
7. [claude-code-plugins](https://github.com/rin2yh/claude-code-plugins)でAgent Skillsやhooks等、Claude Code自体を使いやすくするプラグインリポジトリの作成
8. [study-architecture](https://github.com/rin2yh/study-architecture)で非同期処理やイベント設計を検証

Maxプランのおかげで沢山の技術を学び、検証ループを回すことができました。利用量もMaxプランの週次上限に対して概ね8割程消費する日々が続き、当初の課題感を解消できて満足していました。

しかし、ある時期からコストを確認する度に小さな違和感を覚えていきました。その積み重ねが、解約の一途を辿ることになります。

## Maxプランを解約した理由

Maxプランには満足していました。しかし、3ヶ月目から月額の請求書を確認する度に「収入に対してコストを支払い過ぎているのではないか」と感じるようになりました。

ありがたいことに、私は副業で月3〜4万円程の収入を得ています。一方、Maxプランは月額$100（約15,600円[^1]）です。Maxプランの料金は、副業収入の約39〜52%に当たります。利便性と利用量を考えても、この割合の固定費を毎月支払い続けることには納得がいきませんでした。企業レベルのビジネスであれば選択肢になり得ます。しかし、月3〜4万円程の収入に対する固定費としては厳しいと感じました。

また、4ヶ月目頃から、私がアイデアを形にする速度とClaude Codeが実装する速度の均衡が完全に崩れてしまいました。私の実力不足ですが、Claude Codeで作るものがない状況に陥ってしまったわけです。

Maxプランに変更することで、当初の課題であった「Proプランの利用制限を気にせずClaude Codeを使いたい」は解消されました。その一方で、Maxプランを使い続ける理由がなくなりました。

## 解約後のプランを検討する

Maxプランを解約した後の選択肢として、次の3つを検討しました。

1. Proプランに戻す
2. Proプランに戻し、必要な時だけ追加利用する
3. Proプランと[ChatGPT](https://chatgpt.com/) Plusプランを併用する

Proプランだけに戻す場合、以前と同じく利用制限に当たる可能性が高いです。必要な時だけ追加利用する方法であれば、Maxプラン程の固定費を支払わずにClaude Codeを使い続けることが可能です。

加えて、ChatGPT Plusプランを併用する選択肢も検討しました。過去にChatGPTを無料の範囲で利用したことがあります。ChatGPTユーザの投稿等を見て、汎用性が高そうに感じたことも候補に加えた理由です。ChatGPTはコーディング以外にも画像生成やGPTs等のカスタム機能があり、幅広く利用できます。

また、ChatGPTでのチャットと[Codex](https://openai.com/codex/)の利用量は別枠です。Codexを使ってもチャット側の利用量を消費しない点にかなり惹かれました。[Claude](https://claude.ai/)ではチャットとClaude Codeの利用量が同一のため、常々チャットでの相談がし辛いと思っていました。

最終的に、Claude単体利用時にうっすら感じていた課題も解消できるProプラン＋ChatGPT Plusプランの併用を選択しました。

## ClaudeとChatGPTを併用してみて

9月に入り、併用して大体1ヶ月経ちました。副業では引き続きClaude Codeを使い、利用制限に当たった場合はCodexに切り替えています。それでもCodexの利用量は2〜4割程に収まるため、残りは勉強やOSSに回しています。

また、少しでも気になったことをChatGPTのチャット機能で相談できるようになりました。キーワード検索すらできない未知の領域への調査や小さな相談がしやすくなり、普段の生活がより快適になりました。日常会話レベルの英語も勉強し始めたので、日々GPTsが活躍しています。

用途に応じて使い分けることで、それぞれのサービスを無駄なく活用できています。たまに画像生成で遊ぶこともでき、小さな楽しみも増えました。

## おわりに

Maxプランには半年間かなりお世話になり、契約したこと自体は良い判断だったと思っています。Proプランの利用制限を気にする必要がなくなり、様々なことに挑戦できました。

利用状況やコストを見直した結果、今の私にはProプラン＋ChatGPT Plusプランの併用が合っていると判断しました。用途に応じて使い分けることで、より適材適所といった気配を感じています。現状では、Maxプラン単体だけだった時よりも満足感が高いです。

今回のことから、AIも他の技術やツール選定同様にスタイル/利用状況とコスパを併せて考えることが大切だと学びました（AIは人の能力を拡張するツールですし）。もし似たような悩みがあったり、併用するか迷ったりしている方がいたら参考になると幸いです。

最後まで読んでいただき、ありがとうございました！

## 参考文献

### サービス・料金

- [Claude](https://claude.ai/)
- [Claude Code](https://claude.com/code)
- [Choosing a Claude Plan - Anthropic Help Center](https://support.anthropic.com/en/articles/11049762-choosing-a-claude-ai-plan)
- [Using Claude Code with your Pro or Max plan - Anthropic Help Center](https://support.anthropic.com/en/articles/11145838-using-claude-code-with-your-pro-or-max-plan)
- [ChatGPT](https://chatgpt.com/)
- [Codex](https://openai.com/codex/)
- [ChatGPTプランでCodexを使う - OpenAI Help Center](https://help.openai.com/ja-jp/articles/11369540)

### 本稿で触れた取り組み

- [pumlv](https://github.com/rin2yh/pumlv)
  - [Javaなしで安全に使えるPlantUMLビューア「pumlv」](https://zenn.dev/rinrin_yuuki/articles/9b69cca81875f6)
- [gostty](https://github.com/rin2yh/gostty)
  - [GuiguiとGhosttyを組み合わせてターミナルGUIを作ってみた](https://zenn.dev/rinrin_yuuki/articles/448d45e7df01ee)
- [dinosaur-game](https://github.com/rin2yh/dinosaur-game)
- [games](https://github.com/rin2yh/games)
- [dotfiles](https://github.com/rin2yh/dotfiles)
- [slides](https://github.com/rin2yh/slides)
- [hugo-theme-stack-liquid-glass](https://github.com/rin2yh/hugo-theme-stack-liquid-glass)
- [claude-code-plugins](https://github.com/rin2yh/claude-code-plugins)
- [study-architecture](https://github.com/rin2yh/study-architecture)

[^1]: 2026年9月9日時点の為替レート、$1 = 約156円として換算しています。
