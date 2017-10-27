package trainer

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTmpFile(t *testing.T, path string, data []byte) {
	err := ioutil.WriteFile(path, data, 0644)
	if err != nil {
		t.Fatalf("error writing file: %s", err)
	}
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
	} {
		filePath := path.Join(tmpDir, test.fileName)
		createTmpFile(t, filePath, []byte(test.data))

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
