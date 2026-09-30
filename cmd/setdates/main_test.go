package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSetDates(t *testing.T) {
	type input struct {
		files string
		now   time.Time
	}
	type want struct {
		files   string
		updated []string
		err     bool
	}
	now := time.Date(2026, 9, 30, 12, 34, 56, 0, time.FixedZone("JST", 9*60*60))
	tests := []struct {
		name  string
		input input
		want  want
	}{
		{
			name:  "公開済みの記事だけに日時を追加する",
			input: input{files: "testdata/input", now: now},
			want:  want{files: "testdata/golden", updated: []string{"implicit.md", "multiline.md", "nested.md", "new/index.md", "windows.md"}},
		},
		{
			name:  "再実行しても公開日時を変更しない",
			input: input{files: "testdata/golden", now: now.Add(24 * time.Hour)},
			want:  want{files: "testdata/golden"},
		},
		{
			name:  "front matterがない場合はどの記事も変更しない",
			input: input{files: "testdata/invalid/missing", now: now},
			want:  want{files: "testdata/invalid/missing", err: true},
		},
		{
			name:  "front matterが閉じていない場合はどの記事も変更しない",
			input: input{files: "testdata/invalid/unclosed", now: now},
			want:  want{files: "testdata/invalid/unclosed", err: true},
		},
		{
			name:  "TOMLが不正な場合はどの記事も変更しない",
			input: input{files: "testdata/invalid/invalid-toml", now: now},
			want:  want{files: "testdata/invalid/invalid-toml", err: true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.CopyFS(root, os.DirFS(tt.input.files)); err != nil {
				t.Fatal(err)
			}
			paths, err := setDates(root, tt.input.now)
			if (err != nil) != tt.want.err {
				t.Fatalf("setDates() error = %v, want error = %v", err, tt.want.err)
			}
			var names []string
			for _, path := range paths {
				name, err := filepath.Rel(root, path)
				if err != nil {
					t.Fatal(err)
				}
				names = append(names, filepath.ToSlash(name))
			}
			if diff := cmp.Diff(tt.want.updated, names); diff != "" {
				t.Errorf("updated paths (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(readFiles(t, tt.want.files), readFiles(t, root)); diff != "" {
				t.Errorf("files (-want +got):\n%s", diff)
			}
		})
	}
}

func readFiles(t *testing.T, root string) map[string]string {
	t.Helper()
	files := make(map[string]string)
	dir := os.DirFS(root)
	err := fs.WalkDir(dir, ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(dir, name)
		files[name] = string(content)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	return files
}
