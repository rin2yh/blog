package slug

import (
	"strings"
	"testing"
)

func TestGenerate_Length(t *testing.T) {
	s, err := Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	if len(s) != Length {
		t.Errorf("Generate() returned slug of length %d, want %d", len(s), Length)
	}
}

func TestGenerate_Charset(t *testing.T) {
	s, err := Generate()
	if err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	for i, c := range s {
		if !strings.ContainsRune(Charset, c) {
			t.Errorf("Generate() returned invalid character %q at position %d", c, i)
		}
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	const iterations = 100
	slugs := make(map[string]bool)

	for i := 0; i < iterations; i++ {
		s, err := Generate()
		if err != nil {
			t.Fatalf("Generate() returned error on iteration %d: %v", i, err)
		}

		if slugs[s] {
			t.Errorf("Generate() returned duplicate slug %q on iteration %d", s, i)
		}
		slugs[s] = true
	}
}
