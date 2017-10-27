package trainer

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
)

const (
	extGpx fileExt = ".gpx"
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
	ext, err := getFileExt(fileName)
	if err != nil {
		return nil, err
	}
	data, err := getFileContents(fileName)
	if err != nil {
		return nil, err
	}
	provider, err := getFileActivityProvider(ext, data)
	if err != nil {
		return nil, err
	}
	return buildActivityFromFile(fileName, provider, data), nil
}

func getFileContents(name string) (b []byte, err error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func getFileExt(fileName string) (ext fileExt, err error) {
	switch fileExt(path.Ext(fileName)) {
	case extGpx:
		ext = extGpx
	default:
		err = ErrUnknownExtension
	}
	return
}

func getFileActivityProvider(ext fileExt, data []byte) (p activityProvider, err error) {
	switch ext {
	case extGpx:
		p, err = newGpx(data)
	}
	return
}

func buildActivityFromFile(fileName string, provider activityProvider, data []byte) *Activity {
	metadata := provider.Metadata()
	metadata.DataSource = newDataSource(FileDataSource, fileName)
	datapoints := provider.DataPoints()
	datapoints.process()
	return &Activity{
		Data:       data,
		metadata:   metadata,
		datapoints: datapoints,
	}
}
