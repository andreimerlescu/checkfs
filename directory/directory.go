package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Options struct {
	ReadOnly       bool   // Check if the directory is read-only
	RequireWrite   bool   // Check if the directory is writable
	RequireOwner   string // Check if the directory has a specific owner
	RequireGroup   string // Check if the directory has a specific group
	RequireBaseDir string // Check if the directory is inside a specific base directory

	// NEW OPTIONS

	CreatedBefore  time.Time // Check file creation time
	ModifiedBefore time.Time // Check file modified time
	RequireExt     string    // Check if the file is of an extension
	RequirePrefix  string    // Check if the file name begins with a prefix

	// user intends to create the directory, so if true, verify that we can create a directory in the parent of the path
	WillCreate bool

	// if true, then require the directory to exist ; combining with WillCreate means Exists requires the Create to be
	// successful - the script should (only if it doesn't exist) try to create the file with a random number in it, then
	// delete the file - if both operations succeed, then Exists is true when WillCreate is true
	Exists bool
}

// Directory performs the directory checks
func Directory(path string, opts Options) error {
	// Handle WillCreate logic first
	if opts.WillCreate {
		parentDir := filepath.Dir(path)
		// Check if parent directory exists and is writable
		parentInfo, err := os.Stat(parentDir)
		if err != nil {
			return fmt.Errorf("failed to access parent directory %s: %w", parentDir, err)
		}
		if !parentInfo.IsDir() {
			return fmt.Errorf("parent path is not a directory: %s", parentDir)
		}
		if parentInfo.Mode().Perm()&0200 == 0 {
			return fmt.Errorf("parent directory not writable: %s", parentDir)
		}
	}

	// Get directory info
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if opts.Exists && !opts.WillCreate {
				return fmt.Errorf("directory does not exist: %s", path)
			}
			return nil
		}
		return fmt.Errorf("failed to stat directory %s: %w", path, err)
	}

	// Directory exists - check if we explicitly don't want it to
	if !opts.Exists && !opts.WillCreate {
		return fmt.Errorf("directory exists but was expected not to exist: %s", path)
	}

	// Check if path is a directory
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}

	// Check creation time
	if !opts.CreatedBefore.IsZero() {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("unable to get detailed directory stats for %s", path)
		}
		createTime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
		if createTime.After(opts.CreatedBefore) {
			return fmt.Errorf("directory created after specified time: %s", path)
		}
	}

	// Check modification time
	if !opts.ModifiedBefore.IsZero() && info.ModTime().After(opts.ModifiedBefore) {
		return fmt.Errorf("directory modified after specified time: %s", path)
	}

	// Check directory prefix
	if opts.RequirePrefix != "" {
		basename := filepath.Base(path)
		if !strings.HasPrefix(basename, opts.RequirePrefix) {
			return fmt.Errorf("incorrect directory prefix for %s: expected prefix %s",
				path, opts.RequirePrefix)
		}
	}

	// Check if directory is inside the required base directory
	if opts.RequireBaseDir != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
		}
		absBaseDir, err := filepath.Abs(opts.RequireBaseDir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for base dir %s: %w", opts.RequireBaseDir, err)
		}
		if !strings.HasPrefix(absPath, absBaseDir) {
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
