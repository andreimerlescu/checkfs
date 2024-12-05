package go_checkfs

import (
	"github.com/andreimerlescu/go-checkfs/directory"
	"github.com/andreimerlescu/go-checkfs/file"
)

func File(path string, opts file.Options) error {
	return file.File(path, opts)
}
func Directory(path string, opts directory.Options) error {
	return directory.Directory(path, opts)
}
