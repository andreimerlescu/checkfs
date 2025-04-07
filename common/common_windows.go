//go:build windows
// +build windows

package common

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// GetCreationTime retrieves the creation time of a file or directory on Windows
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

// GetOwnerAndGroup is not supported on Windows
func GetOwnerAndGroup(path string) (uid, gid string, err error) {
	return "", "", fmt.Errorf("owner and group checks are not supported on Windows: %s", path)
}
