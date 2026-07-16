package netutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Shared client. No global timeout so large downloads can take as long as needed.
// Timeout is enforced by the ResponseHeaderTimeout and Context.
var defaultClient = &http.Client{
	Transport: &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       30 * time.Second,
	},
}

// ProgressCallback is called periodically with the downloaded bytes and total bytes.
type ProgressCallback func(downloaded, total int64)

// DownloadFile downloads a file from url to dest, with support for resuming (Range requests)
// and exponential backoff retries.
func DownloadFile(ctx context.Context, url string, dest string, onProgress ProgressCallback) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	// Skip if final file already exists
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	tempDest := dest + ".tmp"
	maxRetries := 5
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := func() error {
			var startBytes int64 = 0

			// Check if temp file exists to resume
			if info, err := os.Stat(tempDest); err == nil {
				startBytes = info.Size()
			}

			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				return err
			}

			if startBytes > 0 {
				req.Header.Set("Range", fmt.Sprintf("bytes=%d-", startBytes))
			}

			resp, err := defaultClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			var out *os.File

			// If server supports range requests, it returns 206 Partial Content
			if resp.StatusCode == http.StatusPartialContent {
				out, err = os.OpenFile(tempDest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					return err
				}
			} else if resp.StatusCode == http.StatusOK {
				// Server didn't respect Range or we didn't send it, start from scratch
				startBytes = 0
				out, err = os.Create(tempDest)
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("unexpected status %d", resp.StatusCode)
			}
			defer out.Close()

			var totalBytes int64 = resp.ContentLength
			if totalBytes > 0 {
				totalBytes += startBytes // True total size
			}

			// Copy with progress reporting
			buf := make([]byte, 32*1024)
			var written int64 = startBytes

			for {
				nr, readErr := resp.Body.Read(buf)
				if nr > 0 {
					nw, writeErr := out.Write(buf[:nr])
					if nw > 0 {
						written += int64(nw)
						if onProgress != nil {
							onProgress(written, totalBytes)
						}
					}
					if writeErr != nil {
						return writeErr
					}
				}
				if readErr == io.EOF {
					break
				}
				if readErr != nil {
					return readErr
				}
			}

			return nil
		}()

		if err == nil {
			if err := os.Rename(tempDest, dest); err != nil {
				return fmt.Errorf("failed to rename temp file: %w", err)
			}
			return nil
		}

		lastErr = err

		// Exponential backoff: 1s, 2s, 4s, 8s
		select {
		case <-ctx.Done():
			return ctx.Err() // Abort immediately if cancelled
		case <-time.After(time.Duration(1<<i) * time.Second):
			// Retry after delay
		}
	}

	return fmt.Errorf("failed to download after %d retries, last error: %w", maxRetries, lastErr)
}
