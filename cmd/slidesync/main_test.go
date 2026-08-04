package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestParseDeck(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    deck
	}{
		{
			// 読まないキーが前後にあっても、引用符付きの値があっても拾えること。
			name:    "実際のデッキ",
			content: "---\nmarp: true\nfooter: '#goconnect'\ntitle: PlantUMLビューア「pumlv」\ndate: 2026-07-29\n---\n\n# 本文\n",
			want:    deck{name: "deck", title: "PlantUMLビューア「pumlv」", date: "2026-07-29"},
		},
		{
			// コロンを含むタイトルはYAMLではクォートが要る。unquoteはそのために居る。
			name:    "クォート付きタイトル",
			content: "---\ntitle: 'クォート付き: タイトル'\ndate: 2026-01-02\n---\n",
			want:    deck{name: "deck", title: "クォート付き: タイトル", date: "2026-01-02"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDeck("deck.md", tt.content)
			if err != nil {
				t.Fatalf("parseDeck() returned error: %v", err)
			}
			if got != tt.want {
				t.Errorf("parseDeck() = %+v, want %+v", got, tt.want)
			}
		})
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
	d := deck{name: "pumlv-go", title: "PlantUMLビューア「pumlv」", date: "2026-07-29"}

	got := renderPost(d, "https://slides.rin2yh.com/pumlv-go/")
	want := `+++
date = '2026-07-29T00:00:00+09:00'
draft = false
title = 'PlantUMLビューア「pumlv」'
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

// 掲載済みのデッキを作り直さないこと、そのとき標準出力が空になることを確かめる。
// workflowはこの出力の有無でコミットするかどうかを決める。
func TestRun_Idempotent(t *testing.T) {
	repo := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repo, "slides"), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "---\ntitle: テスト発表\ndate: 2026-07-29\n---\n"
	if err := os.WriteFile(filepath.Join(repo, "slides", "pumlv-go.md"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	contentDir := t.TempDir()

	runOnce := func() string {
		t.Helper()
		var out bytes.Buffer
		if err := run(repo, defaultBaseURL, contentDir, &out); err != nil {
			t.Fatalf("run() returned error: %v", err)
		}
		return out.String()
	}
	countPosts := func() int {
		t.Helper()
		entries, err := os.ReadDir(contentDir)
		if err != nil {
			t.Fatal(err)
		}
		return len(entries)
	}

	want := "- [テスト発表](https://slides.rin2yh.com/pumlv-go/) (2026-07-29)\n"
	if got := runOnce(); got != want {
		t.Errorf("run() wrote %q, want %q", got, want)
	}
	if got := countPosts(); got != 1 {
		t.Fatalf("run() created %d posts, want 1", got)
	}

	if got := runOnce(); got != "" {
		t.Errorf("second run() wrote %q, want nothing", got)
	}
	if got := countPosts(); got != 1 {
		t.Errorf("second run() created %d posts in total, want 1", got)
	}
}
