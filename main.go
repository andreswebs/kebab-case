package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// run executes all tasks and returns an OS exit code and error
func run() (code int, err error) {

	if len(os.Args) < 2 {
		return 1, errors.New("requires an argument")
	}

	start := os.Args[1]
	
	err = processFile(start)
	if err != nil {
		return 1, err
	}

	return

}

func main() {
	code, err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
	os.Exit(code)
}

func processFile(name string) (err error) {

	fullpath, err := filepath.Abs(name)
	if err != nil {
		return
	}

	newName, err := rename(fullpath)
	if err != nil {
		return
	}

	dir, err := isDir(newName)
	if err != nil {
		return
	}

	if !dir {
		return
	}

	filenames, dirnames, err := sift(newName)
	if err != nil {
		return
	}

	for _, n := range filenames {
		_, err = rename(n)
		if err != nil {
			return err
		}
	}

	var newDirnames []string

	for _, d := range dirnames {
		newD, err := rename(d)
		if err != nil {
			return err
		}
		newDirnames = append(newDirnames, newD)
	}

	for _, d := range newDirnames {
		processFile(d)
	}

	return nil

}

func rename(prevName string) (newName string, err error) {
	newName = getNewPath(prevName)
	if newName != prevName {
		return newName, os.Rename(prevName, newName)
	}
	return
}

func getNewPath(path string) string {

	// get the dir
	dir, fileName := filepath.Split(path)

	// get the ext
	ext := filepath.Ext(path)

	// get the base
	baseName := getBasename(fileName)

	newName := format(baseName) + ext
	newPath := filepath.Join(dir, newName)

	return newPath

}

func getBasename(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func format(s string) string {

	camel := regexp.MustCompile(`([a-z])([A-Z])`)
	spaces := regexp.MustCompile(`(\s+-\s+|\s+)`)
	quotes := regexp.MustCompile(`('|")`)
	parentheses := regexp.MustCompile(`(\(|\)|{|})`)
	accentsA := regexp.MustCompile(`(ã|á|à)`)
	accentsE := regexp.MustCompile(`(é|è)`)
	accentsI := regexp.MustCompile(`(í|ì)`)
	accentsO := regexp.MustCompile(`(ó|ò)`)
	accentsU := regexp.MustCompile(`(ú|ù)`)
	accentsN := regexp.MustCompile(`(ñ)`)
	accentsC := regexp.MustCompile(`(ç|ć)`)
	other := regexp.MustCompile(`(,|;|:|<|>|\?|!|@|#|\$|%|\^|\*|\+|=|~)`)
	multidashes := regexp.MustCompile(`(--+|__+)`)
	trailingDashes := regexp.MustCompile(`(^-|-$)`)

	s = strings.ToValidUTF8(s, "")
	s = strings.TrimSpace(s)
	s = quotes.ReplaceAllString(s, "")
	s = accentsA.ReplaceAllString(s, "a")
	s = accentsE.ReplaceAllString(s, "e")
	s = accentsI.ReplaceAllString(s, "i")
	s = accentsO.ReplaceAllString(s, "o")
	s = accentsU.ReplaceAllString(s, "u")
	s = accentsN.ReplaceAllString(s, "n")
	s = accentsC.ReplaceAllString(s, "c")
	s = parentheses.ReplaceAllString(s, "-")
	s = other.ReplaceAllString(s, "-")
	s = trailingDashes.ReplaceAllString(s, "")
	s = camel.ReplaceAllString(s, "$1-$2")
	s = spaces.ReplaceAllString(s, "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, "-.", ".")
	s = strings.ReplaceAll(s, ".-", "-")
	s = strings.ReplaceAll(s, "..", ".")
	s = multidashes.ReplaceAllString(s, "-")
	s = strings.ToLower(s)

	return s

}

// isDir determines if a file represented
// by `path` is a directory or not
func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), err
}

// sift inspects a directory and returns a list of filenames, a list of dirnames and an error
func sift(dir string) (filenames, dirnames []string, err error) {

	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {

		name, err := filepath.Abs(file.Name())
		if err != nil {
			return filenames, dirnames, err
		}

		if file.IsDir() {
			dirnames = append(dirnames, name)
		} else {
			filenames = append(filenames, name)
		}

	}

	return

}
