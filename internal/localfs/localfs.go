package localfs

import (
	"io/fs"
	"os"
	"path"

	"github.com/davidjspooner/dsrepo/internal/repository"
)

type FS struct {
	base string
}

var _ repository.FS = &FS{}

func (localFs *FS) Create(name string) (fs.File, error) {
	fullname := path.Join(localFs.base, name)
	return os.Create(fullname)
}

func (localFs *FS) CreateDir(name string) error {
	fullname := path.Join(localFs.base, name)
	return os.Mkdir(fullname, 0755)
}

func (localFs *FS) Remove(name string) error {
	fullname := path.Join(localFs.base, name)
	return os.Remove(fullname)
}

func (localFs *FS) Open(name string) (fs.File, error) {
	fullname := path.Join(localFs.base, name)
	return os.Open(fullname)
}

func (localFs *FS) Stat(name string) (fs.FileInfo, error) {
	fullname := path.Join(localFs.base, name)
	return os.Stat(fullname)
}

func (localFs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	fullname := path.Join(localFs.base)
	entries, err := os.ReadDir(fullname)
	return entries, err
}

func (localFs *FS) Sub(name string) (fs.FS, error) {
	fullname := path.Join(localFs.base, name)
	return NewFS(fullname)
}

func NewFS(fullname string) (*FS, error) {
	stat, err := os.Stat(fullname)
	if err != nil {
		return nil, err
	}
	if !stat.IsDir() {
		return nil, fs.ErrNotExist
	}
	return &FS{base: fullname}, nil
}
