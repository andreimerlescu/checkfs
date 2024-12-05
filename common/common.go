package common

import (
	"fmt"
	"path/filepath"
	"strings"
)

// IsPathInBase checks if a path is within the base directory
func IsPathInBase(path, baseDir string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path cannot be empty")
	}
	if baseDir == "" {
		return false, fmt.Errorf("base directory cannot be empty")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of %s: %w", path, err)
	}
	absBaseDir, err := filepath.Abs(baseDir)
	if err != nil {
		return false, fmt.Errorf("failed to get absolute path of base directory %s: %w", baseDir, err)
	}
	rel, err := filepath.Rel(absBaseDir, absPath)
	if err != nil {
		return false, fmt.Errorf("failed to get relative path: %w", err)
	}
	return !RelStartsWithParent(rel), nil
}

// RelStartsWithParent checks if a relative path escapes the base directory
func RelStartsWithParent(rel string) bool {
	// Normalize the path for consistent comparison
	rel = filepath.Clean(rel)
	return strings.HasPrefix(rel, "..") && (len(rel) == 2 || strings.HasPrefix(rel[2:], string(filepath.Separator)))
}
