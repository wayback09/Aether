package instance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"Aether/pkg/fs"
)

type Instance struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Version    string `json:"version"`
	Loader     string `json:"loader"`
	Memory     string `json:"memory"`
	LastPlayed string `json:"lastPlayed"`
	Installed  bool   `json:"installed"`
}

// GetInstances returns a list of instances parsed from the disk
func GetInstances() []Instance {
	instancesDir := filepath.Join(fs.GetDataDir(), "instances")
	instances := []Instance{} // Initialize as empty slice, not nil

	entries, err := os.ReadDir(instancesDir)
	if err != nil {
		return instances
	}

	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		wg.Add(1)
		go func(entry os.DirEntry) {
			defer wg.Done()

			manifestPath := filepath.Join(instancesDir, entry.Name(), "instance.json")
			data, err := os.ReadFile(manifestPath)
			if err != nil {
				return
			}

			var inst Instance
			if err := json.Unmarshal(data, &inst); err != nil {
				return
			}

			if inst.ID == "" {
				inst.ID = entry.Name()
			}
			// Check if client.jar exists
			jarPath := filepath.Join(instancesDir, entry.Name(), "bin", inst.Version+".jar")
			if _, err := os.Stat(jarPath); err == nil {
				inst.Installed = true
			}

			mu.Lock()
			instances = append(instances, inst)
			mu.Unlock()
		}(entry)
	}

	wg.Wait()
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

// UpdateInstance saves the modified instance data to disk
func UpdateInstance(inst *Instance) error {
	manifestPath := filepath.Join(fs.GetDataDir(), "instances", inst.ID, "instance.json")
	data, err := json.MarshalIndent(inst, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(manifestPath, data, 0644)
}

// DeleteInstance permanently removes an instance directory from disk
func DeleteInstance(id string) error {
	instanceDir := filepath.Join(fs.GetDataDir(), "instances", id)
	return os.RemoveAll(instanceDir)
}
