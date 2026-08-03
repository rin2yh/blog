package main

import (
	"os"
	"path/filepath"
	"reflect"
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

func TestReadDecks_SortedByDateDesc(t *testing.T) {
	repo := t.TempDir()
	slidesDir := filepath.Join(repo, "slides")
	if err := os.MkdirAll(slidesDir, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, date string) {
		body := "---\ntitle: " + name + "\ndate: " + date + "\n---\n"
		if err := os.WriteFile(filepath.Join(slidesDir, name+".md"), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("old", "2026-01-01")
	write("new", "2026-07-29")

	decks, err := readDecks(repo)
	if err != nil {
		t.Fatalf("readDecks() returned error: %v", err)
	}

	got := []string{decks[0].name, decks[1].name}
	want := []string{"new", "old"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("readDecks() order = %v, want %v", got, want)
	}
}

func TestCollectExternalURLs(t *testing.T) {
	dir := t.TempDir()
	postDir := filepath.Join(dir, "2026", "a-post")
	if err := os.MkdirAll(postDir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "+++\ndate = '2026-07-04T12:57:21+09:00'\nexternalUrl = 'https://slides.rin2yh.com/pumlv-go/'\n+++\n"
	if err := os.WriteFile(filepath.Join(postDir, "index.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}

	urls, err := collectExternalURLs(dir)
	if err != nil {
		t.Fatalf("collectExternalURLs() returned error: %v", err)
	}

	if !urls["https://slides.rin2yh.com/pumlv-go"] {
		t.Errorf("collectExternalURLs() = %v, want it to contain the pumlv-go URL", urls)
	}
}

func TestMissingDecks(t *testing.T) {
	decks := []deck{
		{name: "pumlv-go", title: "既に記事がある", date: "2026-07-29"},
		{name: "new-talk", title: "まだ記事が無い", date: "2026-08-01"},
	}
	// 末尾スラッシュ無しで記録されていても既出と判定できること。
	published := map[string]bool{"https://slides.rin2yh.com/pumlv-go": true}

	missing := missingDecks(decks, published, defaultBaseURL)

	if len(missing) != 1 || missing[0].name != "new-talk" {
		t.Errorf("missingDecks() = %+v, want only new-talk", missing)
	}
}

func TestDeckURL(t *testing.T) {
	want := "https://slides.rin2yh.com/pumlv-go/"
	if got := deckURL("https://slides.rin2yh.com/", "pumlv-go"); got != want {
		t.Errorf("deckURL() = %q, want %q", got, want)
	}
	if got := deckURL("https://slides.rin2yh.com", "pumlv-go"); got != want {
		t.Errorf("deckURL() = %q, want %q", got, want)
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

func TestTOMLString_Apostrophe(t *testing.T) {
	want := `"Go's coverage"`
	if got := tomlString("Go's coverage"); got != want {
		t.Errorf("tomlString() = %s, want %s", got, want)
	}
}

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

	if err := run(repo, defaultBaseURL, contentDir); err != nil {
		t.Fatalf("run() returned error: %v", err)
	}
	after1, err := os.ReadDir(contentDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(after1) != 1 {
		t.Fatalf("run() created %d posts, want 1", len(after1))
	}

	if err := run(repo, defaultBaseURL, contentDir); err != nil {
		t.Fatalf("second run() returned error: %v", err)
	}
	after2, err := os.ReadDir(contentDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(after2) != 1 {
		t.Errorf("second run() created %d posts in total, want 1", len(after2))
	}
}
