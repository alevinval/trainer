package home

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestSyncNothingToSynchronise(t *testing.T) {
	srcPath, dstPath := createSyncFolders()
	defer cleanUpFolders(srcPath, dstPath)

	Sync(dstPath)
}

func TestSyncOneFileToSynchronise(t *testing.T) {
	srcPath, dstPath := createSyncFolders()
	defer cleanUpFolders(srcPath, dstPath)

	addMockFile(dstPath, "file-1.FIT")

	Sync(dstPath)

	assertFileInPath(t, srcPath, "file-1.FIT")
}

func assertFileInPath(t *testing.T, path, filename string) {
	expectedPath := filepath.Join(path, filename)
	_, err := os.Open(expectedPath)

	if err != nil {
		t.Fatalf("Expected synchronised file in: %s, but got error: %s\n", expectedPath, err)
	}
}

func createSyncFolders() (srcPath, dstPath string) {
	rootPath, err := ioutil.TempDir("", "tmp-sync-test")
	if err != nil {
		log.Fatalf("error creating temp dir: %s\n", err)
	}

	srcPath = filepath.Join(rootPath, "src")
	dstPath = filepath.Join(rootPath, "dst")
	os.MkdirAll(srcPath, os.ModePerm)
	os.MkdirAll(dstPath, os.ModePerm)

	// Hijack original home activities path
	ActivitiesPath = srcPath
	return
}

func cleanUpFolders(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
}

func addMockFile(rootPath, filename string) {
	srcPath := filepath.Join("..", "adapter", "testdata", "sample.fit")
	dstPath := filepath.Join(rootPath, filename)

	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("error opening srcFile: %s\n", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("error reating dstFile: %s\n", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatalf("error copying file: %s\n", err)
	}
}
