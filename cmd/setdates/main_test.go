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
	root := t.TempDir()
	if err := os.CopyFS(root, os.DirFS("testdata/input")); err != nil {
		t.Fatal(err)
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
	compareFiles(t, root, "testdata/golden")

	paths, err = setDates(root, now.Add(24*time.Hour))
	if err != nil || len(paths) != 0 {
		t.Fatalf("second run = %v, %v; want no updates", paths, err)
	}
	compareFiles(t, root, "testdata/golden")
}

func TestInvalidFrontMatterDoesNotPartiallyWrite(t *testing.T) {
	cases, err := os.ReadDir("testdata/invalid")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		t.Run(tc.Name(), func(t *testing.T) {
			fixture := filepath.Join("testdata/invalid", tc.Name())
			root := t.TempDir()
			if err := os.CopyFS(root, os.DirFS(fixture)); err != nil {
				t.Fatal(err)
			}
			if _, err := setDates(root, time.Now()); err == nil {
				t.Fatal("expected invalid front matter to fail")
			}
			compareFiles(t, root, fixture)
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
