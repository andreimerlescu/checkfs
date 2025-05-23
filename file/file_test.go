package file

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFile(t *testing.T) {
	// Setup test directory and files
	dir := t.TempDir()
	regularFile := filepath.Join(dir, "regular.txt")
	prefixFile := filepath.Join(dir, "prefix_test.txt")
	largeFile := filepath.Join(dir, "large.txt")
	permFile := filepath.Join(dir, "perm.txt")

	// Create regular test file
	if err := os.WriteFile(regularFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create file with prefix
	if err := os.WriteFile(prefixFile, []byte("prefix test content"), 0644); err != nil {
		t.Fatalf("Failed to create prefix test file: %v", err)
	}

	// Create large test file
	largeBuf := make([]byte, 1024*1024) // 1MB file
	if err := os.WriteFile(largeFile, largeBuf, 0644); err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// Create file for permission tests
	if err := os.WriteFile(permFile, []byte("perm test"), 0644); err != nil {
		t.Fatalf("Failed to create perm test file: %v", err)
	}

	// Create symlink for testing
	symlinkPath := filepath.Join(dir, "symlink.txt")
	if err := os.Symlink(regularFile, symlinkPath); err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	now := time.Now()
	futureTime := now.Add(24 * time.Hour)
	pastTime := now.Add(-24 * time.Hour)

	tests := []struct {
		name    string
		path    string
		opts    Options
		wantErr bool
	}{
		{"Valid regular file", regularFile, Options{}, false},
		{"Non-existent file with Exists=false", filepath.Join(dir, "nonexistent.txt"), Options{Exists: false}, false},
		{"Non-existent file with Exists=true", filepath.Join(dir, "nonexistent.txt"), Options{Exists: true}, true},
		{"Directory path", dir, Options{}, true},

		// Base directory tests
		{"Valid base directory", regularFile, Options{RequireBaseDir: dir}, false},
		{"Invalid base directory", regularFile, Options{RequireBaseDir: "/invalid"}, true},

		// File extension tests
		{"Valid extension", regularFile, Options{RequireExt: ".txt"}, false},
		{"Invalid extension", regularFile, Options{RequireExt: ".doc"}, true},

		// Prefix tests
		{"Valid prefix", prefixFile, Options{RequirePrefix: "prefix"}, false},
		{"Invalid prefix", regularFile, Options{RequirePrefix: "prefix"}, true},

		// Time-based tests
		{"Valid creation time", regularFile, Options{CreatedBefore: futureTime}, false},
		{"Invalid creation time", regularFile, Options{CreatedBefore: pastTime}, true},
		{"Valid modification time", regularFile, Options{ModifiedBefore: futureTime}, false},
		{"Invalid modification time", regularFile, Options{ModifiedBefore: pastTime}, true},

		// Size tests
		{"Valid exact size", regularFile, Options{IsSize: int64(len("test content"))}, false},
		{"Invalid exact size", regularFile, Options{IsSize: 1000}, true},
		{"Valid size less than", regularFile, Options{IsLessThan: 1000}, false},
		{"Invalid size less than", largeFile, Options{IsLessThan: 1000}, true},
		{"Valid size greater than", largeFile, Options{IsGreaterThan: 1000}, false},
		{"Invalid size greater than", regularFile, Options{IsGreaterThan: 1000}, true},

		// Name length tests
		{"Valid base name length", regularFile, Options{IsBaseNameLen: len("regular.txt")}, false},
		{"Invalid base name length", regularFile, Options{IsBaseNameLen: 5}, true},

		// Permission tests
		{"Valid read-only", regularFile, Options{ReadOnly: true}, true},           // 0644 has write bits
		{"Valid write required", regularFile, Options{RequireWrite: true}, false}, // 0644 has write
		{"Valid write-only", regularFile, Options{WriteOnly: true}, true},         // 0644 has read bits

		// File mode tests
		{"Valid file mode", regularFile, Options{IsFileMode: 0644}, false},
		{"Invalid file mode", regularFile, Options{IsFileMode: 0600}, true},

		// Symlink tests
		{"Valid symlink", symlinkPath, Options{}, false},
		{"Symlink with valid base dir", symlinkPath, Options{RequireBaseDir: dir}, false},

		// Combined options tests
		{"Multiple valid conditions", regularFile, Options{
			RequireExt:     ".txt",
			RequireBaseDir: dir,
			IsLessThan:     1000,
			RequireWrite:   true,
		}, false},
		{"Multiple conditions with one invalid", regularFile, Options{
			RequireExt:     ".txt",
			RequireBaseDir: dir,
			IsLessThan:     10, // This should fail
			RequireWrite:   true,
		}, true},

		// New MorePermissiveThan tests
		{"MorePermissiveThan 0444 with 0644", permFile, Options{MorePermissiveThan: 0444}, false},
		{"MorePermissiveThan 0444 with 0400", permFile, Options{MorePermissiveThan: 0444}, true},  // Set perms later
		{"MorePermissiveThan 0444 with 0744", permFile, Options{MorePermissiveThan: 0444}, false}, // Set perms later
		{"MorePermissiveThan 0644 with 0644", permFile, Options{MorePermissiveThan: 0644}, false},
		{"MorePermissiveThan 0644 with 0444", permFile, Options{MorePermissiveThan: 0644}, true}, // Set perms later

		// New LessPermissiveThan tests
		{"LessPermissiveThan 0400 with 0600", permFile, Options{LessPermissiveThan: 0400}, true}, // Set perms later
		{"LessPermissiveThan 0777 with 0644", permFile, Options{LessPermissiveThan: 0777}, false},
		{"LessPermissiveThan 0644 with 0744", permFile, Options{LessPermissiveThan: 0644}, true}, // Set perms later
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Adjust permissions for specific tests
			switch tt.name {
			case "MorePermissiveThan 0444 with 0400":
				os.Chmod(permFile, 0400)
			case "MorePermissiveThan 0444 with 0744":
				os.Chmod(permFile, 0744)
			case "MorePermissiveThan 0644 with 0444":
				os.Chmod(permFile, 0444)
			case "LessPermissiveThan 0400 with 0400":
				os.Chmod(permFile, 0400)
			case "LessPermissiveThan 0400 with 0600":
				os.Chmod(permFile, 0600)
			case "LessPermissiveThan 0644 with 0744":
				os.Chmod(permFile, 0744)
			}
			err := File(tt.path, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("File() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Reset permFile to 0644 after each test
			if tt.path == permFile {
				os.Chmod(permFile, 0644)
			}
		})
	}
}

func BenchmarkFile(b *testing.B) {
	dir := b.TempDir()
	filePath := filepath.Join(dir, "benchmark.txt")
	if err := os.WriteFile(filePath, []byte("benchmark content"), 0644); err != nil {
		b.Fatalf("Failed to create benchmark file: %v", err)
	}

	cases := []struct {
		name string
		opts Options
	}{
		{"BasicChecks", Options{RequireWrite: true}},
		{"ExtensiveChecks", Options{
			RequireExt:     ".txt",
			RequireBaseDir: dir,
			IsLessThan:     1000,
			RequireWrite:   true,
			ReadOnly:       false,
		}},
		{"PermissiveChecks", Options{
			MorePermissiveThan: 0444,
			LessPermissiveThan: 0777,
		}},
	}

	for _, bc := range cases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = File(filePath, bc.opts)
			}
		})
	}
}
