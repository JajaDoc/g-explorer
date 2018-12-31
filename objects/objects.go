package objects

import (
	"os"
	"io/ioutil"
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
	return objects, nil
}
