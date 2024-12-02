package repository

import "io/fs"

type FS interface {
	fs.ReadDirFS
	fs.StatFS
	fs.SubFS
	Create(name string) (fs.File, error)
	CreateDir(name string) error
	Remove(name string) error
}
