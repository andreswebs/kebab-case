package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// ProcessFile
// recursively renames files and directories using the defined format
func ProcessFile(name string, wg *sync.WaitGroup, errCh chan<- error) {
	var err error

	fullpath, err := filepath.Abs(name)
	if err != nil {
		errCh <- err
		return
	}

	newname := getNewname(fullpath)
	err = rename(fullpath, newname)
	if err != nil {
		errCh <- err
		return
	}

	dirYes, err := isDir(newname)
	if err != nil {
		errCh <- err
		return
	}

	if !dirYes {
		return
	}

	filenames, dirnames, err := sift(newname)
	if err != nil {
		errCh <- err
		return
	}

	// rename files
	for _, n := range filenames {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			ProcessFile(n, wg, errCh)
		}(n)
	}

	// rename dirs
	for _, d := range dirnames {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			ProcessFile(d, wg, errCh)
		}(d)
	}

	return
}

// rename moves a file to new name if new name is different from previous name
func rename(prevname, newname string) (err error) {
	if newname == prevname {
		return
	}
	err = os.Rename(prevname, newname)
	return
}

// getNewname gets new full file name renamed using the defined format
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

// getBasename gets file name without extension
func getBasename(filename string) string {
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

// isDir checks if a file represented by `path` is a directory
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
