package java

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FindJava locates a working Java binary on the system.
// It checks JAVA_HOME first, then falls back to PATH.
func FindJava() (string, error) {
	// 1. Check JAVA_HOME
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome != "" {
		candidate := filepath.Join(javaHome, "bin", "java.exe")
		if _, err := os.Stat(candidate); err == nil {
			if validateJava(candidate) {
				return candidate, nil
			}
		}
		// Also try without .exe for cross-platform
		candidate = filepath.Join(javaHome, "bin", "java")
		if _, err := os.Stat(candidate); err == nil {
			if validateJava(candidate) {
				return candidate, nil
			}
		}
	}

	// 2. Check PATH via exec.LookPath
	pathJava, err := exec.LookPath("java")
	if err == nil && validateJava(pathJava) {
		return pathJava, nil
	}

	// 3. Check common Windows install locations
	programFiles := os.Getenv("ProgramFiles")
	if programFiles != "" {
		candidates := findJavaInDir(programFiles)
		for _, c := range candidates {
			if validateJava(c) {
				return c, nil
			}
		}
	}

	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	if programFilesX86 != "" {
		candidates := findJavaInDir(programFilesX86)
		for _, c := range candidates {
			if validateJava(c) {
				return c, nil
			}
		}
	}

	return "", fmt.Errorf("java not found: please install Java or set JAVA_HOME")
}

// GetVersion returns the major version of a Java binary (e.g. 21 for Java 21)
func GetVersion(javaPath string) (string, error) {
	cmd := exec.Command(javaPath, "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	// Parse output like: openjdk version "21.0.1" or java version "1.8.0_301"
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0]), nil
	}
	return "unknown", nil
}

// validateJava checks if a java binary actually runs
func validateJava(path string) bool {
	cmd := exec.Command(path, "-version")
	err := cmd.Run()
	return err == nil
}

// findJavaInDir searches for java.exe in common Java installation subdirectories
func findJavaInDir(baseDir string) []string {
	var results []string

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return results
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.Contains(name, "java") || strings.Contains(name, "jdk") || strings.Contains(name, "jre") || strings.Contains(name, "adoptium") || strings.Contains(name, "zulu") || strings.Contains(name, "graalvm") {
			candidate := filepath.Join(baseDir, entry.Name(), "bin", "java.exe")
			if _, err := os.Stat(candidate); err == nil {
				results = append(results, candidate)
			}
		}
	}

	return results
}
