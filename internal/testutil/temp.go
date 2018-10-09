package testutil

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Temp holds a temporary directory and provides helper functions
// to create files on it, and remove them.
type Temp struct {
	tmpDir string
}

// NewTemp returns a temp instance with helpers to create test files
// and clean them up after a run.
func NewTemp() *Temp {
	tmpDir, err := ioutil.TempDir("", "tmp-test-files")
	if err != nil {
		log.Fatalf("error creating temporary directory: %s", err)
	}
	return &Temp{
		tmpDir: tmpDir,
	}
}

// Create file with the specified contents in the temporary folder.
func (t *Temp) Create(fname string, data []byte) string {
	fullPath := filepath.Join(t.tmpDir, fname)
	err := ioutil.WriteFile(fullPath, data, 0644)
	if err != nil {
		log.Fatalf("error creating temporary file: %s", err)
	}
	return fullPath
}

// CreateGzip file with the specified contents in the temporary folder.
func (t *Temp) CreateGzip(fname string, data []byte) string {
	var gzipData bytes.Buffer
	w := gzip.NewWriter(&gzipData)
	_, err := w.Write(data)
	if err != nil {
		log.Fatalf("error writing gzip file: %s", err)
	}
	w.Close()
	return t.Create(fname, gzipData.Bytes())
}

// Remove temporary folder and all its contents.
func (t *Temp) Remove() {
	os.RemoveAll(t.tmpDir)
}
