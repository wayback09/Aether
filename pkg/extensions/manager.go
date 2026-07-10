package extensions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"Aether/pkg/fs"
	"Aether/pkg/instance"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Manager handles the lifecycle of all extensions
type Manager struct {
	ctx              context.Context
	serverURL        string
	LoadedExtensions map[string]Extension
	sandboxes        map[string]*Sandbox
	SidebarPages     []map[string]interface{}
	ModLoaders       map[string]ModLoaderConfig
}

// ModLoaderConfig holds info about a registered mod loader
type ModLoaderConfig struct {
	ID          string
	Name        string
	Description string
	ExtensionID string
	Callback    func(ctx map[string]interface{}) (map[string]interface{}, error)
}

// GlobalManager holds the singleton instance for UI access
var GlobalManager *Manager

// NewManager creates a new extension manager
func NewManager(ctx context.Context) *Manager {
	return &Manager{
		ctx:              ctx,
		LoadedExtensions: make(map[string]Extension),
		sandboxes:        make(map[string]*Sandbox),
		SidebarPages:     make([]map[string]interface{}, 0),
		ModLoaders:       make(map[string]ModLoaderConfig),
	}
}

// LoadAll scans the directory, parses manifests, and executes the JS isolate
func (m *Manager) LoadAll() error {
	// Start the local extension server
	server := NewServer()
	url, err := server.Start()
	if err != nil {
		fmt.Printf("[Extensions] Warning: Failed to start UI server: %v\n", err)
	}
	m.serverURL = url
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

			// Determine trust level: Check against gallery API mock
			trust := "local"
			for _, ge := range GetGalleryExtensions() {
				if ge.ID == manifest.ID {
					trust = ge.Trust
					break
				}
			}

			iconUrl := ""
			if manifest.Icon != "" {
				iconUrl = fmt.Sprintf("%s/%s/%s", m.serverURL, manifest.ID, manifest.Icon)
			}

			ext := Extension{
				ID:          manifest.ID,
				Name:        manifest.Name,
				Version:     manifest.Version,
				Author:      manifest.Author,
				Description: manifest.Description,
				Status:      "Running",
				Memory:      "0MB",
				CPU:         "0%",
				Trust:       trust,
				IconURL:     iconUrl,
			}

			sandbox := NewSandbox(
				m.ctx, manifest, m.serverURL,
				func(payload map[string]interface{}) {
					m.SidebarPages = append(m.SidebarPages, payload)
				},
				func(config ModLoaderConfig) {
					m.ModLoaders[config.ID] = config
				},
				func() []InstanceInfo {
					all := instance.GetInstances()
					var out []InstanceInfo
					for _, inst := range all {
						out = append(out, InstanceInfo{
							ID:      inst.ID,
							Name:    inst.Name,
							Version: inst.Version,
							Loader:  inst.Loader,
						})
					}
					return out
				},
				func(instanceID, jarName, downloadURL string) (string, error) {
					jarName = filepath.Base(jarName)
					modsDir := filepath.Join(fs.GetDataDir(), "instances", instanceID, "mods")
					if err := os.MkdirAll(modsDir, 0755); err != nil {
						return "", err
					}
					destPath := filepath.Join(modsDir, jarName)
					resp, err := http.Get(downloadURL)
					if err != nil {
						return "", err
					}
					defer resp.Body.Close()
					out, err := os.Create(destPath)
					if err != nil {
						return "", err
					}
					defer out.Close()
					if _, err = io.Copy(out, resp.Body); err != nil {
						return "", err
					}
					return destPath, nil
				},
			)
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

// HandleIPCMessage routes a message from the UI iframe to the extension sandbox
func (m *Manager) HandleIPCMessage(extID string, payload map[string]interface{}) {
	sandbox, ok := m.sandboxes[extID]
	if !ok {
		fmt.Printf("[Manager] IPC message for unknown extension: %s\n", extID)
		return
	}

	// If the sandbox has a registered listener, call it
	if sandbox.onMessageCallback != nil {
		result, err := sandbox.onMessageCallback(payload)
		if err != nil {
			fmt.Printf("[Manager] IPC callback error for %s: %v\n", extID, err)
			return
		}
		// Emit the response back to the frontend as a Wails event
		runtime.EventsEmit(m.ctx, "extension:message:"+extID, result)
	}
}

// GetSidebarPages returns all registered sidebar pages
func (m *Manager) GetSidebarPages() []map[string]interface{} {
	return m.SidebarPages
}
