package extensions

// Package extensions handles the Aether extension lifecycle


// Extension represents the UI representation of a Manifest
type Extension struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Author  string `json:"author"`
	Status  string `json:"status"`
	Memory  string `json:"memory"`
	CPU     string `json:"cpu"`
}

// GetExtensions returns the list of loaded extensions from the GlobalManager
func GetExtensions() []Extension {
	if GlobalManager != nil {
		return GlobalManager.GetExtensions()
	}
	return []Extension{}
}
