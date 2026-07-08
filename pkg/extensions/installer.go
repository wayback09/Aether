package extensions

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"Aether/pkg/fs"
)

// InstallFromZip extracts a zip file to a temporary location, validates the manifest,
// and moves it to the extensions directory under its ID.
func InstallFromZip(zipPath string) error {
	extDir := filepath.Join(fs.GetDataDir(), "extensions")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp(extDir, "installing-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up temp dir in case of failure

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	// Extract all files
	for _, f := range r.File {
		// Prevent ZipSlip vulnerability
		if strings.Contains(f.Name, "..") {
			continue
		}

		fpath := filepath.Join(tempDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	// Some zips contain a single root folder (e.g. 'my-extension/manifest.json')
	// Let's find the manifest.json
	manifestPath := ""
	rootDir := tempDir

	err = filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "manifest.json" {
			manifestPath = path
			rootDir = filepath.Dir(path)
			return filepath.SkipDir // Stop walking
		}
		return nil
	})

	if err != nil {
		return err
	}

	if manifestPath == "" {
		return fmt.Errorf("invalid extension: manifest.json not found in zip")
	}

	// Parse manifest
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("failed to parse manifest: %w", err)
	}

	if manifest.ID == "" {
		return fmt.Errorf("invalid manifest: missing 'id'")
	}

	// Target directory
	targetDir := filepath.Join(extDir, manifest.ID)
	
	// Remove old version if it exists
	os.RemoveAll(targetDir)

	// Move the rootDir to the targetDir
	if err := os.Rename(rootDir, targetDir); err != nil {
		// Fallback for cross-device rename issues if needed, though they are in the same folder here
		return fmt.Errorf("failed to move extension to final directory: %w", err)
	}

	return nil
}
