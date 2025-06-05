package main

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestProcessFile_SingleFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "TestFile.txt")
	err := os.WriteFile(filePath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ProcessFile(filePath, &wg, errCh)
	}()
	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			t.Fatalf("unexpected error: %v", e)
		}
	}

	newFilePath := filepath.Join(tempDir, "test-file.txt")
	if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist, but it does not", newFilePath)
	}
}

func TestProcessFile_DirectoryWithFiles(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "TestDir")
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("failed to create test dir: %v", err)
	}
	filePath1 := filepath.Join(dirPath, "FileOne.txt")
	filePath2 := filepath.Join(dirPath, "FileTwo.txt")
	os.WriteFile(filePath1, []byte("content one"), 0644)
	os.WriteFile(filePath2, []byte("content two"), 0644)

	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ProcessFile(dirPath, &wg, errCh)
	}()
	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			t.Fatalf("unexpected error: %v", e)
		}
	}

	newDirPath := filepath.Join(tempDir, "test-dir")
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist, but it does not", newDirPath)
	}

	newFilePath1 := filepath.Join(newDirPath, "file-one.txt")
	newFilePath2 := filepath.Join(newDirPath, "file-two.txt")
	if _, err := os.Stat(newFilePath1); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist, but it does not", newFilePath1)
	}
	if _, err := os.Stat(newFilePath2); os.IsNotExist(err) {
		t.Errorf("expected file %s to exist, but it does not", newFilePath2)
	}
}

func TestProcessFile_EmptyDir(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "EmptyDir")
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("failed to create empty dir: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 10)
	wg.Add(1)
	go func() {
		defer wg.Done()
		ProcessFile(dirPath, &wg, errCh)
	}()
	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			t.Fatalf("unexpected error: %v", e)
		}
	}

	newDirPath := filepath.Join(tempDir, "empty-dir")
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		t.Errorf("expected directory %s to exist, but it does not", newDirPath)
	}
}
