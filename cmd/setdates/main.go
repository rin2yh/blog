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
	paths, err := collectArticlePaths("content/post")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	now := time.Now().In(zone)
	for _, path := range paths {
		updated, err := setDate(path, now)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			os.Exit(1)
		}
		if updated {
			fmt.Println(path)
		}
	}
}

func collectArticlePaths(root string) ([]string, error) {
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

func needsDate(metadata map[string]any) bool {
	_, exists := metadata["date"]
	return !exists && metadata["draft"] != true
}

func setDate(path string, now time.Time) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	frontMatter, metadata, err := readFrontMatter(reader)
	if err != nil {
		return false, err
	}
	if !needsDate(metadata) {
		return false, nil
	}

	content := io.MultiReader(bytes.NewReader(addDate(frontMatter, now)), reader)
	if err := replaceFile(file, content); err != nil {
		return false, err
	}
	return true, nil
}

func addDate(frontMatter []byte, now time.Time) []byte {
	opening, rest, _ := bytes.Cut(frontMatter, []byte("\n"))
	newline := "\n"
	if bytes.HasSuffix(opening, []byte("\r")) {
		newline = "\r\n"
	}
	return fmt.Appendf(nil, "%s\ndate = '%s'%s%s", opening, now.Format(time.RFC3339), newline, rest)
}

// replaceFile streams content to a temporary file and replaces the original only after success.
func replaceFile(file *os.File, content io.Reader) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	temp, err := os.CreateTemp(filepath.Dir(file.Name()), ".setdates-*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(temp.Name())
	defer temp.Close()

	if err := temp.Chmod(info.Mode().Perm()); err != nil {
		return err
	}
	if _, err := io.Copy(temp, content); err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return os.Rename(temp.Name(), file.Name())
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
