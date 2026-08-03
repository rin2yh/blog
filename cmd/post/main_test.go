package main

import (
	"reflect"
	"testing"
)

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
