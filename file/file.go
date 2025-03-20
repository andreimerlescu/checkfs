package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type Options struct {
	CreatedBefore      time.Time   // Check file creation time
	ModifiedBefore     time.Time   // Check file modified time
	IsLessThan         int64       // Check if the size is less than
	IsSize             int64       // Check the file size
	IsGreaterThan      int64       // Check if the size is greater than
	RequireExt         string      // Check if the file is of an extension
	RequirePrefix      string      // Check if the file name begins with a prefix
	RequireOwner       string      // Check if the file has a specific owner
	RequireGroup       string      // Check if the file has a specific group
	RequireBaseDir     string      // Check if the file is inside a specific base directory
	IsFileMode         os.FileMode // Check the os.FileMode value
	MorePermissiveThan os.FileMode // Check if mode is at least this permissive (e.g., >= 0444)
	LessPermissiveThan os.FileMode // Check if mode is less permissive than this (e.g., <= 0400)
	IsBaseNameLen      int         // Check if the file name length
	RequireWrite       bool        // Check if the file is writable
	ReadOnly           bool        // Check if the file is read-only
	WriteOnly          bool        // Check if the file is write-only
	Exists             bool        // Check if the file exists
}

// File performs the file checks
func File(path string, opts Options) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if opts.Exists {
				return fmt.Errorf("file does not exist: %s", path)
			}
			return nil
		}
		return fmt.Errorf("failed to stat file %s: %w", path, err)
	}

	// Check if file is a regular file
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", path)
	}

	// Check file creation time
	if !opts.CreatedBefore.IsZero() {
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("unable to get detailed file stats for %s", path)
		}
		createTime := time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
		if createTime.After(opts.CreatedBefore) {
			return fmt.Errorf("file created after specified time: %s", path)
		}
	}

	// Check modification time
	if !opts.ModifiedBefore.IsZero() && info.ModTime().After(opts.ModifiedBefore) {
		return fmt.Errorf("file modified after specified time: %s", path)
	}

	// Check file extension
	if opts.RequireExt != "" {
		ext := filepath.Ext(path)
		if ext != opts.RequireExt {
			return fmt.Errorf("incorrect file extension for %s: expected %s, got %s",
				path, opts.RequireExt, ext)
		}
	}

	// Check file prefix
	if opts.RequirePrefix != "" {
		basename := filepath.Base(path)
		if !strings.HasPrefix(basename, opts.RequirePrefix) {
			return fmt.Errorf("incorrect file prefix for %s: expected prefix %s",
				path, opts.RequirePrefix)
		}
	}

	// Check base directory
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
			return &ErrCheckBadBaseDir{Path: path, BaseDir: opts.RequireBaseDir}
		}
	}

	// Check file size constraints
	size := info.Size()
	if opts.IsSize != 0 && size != opts.IsSize {
		return fmt.Errorf("incorrect file size for %s: expected %d, got %d",
			path, opts.IsSize, size)
	}
	if opts.IsLessThan != 0 && size >= opts.IsLessThan {
		return fmt.Errorf("file size %d is not less than %d: %s",
			size, opts.IsLessThan, path)
	}
	if opts.IsGreaterThan != 0 && size <= opts.IsGreaterThan {
		return fmt.Errorf("file size %d is not greater than %d: %s",
			size, opts.IsGreaterThan, path)
	}

	// Check base name length
	if opts.IsBaseNameLen != 0 {
		basename := filepath.Base(path)
		if len(basename) != opts.IsBaseNameLen {
			return fmt.Errorf("incorrect base name length for %s: expected %d, got %d",
				path, opts.IsBaseNameLen, len(basename))
		}
	}

	// Check file mode
	mode := info.Mode()
	if opts.IsFileMode != 0 && mode != opts.IsFileMode {
		return fmt.Errorf("incorrect file mode for %s: expected %s, got %s",
			path, opts.IsFileMode, mode)
	}

	// Check more permissive than
	if opts.MorePermissiveThan != 0 {
		perms := mode.Perm()
		required := opts.MorePermissiveThan
		if perms&required != required {
			return fmt.Errorf("file mode for %s is less permissive than required: expected at least %o, got %o",
				path, required, perms)
		}
	}

	// Check less permissive than
	if opts.LessPermissiveThan != 0 {
		perms := mode.Perm()
		limit := opts.LessPermissiveThan
		if perms&^limit != 0 {
			return fmt.Errorf("file mode for %s is more permissive than allowed: expected at most %o, got %o",
				path, limit, perms)
		}
	}

	// Check permissions
	if opts.ReadOnly && mode.Perm()&0222 != 0 {
		return &ErrCheckOpenPermissions{Path: path}
	}
	if opts.WriteOnly && mode.Perm()&0444 != 0 {
		return fmt.Errorf("file has read permissions when write-only required: %s", path)
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
