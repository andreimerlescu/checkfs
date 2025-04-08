//go:build darwin

package common

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

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

// GetOwnerAndGroup retrieves the owner UID and group GID of a file or directory on Darwin
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

// GetCreationTime retrieves the creation time of a file or directory on Darwin
func GetCreationTime(path string) (time.Time, error) {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to stat %s: %w", path, err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}, fmt.Errorf("unable to get detailed stats for %s", path)
	}
	return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec), nil
}
