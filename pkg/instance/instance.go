package instance

import (
	"encoding/json"
	"os"
	"path/filepath"

	"Aether/pkg/fs"
)

type Instance struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Loader     string `json:"loader"`
	Memory     string `json:"memory"`
	LastPlayed string `json:"lastPlayed"`
}

// GetInstances returns a list of instances parsed from the disk
func GetInstances() []Instance {
	instancesDir := filepath.Join(fs.GetDataDir(), "instances")
	var instances []Instance

	entries, err := os.ReadDir(instancesDir)
	if err != nil {
		return []Instance{}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			manifestPath := filepath.Join(instancesDir, entry.Name(), "instance.json")
			data, err := os.ReadFile(manifestPath)
			if err == nil {
				var inst Instance
				if err := json.Unmarshal(data, &inst); err == nil {
					if inst.ID == "" {
						inst.ID = entry.Name()
					}
					instances = append(instances, inst)
				}
			}
		}
	}
	return instances
}

// GetActiveInstance returns the first loaded instance
func GetActiveInstance() *Instance {
	instances := GetInstances()
	if len(instances) > 0 {
		return &instances[0]
	}
	return nil
}
