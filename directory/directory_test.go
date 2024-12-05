package directory

import (
	"os"
	"testing"
)

func TestDirectory(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		opts    Options
		wantErr bool
	}{
		{"Valid directory", dir, Options{}, false},
		{"Non-directory path", dir + "/file.txt", Options{}, true},
		{"Invalid base directory", dir, Options{RequireBaseDir: "/invalid"}, true},
	}

	wErr := os.WriteFile(dir+"/file.txt", []byte("test"), 0644)
	if wErr != nil {
		t.Errorf("Error writing to directory: %v", wErr)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Directory(tt.path, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Directory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkDirectory(b *testing.B) {
	dir := b.TempDir()
	opts := Options{ReadOnly: false, RequireWrite: true}

	for i := 0; i < b.N; i++ {
		_ = Directory(dir, opts)
	}
}
