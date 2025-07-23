package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindGoPackages finds all directories containing Go source files.
func FindGoPackages(rootDir string, excludeDirs []string) ([]string, error) {
	var packages []string

	excludeMap := make(map[string]bool)
	for _, dir := range excludeDirs {
		excludeMap[dir] = true
	}

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Skip excluded directories
		if shouldSkipDir(path, excludeMap) {
			return filepath.SkipDir
		}

		// Check if directory contains Go files
		hasGoFiles, err := hasGoSourceFiles(path)
		if err != nil {
			return fmt.Errorf("check Go files in %q: %w", path, err)
		}

		if hasGoFiles {
			packages = append(packages, path)
		}

		return nil
	})

	return packages, err
}

// hasGoSourceFiles checks if a directory contains Go source files (excluding test files).
func hasGoSourceFiles(dir string) (bool, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			return true, nil
		}
	}

	return false, nil
}

// shouldSkipDir determines if a directory should be skipped during traversal.
func shouldSkipDir(path string, excludeMap map[string]bool) bool {
	base := filepath.Base(path)

	// Check exact matches
	if excludeMap[base] {
		return true
	}

	// Check common patterns
	if strings.HasPrefix(base, ".") && base != "." {
		return true
	}

	if strings.HasSuffix(base, "_test") {
		return true
	}

	return false
}
