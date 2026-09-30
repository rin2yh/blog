// Command setdates records missing publication dates without rewriting article content.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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

	// Validate all front matter before writing. Keep only paths in memory.
	var pending []string
	for _, path := range paths {
		needed, err := needsPublicationDate(path)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if needed {
			pending = append(pending, path)
		}
	}

	var changed []string
	for _, path := range pending {
		updated, err := writePublicationDate(path, now)
		if err != nil {
			return changed, fmt.Errorf("%s: %w", path, err)
		}
		if updated {
			changed = append(changed, path)
		}
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

func needsPublicationDate(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	_, metadata, err := readFrontMatter(bufio.NewReader(file))
	if err != nil {
		return false, err
	}
	return needsDate(metadata), nil
}

func needsDate(metadata map[string]any) bool {
	_, exists := metadata["date"]
	return !exists && metadata["draft"] != true
}

func writePublicationDate(path string, now time.Time) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	frontMatter, metadata, err := readFrontMatter(reader)
	if err != nil || !needsDate(metadata) {
		return false, err
	}

	temp, err := os.CreateTemp(filepath.Dir(path), ".setdates-*.tmp")
	if err != nil {
		return false, err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()
	info, err := file.Stat()
	if err != nil {
		return false, err
	}
	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		return false, err
	}
	if err := withDate(temp, frontMatter, now); err != nil {
		return false, err
	}
	// Copy the body unchanged without loading it into memory.
	if _, err := io.Copy(temp, reader); err != nil {
		return false, err
	}
	if err := temp.Close(); err != nil {
		return false, err
	}
	if err := os.Rename(temp.Name(), path); err != nil {
		return false, err
	}
	return true, nil
}

func withDate(writer io.Writer, frontMatter []byte, now time.Time) error {
	opening, rest, _ := bytes.Cut(frontMatter, []byte("\n"))
	newline := "\n"
	if bytes.HasSuffix(opening, []byte("\r")) {
		newline = "\r\n"
	}
	_, err := fmt.Fprintf(writer, "%s\ndate = '%s'%s%s", opening, now.Format(time.RFC3339), newline, rest)
	return err
}

func readFrontMatter(reader *bufio.Reader) ([]byte, map[string]any, error) {
	opening, err := reader.ReadString('\n')
	if strings.TrimSpace(opening) != "+++" {
		return nil, nil, fmt.Errorf("expected TOML front matter")
	}
	var header bytes.Buffer
	for err == nil {
		var line string
		line, err = reader.ReadString('\n')
		if strings.TrimSpace(line) == "+++" {
			var metadata map[string]any
			if err := toml.Unmarshal(header.Bytes(), &metadata); err != nil {
				return nil, nil, err
			}
			return []byte(opening + header.String() + line), metadata, nil
		}
		header.WriteString(line)
	}
	if err != io.EOF {
		return nil, nil, err
	}
	return nil, nil, fmt.Errorf("unclosed front matter")
}
