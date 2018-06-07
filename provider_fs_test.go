package trainer

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTmpFile(t *testing.T, path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("error writing file: %s", err)
	}
}

func createGzipTmpFile(t *testing.T, path string, data []byte) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		t.Fatalf("error writing gzip file: %s", err)
	}
	w.Close()
	createTmpFile(t, path, b.Bytes())
}

func createGpxActivityAsText() string {
	trackpoint := func(lat, lon float64) string {
		return fmt.Sprintf(`<trkpt lat="%f" lon="%f">
		<ele>%f</ele>
		<time>2017-06-19T16:49:35.000Z</time>
		<extensions>
			<ns3:TrackPointExtension>
				<ns3:hr>94</ns3:hr>
				<ns3:cad>85</ns3:cad>
			</ns3:TrackPointExtension>
		</extensions>
		</trkpt>`, lat, lon, 100.0)
	}

	trackpoints := []string{
		trackpoint(1.0, 1.0),
		trackpoint(2.0, 2.0),
		trackpoint(3.0, 3.0),
	}
	trackpointsStr := strings.Join(trackpoints, "")

	return fmt.Sprintf(`<gpx creator="Garmin Connect" version="1.1"
	xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/11.xsd"
	xmlns="http://www.topografix.com/GPX/1/1"
	xmlns:ns3="http://www.garmin.com/xmlschemas/TrackPointExtension/v1"
	xmlns:ns2="http://www.garmin.com/xmlschemas/GpxExtensions/v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<trk>
		<name>Some activity name</name>
		<type>running</type>
		<trkseg>
			%s
		</trkseg>
	</trk>
	</gpx>`, trackpointsStr)
}

func TestOpenFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "trainer-activities-test")
	if err != nil {
		t.Fatalf("error creating temporary directory: %s", err)
	}
	defer os.RemoveAll(tmpDir)

	for _, test := range []struct {
		fileName      string
		data          string
		err           error
		dataPointsLen int
	}{
		{"wrong-extension-1", "", ErrUnknownExtension, 0},
		{"wrong-extension-2.tmp", "", ErrUnknownExtension, 0},
		{"wrong-extension-3.gpx.txt", "", ErrUnknownExtension, 0},

		{"right-extension-no-data.gpx", "", io.EOF, 0},
		{"right-extension-invalid-data.gpx", "invalid blob", io.EOF, 0},

		{"right-extension-valid-data.gpx", "<xml></xml>", nil, 0},

		{"right-file-1.gpx", createGpxActivityAsText(), nil, 3},
		{"right-file-2-compressed.gpx.gz", createGpxActivityAsText(), nil, 3},
	} {
		filePath := path.Join(tmpDir, test.fileName)
		if path.Ext(filePath) == ".gz" {
			createGzipTmpFile(t, filePath, []byte(test.data))
		} else {
			createTmpFile(t, filePath, []byte(test.data))
		}

		activity, err := OpenFile(filePath)
		if !assert.Equal(t, test.err, err) {
			return
		}
		if err == nil {
			metadata := activity.Metadata()
			assert.Equal(t, newDataSource(FileDataSource, filePath), metadata.DataSource)
			assert.Equal(t, test.dataPointsLen, len(activity.DataPoints()))
		}
	}
}

func TestOpenFileWithMissingFile(t *testing.T) {
	_, err := OpenFile("some-path.gpx")
	assert.NotNil(t, err)
}
