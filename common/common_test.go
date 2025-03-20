package common

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommonUtils(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	os.WriteFile(file, []byte("test"), 0644)

	t.Run("GetCreationTime", func(t *testing.T) {
		ctime, err := GetCreationTime(file)
		if err != nil {
			t.Errorf("GetCreationTime failed: %v", err)
		}
		if ctime.IsZero() {
			t.Error("Expected non-zero creation time")
		}
	})

	t.Run("HasPermissions", func(t *testing.T) {
		ok, err := HasPermissions(file, 0444)
		if err != nil || !ok {
			t.Errorf("HasPermissions(0444) failed: %v, got %v", err, ok)
		}
	})

	t.Run("IsMorePermissiveThan", func(t *testing.T) {
		ok, err := IsMorePermissiveThan(file, 0444)
		if err != nil || !ok {
			t.Errorf("IsMorePermissiveThan(0444) failed: %v, got %v", err, ok)
		}
		ok, err = IsMorePermissiveThan(file, 0666)
		if err != nil || ok {
			t.Errorf("IsMorePermissiveThan(0666) should fail: %v, got %v", err, ok)
		}
	})

	t.Run("IsLessPermissiveThan", func(t *testing.T) {
		ok, err := IsLessPermissiveThan(file, 0777)
		if err != nil || !ok {
			t.Errorf("IsLessPermissiveThan(0777) failed: %v, got %v", err, ok)
		}
		ok, err = IsLessPermissiveThan(file, 0400)
		if err != nil || ok {
			t.Errorf("IsLessPermissiveThan(0400) should fail: %v, got %v", err, ok)
		}
	})

	t.Run("SanitizePath", func(t *testing.T) {
		clean, err := SanitizePath("/dir//file/../test")
		if err != nil || clean != "/dir/test" {
			t.Errorf("SanitizePath failed: %v, got %v", err, clean)
		}
	})
}

func TestIsPathInBase(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		baseDir string
		want    bool
		wantErr bool
	}{
		{"Valid path in base", "/tmp/test/file.txt", "/tmp/test", true, false},
		{"Path outside base", "/tmp/other/file.txt", "/tmp/test", false, false},
		{"Path escaping base", "/tmp/test/../file.txt", "/tmp/test", false, false},
		{"Empty path", "", "/tmp/test", false, true},                    // Updated case
		{"Empty base directory", "/tmp/test/file.txt", "", false, true}, // Updated case
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsPathInBase(tt.path, tt.baseDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsPathInBase() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("IsPathInBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRelStartsWithParent(t *testing.T) {
	tests := []struct {
		name string
		rel  string
		want bool
	}{
		{"Relative path escapes", "../file.txt", true},
		{"Relative path inside", "subdir/file.txt", false},
		{"Current directory", "./file.txt", false},
		{"Escaping with separator", "../../file.txt", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RelStartsWithParent(tt.rel); got != tt.want {
				t.Errorf("RelStartsWithParent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkIsPathInBase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = IsPathInBase("/tmp/test/file.txt", "/tmp/test")
	}
}

func BenchmarkRelStartsWithParent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = RelStartsWithParent("../file.txt")
	}
}
