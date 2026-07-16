package extensions

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"Aether/pkg/fs"
	"Aether/pkg/netutil"
	"github.com/dop251/goja"
)

// Sandbox represents an isolated JavaScript environment for an extension
type Sandbox struct {
	ctx                context.Context
	vm                 *goja.Runtime
	manifest           Manifest
	onMessageCallback  func(map[string]interface{}) (map[string]interface{}, error)
}

// InstanceInfo is a minimal view of an instance passed into the sandbox
type InstanceInfo struct {
	ID      string
	Name    string
	Version string
	Loader  string
}

// NewSandbox creates a new Goja isolate restricted by the given manifest.
// emit is an optional function for broadcasting events (e.g. runtime.EventsEmit);
// pass nil to disable event broadcasting (useful in tests).
func NewSandbox(
	ctx context.Context,
	manifest Manifest,
	serverURL string,
	onSidebarPage func(map[string]interface{}),
	onModLoader func(ModLoaderConfig),
	listInstances func() []InstanceInfo,
	installMod func(instanceID, jarName, downloadURL string) (string, error),
	emit func(ctx context.Context, event string, data ...interface{}),
) *Sandbox {
	if emit == nil {
		emit = func(_ context.Context, _ string, _ ...interface{}) {}
	}
	vm := goja.New()
	
	// Create the secure Aether bridge object
	aetherObj := vm.NewObject()
	
	// URL Whitelist Helper
	isAllowedURL := func(target string) bool {
		// Example implementation (simplified)
		for _, host := range manifest.Hosts {
			if strings.Contains(target, host) {
				return true
			}
		}
		return false
	}

	// Capability: ui:sidebar
	if manifest.HasPermission("ui:sidebar") {
		uiObj := vm.NewObject()
		uiObj.Set("registerSidebarPage", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				arg := call.Argument(0).Export().(map[string]interface{})
				
				// Reconstruct the URL using the local server
				relURL := arg["url"].(string)
				fullURL := fmt.Sprintf("%s/%s/%s", serverURL, manifest.ID, relURL)
				
				payload := map[string]interface{}{
					"extensionId": manifest.ID,
					"id":          arg["id"],
					"label":       arg["label"],
					"url":         fullURL,
				}
				
				if onSidebarPage != nil {
					onSidebarPage(payload)
				}
				
				emit(ctx, "extension:sidebar:add", payload)
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
	
	// Capability: network:http
	if manifest.HasPermission("network:http") {
		httpObj := vm.NewObject()
		httpObj.Set("get", func(call goja.FunctionCall) goja.Value {
			targetURL := call.Argument(0).String()
			if !isAllowedURL(targetURL) {
				panic(vm.NewGoError(fmt.Errorf("access denied to URL: %s", targetURL)))
			}
			
			resp, err := http.Get(targetURL)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			defer resp.Body.Close()
			
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(string(body))
		})
		aetherObj.Set("http", httpObj)
	}
	
	// Capability: fs:download
	if manifest.HasPermission("fs:download") {
		fsObj := vm.NewObject()
		fsObj.Set("download", func(call goja.FunctionCall) goja.Value {
			targetURL := call.Argument(0).String()
			destPath := call.Argument(1).String()
			
			if !isAllowedURL(targetURL) {
				panic(vm.NewGoError(fmt.Errorf("access denied to URL: %s", targetURL)))
			}
			
			// Force destPath to be within the instances/libraries folder
			safePath := filepath.Join(fs.GetDataDir(), "libraries", filepath.Clean(destPath))
			
			if err := netutil.DownloadFile(ctx, targetURL, safePath, nil); err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(safePath)
		})
		aetherObj.Set("fs", fsObj)
	}
	
	// Capability: instances:patch
	if manifest.HasPermission("instances:patch") {
		instancesObj := vm.NewObject()

		instancesObj.Set("list", func(call goja.FunctionCall) goja.Value {
			if listInstances == nil {
				return vm.ToValue([]interface{}{})
			}
			all := listInstances()
			var result []map[string]interface{}
			for _, inst := range all {
				result = append(result, map[string]interface{}{
					"id":      inst.ID,
					"name":    inst.Name,
					"version": inst.Version,
					"loader":  inst.Loader,
				})
			}
			return vm.ToValue(result)
		})

		instancesObj.Set("installMod", func(call goja.FunctionCall) goja.Value {
			instanceID := call.Argument(0).String()
			jarName := call.Argument(1).String()
			downloadURL := call.Argument(2).String()

			if !isAllowedURL(downloadURL) {
				panic(vm.NewGoError(fmt.Errorf("access denied to URL: %s", downloadURL)))
			}

			if installMod == nil {
				panic(vm.NewGoError(fmt.Errorf("installMod not available")))
			}

			path, err := installMod(instanceID, jarName, downloadURL)
			if err != nil {
				panic(vm.NewGoError(err))
			}
			return vm.ToValue(path)
		})

		aetherObj.Set("instances", instancesObj)
	}

	// Capability: launcher:modloader
	if manifest.HasPermission("launcher:modloader") {
		launcherObj := vm.NewObject()
		launcherObj.Set("registerModLoader", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				arg := call.Argument(0).Export().(map[string]interface{})
				
				config := ModLoaderConfig{
					ID:          arg["id"].(string),
					Name:        arg["name"].(string),
					Description: arg["description"].(string),
					ExtensionID: manifest.ID,
				}
				
				// Extract the JS callback function safely
				if cb, ok := goja.AssertFunction(call.Argument(0).ToObject(vm).Get("onLaunch")); ok {
					config.Callback = func(ctx map[string]interface{}) (map[string]interface{}, error) {
						val, err := cb(goja.Undefined(), vm.ToValue(ctx))
						if err != nil {
							return nil, err
						}
						return val.Export().(map[string]interface{}), nil
					}
				}
				
				if onModLoader != nil {
					onModLoader(config)
				}
				fmt.Printf("[Sandbox:%s] Registered mod loader: %s\n", manifest.ID, config.ID)
			}
			return goja.Undefined()
		})
		aetherObj.Set("launcher", launcherObj)
	}

	// Capability: skin:export
	if manifest.HasPermission("skin:export") {
		skinsObj := vm.NewObject()
		skinsObj.Set("export", func(call goja.FunctionCall) goja.Value {
			b64Data := call.Argument(0).String()
			filename := call.Argument(1).String()
			if filename == "" {
				filename = "skin.png"
			}

			skinsDir := filepath.Join(fs.GetDataDir(), "skins")
			if err := os.MkdirAll(skinsDir, 0755); err != nil {
				panic(vm.NewGoError(err))
			}

			safePath := filepath.Join(skinsDir, filepath.Clean(filename))
			data, err := base64.StdEncoding.DecodeString(b64Data)
			if err != nil {
				panic(vm.NewGoError(fmt.Errorf("invalid base64 data: %w", err)))
			}

			if err := os.WriteFile(safePath, data, 0644); err != nil {
				panic(vm.NewGoError(err))
			}

			return vm.ToValue(safePath)
		})
		aetherObj.Set("skins", skinsObj)
	}

	// Capability: ui:sidebar - also inject IPC methods
	if manifest.HasPermission("ui:sidebar") {
		// Grab the existing uiObj that was created earlier (or create if somehow nil)
		uiIPC := vm.NewObject()
		
		// Sandbox not yet created. Store callback via pointer closure to allow
		// execution later.
		var jsMessageHandler goja.Callable

		uiIPC.Set("onMessage", func(call goja.FunctionCall) goja.Value {
			if fn, ok := goja.AssertFunction(call.Argument(0)); ok {
				jsMessageHandler = fn
			}
			return goja.Undefined()
		})

		uiIPC.Set("postMessage", func(call goja.FunctionCall) goja.Value {
			if len(call.Arguments) > 0 {
				emit(ctx, "extension:message:"+manifest.ID, call.Argument(0).Export())
			}
			return goja.Undefined()
		})

		// Merge into existing uiObj (which already has registerSidebarPage)
		if existingUI := aetherObj.Get("ui"); existingUI != nil {
			existingUI.(*goja.Object).Set("onMessage", uiIPC.Get("onMessage"))
			existingUI.(*goja.Object).Set("postMessage", uiIPC.Get("postMessage"))
		} else {
			aetherObj.Set("ui", uiIPC)
		}

		// Build sandbox early to assign the deferred callback closure.
		_ = jsMessageHandler // Captured in callback below

		// Inject the bridge into the global scope
		vm.Set("Aether", aetherObj)

		sb := &Sandbox{
			ctx:      ctx,
			vm:       vm,
			manifest: manifest,
		}
		sb.onMessageCallback = func(payload map[string]interface{}) (map[string]interface{}, error) {
			if jsMessageHandler == nil {
				return nil, fmt.Errorf("no onMessage handler registered")
			}
			val, err := jsMessageHandler(goja.Undefined(), vm.ToValue(payload))
			if err != nil {
				return nil, err
			}
			if exported := val.Export(); exported != nil {
				if m, ok := exported.(map[string]interface{}); ok {
					return m, nil
				}
			}
			return map[string]interface{}{}, nil
		}
		return sb
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
