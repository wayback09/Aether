package fs

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ContainedPath resolves relativePath below root and rejects path traversal.
func ContainedPath(root, relativePath string) (string, error) {
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(filepath.Join(rootAbs, relativePath))
	if err != nil {
		return "", err
	}

	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes root directory")
	}
	return targetAbs, nil
}
