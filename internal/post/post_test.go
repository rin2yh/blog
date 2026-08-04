package post

import (
	"strings"
	"testing"
)

func TestDirName(t *testing.T) {
	name, err := DirName("my-post")
	if err != nil {
		t.Fatalf("DirName() returned error: %v", err)
	}

	const prefix = "my-post-"
	if !strings.HasPrefix(name, prefix) {
		t.Fatalf("DirName() = %q, want it to start with %q", name, prefix)
	}

	s := strings.TrimPrefix(name, prefix)
	if len(s) != slugLength {
		t.Errorf("DirName() slug is %d characters, want %d", len(s), slugLength)
	}
	for i, c := range s {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("DirName() slug has invalid character %q at position %d", c, i)
		}
	}
}
