+++
date = '2024-02-27T20:07:07+09:00'
draft = false
title = 'gh CLIの学習めも'
categories = []
tags = ["GitHub", "CLI", "メモ"]
+++

## ざっくり所感
ghコマンドはgitコマンドと比べると、GitHubの操作に長けているらしい。

> 詳細: https://docs.github.com/ja/github-cli/github-cli/about-github-cli#%E3%82%B3%E3%83%9E%E3%83%B3%E3%83%89-%E3%83%A9%E3%82%A4%E3%83%B3%E3%81%A7%E3%81%AE-github-cli-%E3%81%A8-git-%E3%81%AE%E9%81%95%E3%81%84%E3%81%AF%E4%BD%95%E3%81%A7%E3%81%99%E3%81%8B

`gh repo clone OWNER/REPO`でリポジトリをクローンできる（`git clone URL`相当）。

参考記事:

- https://qiita.com/ryo2132/items/2a29dd7b1627af064d7b
- https://zenn.dev/fusic/articles/336c5192d2f162

チェックリスト（個人メモ）:

- [ ] 便利なエイリアス、スクリプト作る

## エイリアスがどこまで便利になるか
### 個人的によく使うエイリアス

```zsh
alias status="git status"
alias add="git add"
alias commit="git commit"
alias push="git push"
alias pull="git pull"
alias fetch="git fetch"
alias merge="git merge"
alias branch="git branch"
alias clone="git clone"
alias diff="git diff"
alias pushoh="git push origin HEAD"
alias switch="git switch"
alias stash="git stash"
```

`git checkout` or `git switch` 的なコマンド:

```sh
gh pr checkout {<number> | <url> | <branch>} [flags]
```

`git switch -c <branch>`相当:

```sh
# pr作成
gh pr create new-branch

# prのブランチに入る
gh pr checkout new-branch
```

## 結論

ghはgitを完全に置き換えるわけではなく、GitHub上の操作を単純化するツールである。最初は対応表を作るつもりだったが、すべて対応しているわけでもなく技術的に異なる点があるため、表は作らないことにした。
