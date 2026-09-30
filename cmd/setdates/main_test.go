package main

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestSetDates(t *testing.T) {
	root := t.TempDir()
	originals := map[string]string{
		"existing.md":  "+++\ndate = '2020-01-02T03:04:05+09:00'\ndraft = false\n+++\nBody\n",
		"draft.md":     "+++\ndraft = true\n+++\nBody\n",
		"_index.md":    "+++\ntitle = 'Posts'\n+++\n",
		"new/index.md": "+++\ndraft = false\ntitle = 'New'\n+++\n\nBody +++\n",
		"implicit.md":  "+++\ntitle = 'Implicit'\n+++\n",
		"windows.md":   "+++\r\ndraft = false\r\n+++\r\nBody\r\n",
		"nested.md":    "+++\n[params]\ndate = 'custom'\n+++\n",
		"multiline.md": "+++\ntitle = '''\ndate = 'inside a string'\ndraft = true\n'''\n+++\n本文\n",
		"quoted.md":    "+++\n'date' = 2020-01-02T03:04:05+09:00\n+++\n",
		"ignored.txt":  "not an article",
	}
	for name, content := range originals {
		path := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		writeFile(t, path, []byte(content))
	}
	now := time.Date(2026, 9, 30, 12, 34, 56, 0, time.FixedZone("JST", 9*60*60))
	paths, err := setDates(root, now)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, path := range paths {
		name, err := filepath.Rel(root, path)
		if err != nil {
			t.Fatal(err)
		}
		names = append(names, filepath.ToSlash(name))
	}
	wantNames := []string{"implicit.md", "multiline.md", "nested.md", "new/index.md", "windows.md"}
	if !reflect.DeepEqual(names, wantNames) {
		t.Fatalf("updated = %v, want %v", names, wantNames)
	}
	for name, original := range originals {
		t.Run(name, func(t *testing.T) {
			want := []byte(original)
			for _, changed := range wantNames {
				if name == changed {
					newline := "\n"
					if bytes.Contains(want, []byte("\r\n")) {
						newline = "\r\n"
					}
					want = bytes.Replace(want, []byte("+++"+newline), []byte("+++"+newline+"date = '2026-09-30T12:34:56+09:00'"+newline), 1)
				}
			}
			if got := readFile(t, filepath.Join(root, name)); !bytes.Equal(got, want) {
				t.Errorf("content = %q, want %q", got, want)
			}
		})
	}
	paths, err = setDates(root, now.Add(24*time.Hour))
	if err != nil || len(paths) != 0 {
		t.Fatalf("second run = %v, %v; want no updates", paths, err)
	}
}

func TestInvalidFrontMatterDoesNotPartiallyWrite(t *testing.T) {
	for _, tc := range []struct {
		name    string
		content string
	}{
		{"missing", "---\ntitle: bad\n---\n"},
		{"unclosed", "+++\ntitle = 'bad'\n"},
		{"invalid TOML", "+++\ndraft = !!!\n+++\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			original := []byte("+++\ndraft = false\n+++\n")
			path := filepath.Join(root, "a.md")
			writeFile(t, path, original)
			writeFile(t, filepath.Join(root, "z.md"), []byte(tc.content))
			if _, err := setDates(root, time.Now()); err == nil {
				t.Fatal("expected invalid front matter to fail")
			}
			if got := readFile(t, path); !bytes.Equal(got, original) {
				t.Errorf("valid article changed: %q", got)
			}
		})
	}
}

func writeFile(t *testing.T, path string, content []byte) {
	t.Helper()
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return content
}
