package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"blog/internal/post"
)

func buildHugoArgs(title, kind string) []string {
	args := []string{"new"}
	if kind != "" {
		args = append(args, "--kind", kind)
	}
	args = append(args, fmt.Sprintf("%s/%s/index.md", post.Section, title))
	return args
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/post <title> [kind]")
		fmt.Fprintln(os.Stderr, "  title: 記事タイトル（英数字ハイフン、例: my-new-post）")
		fmt.Fprintln(os.Stderr, "  kind:  アーキタイプ（省略可、例: external）")
		os.Exit(1)
	}

	title := os.Args[1]
	kind := ""
	if len(os.Args) >= 3 {
		kind = os.Args[2]
	}

	newFilename, err := post.DirName(title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating slug: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("hugo", buildHugoArgs(title, kind)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running hugo: %v\n", err)
		os.Exit(1)
	}

	oldPath := filepath.Join(post.Dir, title)
	newPath := filepath.Join(post.Dir, newFilename)

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created: %s\n", newPath)
	fmt.Printf("URL: /%s/%s/\n", post.Section, newFilename)
}
