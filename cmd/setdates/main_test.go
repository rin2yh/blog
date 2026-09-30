package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/sebdah/goldie/v2"
)

func TestSetDate(t *testing.T) {
	type input struct {
		file string
		now  time.Time
		path string
	}
	type want struct {
		file string
		err  string
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
			want:  want{file: "testdata/golden/existing.md"},
		},
		{
			name:  "下書きは変更しない",
			input: input{file: "testdata/input/draft.md", now: now},
			want:  want{file: "testdata/golden/draft.md"},
		},
		{
			name:  "ネストした記事に公開日を追加する",
			input: input{file: "testdata/input/new/index.md", now: now, path: "new/index.md"},
			want:  want{file: "testdata/golden/new/index.md"},
		},
		{
			name:  "draft未指定の記事に公開日を追加する",
			input: input{file: "testdata/input/implicit.md", now: now},
			want:  want{file: "testdata/golden/implicit.md"},
		},
		{
			name:  "CRLFを維持する",
			input: input{file: "testdata/input/windows.md", now: now},
			want:  want{file: "testdata/golden/windows.md"},
		},
		{
			name:  "paramsのdateは公開日と扱わない",
			input: input{file: "testdata/input/nested.md", now: now},
			want:  want{file: "testdata/golden/nested.md"},
		},
		{
			name:  "複数行文字列内のdateとdraftは判定に使わない",
			input: input{file: "testdata/input/multiline.md", now: now},
			want:  want{file: "testdata/golden/multiline.md"},
		},
		{
			name:  "引用符付きのdateも維持する",
			input: input{file: "testdata/input/quoted.md", now: now},
			want:  want{file: "testdata/golden/quoted.md"},
		},
		{
			name:  "再実行しても公開日時を変更しない",
			input: input{file: "testdata/golden/implicit.md", now: now.Add(24 * time.Hour)},
			want:  want{file: "testdata/golden/implicit.md"},
		},
		{
			name:  "front matterがない場合はエラー",
			input: input{file: "testdata/invalid/missing/z.md", now: now},
			want:  want{file: "testdata/invalid/missing/z.md", err: "expected TOML front matter"},
		},
		{
			name:  "front matterが閉じていない場合はエラー",
			input: input{file: "testdata/invalid/unclosed/z.md", now: now},
			want:  want{file: "testdata/invalid/unclosed/z.md", err: "unclosed front matter"},
		},
		{
			name:  "TOMLが不正な場合はエラー",
			input: input{file: "testdata/invalid/invalid-toml/z.md", now: now},
			want:  want{file: "testdata/invalid/invalid-toml/z.md", err: "toml: incomplete number"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, path := prepareArticle(t, tt.input.file, tt.input.path)

			_, err := setDate(path, tt.input.now)
			var message string
			if err != nil {
				message = err.Error()
			}
			if diff := cmp.Diff(tt.want.err, message); diff != "" {
				t.Errorf("error (-want +got):\n%s", diff)
			}
			golden := goldie.New(t, goldie.WithFixtureDir("."), goldie.WithNameSuffix(""))
			golden.Assert(t, tt.want.file, readTestFile(t, path))
		})
	}
}

func TestCollectArticlePaths(t *testing.T) {
	type input struct {
		file string
		path string
	}
	tests := []struct {
		name  string
		input input
		want  []string
	}{
		{
			name:  "Markdown記事を対象にする",
			input: input{file: "testdata/input/implicit.md"},
			want:  []string{"implicit.md"},
		},
		{
			name:  "ネストした記事を対象にする",
			input: input{file: "testdata/input/new/index.md", path: "new/index.md"},
			want:  []string{"new/index.md"},
		},
		{
			name:  "セクションページを除外する",
			input: input{file: "testdata/input/_index.md"},
			want:  nil,
		},
		{
			name:  "Markdown以外を除外する",
			input: input{file: "testdata/input/ignored.txt"},
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, _ := prepareArticle(t, tt.input.file, tt.input.path)
			paths, err := collectArticlePaths(root)
			if err != nil {
				t.Fatal(err)
			}
			var want []string
			for _, path := range tt.want {
				want = append(want, filepath.Join(root, path))
			}
			if diff := cmp.Diff(want, paths); diff != "" {
				t.Errorf("paths (-want +got):\n%s", diff)
			}
		})
	}
}

func prepareArticle(t *testing.T, fixture, name string) (root, path string) {
	t.Helper()
	root = t.TempDir()
	if name == "" {
		name = filepath.Base(fixture)
	}
	path = filepath.Join(root, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, readTestFile(t, fixture), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, path
}

func readTestFile(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}
