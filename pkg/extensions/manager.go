package extensions

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"Aether/pkg/fs"
)

// Manager handles the lifecycle of all extensions
type Manager struct {
	LoadedExtensions map[string]Extension
	sandboxes        map[string]*Sandbox
}

// GlobalManager holds the singleton instance for UI access
var GlobalManager *Manager

// NewManager creates a new extension manager
func NewManager() *Manager {
	return &Manager{
		LoadedExtensions: make(map[string]Extension),
		sandboxes:        make(map[string]*Sandbox),
	}
}

// LoadAll scans the directory, parses manifests, and executes the JS isolate
func (m *Manager) LoadAll() error {
	extDir := filepath.Join(fs.GetDataDir(), "extensions")
	entries, err := os.ReadDir(extDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(extDir, entry.Name(), "manifest.json")
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				continue
			}

			var manifest Manifest
			if err := json.Unmarshal(data, &manifest); err != nil {
				continue
			}

			if manifest.ID == "" {
				manifest.ID = entry.Name()
			}

			ext := Extension{
				ID:      manifest.ID,
				Name:    manifest.Name,
				Version: manifest.Version,
				Author:  manifest.Author,
				Status:  "Running",
				Memory:  "0MB",
				CPU:     "0%",
			}

			sandbox := NewSandbox(manifest)
			m.sandboxes[manifest.ID] = sandbox

			if manifest.Main != "" {
				scriptPath := filepath.Join(extDir, entry.Name(), manifest.Main)
				scriptData, err := os.ReadFile(scriptPath)
				if err == nil {
					if err := sandbox.Execute(string(scriptData)); err != nil {
						fmt.Printf("[Manager] Failed to execute %s for %s: %v\n", manifest.Main, manifest.ID, err)
						ext.Status = "Error"
					} else {
						fmt.Printf("[Manager] Successfully loaded extension isolate: %s\n", manifest.ID)
					}
				} else {
					fmt.Printf("[Manager] Missing main script %s for %s\n", manifest.Main, manifest.ID)
					ext.Status = "Error (Missing Main)"
				}
			}

			m.LoadedExtensions[manifest.ID] = ext
		}
	}
	return nil
}

// GetExtensions returns the list of loaded extensions
func (m *Manager) GetExtensions() []Extension {
	var list []Extension
	for _, ext := range m.LoadedExtensions {
		list = append(list, ext)
	}
	return list
}
