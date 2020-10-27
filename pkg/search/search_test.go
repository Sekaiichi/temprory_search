package search

import (
	"context"
	"testing"
)

func TestSearch_All(t *testing.T) {
	root := context.Background()
	files := make([]string, 0)
	files = append(files, "lorem.txt")
	files = append(files, "go.txt")
	files = append(files, "php.txt")

	results := All(root, "concurrency", files)

	if len(results) != 0 {
		t.Error("Error")
	}
}

func TestSearch_Any(t *testing.T) {
	root := context.Background()
	files := make([]string, 0)
	files = append(files, "lorem.txt")
	files = append(files, "go.txt")
	files = append(files, "php.txt")

	results := Any(root, "concurrency", files)

	if len(results) > 1 {
		t.Error("Error")
	}
}
