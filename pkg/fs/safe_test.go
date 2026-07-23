package fs

import (
	"path/filepath"
	"testing"
)

func TestContainedPath(t *testing.T) {
	root := filepath.Join("test-data", "instances")

	inside, err := ContainedPath(root, filepath.Join("survival", "instance.json"))
	if err != nil || !filepath.IsAbs(inside) {
		t.Fatalf("expected an absolute contained path, got %q, %v", inside, err)
	}

	for _, relative := range []string{"..", filepath.Join("survival", "..", "..", "outside.txt")} {
		if _, err := ContainedPath(root, relative); err == nil {
			t.Fatalf("expected path traversal to be rejected: %q", relative)
		}
	}
}
