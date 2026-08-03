package main

import (
	"fmt"
	"os"
	"os/exec"

	"blog/internal/slug"
)

const contentDir = "content/post"

func buildHugoArgs(title, kind string) []string {
	args := []string{"new"}
	if kind != "" {
		args = append(args, "--kind", kind)
	}
	args = append(args, fmt.Sprintf("post/%s/index.md", title))
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

	s, err := slug.Generate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating slug: %v\n", err)
		os.Exit(1)
	}

	args := buildHugoArgs(title, kind)
	cmd := exec.Command("hugo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running hugo: %v\n", err)
		os.Exit(1)
	}

	oldPath := fmt.Sprintf("%s/%s", contentDir, title)
	newFilename := fmt.Sprintf("%s-%s", title, s)
	newPath := fmt.Sprintf("%s/%s", contentDir, newFilename)

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created: %s\n", newPath)
	fmt.Printf("URL: /post/%s/\n", newFilename)
}
