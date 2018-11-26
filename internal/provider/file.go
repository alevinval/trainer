package provider

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/alevinval/trainer/internal/adapter"
	"github.com/alevinval/trainer/internal/trainer"
)

type (
	adapterConstructor func(data []byte) (provider trainer.ActivityProvider, err error)
)

var (
	// ErrExtensionNotSupported is returned when the file extension is not supported
	ErrExtensionNotSupported = errors.New("unrecognized file extension")

	extToAdapter = map[string]adapterConstructor{
		".gpx": adapter.Gpx,
		".fit": adapter.Fit,
	}
)

// File reads a file content and returns an ActivityProvider
func File(filePath string) (provider trainer.ActivityProvider, err error) {
	ext, isGzip := isGzip(filePath)
	if !isExtSupported(ext) {
		return nil, ErrExtensionNotSupported
	}
	r, err := getReaderForFile(filePath, isGzip)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	provider, err = extToAdapter[ext](data)
	if err != nil {
		return nil, err
	}
	provider.Metadata().DataSource = trainer.DataSource{
		Type: trainer.FileDataSource,
		Name: filePath,
	}
	return provider, nil
}

// isGzip returns the extension of a file and whether its zipped or not.
func isGzip(name string) (ext string, isGzip bool) {
	name = strings.ToLower(name)
	if strings.HasSuffix(name, ".gz") {
		nameWithoutGz := strings.TrimSuffix(name, ".gz")
		ext := path.Ext(nameWithoutGz)
		return ext, true
	}
	return path.Ext(name), false
}

// isExtSupported returns if an adapter exists for the given extension.
func isExtSupported(ext string) (supported bool) {
	_, ok := extToAdapter[ext]
	return ok
}

// getReaderForFile returns a reader for a file, supports zipped files.
func getReaderForFile(name string, isGzip bool) (r io.Reader, err error) {
	r, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	if !isGzip {
		return r, err
	}
	return gzip.NewReader(r)
}
