package file

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	dir := t.TempDir()
	filePath := dir + "/file.txt"
	wErr := os.WriteFile(filePath, []byte("test"), 0644)
	if wErr != nil {
		t.Errorf("Error writing to file: %v", wErr)
	}

	tests := []struct {
		name    string
		path    string
		opts    Options
		wantErr bool
	}{
		{"Valid file", filePath, Options{}, false},
		{"Non-file path", dir, Options{}, true},
		{"Invalid base directory", filePath, Options{RequireBaseDir: "/invalid"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := File(tt.path, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("File() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkFile(b *testing.B) {
	dir := b.TempDir()
	filePath := dir + "/file.txt"
	wErr := os.WriteFile(filePath, []byte("test"), 0644)
	if wErr != nil {
		b.Errorf("Error writing to file: %v", wErr)
	}
	opts := Options{ReadOnly: false, RequireWrite: true}

	for i := 0; i < b.N; i++ {
		_ = File(filePath, opts)
	}
}
