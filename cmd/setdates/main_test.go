package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSetDates(t *testing.T) {
	type input struct {
		file string
		now  time.Time
		path string
	}
	type want struct {
		file    string
		updated bool
		err     bool
	}
	now := time.Date(2026, 9, 30, 12, 34, 56, 0, time.FixedZone("JST", 9*60*60))
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "既存の公開日は変更しない",
			input: input{file: "testdata/input/existing.md", now: now},
			want:  want{file: "testdata/golden/existing.md", updated: false},
		},
		{
			name:  "下書きは変更しない",
			input: input{file: "testdata/input/draft.md", now: now},
			want:  want{file: "testdata/golden/draft.md", updated: false},
		},
		{
			name:  "セクションページは変更しない",
			input: input{file: "testdata/input/_index.md", now: now},
			want:  want{file: "testdata/golden/_index.md", updated: false},
		},
		{
			name:  "ネストした記事に公開日を追加する",
			input: input{file: "testdata/input/new/index.md", now: now, path: "new/index.md"},
			want:  want{file: "testdata/golden/new/index.md", updated: true},
		},
		{
			name:  "draft未指定の記事に公開日を追加する",
			input: input{file: "testdata/input/implicit.md", now: now},
			want:  want{file: "testdata/golden/implicit.md", updated: true},
		},
		{
			name:  "CRLFを維持する",
			input: input{file: "testdata/input/windows.md", now: now},
			want:  want{file: "testdata/golden/windows.md", updated: true},
		},
		{
			name:  "paramsのdateは公開日と扱わない",
			input: input{file: "testdata/input/nested.md", now: now},
			want:  want{file: "testdata/golden/nested.md", updated: true},
		},
		{
			name:  "複数行文字列内のdateとdraftは判定に使わない",
			input: input{file: "testdata/input/multiline.md", now: now},
			want:  want{file: "testdata/golden/multiline.md", updated: true},
		},
		{
			name:  "引用符付きのdateも維持する",
			input: input{file: "testdata/input/quoted.md", now: now},
			want:  want{file: "testdata/golden/quoted.md", updated: false},
		},
		{
			name:  "Markdown以外は変更しない",
			input: input{file: "testdata/input/ignored.txt", now: now},
			want:  want{file: "testdata/golden/ignored.txt", updated: false},
		},
		{
			name:  "再実行しても公開日時を変更しない",
			input: input{file: "testdata/golden/implicit.md", now: now.Add(24 * time.Hour)},
			want:  want{file: "testdata/golden/implicit.md"},
		},
		{
			name:  "front matterがない場合はエラー",
			input: input{file: "testdata/invalid/missing/z.md", now: now},
			want:  want{file: "testdata/invalid/missing/z.md", err: true},
		},
		{
			name:  "front matterが閉じていない場合はエラー",
			input: input{file: "testdata/invalid/unclosed/z.md", now: now},
			want:  want{file: "testdata/invalid/unclosed/z.md", err: true},
		},
		{
			name:  "TOMLが不正な場合はエラー",
			input: input{file: "testdata/invalid/invalid-toml/z.md", now: now},
			want:  want{file: "testdata/invalid/invalid-toml/z.md", err: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			name := tt.input.path
			if name == "" {
				name = filepath.Base(tt.input.file)
			}
			path := filepath.Join(root, name)
			content, err := os.ReadFile(tt.input.file)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatal(err)
			}

			paths, err := setDates(root, tt.input.now)
			if (err != nil) != tt.want.err {
				t.Fatalf("setDates() error = %v, want error = %v", err, tt.want.err)
			}
			var wantPaths []string
			if tt.want.updated {
				wantPaths = []string{path}
			}
			if diff := cmp.Diff(wantPaths, paths); diff != "" {
				t.Errorf("updated paths (-want +got):\n%s", diff)
			}
			expected, err := os.ReadFile(tt.want.file)
			if err != nil {
				t.Fatal(err)
			}
			actual, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(string(expected), string(actual)); diff != "" {
				t.Errorf("content (-want +got):\n%s", diff)
			}
		})
	}
}
