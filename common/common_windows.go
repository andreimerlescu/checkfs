//go:build windows

package common

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// HasPermissions checks if a file or directory has at least the specified permissions
// On Windows, this is simplified due to NTFS permissions not mapping directly to Unix modes
func HasPermissions(path string, perms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	actualPerms := info.Mode().Perm()
	// Windows often reports broader perms; focus on read/write bits
	return actualPerms&perms&0666 != 0, nil // Ignore execute bits, focus on read/write
}

// IsMorePermissiveThan checks if a file or directory’s permissions are at least as permissive as the given mode
// Adjusted for Windows behavior where strict Unix perms aren't enforced
func IsMorePermissiveThan(path string, minPerms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	perms := info.Mode().Perm()
	// On Windows, assume read/write perms are broader; mask to relevant bits
	return perms&0444 >= minPerms&0444, nil // Focus on read bits as a minimum
}

func GetOwnerAndGroup(path string) (uid, gid string, err error) {
	return "", "", fmt.Errorf("owner and group checks are not supported on Windows: %s", path)
}

func GetCreationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	if stat, ok := info.Sys().(*syscall.Win32FileAttributeData); ok {
		return time.Unix(0, stat.CreationTime.Nanoseconds()), nil
	}
	return time.Time{}, fmt.Errorf("unable to get creation time for %s on Windows", path)
}

// IsLessPermissiveThan checks if a file or directory’s permissions are no more permissive than the given mode
// Adjusted for Windows behavior
func IsLessPermissiveThan(path string, maxPerms os.FileMode) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	perms := info.Mode().Perm()
	// Windows perms are often broader; check if within maxPerms bounds
	return perms&0666 <= maxPerms&0666, nil // Focus on read/write bits
}
