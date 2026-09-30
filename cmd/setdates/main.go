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
	paths, err := articlePaths(root)
	if err != nil {
		return nil, err
	}

	type update struct {
		path    string
		content []byte
	}
	var updates []update
	// Prepare every update before writing, so invalid front matter changes no files.
	for _, path := range paths {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		updated, err := withDate(content, now)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if updated != nil {
			updates = append(updates, update{path: path, content: updated})
		}
	}

	var changed []string
	for _, update := range updates {
		if err := os.WriteFile(update.path, update.content, 0o644); err != nil {
			return changed, err
		}
		changed = append(changed, update.path)
	}
	return changed, nil
}

func articlePaths(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !entry.IsDir() && filepath.Ext(path) == ".md" && entry.Name() != "_index.md" {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}

func withDate(content []byte, now time.Time) ([]byte, error) {
	metadata, err := parseFrontMatter(content)
	if err != nil {
		return nil, err
	}
	if _, exists := metadata["date"]; exists || metadata["draft"] == true {
		return nil, nil
	}

	opening, rest, _ := bytes.Cut(content, []byte("\n"))
	newline := "\n"
	if bytes.HasSuffix(opening, []byte("\r")) {
		newline = "\r\n"
	}
	return fmt.Appendf(nil, "%s\ndate = '%s'%s%s", opening, now.Format(time.RFC3339), newline, rest), nil
}

func parseFrontMatter(content []byte) (map[string]any, error) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	if string(bytes.TrimSpace(lines[0])) != "+++" {
		return nil, fmt.Errorf("expected TOML front matter")
	}
	for i := 1; i < len(lines); i++ {
		if string(bytes.TrimSpace(lines[i])) != "+++" {
			continue
		}
		var metadata map[string]any
		if err := toml.Unmarshal(bytes.Join(lines[1:i], nil), &metadata); err != nil {
			return nil, err
		}
		return metadata, nil
	}
	return nil, fmt.Errorf("unclosed front matter")
}
