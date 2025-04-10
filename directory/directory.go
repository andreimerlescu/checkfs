package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreimerlescu/checkfs/common"
)

type CreateKind int8

const (

	// NoAction CreateKind Skips Create in Directory in Options
	NoAction CreateKind = iota

	// IfNotExists CreateKind reads err := os.Stat, then if os.IsNotExists(err) is true and Create.Kind is IfNotExist,
	// the Create.Run() is called to create the directory in the Create.Path
	IfNotExists CreateKind = iota

	// IfExists CreateKind relies on os.Stat where os.IsNotExists(err) is false ; meaning the path exists; if the
	// Create.Kind is IfExists then checkfs will delete the path first, then create a new directory at the path in
	// Create.Path
	IfExists CreateKind = iota
)

// Create defines a New Directory that is a CreateKind (default NoAction), options include:
// - IfNotExists
// - IfExists
// Properties in the Create struct dictate the runtime of the Create.Run() method
type Create struct {
	Kind     CreateKind  // Kind requires either CreateFileIfNotExists or IfNotExists CreateKind
	FileMode os.FileMode // FileMode allows you to set os.ModePerm etc.
	Path     string      // Path stores where the resource will be created
}

// directory will consume a pointer to Create and apply the policy against the host
func (create *Create) directory() error {
	if create.Kind != IfNotExists {
		return nil
	}
	defer func() { create.Kind = NoAction }()
	return os.MkdirAll(create.Path, create.FileMode)
}

// replaceDirectory  will consume a pointer to Create an apply the policy against the host
func (create *Create) replaceDirectory() error {
	if create.Kind != IfExists {
		return nil
	}
	err := os.RemoveAll(create.Path)
	if err != nil {
		return fmt.Errorf("could not remove directory: %w", err)
	}
	create.Kind = IfNotExists
	return create.directory()
}

// Run will read the Create.Kind and switch between IfExists and IfNotExists to run either createDirectory or
// replaceDirectory internally.
func (create *Create) Run() error {
	switch create.Kind {
	case IfExists:
		return create.replaceDirectory()
	case IfNotExists:
		return create.directory()
	default:
		return fmt.Errorf("create kind not supported: %v", create.Kind)
	}
}

type Options struct {
	CreatedBefore      time.Time   // Check directory creation time
	ModifiedBefore     time.Time   // Check directory modified time
	RequireOwner       string      // Check if the directory has a specific owner
	RequireGroup       string      // Check if the directory has a specific group
	RequireBaseDir     string      // Check if the directory is inside a specific base directory
	RequireExt         string      // Check if the directory has an extension (unlikely, but included for parity)
	RequirePrefix      string      // Check if the directory name begins with a prefix
	MorePermissiveThan os.FileMode // Check if mode is at least this permissive (e.g., >= 0444)
	LessPermissiveThan os.FileMode // Check if mode is less permissive than this (e.g., <= 0400)
	ReadOnly           bool        // Check if the directory is read-only
	RequireWrite       bool        // Check if the directory is writable
	WillCreate         bool        // User intends to create the directory, so if true, verify that we can create a directory in the parent of the path
	Create             Create      // user intends to create the directory
	Exists             bool        // If true, require the directory to exist; combining with WillCreate means Exists requires the Create to be successful
}

// Directory performs the directory checks
func Directory(path string, opts Options) error {

	// Handle WillCreate logic first
	if opts.WillCreate {
		if opts.Create.Kind == NoAction {
			opts.Create.Kind = IfNotExists
		}
		parentDir := filepath.Dir(path)
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

	if opts.Create.Kind != NoAction && len(opts.Create.Path) == 0 {
		opts.Create.Path = path
	}

	// Get directory info
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if !opts.Exists && opts.Create.Kind == NoAction {
				return nil
			}
			if opts.Create.Kind == IfNotExists {
				return opts.Create.Run()
			}
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

	if opts.Exists && opts.Create.Kind == IfExists {
		return opts.Create.Run()
	}

	// Check creation time
	if !opts.CreatedBefore.IsZero() {
		createTime, err := common.GetCreationTime(path)
		if err != nil {
			return fmt.Errorf("failed to get creation time for %s: %w", path, err)
		}
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
		isInBase, err := common.IsPathInBase(path, opts.RequireBaseDir)
		if err != nil {
			return fmt.Errorf("failed to check base directory for %s: %w", path, err)
		}
		if !isInBase {
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

	// Check more permissive than
	if opts.MorePermissiveThan != 0 {
		isMorePermissive, err := common.IsMorePermissiveThan(path, opts.MorePermissiveThan)
		if err != nil {
			return fmt.Errorf("failed to check permissions for %s: %w", path, err)
		}
		if !isMorePermissive {
			return fmt.Errorf("directory mode for %s is less permissive than required: expected at least %o, got %o",
				path, opts.MorePermissiveThan, mode.Perm())
		}
	}

	// Check less permissive than
	if opts.LessPermissiveThan != 0 {
		isLessPermissive, err := common.IsLessPermissiveThan(path, opts.LessPermissiveThan)
		if err != nil {
			return fmt.Errorf("failed to check permissions for %s: %w", path, err)
		}
		if !isLessPermissive {
			return fmt.Errorf("directory mode for %s is more permissive than allowed: expected at most %o, got %o",
				path, opts.LessPermissiveThan, mode.Perm())
		}
	}

	// Check owner and group
	if opts.RequireOwner != "" || opts.RequireGroup != "" {
		uid, gid, err := common.GetOwnerAndGroup(path)
		if err != nil {
			return fmt.Errorf("failed to get owner/group for %s: %w", path, err)
		}
		if opts.RequireOwner != "" && uid != opts.RequireOwner {
			return &ErrCheckDirBadOwner{Path: path, Expected: opts.RequireOwner, Actual: uid}
		}
		if opts.RequireGroup != "" && gid != opts.RequireGroup {
			return &ErrCheckDirBadGroup{Path: path, Expected: opts.RequireGroup, Actual: gid}
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
