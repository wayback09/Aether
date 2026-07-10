package java

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"Aether/pkg/fs"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// adoptiumAsset is a partial unmarshal of the Adoptium API response.
type adoptiumAsset struct {
	Binary struct {
		Package struct {
			Link     string `json:"link"`
			Checksum string `json:"checksum"`
			Size     int64  `json:"size"`
		} `json:"package"`
	} `json:"binary"`
}

// GetManagedJavaPath returns the path to the java(.exe) binary for the
// given major version inside Aether's managed runtimes directory.
func GetManagedJavaPath(majorVersion int) string {
	base := getManagedJavaDir(majorVersion)
	if runtime.GOOS == "windows" {
		return filepath.Join(base, "bin", "java.exe")
	}
	return filepath.Join(base, "bin", "java")
}

// IsManagedJavaInstalled returns true if we already have a managed JRE for this version.
func IsManagedJavaInstalled(majorVersion int) bool {
	path := GetManagedJavaPath(majorVersion)
	_, err := os.Stat(path)
	return err == nil
}

// getManagedJavaDir returns the directory where the managed JRE is stored.
func getManagedJavaDir(majorVersion int) string {
	return filepath.Join(fs.GetDataDir(), "runtimes", fmt.Sprintf("java-%d", majorVersion))
}

// DownloadJava downloads and installs a managed JRE from Adoptium for the
// given major version. Emits progress events on ctx via Wails.
func DownloadJava(ctx context.Context, majorVersion int) error {
	os_ := "windows"
	arch := "x64"
	if runtime.GOARCH == "arm64" {
		arch = "aarch64"
	}

	apiURL := fmt.Sprintf(
		"https://api.adoptium.net/v3/assets/latest/%d/hotspot?os=%s&architecture=%s&image_type=jre",
		majorVersion, os_, arch,
	)

	wailsruntime.EventsEmit(ctx, "java:status", map[string]interface{}{
		"version": majorVersion,
		"phase":   "fetching",
		"message": fmt.Sprintf("Fetching Java %d download info...", majorVersion),
	})

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("adoptium API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("adoptium API returned status %s", resp.Status)
	}

	var assets []adoptiumAsset
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return fmt.Errorf("failed to parse adoptium API response: %w", err)
	}
	if len(assets) == 0 {
		return fmt.Errorf("no JRE packages found for Java %d on %s/%s", majorVersion, os_, arch)
	}

	downloadURL := assets[0].Binary.Package.Link
	totalSize := assets[0].Binary.Package.Size

	wailsruntime.EventsEmit(ctx, "java:status", map[string]interface{}{
		"version": majorVersion,
		"phase":   "downloading",
		"message": fmt.Sprintf("Downloading Java %d...", majorVersion),
		"total":   totalSize,
	})

	// Download to a temp file
	tmpFile, err := os.CreateTemp("", fmt.Sprintf("aether-java-%d-*.zip", majorVersion))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	dlResp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to start download: %w", err)
	}
	defer dlResp.Body.Close()

	// Copy with progress reporting
	written, err := copyWithProgress(ctx, tmpFile, dlResp.Body, totalSize, majorVersion)
	if err != nil {
		return fmt.Errorf("download failed after %d bytes: %w", written, err)
	}
	tmpFile.Close()

	wailsruntime.EventsEmit(ctx, "java:status", map[string]interface{}{
		"version": majorVersion,
		"phase":   "extracting",
		"message": fmt.Sprintf("Extracting Java %d...", majorVersion),
	})

	destDir := getManagedJavaDir(majorVersion)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create runtime dir: %w", err)
	}

	if err := extractZip(tmpFile.Name(), destDir); err != nil {
		return fmt.Errorf("failed to extract JRE: %w", err)
	}

	wailsruntime.EventsEmit(ctx, "java:status", map[string]interface{}{
		"version": majorVersion,
		"phase":   "done",
		"message": fmt.Sprintf("Java %d ready.", majorVersion),
	})

	return nil
}

// copyWithProgress copies from src to dst and emits periodic progress events.
func copyWithProgress(ctx context.Context, dst io.Writer, src io.Reader, total int64, majorVersion int) (int64, error) {
	buf := make([]byte, 32*1024)
	var written int64
	for {
		nr, err := src.Read(buf)
		if nr > 0 {
			nw, werr := dst.Write(buf[:nr])
			written += int64(nw)
			if werr != nil {
				return written, werr
			}
			// Emit progress every ~2 MB
			if total > 0 && written%(2*1024*1024) < int64(nr) {
				pct := int(written * 100 / total)
				wailsruntime.EventsEmit(ctx, "java:status", map[string]interface{}{
					"version":  majorVersion,
					"phase":    "downloading",
					"message":  fmt.Sprintf("Downloading Java %d... %d%%", majorVersion, pct),
					"progress": pct,
					"total":    total,
				})
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

// extractZip extracts a zip archive into destDir, stripping the top-level
// directory that Adoptium bundles include (e.g. jdk-21.0.1-jre/).
func extractZip(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Determine the common prefix to strip (the bundled root folder name)
	prefix := ""
	if len(r.File) > 0 {
		parts := strings.SplitN(r.File[0].Name, "/", 2)
		if len(parts) > 1 {
			prefix = parts[0] + "/"
		}
	}

	for _, f := range r.File {
		// Strip the top-level directory
		relPath := strings.TrimPrefix(f.Name, prefix)
		if relPath == "" {
			continue
		}

		target := filepath.Join(destDir, filepath.FromSlash(relPath))

		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
