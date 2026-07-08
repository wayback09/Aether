package extensions

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Sandbox represents an isolated JavaScript environment for an extension
type Sandbox struct {
	ctx      context.Context
	vm       *goja.Runtime
	manifest Manifest
}

// NewSandbox creates a new Goja isolate restricted by the given manifest
func NewSandbox(ctx context.Context, manifest Manifest, serverURL string, onSidebarPage func(map[string]interface{})) *Sandbox {
	vm := goja.New()
	
	// Create the secure Aether bridge object
	aetherObj := vm.NewObject()
	
	// Capability: ui:sidebar
	if manifest.HasPermission("ui:sidebar") {
		uiObj := vm.NewObject()
		uiObj.Set("registerSidebarPage", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				arg := call.Argument(0).Export().(map[string]interface{})
				
				// Reconstruct the URL using the local server
				// E.g., if arg["url"] is "ui/index.html", we prepend "http://127.0.0.1:port/modrinth/ui/index.html"
				relURL := arg["url"].(string)
				fullURL := fmt.Sprintf("%s/%s/%s", serverURL, manifest.ID, relURL)
				
				payload := map[string]interface{}{
					"extensionId": manifest.ID,
					"id":          arg["id"],
					"label":       arg["label"],
					"url":         fullURL,
				}
				
				fmt.Printf("[Sandbox:%s] Registered sidebar page: %+v\n", manifest.ID, payload)
				
				if onSidebarPage != nil {
					onSidebarPage(payload)
				}
				
				runtime.EventsEmit(ctx, "extension:sidebar:add", payload)
			}
			return goja.Undefined()
		})
		aetherObj.Set("ui", uiObj)
	}

	// Capability: ui:dialogs
	if manifest.HasPermission("ui:dialogs") {
		uiObj := aetherObj.Get("ui")
		if uiObj == nil {
			uiObj = vm.NewObject()
			aetherObj.Set("ui", uiObj)
		}
		uiObj.(*goja.Object).Set("openDialog", func(call goja.FunctionCall) goja.Value {
			fmt.Printf("[Sandbox:%s] Opened dialog\n", manifest.ID)
			return goja.Undefined()
		})
	}
	
	// Inject the bridge into the global scope
	vm.Set("Aether", aetherObj)
	
	return &Sandbox{
		ctx:      ctx,
		vm:       vm,
		manifest: manifest,
	}
}

// Execute runs a JS script inside the sandbox
func (s *Sandbox) Execute(script string) error {
	_, err := s.vm.RunString(script)
	if err != nil {
		return fmt.Errorf("sandbox execution error: %w", err)
	}
	return nil
}
