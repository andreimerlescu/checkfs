package file

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
	// NoAction will perform no action against the Create structure
	NoAction CreateKind = iota

	// IfNotExists CreateKind will perform an action on the Create structure if the path doesn't exist
	IfNotExists CreateKind = iota

	// IfExists CreateKind will perform an action on the Create structure if the path exists
	// This is intended to be a DESTRUCTIVE act when used since it removes the file first before Create.Run() is called.
	IfExists CreateKind = iota
)

// Create is used to describe the File you wish to Create, you are not required to set the Path,
// but you can if you wish to change it
type Create struct {
	Path     string      // Path stores where the resource will be created
	Kind     CreateKind  // Kind requires either IfNotExists or another CreateKind
	FileMode os.FileMode // FileMode allows you to set os.ModePerm etc.
	OpenFlag int         // OpenFlag allows you to use os.O_CREATE|os.O_TRUNC|os.O_WRONLY
	Size     int64       // Size allows you to fill a file with zeros, throws error if applied to a directory
}

// NewCreate allows you to stack the .Run() call
//
// Example:
//
//			err := file.NewCreate(file.Create{
//				Kind: file.IfNotExists,
//				Path: "/opt/test.txt",
//	  		OpenFlag: os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
//				FileMode: 0644,
//			}).Run()
func NewCreate(create *Create) *Create {
	return &Create{}
}

const (
	KB = 1 << (10 * iota)
	MB
	GB
	TB
)

func (create *Create) file() error {
	if create.Kind != IfNotExists {
		return nil
	}
	defer func() { create.Kind = NoAction }()
	theFile, err := os.OpenFile(create.Path, create.OpenFlag, create.FileMode)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}
	defer theFile.Close()

	if create.Size > TB {
		return fmt.Errorf("file size too big (max 1TB): %d", create.Size)
	}

	if create.Size > 0 {
		b := make([]byte, create.Size)
		for i := int64(0); i < create.Size; i++ {
			b[i] = byte(i)
		}
		_, err := theFile.Seek(0, 0)
		if err != nil {
			return err
		}
		bytesWritten, err := theFile.Write(b)
		if err != nil {
			return fmt.Errorf("could not write to file: %w", err)
		}
		if bytesWritten != len(b) {
			return fmt.Errorf("didnt write %d of %d to file", bytesWritten, create.Size)
		}
	}

	return nil
}

func (create *Create) replaceFile() error {
	if create.Kind != IfExists {
		return nil
	}
	err := os.Remove(create.Path)
	if err != nil {
		return fmt.Errorf("could not remove file: %w", err)
	}
	create.Kind = IfNotExists
	return create.file()
}

func (create *Create) Run() error {
	switch create.Kind {
	case IfExists:
		return create.replaceFile()
	case IfNotExists:
		return create.file()
	default:
		return fmt.Errorf("create kind not supported: %v", create.Kind)
	}
}

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
	Create             Create      // Allow the user to create the file
}

// File performs the file checks
func File(path string, opts Options) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			if opts.Create.Kind == IfNotExists {
				if len(opts.Create.Path) == 0 {
					opts.Create.Path = path
				}
				return opts.Create.Run()
			}
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
		createTime, err := common.GetCreationTime(path)
		if err != nil {
			return fmt.Errorf("failed to get creation time for %s: %w", path, err)
		}
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
		isInBase, err := common.IsPathInBase(path, opts.RequireBaseDir)
		if err != nil {
			return fmt.Errorf("failed to check base directory for %s: %w", path, err)
		}
		if !isInBase {
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
		isMorePermissive, err := common.IsMorePermissiveThan(path, opts.MorePermissiveThan)
		if err != nil {
			return fmt.Errorf("failed to check permissions for %s: %w", path, err)
		}
		if !isMorePermissive {
			return fmt.Errorf("file mode for %s is less permissive than required: expected at least %o, got %o",
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
			return fmt.Errorf("file mode for %s is more permissive than allowed: expected at most %o, got %o",
				path, opts.LessPermissiveThan, mode.Perm())
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

	// Check owner and group
	if opts.RequireOwner != "" || opts.RequireGroup != "" {
		uid, gid, err := common.GetOwnerAndGroup(path)
		if err != nil {
			return fmt.Errorf("failed to get owner/group for %s: %w", path, err)
		}
		if opts.RequireOwner != "" && uid != opts.RequireOwner {
			return &ErrCheckBadOwner{Path: path, Expected: opts.RequireOwner, Actual: uid}
		}
		if opts.RequireGroup != "" && gid != opts.RequireGroup {
			return &ErrCheckBadGroup{Path: path, Expected: opts.RequireGroup, Actual: gid}
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
