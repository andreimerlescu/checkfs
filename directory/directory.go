package directory

import (
	"fmt"
	"github.com/andreimerlescu/go-checkfs/common"
	"os"
	"syscall"
)

type Options struct {
	ReadOnly       bool   // Check if the directory is read-only
	RequireWrite   bool   // Check if the directory is writable
	RequireOwner   string // Check if the directory has a specific owner
	RequireGroup   string // Check if the directory has a specific group
	RequireBaseDir string // Check if the directory is inside a specific base directory
}

// Directory performs the directory checks
func Directory(path string, opts Options) error {
	// Get directory info
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat directory %s: %w", path, err)
	}

	// Check if path is a directory
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	// Check if directory is inside the required base directory
	if opts.RequireBaseDir != "" {
		inBase, err := common.IsPathInBase(path, opts.RequireBaseDir)
		if err != nil {
			return err
		}
		if !inBase {
			return &ErrCheckDirBadBaseDir{Path: path, BaseDir: opts.RequireBaseDir}
		}
	}

	// Get directory permissions
	mode := info.Mode()
	if opts.ReadOnly && mode.Perm()&0222 != 0 {
		return &ErrCheckDirOpenPermissions{Path: path}
	}
	if opts.RequireWrite && mode.Perm()&0200 == 0 {
		return &ErrCheckDirNoWritePermissions{Path: path}
	}

	// Get directory owner and group
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if opts.RequireOwner != "" {
			owner := fmt.Sprint(stat.Uid)
			if owner != opts.RequireOwner {
				return &ErrCheckDirBadOwner{Path: path, Expected: opts.RequireOwner, Actual: owner}
			}
		}
		if opts.RequireGroup != "" {
			group := fmt.Sprint(stat.Gid)
			if group != opts.RequireGroup {
				return &ErrCheckDirBadGroup{Path: path, Expected: opts.RequireGroup, Actual: group}
			}
		}
	}

	return nil
}

type ErrCheckDirOpenPermissions struct{ Path string }
type ErrCheckDirNoWritePermissions struct{ Path string }
type ErrCheckDirBadOwner struct{ Path, Expected, Actual string }
type ErrCheckDirBadGroup struct{ Path, Expected, Actual string }

type ErrCheckDirBadBaseDir struct{ Path, BaseDir string }

func (e *ErrCheckDirOpenPermissions) Error() string {
	return fmt.Sprintf("permissions too open: %s", e.Path)
}
func (e *ErrCheckDirNoWritePermissions) Error() string {
	return fmt.Sprintf("no write permission: %s", e.Path)
}
func (e *ErrCheckDirBadOwner) Error() string {
	return fmt.Sprintf("bad owner for %s: expected %s, got %s", e.Path, e.Expected, e.Actual)
}
func (e *ErrCheckDirBadGroup) Error() string {
	return fmt.Sprintf("bad group for %s: expected %s, got %s", e.Path, e.Expected, e.Actual)
}

func (e *ErrCheckDirBadBaseDir) Error() string {
	return fmt.Sprintf("directory %s is not in required base directory %s", e.Path, e.BaseDir)
}
