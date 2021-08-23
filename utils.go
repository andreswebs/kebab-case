package main

import (
	"os"
	"path/filepath"
	"strings"
)

// ProcessFile recursively renames files and directories using the defined format
func ProcessFile(name string) (err error) {

	fullpath, err := filepath.Abs(name)
	if err != nil {
		return
	}

	newname := getNewname(fullpath)
	err = rename(fullpath, newname)
	if err != nil {
		return
	}

	dir, err := isDir(newname)
	if err != nil {
		return
	}

	if !dir {
		return
	}

	filenames, dirnames, err := sift(newname)
	if err != nil {
		return
	}

	// rename files
	for _, n := range filenames {
		newname := getNewname(n)
		err = rename(n, newname)
		if err != nil {
			return
		}
	}

	var newDirnames []string

	// rename dirs
	for _, d := range dirnames {
		newdirname := getNewname(d)
		err = rename(d, newdirname)
		if err != nil {
			return
		}
		newDirnames = append(newDirnames, newdirname)
	}

	for _, d := range newDirnames {
		ProcessFile(d)
	}

	return

}

// rename mv file to new name if new name is different from previous name
func rename(prevname, newname string) (err error) {
	if newname == prevname {
		return
	}
	err = os.Rename(prevname, newname)
	return
}

// getNewname get new full file name renamed using the defined format
func getNewname(name string) (newname string) {

	// get the dir
	dir, filename := filepath.Split(name)

	// get the ext
	ext := filepath.Ext(name)

	// get the base
	basename := getBasename(filename)

	// rename
	newbasename := Format(basename) + ext
	newname = filepath.Join(dir, newbasename)

	return

}

// getBasename get file name without extension
func getBasename(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// isDir check if a file represented by `path` is a directory
func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// sift inspects a directory and returns a list of filenames, a list of dirnames and an error
func sift(dir string) (filenames, dirnames []string, err error) {

	fullpath, err := filepath.Abs(dir)
	if err != nil {
		return
	}

	files, err := os.ReadDir(fullpath)
	if err != nil {
		return
	}

	for _, file := range files {

		name := filepath.Join(fullpath, file.Name())

		if file.IsDir() {
			dirnames = append(dirnames, name)
		} else {
			filenames = append(filenames, name)
		}

	}

	return

}
