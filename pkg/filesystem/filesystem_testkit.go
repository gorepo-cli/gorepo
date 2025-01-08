package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type MockMethods interface {
	Methods
	Output() map[string][]byte
}

type MockFileSystem struct {
	Files map[string][]byte
	Wd    string
}

var _ MockMethods = &MockFileSystem{}

func NewMockFilesystem(files map[string][]byte, wd string) *MockFileSystem {
	return &MockFileSystem{
		Files: files,
		Wd:    wd,
	}
}

func (m *MockFileSystem) Exists(absolutePath string) bool {
	_, exists := m.Files[absolutePath]
	return exists
}

func (m *MockFileSystem) Read(absolutePath string) ([]byte, error) {
	if data, exists := m.Files[absolutePath]; exists {
		return data, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystem) Write(absolutePath string, content []byte) error {
	m.Files[absolutePath] = content
	return nil
}

func (m *MockFileSystem) Walk(rootPath string, walkFn filepath.WalkFunc) error {
	dirs := make(map[string]bool)
	filesInDir := make(map[string][]string)

	rootPath = filepath.Clean(rootPath)
	dirs[rootPath] = true

	for path := range m.Files {
		if !strings.HasPrefix(filepath.Clean(path), rootPath) {
			continue
		}

		dir := filepath.Dir(path)
		dirs[dir] = true
		filesInDir[dir] = append(filesInDir[dir], path)
	}

	var dirList []string
	for d := range dirs {
		dirList = append(dirList, d)
	}
	dirList = sortByPathLength(dirList)

	for _, d := range dirList {
		dInfo := mockFileInfo{
			name:  filepath.Base(d),
			size:  0,
			mode:  fs.ModeDir | 0755,
			isDir: true,
			sys:   nil,
		}

		err := walkFn(d, dInfo, nil)
		if err != nil {
			if err == filepath.SkipDir {
				continue
			}
			return err
		}

		for _, filePath := range filesInDir[d] {
			info := mockFileInfo{
				name:  filepath.Base(filePath),
				size:  int64(len(m.Files[filePath])),
				mode:  0644,
				isDir: false,
				sys:   nil,
			}
			err := walkFn(filePath, info, nil)
			if err != nil {
				if err == filepath.SkipDir {
					continue
				}
				return err
			}
		}
	}
	return nil
}

// Helper function to sort by path length so that shorter directories
// (closer to root) come first. This ensures we walk parent directories
// before children, mimicking the behavior of filepath.Walk.
func sortByPathLength(paths []string) []string {
	sorted := append([]string{}, paths...)
	// Sort by length, then lexicographically to have a deterministic order
	// if needed.
	sort.Slice(sorted, func(i, j int) bool {
		if len(sorted[i]) == len(sorted[j]) {
			return sorted[i] < sorted[j]
		}
		return len(sorted[i]) < len(sorted[j])
	})
	return sorted
}

type mockFileInfo struct {
	name  string
	size  int64
	mode  fs.FileMode
	isDir bool
	sys   interface{}
}

func (m mockFileInfo) Name() string       { return m.name }
func (m mockFileInfo) Size() int64        { return m.size }
func (m mockFileInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFileInfo) IsDir() bool        { return m.isDir }
func (m mockFileInfo) Sys() interface{}   { return m.sys }
func (m mockFileInfo) ModTime() time.Time { return time.Time{} }

type mockTime struct{}

func (mockTime) Unix() int64            { return 0 }
func (mockTime) String() string         { return "mockTime" }
func (mockTime) IsZero() bool           { return true }
func (mockTime) Before(t mockTime) bool { return false }

func (m *MockFileSystem) GetWd() (string, error) {
	return m.Wd, nil
}

func (m *MockFileSystem) Output() map[string][]byte {
	return m.Files
}
