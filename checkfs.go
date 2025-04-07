// Package checkfs provides utilities to validate and manage files and directories
// with customizable options for permissions, ownership, and creation behavior.
package checkfs

import (
	"github.com/andreimerlescu/checkfs/directory"
	"github.com/andreimerlescu/checkfs/file"
)

// File will use the file package to validate the file.Options passed into the path
func File(path string, opts file.Options) error {
	return file.File(path, opts)
}

// Directory will use the directory package to validate the directory.Options passed into the path
func Directory(path string, opts directory.Options) error {
	return directory.Directory(path, opts)
}
