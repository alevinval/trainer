package trainer

import (
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	extGpx fileExt = ".gpx"
	extFit fileExt = ".fit"
)

var (
	// ErrUnknownExtension is returned when the file extension is not recognized.
	ErrUnknownExtension = errors.New("unrecognized file extension")
)

type (
	fileExt string
)

// OpenFile reads a file content and returns an Activity.
func OpenFile(fileName string) (a *Activity, err error) {
	ext, isGzip, err := getFileExt(fileName)
	if err != nil {
		return nil, err
	}
	data, err := getFileContents(fileName, isGzip)
	if err != nil {
		return nil, err
	}
	provider, err := getFileActivityProvider(ext, data)
	if err != nil {
		return nil, err
	}
	return buildActivityFromFile(fileName, provider, data), nil
}

func getFileContents(name string, isGzip bool) (b []byte, err error) {
	var src io.Reader

	src, err = os.Open(name)
	if err != nil {
		return nil, err
	}
	if isGzip {
		src, err = gzip.NewReader(src)
		if err != nil {
			return nil, err
		}
	}
	return ioutil.ReadAll(src)
}

func getFileExt(fileName string) (ext fileExt, isGzip bool, err error) {
	extStr := path.Ext(fileName)
	if extStr == ".gz" {
		extStr = path.Ext(strings.TrimSuffix(fileName, ".gz"))
		isGzip = true
	}
	switch fileExt(extStr) {
	case extGpx:
		ext = extGpx
	case extFit:
		ext = extFit
	default:
		err = ErrUnknownExtension
	}
	return
}

func getFileActivityProvider(ext fileExt, data []byte) (p activityProvider, err error) {
	switch ext {
	case extGpx:
		p, err = newGpxAdapter(data)
	case extFit:
		p, err = newFitAdapter(data)
	}
	return
}

func buildActivityFromFile(fileName string, provider activityProvider, data []byte) *Activity {
	metadata := provider.Metadata()
	metadata.DataSource = newDataSource(FileDataSource, fileName)
	datapoints := provider.DataPoints()
	datapoints.process()
	return &Activity{
		rawData:    data,
		metadata:   metadata,
		datapoints: datapoints,
	}
}
