package file

import (
	"fmt"
	"github.com/andreimerlescu/go-checkfs/common"
	"os"
	"syscall"
)

type Options struct {
	ReadOnly       bool   // Check if the file is read-only
	RequireWrite   bool   // Check if the file is writable
	RequireOwner   string // Check if the file has a specific owner
	RequireGroup   string // Check if the file has a specific group
	RequireBaseDir string // Check if the file is inside a specific base directory
}

// File performs the file checks
func File(path string, opts Options) error {
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// Check if file is a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", path)
	}

	// Check if file is inside the required base directory
	if opts.RequireBaseDir != "" {
		inBase, err := common.IsPathInBase(path, opts.RequireBaseDir)
		if err != nil {
			return err
		}
		if !inBase {
			return &ErrCheckBadBaseDir{Path: path, BaseDir: opts.RequireBaseDir}
		}
	}

	// Get file permissions
	mode := info.Mode()
	if opts.ReadOnly && mode.Perm()&0222 != 0 {
		return &ErrCheckOpenPermissions{Path: path}
	}
	if opts.RequireWrite && mode.Perm()&0200 == 0 {
		return &ErrCheckNoWritePermissions{Path: path}
	}

	// Get file owner and group
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if opts.RequireOwner != "" {
			owner := fmt.Sprint(stat.Uid)
			if owner != opts.RequireOwner {
				return &ErrCheckBadOwner{Path: path, Expected: opts.RequireOwner, Actual: owner}
			}
		}
		if opts.RequireGroup != "" {
			group := fmt.Sprint(stat.Gid)
			if group != opts.RequireGroup {
				return &ErrCheckBadGroup{Path: path, Expected: opts.RequireGroup, Actual: group}
			}
		}
	}

	return nil
}

type ErrCheckOpenPermissions struct{ Path string }
type ErrCheckNoWritePermissions struct{ Path string }
type ErrCheckBadOwner struct{ Path, Expected, Actual string }
type ErrCheckBadGroup struct{ Path, Expected, Actual string }

type ErrCheckBadBaseDir struct{ Path, BaseDir string }

func (e *ErrCheckOpenPermissions) Error() string {
	return fmt.Sprintf("permissions too open: %s", e.Path)
}
func (e *ErrCheckNoWritePermissions) Error() string {
	return fmt.Sprintf("no write permission: %s", e.Path)
}
func (e *ErrCheckBadOwner) Error() string {
	return fmt.Sprintf("bad owner for %s: expected %s, got %s", e.Path, e.Expected, e.Actual)
}
func (e *ErrCheckBadGroup) Error() string {
	return fmt.Sprintf("bad group for %s: expected %s, got %s", e.Path, e.Expected, e.Actual)
}

func (e *ErrCheckBadBaseDir) Error() string {
	return fmt.Sprintf("file %s is not in required base directory %s", e.Path, e.BaseDir)
}
