package main

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"Aether/pkg/auth"
	"Aether/pkg/extensions"
	"Aether/pkg/fs"
	"Aether/pkg/instance"
	"Aether/pkg/mojang"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fs.EnsureDirectories()
	
	// Initialize and load all extensions into their isolates
	extensions.GlobalManager = extensions.NewManager(ctx)
	extensions.GlobalManager.LoadAll()
}

// GetInstances returns all installed instances
func (a *App) GetInstances() []instance.Instance {
	return instance.GetInstances()
}

// GetActiveInstance returns the currently selected instance
func (a *App) GetActiveInstance() *instance.Instance {
	return instance.GetActiveInstance()
}

// --- Extension Methods ---

func (a *App) GetExtensions() []extensions.Extension {
	return extensions.GetExtensions()
}

func (a *App) GetExtensionSidebarPages() []map[string]interface{} {
	if extensions.GlobalManager != nil {
		return extensions.GlobalManager.GetSidebarPages()
	}
	return []map[string]interface{}{}
}

func (a *App) SelectAndInstallExtension() (bool, error) {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Extension Zip",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Aether Extensions (*.zip)",
				Pattern:     "*.zip",
			},
		},
	})
	if err != nil {
		return false, err
	}
	if file == "" {
		// User cancelled
		return false, nil
	}

	if err := extensions.InstallFromZip(file); err != nil {
		return false, err
	}

	// Reload all extensions dynamically!
	if extensions.GlobalManager != nil {
		extensions.GlobalManager.LoadAll()
	}

	return true, nil
}


// --- Auth Methods ---

// GetActiveAccount returns the currently active account
func (a *App) GetActiveAccount() *auth.Account {
	return auth.GetActiveAccount()
}

// GetAccounts returns all saved accounts
func (a *App) GetAccounts() []auth.Account {
	return auth.GetAccounts()
}

// LoginOffline creates or switches to an offline account with the given username
func (a *App) LoginOffline(username string) (auth.Account, error) {
	return auth.AddOfflineAccount(username)
}

// SetActiveAccount sets the active account by ID
func (a *App) SetActiveAccount(id string) error {
	return auth.SetActiveAccount(id)
}

// LaunchInstance starts the specified instance
func (a *App) LaunchInstance(id string) error {
	instances := instance.GetInstances()
	var target *instance.Instance
	for i := range instances {
		if instances[i].ID == id {
			target = &instances[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("instance not found: %s", id)
	}
	return instance.Launch(a.ctx, target)
}

// InstallInstance triggers the Mojang download pipeline
func (a *App) InstallInstance(id string) error {
	instances := instance.GetInstances()
	var target *instance.Instance
	for i := range instances {
		if instances[i].ID == id {
			target = &instances[i]
			break
		}
	}
	if target == nil {
		return fmt.Errorf("instance not found: %s", id)
	}

	info, err := mojang.GetVersionInfo(target.Version)
	if err != nil {
		return fmt.Errorf("failed to get version info: %w", err)
	}

	basePath := filepath.Join(fs.GetDataDir(), "instances", target.ID)
	assetsDir := fs.GetAssetsDir()
	engine := mojang.NewDownloadEngine(a.ctx, target.ID, basePath)

	go func() {
		if err := engine.Install(info, assetsDir); err != nil {
			fmt.Printf("Install error: %v\n", err)
			runtime.EventsEmit(a.ctx, "instance:state", map[string]interface{}{
				"id":    target.ID,
				"state": "Error",
			})
		} else {
			runtime.EventsEmit(a.ctx, "instance:state", map[string]interface{}{
				"id":    target.ID,
				"state": "Idle",
			})
		}
	}()

	return nil
}

// GetAvailableVersions fetches stable releases from Mojang
func (a *App) GetAvailableVersions() ([]string, error) {
	manifest, err := mojang.GetVersionManifest()
	if err != nil {
		return nil, err
	}

	var releases []string
	for _, v := range manifest.Versions {
		if v.Type == "release" {
			releases = append(releases, v.ID)
		}
	}
	return releases, nil
}

// CreateInstance creates a new instance on disk
func (a *App) CreateInstance(name, version, loader string) error {
	_, err := instance.Create(name, version, loader)
	return err
}
