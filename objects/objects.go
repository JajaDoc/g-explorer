package objects

import (
	"os"
	"io/ioutil"
	"time"
)

type Object struct {
	Info os.FileInfo
}


// CurrentDir returns directory path.
// if occurred error, return err.
func CurrentDir() (dir string, err error) {
	return  os.Getwd()
}

// ChangeDir changes directory and return directory path.
func ChangeDir(path string) (dir string, err error) {
	os.Chdir(path)
	return CurrentDir()
}

func GetObjects(path string) ([]Object, error) {
	files, err :=  ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var objects []Object
	for _, f := range files {
		objects = append(objects, Object{f})
	}

	// parent dir
	parentDir := &ParentDir{
		name:    "..",
		isDir:   true,
		modTime: time.Now(),
		mode:    os.FileMode(os.ModeDir),
		size:    0,
	}
	objects = append([]Object{{parentDir}}, objects...)

	return objects, nil
}

type ParentDir struct {
	name string
	size int64
	mode os.FileMode
	modTime time.Time
	isDir bool
}

func (p *ParentDir) Name() string {
	return p.name
}

func (p *ParentDir) Size() int64 {
	return p.size
}

func (p *ParentDir) Mode() os.FileMode {
	return p.mode
}

func (p ParentDir) ModTime() time.Time {
	return p.modTime
}

func (p *ParentDir) IsDir() bool {
	return p.isDir
}

func (p *ParentDir) Sys() interface{} {
	return nil
}
