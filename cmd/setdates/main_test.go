package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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
			name: "公開済みの記事だけに日時を追加する",
			input: input{
				files: "testdata/input",
				now:   now,
			},
			want: want{
				files:   "testdata/golden",
				updated: []string{"implicit.md", "multiline.md", "nested.md", "new/index.md", "windows.md"},
			},
		},
		{
			name: "再実行しても公開日時を変更しない",
			input: input{
				files: "testdata/golden",
				now:   now.Add(24 * time.Hour),
			},
			want: want{
				files: "testdata/golden",
			},
		},
		{
			name: "front matterがない場合はどの記事も変更しない",
			input: input{
				files: "testdata/invalid/missing",
				now:   now,
			},
			want: want{
				files: "testdata/invalid/missing",
				err:   true,
			},
		},
		{
			name: "front matterが閉じていない場合はどの記事も変更しない",
			input: input{
				files: "testdata/invalid/unclosed",
				now:   now,
			},
			want: want{
				files: "testdata/invalid/unclosed",
				err:   true,
			},
		},
		{
			name: "TOMLが不正な場合はどの記事も変更しない",
			input: input{
				files: "testdata/invalid/invalid-toml",
				now:   now,
			},
			want: want{
				files: "testdata/invalid/invalid-toml",
				err:   true,
			},
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
			if !reflect.DeepEqual(names, tt.want.updated) {
				t.Errorf("updated = %v, want %v", names, tt.want.updated)
			}
			compareFiles(t, root, tt.want.files)
		})
	}
}

func compareFiles(t *testing.T, root, expected string) {
	t.Helper()
	err := fs.WalkDir(os.DirFS(expected), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		want, err := os.ReadFile(filepath.Join(expected, name))
		if err != nil {
			return err
		}
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			return err
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s: content = %q, want %q", name, got, want)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}
