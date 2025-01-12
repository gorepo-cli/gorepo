package filesystem

import (
	"os"
	"path/filepath"
)

type Methods interface {
	Exists(absolutePath string) bool
	Read(absolutePath string) ([]byte, error)
	Write(absolutePath string, content []byte) error
	Walk(rootPath string, walkFn filepath.WalkFunc) error
	MkDir(absolutePath string) error
	GetWd() (absolutePath string, err error)
}

type Filesystem struct{}

var _ Methods = &Filesystem{}

func (fs *Filesystem) Exists(absolutePath string) (exists bool) {
	_, err := os.Stat(absolutePath)
	return err == nil
}

func (fs *Filesystem) Read(absolutePath string) ([]byte, error) {
	return os.ReadFile(absolutePath)
}

func (fs *Filesystem) Write(absolutePath string, content []byte) (err error) {
	return os.WriteFile(absolutePath, content, 0644)
}

func (fs *Filesystem) Walk(rootPath string, walkFn filepath.WalkFunc) (err error) {
	return filepath.Walk(rootPath, walkFn)
}

func (fs *Filesystem) MkDir(absolutePath string) (err error) {
	return os.MkdirAll(absolutePath, 0755)
}

func (fs *Filesystem) GetWd() (dir string, err error) {
	return os.Getwd()
}
