package extensions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

const galleryIndexURL = "https://raw.githubusercontent.com/wayback09/Aether-Extensions/main/index.json"

// GalleryExtension represents an extension in the Aether Registry.
// Trust tier is assigned by the Aether team in the registry — never by the extension itself.
type GalleryExtension struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Version     string `json:"version"`
	Trust       string `json:"trust"`
	URL         string `json:"url"`
}

var (
	galleryCache    []GalleryExtension
	galleryCacheAt  time.Time
	galleryCacheMu  sync.Mutex
	galleryCacheTTL = 5 * time.Minute
)

// GetGalleryExtensions returns the live registry from GitHub, with a 5-minute in-memory cache.
// Gracefully falls back to an empty slice if the user is offline or GitHub is unreachable.
func GetGalleryExtensions() []GalleryExtension {
	galleryCacheMu.Lock()
	defer galleryCacheMu.Unlock()

	if galleryCache != nil && time.Since(galleryCacheAt) < galleryCacheTTL {
		return galleryCache
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(galleryIndexURL)
	if err != nil {
		fmt.Printf("[Gallery] Could not fetch registry (offline?): %v\n", err)
		return galleryCache // return stale cache rather than nothing
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[Gallery] Registry returned HTTP %d\n", resp.StatusCode)
		return galleryCache
	}

	var entries []GalleryExtension
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		fmt.Printf("[Gallery] Failed to parse registry: %v\n", err)
		return galleryCache
	}

	galleryCache = entries
	galleryCacheAt = time.Now()
	fmt.Printf("[Gallery] Fetched %d extensions from registry\n", len(entries))
	return galleryCache
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
