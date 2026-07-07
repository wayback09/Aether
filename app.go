package main

import (
	"context"
	"fmt"

	"Aether/pkg/extensions"
	"Aether/pkg/fs"
	"Aether/pkg/instance"
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
	extensions.GlobalManager = extensions.NewManager()
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

// GetExtensions returns all installed extensions
func (a *App) GetExtensions() []extensions.Extension {
	return extensions.GetExtensions()
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
