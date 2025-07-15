package ocifs

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/nlepage/go-tarfs"
	"io/fs"
	"net/url"
)

type ociFS struct {
	ctx context.Context

	repo *url.URL
	img  *v1.Image

	tarFS fs.FS
}

var FS = fsimpl.FSProviderFunc(New, "oci")

var (
	_ fs.FS = (*ociFS)(nil)
	// TODO: add WithAuthenticator
)

func (f *ociFS) URL() string {
	return f.repo.String()
}

func New(u *url.URL) (fs.FS, error) {
	repoUrl := *u

	fsys := &ociFS{
		ctx:  context.Background(),
		repo: &repoUrl,
	}

	img, err := fsys.fetchManifest()
	if err != nil {
		return nil, err
	}

	fsys.img = &img

	// Read content of flattened artifact filesystem
	rc := mutate.Extract(img)
	fsys.tarFS, err = tarfs.New(rc)
	if err != nil {
		return nil, err
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

	ociFile, err := f.tarFS.Open(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	return ociFile, nil
}

func (f *ociFS) fetchManifest() (v1.Image, error) {
	ref := fmt.Sprintf("%s%s", f.repo.Host, f.repo.Path)
	rOpts, nOpts, err := RegistryOpts()
	if err != nil {
		return nil, err
	}

	repoRef, err := name.ParseReference(ref, nOpts...)
	if err != nil {
		return nil, err
	}

	// TODO: Add auth options for private image fetching
	// TODO: Add signature verification for manifest
	img, err := remote.Image(repoRef, rOpts...)
	if err != nil {
		return nil, err
	}
	f.img = &img

	return img, err
}
