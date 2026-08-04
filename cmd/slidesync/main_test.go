package main

import (
	"os"
	"path/filepath"
	"testing"
)

const deckFixture = `---
marp: true
theme: dc
paginate: true
size: 16:9
footer: '#goconnect'
title: Javaなしで安全に使えるPlantUMLビューア「pumlv」
description: PlantUMLのビューアを作った話です。
date: 2026-07-29
---

<!-- _class: cover -->

# タイトル
`

func TestParseDeck(t *testing.T) {
	d, err := parseDeck("pumlv-go.md", deckFixture)
	if err != nil {
		t.Fatalf("parseDeck() returned error: %v", err)
	}

	want := deck{
		name:  "pumlv-go",
		title: "Javaなしで安全に使えるPlantUMLビューア「pumlv」",
		date:  "2026-07-29",
	}
	if d != want {
		t.Errorf("parseDeck() = %+v, want %+v", d, want)
	}
}

// コロンを含むタイトルはYAMLではクォートが要る。unquoteはそのために居る。
func TestParseDeck_QuotedTitle(t *testing.T) {
	d, err := parseDeck("quoted.md", "---\ntitle: 'クォート付き: タイトル'\ndate: 2026-01-02\n---\n")
	if err != nil {
		t.Fatalf("parseDeck() returned error: %v", err)
	}

	if d.title != "クォート付き: タイトル" {
		t.Errorf("parseDeck() title = %q, want %q", d.title, "クォート付き: タイトル")
	}
}

func TestParseDeck_Errors(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{"titleが無い", "---\ndate: 2026-01-02\n---\n"},
		{"dateが無い", "---\ntitle: タイトル\n---\n"},
		{"dateの形式が違う", "---\ntitle: タイトル\ndate: 2026/01/02\n---\n"},
		{"frontmatterが無い", "# 見出しだけ\n"},
		{"frontmatterが閉じていない", "---\ntitle: タイトル\ndate: 2026-01-02\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := parseDeck("deck.md", tt.content); err == nil {
				t.Error("parseDeck() returned nil error, want error")
			}
		})
	}
}

// 既存のexternal記事（content/post/why-go-coverage-stmt-fn-*/index.md）と同じ形。
func TestRenderPost(t *testing.T) {
	d := deck{
		name:  "pumlv-go",
		title: "Javaなしで安全に使えるPlantUMLビューア「pumlv」",
		date:  "2026-07-29",
	}

	got := renderPost(d, "https://slides.rin2yh.com/pumlv-go/")
	want := `+++
date = '2026-07-29T00:00:00+09:00'
draft = false
title = 'Javaなしで安全に使えるPlantUMLビューア「pumlv」'
categories = ['tech', 'external']
tags = ['slide']
externalUrl = 'https://slides.rin2yh.com/pumlv-go/'
+++
`
	if got != want {
		t.Errorf("renderPost() =\n%s\nwant\n%s", got, want)
	}
}

// アポストロフィを含むタイトルはTOMLのリテラル文字列に入らない。
func TestTOMLString_Apostrophe(t *testing.T) {
	want := `"Go's coverage"`
	if got := tomlString("Go's coverage"); got != want {
		t.Errorf("tomlString() = %s, want %s", got, want)
	}
}

// 掲載済みのデッキを作り直さないこと、そのときサマリが空になることを確かめる。
// workflowはサマリの中身だけを見てコミットするかどうかを決める。
func TestRun_Idempotent(t *testing.T) {
	repo := t.TempDir()
	slidesDir := filepath.Join(repo, "slides")
	if err := os.MkdirAll(slidesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\ntitle: テスト発表\ndate: 2026-07-29\n---\n"
	if err := os.WriteFile(filepath.Join(slidesDir, "pumlv-go.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	contentDir := filepath.Join(t.TempDir(), "post")
	if err := os.MkdirAll(contentDir, 0o755); err != nil {
		t.Fatal(err)
	}
	summaryPath := filepath.Join(t.TempDir(), "summary.md")

	countPosts := func() int {
		t.Helper()
		entries, err := os.ReadDir(contentDir)
		if err != nil {
			t.Fatal(err)
		}
		return len(entries)
	}
	readSummary := func() string {
		t.Helper()
		b, err := os.ReadFile(summaryPath)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	if err := run(repo, defaultBaseURL, contentDir, summaryPath); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
	if got := countPosts(); got != 1 {
		t.Fatalf("run() created %d posts, want 1", got)
	}
	want := "- [テスト発表](https://slides.rin2yh.com/pumlv-go/) (2026-07-29)\n"
	if got := readSummary(); got != want {
		t.Errorf("run() summary = %q, want %q", got, want)
	}

	if err := run(repo, defaultBaseURL, contentDir, summaryPath); err != nil {
		t.Fatalf("second run() returned error: %v", err)
	}
	if got := countPosts(); got != 1 {
		t.Errorf("second run() created %d posts in total, want 1", got)
	}
	if got := readSummary(); got != "" {
		t.Errorf("second run() summary = %q, want it to be empty", got)
	}
}
