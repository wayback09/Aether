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
	"Aether/pkg/java"
	"Aether/pkg/mojang"
)

type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// to call runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fs.EnsureDirectories()
	
	// Initialize and load all extensions into their isolates
	extensions.GlobalManager = extensions.NewManager(ctx)
	extensions.GlobalManager.LoadAll()

	// Wire the mod loader hook so launcher.go can call extension mod loaders
	// without an import cycle (instance → extensions → instance)
	instance.ModLoaderHook = func(loaderID string, hookCtx map[string]interface{}) (map[string]interface{}, error) {
		if loader, ok := extensions.GlobalManager.ModLoaders[loaderID]; ok {
			return loader.Callback(hookCtx)
		}
		return nil, fmt.Errorf("mod loader '%s' is not installed", loaderID)
	}
}

// GetInstances returns all installed instances
func (a *App) GetInstances() []instance.Instance {
	return instance.GetInstances()
}

// GetActiveInstance returns the currently selected instance
func (a *App) GetActiveInstance() *instance.Instance {
	return instance.GetActiveInstance()
}

func (a *App) GetExtensions() []extensions.Extension {
	return extensions.GetExtensions()
}

// GetExtensionSidebarPages returns sidebar pages contributed by extensions
func (a *App) GetExtensionSidebarPages() []map[string]interface{} {
	if extensions.GlobalManager != nil {
		return extensions.GlobalManager.GetSidebarPages()
	}
	return []map[string]interface{}{}
}

// SendExtensionMessage routes an IPC message from the UI iframe to the extension sandbox
func (a *App) SendExtensionMessage(extID string, payload map[string]interface{}) {
	if extensions.GlobalManager != nil {
		extensions.GlobalManager.HandleIPCMessage(extID, payload)
	}
}

// ModLoaderInfo represents a registered mod loader
type ModLoaderInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetModLoaders returns all mod loaders registered by extensions
func (a *App) GetModLoaders() []ModLoaderInfo {
	var loaders []ModLoaderInfo
	if extensions.GlobalManager != nil {
		for _, loader := range extensions.GlobalManager.ModLoaders {
			loaders = append(loaders, ModLoaderInfo{
				ID:          loader.ID,
				Name:        loader.Name,
				Description: loader.Description,
			})
		}
	}
	return loaders
}

// JavaRuntimeStatus describes the status of a managed Java runtime.
type JavaRuntimeStatus struct {
	Version   int    `json:"version"`
	Installed bool   `json:"installed"`
	Path      string `json:"path"`
}

// GetJavaStatus returns the installation status for each required Java version.
func (a *App) GetJavaStatus() []JavaRuntimeStatus {
	versions := []int{8, 17, 21}
	var statuses []JavaRuntimeStatus
	for _, v := range versions {
		installed := java.IsManagedJavaInstalled(v)
		path := ""
		if installed {
			path = java.GetManagedJavaPath(v)
		}
		statuses = append(statuses, JavaRuntimeStatus{
			Version:   v,
			Installed: installed,
			Path:      path,
		})
	}
	return statuses
}

// DownloadJavaRuntime downloads a managed JRE for the given major version.
func (a *App) DownloadJavaRuntime(version int) error {
	return java.DownloadJava(a.ctx, version)
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

// DownloadAndInstallExtension downloads a remote zip and installs it
func (a *App) DownloadAndInstallExtension(url string) (bool, error) {
	if err := extensions.DownloadAndInstallExtension(url); err != nil {
		return false, err
	}

	if extensions.GlobalManager != nil {
		extensions.GlobalManager.LoadAll()
	}

	return true, nil
}


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

// RemoveAccount removes an account by ID
func (a *App) RemoveAccount(id string) error {
	return auth.RemoveAccount(id)
}

// LoginWithMicrosoft starts the OAuth2 flow and returns immediately.
func (a *App) LoginWithMicrosoft() error {
	go func() {
		var browserPort int
		
		// Start local callback server for OAuth loop.
		// Blocks until completion, but triggers callback once listening.
		code, err := auth.StartCallbackServer(func(port int) {
			browserPort = port
			url := auth.GetMicrosoftAuthURL(port)
			runtime.BrowserOpenURL(a.ctx, url)
		})
		
		if err != nil {
			runtime.EventsEmit(a.ctx, "auth:error", err.Error())
			return
		}

		// Exchange authorization code for token chain.
		acc, err := auth.LoginWithMicrosoft(code, browserPort)
		if err != nil {
			runtime.EventsEmit(a.ctx, "auth:error", err.Error())
			return
		}

		// Save to keyring and set as active profile.
		if err := auth.AddMicrosoftAccount(acc); err != nil {
			runtime.EventsEmit(a.ctx, "auth:error", err.Error())
			return
		}

		// Notify frontend of successful authentication.
		runtime.EventsEmit(a.ctx, "auth:complete")
	}()
	return nil
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

// CreateInstance creates a new instance on disk and returns the created instance
func (a *App) CreateInstance(name, version, loader string) (*instance.Instance, error) {
	inst, err := instance.Create(name, version, loader)
	if err != nil {
		return nil, err
	}
	return inst, nil
}

// UpdateInstance saves changes to an instance
func (a *App) UpdateInstance(inst *instance.Instance) error {
	return instance.UpdateInstance(inst)
}

// DeleteInstance deletes an instance completely
func (a *App) DeleteInstance(id string) error {
	return instance.DeleteInstance(id)
}
