package mojang

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractNatives extracts native DLLs/SOs from native library JARs into the natives directory.
func ExtractNatives(libraries []Library, librariesDir, nativesDir string) error {
	if err := os.MkdirAll(nativesDir, 0755); err != nil {
		return err
	}

	for _, lib := range libraries {
		if !IsLibraryAllowed(lib) {
			continue
		}

		// Only process native libraries (name contains "natives-")
		if !strings.Contains(lib.Name, "natives-") {
			continue
		}

		if lib.Downloads.Artifact.Path == "" {
			continue
		}

		jarPath := filepath.Join(librariesDir, lib.Downloads.Artifact.Path)

		// Skip if JAR doesn't exist
		if _, err := os.Stat(jarPath); os.IsNotExist(err) {
			continue
		}

		if err := extractJar(jarPath, nativesDir); err != nil {
			return fmt.Errorf("failed to extract natives from %s: %w", lib.Name, err)
		}
	}

	return nil
}

// extractJar extracts relevant files from a JAR (zip) into the destination directory.
// Skips META-INF/ and directories.
func extractJar(jarPath, destDir string) error {
	r, err := zip.OpenReader(jarPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Skip directories and META-INF
		if f.FileInfo().IsDir() || strings.HasPrefix(f.Name, "META-INF") {
			continue
		}

		destPath := filepath.Join(destDir, filepath.Base(f.Name))

		// Skip if already extracted
		if _, err := os.Stat(destPath); err == nil {
			continue
		}

		src, err := f.Open()
		if err != nil {
			return err
		}

		dst, err := os.Create(destPath)
		if err != nil {
			src.Close()
			return err
		}

		_, err = io.Copy(dst, src)
		src.Close()
		dst.Close()

		if err != nil {
			return err
		}
	}

	return nil
}
