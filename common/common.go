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

// GetCreationTime retrieves the creation time of a file or directory
func GetCreationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, fmt.Errorf("unable to get detailed stats for %s", path)
	}
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec), nil
}

// HasPermissions checks if a file or directory has at least the specified permissions
func HasPermissions(path string, perms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	actualPerms := info.Mode().Perm()
	return actualPerms&perms == perms, nil
}

// IsMorePermissiveThan checks if a file or directory’s permissions are at least as permissive as the given mode
func IsMorePermissiveThan(path string, minPerms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	perms := info.Mode().Perm()
	return perms&minPerms == minPerms, nil
}

// IsLessPermissiveThan checks if a file or directory’s permissions are no more permissive than the given mode
func IsLessPermissiveThan(path string, maxPerms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	perms := info.Mode().Perm()
	return perms&^maxPerms == 0, nil
}

// GetOwnerAndGroup retrieves the owner UID and group GID of a file or directory
func GetOwnerAndGroup(path string) (uid, gid string, err error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to stat %s: %w", path, err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return "", "", fmt.Errorf("unable to get detailed stats for %s", path)
	}
	return fmt.Sprint(stat.Uid), fmt.Sprint(stat.Gid), nil
}

// SanitizePath removes redundant separators and resolves relative components in a path
func SanitizePath(path string) (string, error) {
	cleaned := filepath.Clean(path)
	if cleaned == "" {
		return "", fmt.Errorf("path cannot be empty after cleaning")
	}
	return cleaned, nil
}
