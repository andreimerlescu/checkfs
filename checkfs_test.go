package checkfs

import (
	"github.com/andreimerlescu/checkfs/directory"
	"github.com/andreimerlescu/checkfs/file"
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	dir := t.TempDir()
	filePath := dir + "/file.txt"
	wErr := os.WriteFile(filePath, []byte("test"), 0644)
	if wErr != nil {
		t.Errorf("Error writing file: %v", wErr)
	}

	err := File(filePath, file.Options{})
	if err != nil {
		t.Errorf("File() error = %v", err)
	}
}

func TestDirectory(t *testing.T) {
	dir := t.TempDir()

	err := Directory(dir, directory.Options{})
	if err == nil {
		t.Errorf("Directory() should have thrown err but got %v", err)
	}
}

func BenchmarkFile(b *testing.B) {
	dir := b.TempDir()
	filePath := dir + "/file.txt"
	wErr := os.WriteFile(filePath, []byte("test"), 0644)
	if wErr != nil {
		b.Fatalf("Error writing file: %v", wErr)
	}

	for i := 0; i < b.N; i++ {
		_ = File(filePath, file.Options{})
	}
}

func BenchmarkDirectory(b *testing.B) {
	dir := b.TempDir()

	for i := 0; i < b.N; i++ {
		_ = Directory(dir, directory.Options{})
	}
}
