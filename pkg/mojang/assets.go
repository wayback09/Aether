package mojang

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"context"
)

// AssetIndexData represents the parsed asset index JSON
type AssetIndexData struct {
	Objects map[string]AssetObject `json:"objects"`
}

type AssetObject struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

// DownloadAssets fetches the asset index and downloads all game assets to a shared directory.
func (e *DownloadEngine) DownloadAssets(ctx context.Context, assetIndex AssetIndex, assetsDir string) error {
	// Download the asset index JSON
	indexDir := filepath.Join(assetsDir, "indexes")
	if err := os.MkdirAll(indexDir, 0755); err != nil {
		return err
	}

	indexPath := filepath.Join(indexDir, assetIndex.ID+".json")
	if err := e.downloadFile(assetIndex.URL, indexPath); err != nil {
		return fmt.Errorf("failed to download asset index: %w", err)
	}

	// Parse the asset index
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	var indexData AssetIndexData
	if err := json.Unmarshal(data, &indexData); err != nil {
		return fmt.Errorf("failed to parse asset index: %w", err)
	}

	// Download all asset objects concurrently
	objectsDir := filepath.Join(assetsDir, "objects")
	totalAssets := len(indexData.Objects)
	var completed int
	var mu sync.Mutex
	var wg sync.WaitGroup
	errors := make(chan error, totalAssets)
	sem := make(chan struct{}, 15) // Limit concurrency

	for name, obj := range indexData.Objects {
		prefix := obj.Hash[:2]
		objPath := filepath.Join(objectsDir, prefix, obj.Hash)

		// Skip if already downloaded
		if info, err := os.Stat(objPath); err == nil && info.Size() == obj.Size {
			mu.Lock()
			completed++
			mu.Unlock()
			continue
		}

		wg.Add(1)
		go func(assetName string, assetObj AssetObject, dest string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			url := fmt.Sprintf("https://resources.download.minecraft.net/%s/%s", assetObj.Hash[:2], assetObj.Hash)
			if err := e.downloadFile(url, dest); err != nil {
				errors <- fmt.Errorf("asset %s: %w", assetName, err)
				return
			}

			mu.Lock()
			completed++
			pct := float64(completed) / float64(totalAssets) * 100
			mu.Unlock()

			// Only emit progress every ~5% to avoid flooding the event bus
			if int(pct)%5 == 0 {
				runtime.EventsEmit(ctx, "instance:progress", map[string]interface{}{
					"id":       e.instance,
					"progress": pct,
					"status":   fmt.Sprintf("Downloading assets (%d/%d)", completed, totalAssets),
				})
			}
		}(name, obj, objPath)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
