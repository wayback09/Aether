package extensions

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// GalleryExtension represents an extension in the Aether Registry.
// Trust tier is assigned server-side by the Aether team — not by the extension manifest.
type GalleryExtension struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Trust       string `json:"trust"`
	DownloadURL string `json:"downloadUrl"`
}

// GetGalleryExtensions fetches the list of available extensions from the Aether Registry.
// Returns an empty slice until the Aether Registry API is live.
func GetGalleryExtensions() []GalleryExtension {
	return []GalleryExtension{}
}

// DownloadAndInstallExtension downloads a zip file from a trusted Registry URL and installs it.
func DownloadAndInstallExtension(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download extension: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	tmpFile, err := os.CreateTemp("", "aether-ext-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}
	tmpFile.Close()

	return InstallFromZip(tmpFile.Name())
}
