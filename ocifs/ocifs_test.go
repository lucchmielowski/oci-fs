package ocifs

import (
	"context"
	"io"
	"io/fs"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			url:     "oci://ghcr.io/test/image:latest",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil != tt.wantErr {
				t.Errorf("url.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got, err := New(u)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("New() = nil, want non-nil")
			}
		})
	}
}

func TestOciFS_WithContext(t *testing.T) {
	u, _ := url.Parse("oci://ghcr.io/test/image:latest")
	fsys := &ociFS{
		ctx:  context.Background(),
		repo: u,
	}

	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		{
			name: "With valid context",
			ctx:  context.Background(),
			want: true,
		},
		{
			name: "With nil context",
			ctx:  nil,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fsys.WithContext(tt.ctx)
			if (got != nil) != tt.want {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOciFS_Open(t *testing.T) {
	u, _ := url.Parse("oci://ghcr.io/test/image:latest")
	fsys := &ociFS{
		ctx:  context.Background(),
		repo: u,
	}

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "Invalid path with ..",
			path:    "../test",
			wantErr: true,
		},
		{
			name:    "Empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := fsys.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOciFile_Operations(t *testing.T) {
	f := &ociFile{
		ctx:  context.Background(),
		name: "test.txt",
		read: true,
	}

	t.Run("Read when already read", func(t *testing.T) {
		buf := make([]byte, 10)
		n, err := f.Read(buf)
		if n != 0 || err != io.EOF {
			t.Errorf("Read() = %v, %v, want 0, EOF", n, err)
		}
	})

	t.Run("Close", func(t *testing.T) {
		err := f.Close()
		if err != nil {
			t.Errorf("Close() error = %v", err)
		}
		if !f.read {
			t.Error("Close() should set read to true")
		}
		if f.tr != nil {
			t.Error("Close() should set tr to nil")
		}
	})
}

func TestOciFS_URL(t *testing.T) {
	testURL := "oci://ghcr.io/test/image:latest"
	u, _ := url.Parse(testURL)
	fsys := &ociFS{
		ctx:  context.Background(),
		repo: u,
	}

	if got := fsys.URL(); got != testURL {
		t.Errorf("URL() = %v, want %v", got, testURL)
	}
}

func TestFS_Implementation(t *testing.T) {
	var _ fs.FS = (*ociFS)(nil)
	var _ fs.File = (*ociFile)(nil)
}
