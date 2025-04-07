//go:build darwin
// +build darwin

package common

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// GetCreationTime retrieves the creation time of a file or directory on Darwin.
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

// GetOwnerAndGroup retrieves the owner UID and group GID of a file or directory on Darwin.
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
