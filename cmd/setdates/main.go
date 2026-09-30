// Command setdates records missing publication dates without rewriting article content.
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"
	_ "time/tzdata"

	"github.com/pelletier/go-toml/v2"
)

func main() {
	zone, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	paths, err := setDates("content/post", time.Now().In(zone))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, path := range paths {
		fmt.Println(path)
	}
}

func setDates(root string, now time.Time) ([]string, error) {
	type update struct {
		path    string
		content []byte
	}
	var updates []update
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".md" || entry.Name() == "_index.md" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		updated, err := withDate(content, now)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if updated != nil {
			updates = append(updates, update{path, updated})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// Validate every front matter block before changing any file.
	var paths []string
	for _, u := range updates {
		if err := os.WriteFile(u.path, u.content, 0o644); err != nil {
			return paths, err
		}
		paths = append(paths, u.path)
	}
	return paths, nil
}

func withDate(content []byte, now time.Time) ([]byte, error) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	if len(lines) == 0 || string(bytes.TrimSpace(lines[0])) != "+++" {
		return nil, fmt.Errorf("expected TOML front matter")
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if string(bytes.TrimSpace(lines[i])) == "+++" {
			end = i
			break
		}
	}
	if end == -1 {
		return nil, fmt.Errorf("unclosed front matter")
	}
	var metadata map[string]any
	if err := toml.Unmarshal(bytes.Join(lines[1:end], nil), &metadata); err != nil {
		return nil, err
	}
	if _, exists := metadata["date"]; exists || metadata["draft"] == true {
		return nil, nil
	}
	newline := "\n"
	if bytes.HasSuffix(lines[0], []byte("\r\n")) {
		newline = "\r\n"
	}
	date := fmt.Sprintf("date = '%s'%s", now.Format(time.RFC3339), newline)
	updated := make([]byte, 0, len(content)+len(date))
	updated = append(updated, lines[0]...)
	updated = append(updated, date...)
	updated = append(updated, content[len(lines[0]):]...)
	return updated, nil
}
