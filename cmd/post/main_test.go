package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGenerateSlug_Length(t *testing.T) {
	slug, err := generateSlug()
	if err != nil {
		t.Fatalf("generateSlug() returned error: %v", err)
	}

	if len(slug) != slugLength {
		t.Errorf("generateSlug() returned slug of length %d, want %d", len(slug), slugLength)
	}
}

func TestGenerateSlug_Charset(t *testing.T) {
	slug, err := generateSlug()
	if err != nil {
		t.Fatalf("generateSlug() returned error: %v", err)
	}

	for i, c := range slug {
		if !strings.ContainsRune(charset, c) {
			t.Errorf("generateSlug() returned invalid character %q at position %d", c, i)
		}
	}
}

func TestGenerateSlug_Uniqueness(t *testing.T) {
	const iterations = 100
	slugs := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		slug, err := generateSlug()
		if err != nil {
			t.Fatalf("generateSlug() returned error on iteration %d: %v", i, err)
		}

		if slugs[slug] {
			t.Errorf("generateSlug() returned duplicate slug %q on iteration %d", slug, i)
		}
		slugs[slug] = true
	}
}

func TestBuildHugoArgs_WithoutKind(t *testing.T) {
	args := buildHugoArgs("my-post", "")
	expected := []string{"new", "post/my-post/index.md"}

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("buildHugoArgs() = %v, want %v", args, expected)
	}
}

func TestBuildHugoArgs_WithKind(t *testing.T) {
	args := buildHugoArgs("my-post", "external")
	expected := []string{"new", "--kind", "external", "post/my-post/index.md"}

	if !reflect.DeepEqual(args, expected) {
		t.Errorf("buildHugoArgs() = %v, want %v", args, expected)
	}
}
