package common

import (
	"testing"
)

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
