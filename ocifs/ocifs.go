package ocifs

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hairyhenderson/go-fsimpl"
	"io"
	"io/fs"
	"net/url"
)

type ociFS struct {
	ctx context.Context

	repo *url.URL
	img  *v1.Image
}

var FS = fsimpl.FSProviderFunc(New, "oci")

var (
	_ fs.FS = (*ociFS)(nil)
	// TODO: _ fs.ReadDirFS = (*ociFS)(nil)

	// TODO: add WithAuthenticator
)

func (o *ociFS) URL() string {
	return o.repo.String()
}

func New(u *url.URL) (fs.FS, error) {
	repoUrl := *u

	fsys := &ociFS{
		ctx:  context.Background(),
		repo: &repoUrl,
	}

	return fsys, nil
}

func (f *ociFS) WithContext(ctx context.Context) fs.FS {
	if ctx == nil {
		return f
	}

	fsys := *f
	fsys.ctx = ctx

	return &fsys
}

func (f *ociFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	if f.img == nil {
		img, err := f.fetchManifest()
		if err != nil {
			return nil, err
		}
		f.img = &img
	}

	ociFile := &ociFile{
		name: name,
		img:  f.img,
	}

	err := ociFile.extractFile(name)
	if err != nil {
		return nil, err
	}
	return ociFile, nil
}

func (f *ociFS) fetchManifest() (v1.Image, error) {
	ref := fmt.Sprintf("%s%s", f.repo.Host, f.repo.Path)
	fmt.Println(ref)
	repoRef, err := name.ParseReference(ref)
	if err != nil {
		return nil, err
	}

	// TODO: Add auth options for private image fetching
	// TODO: Add signature verification for manifest
	img, err := remote.Image(repoRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return nil, err
	}
	f.img = &img

	return img, err
}

type ociFile struct {
	ctx context.Context

	name string
	img  *v1.Image

	fi   fs.FileInfo
	tr   *tar.Reader
	rc   io.ReadCloser
	read bool
}

var _ fs.File = (*ociFile)(nil)

func (f *ociFile) extractFile(name string) error {
	// Read content of flattened artifact filesystem
	f.rc = mutate.Extract(*f.img)
	tr := tar.NewReader(f.rc)

	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}

		if name == hdr.Name {
			f.fi = hdr.FileInfo()
			f.tr = tr
			f.read = false
			return nil
		}

	}

	return nil
}

func (f *ociFile) Stat() (fs.FileInfo, error) {
	return f.fi, nil
}

func (f *ociFile) Read(p []byte) (int, error) {
	if f.read {
		return 0, io.EOF
	}

	n, err := f.tr.Read(p)
	if errors.Is(err, io.EOF) {
		f.read = true
	}
	return n, err
}

func (f *ociFile) Close() error {
	f.read = true
	f.tr = nil
	if f.rc != nil {
		return f.rc.Close()
	}
	return nil
}
