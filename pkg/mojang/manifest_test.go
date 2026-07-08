package mojang

import (
	"testing"
)

func TestGetVersionManifest(t *testing.T) {
	manifest, err := GetVersionManifest()
	if err != nil {
		t.Fatalf("Failed to get manifest: %v", err)
	}

	if manifest.Latest.Release == "" {
		t.Error("Expected Latest.Release to be populated")
	}

	if len(manifest.Versions) == 0 {
		t.Error("Expected versions array to have items")
	}
}

func TestGetVersionInfo(t *testing.T) {
	info, err := GetVersionInfo("1.20.4")
	if err != nil {
		t.Fatalf("Failed to get version info: %v", err)
	}

	if info.Downloads.Client.URL == "" {
		t.Error("Expected client download URL")
	}

	if len(info.Libraries) == 0 {
		t.Error("Expected libraries to be parsed")
	}
}
