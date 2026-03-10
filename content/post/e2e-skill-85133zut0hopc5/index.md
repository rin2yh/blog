+++
date = '2026-03-08T15:58:56+09:00'
draft = false
title = '自律的にE2EをデバッグするSkillを作ってみた'
categories = ['tech']
tags = ['E2E', 'Playwright', 'GitHub Actions', 'Codex']
+++

## 概要
個人開発のプロダクトでE2Eテストを導入しているのですが、CIのスペック問題から何かしらの理由で
よく落ちます。さらに、適当にAIエージェントに依頼してもなかなか解決できなかったり、
解決できたとしても人間よりも時間がかかったりします。
そこで、AIエージェントが自律的にデバッグできるSkillを作成して効率化しました。

実際に作成したSkillを使ってE2Eテストの修正を何度も依頼していますが、
体感としては人間の介入を必要とせずに、3回以内のコミットでバグ修正できています。
（導入前の解決までのコミット数：6回以上、50%近く人間の介入が必要）

作成したSkill：
https://github.com/Geek-Teck-Mentors/trend_diary/blob/main/.claude/skills/e2e-debug-playwright/SKILL.md

ツール・環境
- GitHub Actions
- Playwright

### 手短なまとめ

効果的なSkill作成には、エージェントが確実に遂行できる条件を定義することが大切です。
遂行できる条件とは、以下のようなものです。
1. 終了条件（完了条件）
2. 人間が運用しているレベルの効果的なデバッグプロセス

2についてはコマンドレベルでの具体性までは不要です。

これらの定義をしっかり行った上で、Skill Creatorを使いましょう！

## 作成までの流れ
手順自体は簡単で、以下の通りです。

1. ローカルのCodexでE2Eのデバッグを依頼
2. [Skill Creator](https://github.com/anthropics/skills/blob/main/skills/skill-creator/SKILL.md)でSkillを作成

作成手順よりもデバッグプロセスをどう定義するかがすごく大切で
依頼時にCIがアップロードするartifactを見るように指示しています。
特にPlaywrightではログやトレースが丁寧に残っていて、私のデバッグでも
多用していたため、依頼内容に記載していました。

そのため、その旨がSkillに記載されました。

```md
## 3. artifactで画面状態遷移を読む

必ず見る:
 - `error-context.md`
 - `test-failed-1.png`
 - `video.webm`

確認ポイント:
 - どの画面で止まっているか（`/signup` or `/login` or `/trends`）
 - ボタン状態（例: `ログイン中...` でdisabled）
 - 入力値がDOMに反映されているか

```

また、gh CLIでActionsのartifactの取得やActionsのステータスチェック、ログチェックもできることが
Codexの作業中にわかった（Codexが実行していた）ので、Skillに組み込んでいます。

artifact取得のコマンド
```sh
gh run download <run_id> -n playwright-test-results-diff-<run_id>-1 -D artifacts/run-<run_id>
```

終了条件も明示しています。

> 終了条件:
> - Diff actionがエラー時に確実に失敗する
> - Diff action上のE2Eエラーが消える

## 終わりに

今回、デバッグプロセスをSkill化する中で「そんなghコマンドの使い方が！」という驚きが多く、
ghコマンドの使い方自体をAIエージェントから学ぶことになりました。

みんなも適切にSkillを作って楽しく開発！（某デュエル風）

読んでいただき、ありがとうございました！

完成したSkill：
https://github.com/Geek-Teck-Mentors/trend_diary/blob/main/.claude/skills/e2e-debug-playwright/SKILL.md

## 参考文献
- Skill Creator： https://github.com/anthropics/skills/blob/main/skills/skill-creator/SKILL.md
