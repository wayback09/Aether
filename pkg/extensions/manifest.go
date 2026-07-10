package extensions

// Manifest represents the structure of an extension's manifest.json
type Manifest struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Icon        string   `json:"icon,omitempty"`
	Main        string   `json:"main"`
	API         string   `json:"api"`
	Permissions []string `json:"permissions"`
	Hosts       []string `json:"hosts,omitempty"`
}

// HasPermission checks if the extension has requested a specific capability
func (m *Manifest) HasPermission(perm string) bool {
	for _, p := range m.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}
