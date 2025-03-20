package directory

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDirectory(t *testing.T) {
	baseDir := t.TempDir()
	testDir := filepath.Join(baseDir, "test_directory")
	prefixDir := filepath.Join(baseDir, "prefix_test_dir")
	nonExistentDir := filepath.Join(baseDir, "nonexistent")
	readOnlyDir := filepath.Join(baseDir, "readonly_dir")
	writeableDir := filepath.Join(baseDir, "writeable_dir")
	permDir := filepath.Join(baseDir, "perm_dir")

	// Create test directories with specific permissions
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	if err := os.Mkdir(prefixDir, 0755); err != nil {
		t.Fatalf("Failed to create prefix directory: %v", err)
	}
	if err := os.Mkdir(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create readonly directory: %v", err)
	}
	if err := os.Mkdir(writeableDir, 0755); err != nil {
		t.Fatalf("Failed to create writeable directory: %v", err)
	}
	if err := os.Mkdir(permDir, 0755); err != nil {
		t.Fatalf("Failed to create perm directory: %v", err)
	}

	// Create a test file to validate non-directory path
	testFile := filepath.Join(baseDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
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
		// Basic existence tests
		{"Valid existing directory", testDir, Options{Exists: true}, false},
		{"Non-existent directory with Exists=false", nonExistentDir, Options{Exists: false}, false},
		{"Non-existent directory with Exists=true", nonExistentDir, Options{Exists: true}, true},
		{"Non-directory path", testFile, Options{Exists: true}, true},

		// Creation capability tests
		{"Will create in existing parent", filepath.Join(baseDir, "new_dir"), Options{
			WillCreate: true,
			Exists:     true,
		}, false},
		{"Will create with existing target", testDir, Options{
			WillCreate: true,
			Exists:     true,
		}, false},
		{"Will create without existence check", filepath.Join(baseDir, "just_create"), Options{
			WillCreate: true,
		}, false},
		{"Will create and require existence", filepath.Join(baseDir, "create_exist"), Options{
			WillCreate: true,
			Exists:     true,
		}, false},
		{"Will create without existence check", filepath.Join(baseDir, "just_create"), Options{
			WillCreate: true,
			Exists:     false,
		}, false},

		// Base directory tests
		{"Valid base directory", testDir, Options{Exists: true, RequireBaseDir: baseDir}, false},
		{"Invalid base directory", testDir, Options{Exists: true, RequireBaseDir: "/invalid"}, true},

		// Prefix tests
		{"Valid prefix", prefixDir, Options{Exists: true, RequirePrefix: "prefix"}, false},
		{"Invalid prefix", testDir, Options{Exists: true, RequirePrefix: "prefix"}, true},

		// Time-based tests
		{"Valid creation time", testDir, Options{Exists: true, CreatedBefore: futureTime}, false},
		{"Invalid creation time", testDir, Options{Exists: true, CreatedBefore: pastTime}, true},
		{"Valid modification time", testDir, Options{Exists: true, ModifiedBefore: futureTime}, false},
		{"Invalid modification time", testDir, Options{Exists: true, ModifiedBefore: pastTime}, true},

		// Permission tests
		{"Read-only directory check", readOnlyDir, Options{Exists: true, ReadOnly: true}, false},
		{"Write permission check", writeableDir, Options{Exists: true, RequireWrite: true}, false},
		{"Invalid write permission", readOnlyDir, Options{Exists: true, RequireWrite: true}, true},

		// Owner and group tests
		{"Valid owner", testDir, Options{Exists: true, RequireOwner: fmt.Sprint(os.Getuid())}, false},
		{"Invalid owner", testDir, Options{Exists: true, RequireOwner: "99999"}, true},
		{"Valid group", testDir, Options{Exists: true, RequireGroup: fmt.Sprint(os.Getgid())}, false},
		{"Invalid group", testDir, Options{Exists: true, RequireGroup: "99999"}, true},

		// MorePermissiveThan tests
		{"MorePermissiveThan 0444 with 0755", permDir, Options{Exists: true, MorePermissiveThan: 0444}, false},
		{"MorePermissiveThan 0444 with 0400", permDir, Options{Exists: true, MorePermissiveThan: 0444}, true}, // Set perms later
		{"MorePermissiveThan 0444 with 0744", permDir, Options{Exists: true, MorePermissiveThan: 0444}, false}, // Set perms later
		{"MorePermissiveThan 0644 with 0755", permDir, Options{Exists: true, MorePermissiveThan: 0644}, false},
		{"MorePermissiveThan 0644 with 0444", permDir, Options{Exists: true, MorePermissiveThan: 0644}, true},  // Set perms later

		// LessPermissiveThan tests
		{"LessPermissiveThan 0400 with 0400", permDir, Options{Exists: true, LessPermissiveThan: 0400}, false}, // Set perms later
		{"LessPermissiveThan 0400 with 0755", permDir, Options{Exists: true, LessPermissiveThan: 0400}, true},
		{"LessPermissiveThan 0777 with 0755", permDir, Options{Exists: true, LessPermissiveThan: 0777}, false},
		{"LessPermissiveThan 0755 with 0777", permDir, Options{Exists: true, LessPermissiveThan: 0755}, true},  // Set perms later

		// Multiple conditions
		{"Multiple valid conditions", writeableDir, Options{
			Exists:         true,
			RequireWrite:   true,
			RequirePrefix:  "writeable",
			RequireBaseDir: baseDir,
		}, false},
		{"Multiple conditions with one invalid", writeableDir, Options{
			Exists:         true,
			RequireWrite:   true,
			RequirePrefix:  "invalid",
			RequireBaseDir: baseDir,
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Adjust permissions for specific tests
			switch tt.name {
			case "MorePermissiveThan 0444 with 0400":
				os.Chmod(permDir, 0400)
			case "MorePermissiveThan 0444 with 0744":
				os.Chmod(permDir, 0744)
			case "MorePermissiveThan 0644 with 0444":
				os.Chmod(permDir, 0444)
			case "LessPermissiveThan 0400 with 0400":
				os.Chmod(permDir, 0400)
			case "LessPermissiveThan 0755 with 0777":
				os.Chmod(permDir, 0777)
			}
			err := Directory(tt.path, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Directory() error = %v, wantErr %v", err, tt.wantErr)
			}
			// Reset permDir to 0755 after each test
			if tt.path == permDir {
				os.Chmod(permDir, 0755)
			}
		})
	}
}

func BenchmarkDirectory(b *testing.B) {
	dir := b.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "bench"), 0755); err != nil {
		b.Fatalf("Failed to create benchmark directory: %v", err)
	}

	cases := []struct {
		name string
		opts Options
	}{
		{"BasicChecks", Options{
			Exists:       true,
			RequireWrite: true,
		}},
		{"ExtensiveChecks", Options{
			Exists:         true,
			WillCreate:     true,
			RequireWrite:   true,
			RequirePrefix:  "bench",
			RequireBaseDir: dir,
		}},
		{"PermissiveChecks", Options{
			Exists:             true,
			MorePermissiveThan: 0444,
			LessPermissiveThan: 0777,
		}},
	}

	for _, bc := range cases {
		b.Run(bc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = Directory(filepath.Join(dir, "bench"), bc.opts)
			}
		})
	}
}