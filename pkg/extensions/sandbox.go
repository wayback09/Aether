package extensions

import (
	"fmt"
	"github.com/dop251/goja"
)

// Sandbox represents an isolated JavaScript environment for an extension
type Sandbox struct {
	vm       *goja.Runtime
	manifest Manifest
}

// NewSandbox creates a new V8/Goja isolate restricted by the given manifest
func NewSandbox(manifest Manifest) *Sandbox {
	vm := goja.New()
	
	// Create the secure Aether bridge object
	aetherObj := vm.NewObject()
	
	// Capability: ui:sidebar
	if manifest.HasPermission("ui:sidebar") {
		uiObj := vm.NewObject()
		uiObj.Set("registerSidebarPage", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				arg := call.Argument(0).Export()
				fmt.Printf("[Sandbox:%s] Registered sidebar page: %+v\n", manifest.ID, arg)
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
