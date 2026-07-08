package mojang

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// DownloadEngine handles fetching assets concurrently
type DownloadEngine struct {
	ctx      context.Context
	instance string
	basePath string
}

func NewDownloadEngine(ctx context.Context, instanceID, basePath string) *DownloadEngine {
	return &DownloadEngine{
		ctx:      ctx,
		instance: instanceID,
		basePath: basePath,
	}
}

// Install processes the VersionInfo and downloads everything needed to launch
func (e *DownloadEngine) Install(info *VersionInfo, assetsDir string) error {
	// Count only the libraries we'll actually download + client jar
	allowedLibs := []Library{}
	for _, lib := range info.Libraries {
		if lib.Downloads.Artifact.URL != "" && IsLibraryAllowed(lib) {
			allowedLibs = append(allowedLibs, lib)
		}
	}

	var wg sync.WaitGroup
	errors := make(chan error, len(allowedLibs)+5)

	totalItems := len(allowedLibs) + 1 // +1 for client jar
	var completed int
	var mu sync.Mutex

	reportProgress := func(name string, phase string) {
		mu.Lock()
		completed++
		pct := float64(completed) / float64(totalItems) * 100
		mu.Unlock()

		runtime.EventsEmit(e.ctx, "instance:progress", map[string]interface{}{
			"id":       e.instance,
			"progress": pct,
			"status":   fmt.Sprintf("[%s] %s", phase, name),
		})
	}

	// 1. Download Client Jar
	wg.Add(1)
	go func() {
		defer wg.Done()
		clientPath := filepath.Join(e.basePath, "bin", fmt.Sprintf("%s.jar", info.ID))
		if err := e.downloadFile(info.Downloads.Client.URL, clientPath); err != nil {
			errors <- err
			return
		}
		reportProgress("client.jar", "Core")
	}()

	// 2. Download Libraries concurrently
	sem := make(chan struct{}, 10)

	for _, lib := range allowedLibs {
		wg.Add(1)
		go func(l Library) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			libPath := filepath.Join(e.basePath, "libraries", l.Downloads.Artifact.Path)
			if err := e.downloadFile(l.Downloads.Artifact.URL, libPath); err != nil {
				errors <- err
				return
			}
			reportProgress(l.Name, "Libraries")
		}(lib)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
	}

	// 3. Download Assets
	runtime.EventsEmit(e.ctx, "instance:progress", map[string]interface{}{
		"id":       e.instance,
		"progress": 50,
		"status":   "Downloading game assets...",
	})

	if err := e.DownloadAssets(e.ctx, info.AssetIndex, assetsDir); err != nil {
		return fmt.Errorf("asset download error: %w", err)
	}

	// 4. Extract Native Libraries
	runtime.EventsEmit(e.ctx, "instance:progress", map[string]interface{}{
		"id":       e.instance,
		"progress": 90,
		"status":   "Extracting native libraries...",
	})

	nativesDir := filepath.Join(e.basePath, "natives")
	librariesDir := filepath.Join(e.basePath, "libraries")
	if err := ExtractNatives(info.Libraries, librariesDir, nativesDir); err != nil {
		return fmt.Errorf("native extraction error: %w", err)
	}

	// 5. Download Log4j config file
	if info.Logging.Client.File.URL != "" {
		logConfigPath := filepath.Join(e.basePath, info.Logging.Client.File.ID)
		if err := e.downloadFile(info.Logging.Client.File.URL, logConfigPath); err != nil {
			fmt.Printf("Warning: failed to download log config: %v\n", err)
			// Non-fatal, continue
		}
	}

	// 6. Save version.json to disk (needed by launcher to read mainClass, arguments, etc.)
	versionData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal version info: %w", err)
	}
	versionPath := filepath.Join(e.basePath, "version.json")
	if err := os.WriteFile(versionPath, versionData, 0644); err != nil {
		return fmt.Errorf("failed to save version.json: %w", err)
	}

	runtime.EventsEmit(e.ctx, "instance:progress", map[string]interface{}{
		"id":       e.instance,
		"progress": 100,
		"status":   "Installation Complete",
	})

	return nil
}

func (e *DownloadEngine) downloadFile(url string, dest string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Skip if file already exists (simplistic check for this phase)
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
