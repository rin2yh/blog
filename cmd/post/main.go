package main

import (
	"crypto/rand"
	"fmt"
	"os"
	"os/exec"
)

const (
	slugLength = 14
	charset    = "abcdefghijklmnopqrstuvwxyz0123456789"
	contentDir = "content/post"
)

func generateSlug() (string, error) {
	b := make([]byte, slugLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/post <title>")
		fmt.Fprintln(os.Stderr, "  title: 記事タイトル（英数字ハイフン、例: my-new-post）")
		os.Exit(1)
	}

	title := os.Args[1]
	slug, err := generateSlug()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating slug: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command("hugo", "new", fmt.Sprintf("post/%s/index.md", title))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running hugo: %v\n", err)
		os.Exit(1)
	}

	oldPath := fmt.Sprintf("%s/%s", contentDir, title)
	newFilename := fmt.Sprintf("%s-%s", title, slug)
	newPath := fmt.Sprintf("%s/%s", contentDir, newFilename)

	if err := os.Rename(oldPath, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error renaming file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created: %s\n", newPath)
	fmt.Printf("URL: /post/%s/\n", newFilename)
}
