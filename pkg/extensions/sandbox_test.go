package extensions

import (
	"context"
	"testing"
)

func TestSandboxCapabilities(t *testing.T) {
	// 1. Create a manifest WITH sidebar permission
	manifestAllowed := Manifest{
		ID:          "com.test.allowed",
		Permissions: []string{"ui:sidebar"},
	}
	
	sandbox1 := NewSandbox(
		context.Background(),
		manifestAllowed,
		"http://localhost",
		nil,
		nil,
		nil,
		nil,
		nil, // emit: nil disables Wails event broadcasting in tests
	)
	
	// This should succeed without panic
	script1 := `
		Aether.ui.registerSidebarPage({ id: "test", label: "Test Page", url: "ui/index.html" });
	`
	err := sandbox1.Execute(script1)
	if err != nil {
		t.Fatalf("Expected script to run successfully, got: %v", err)
	}
	
	// 2. Create a manifest WITHOUT any permissions
	manifestDenied := Manifest{
		ID:          "com.test.denied",
		Permissions: []string{},
	}
	
	sandbox2 := NewSandbox(
		context.Background(),
		manifestDenied,
		"http://localhost",
		nil,
		nil,
		nil,
		nil,
		nil, // emit: nil disables Wails event broadcasting in tests
	)
	
	// This should throw a JS error because Aether.ui is undefined
	script2 := `
		Aether.ui.registerSidebarPage({ id: "test", label: "Test Page", url: "ui/index.html" });
	`
	err = sandbox2.Execute(script2)
	if err == nil {
		t.Fatalf("Expected script to fail due to missing capabilities, but it succeeded")
	}
}
