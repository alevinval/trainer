package provider

import (
	"io"
	"io/ioutil"
	"path"
	"testing"

	"github.com/alevinval/trainer/internal/testutil"
	"github.com/alevinval/trainer/internal/trainer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenFile(t *testing.T) {
	tmp := testutil.NewTemp()
	defer tmp.Remove()

	sampleGpx, err := ioutil.ReadFile("../adapter/testdata/sample.gpx")
	require.Nil(t, err)

	sampleFit, err := ioutil.ReadFile("../adapter/testdata/sample.fit")
	require.Nil(t, err)

	for _, test := range []struct {
		fileName      string
		data          string
		err           error
		dataPointsLen int
	}{
		{"wrong-extension-1", "", ErrExtensionNotSupported, 0},
		{"wrong-extension-2.tmp", "", ErrExtensionNotSupported, 0},
		{"wrong-extension-3.gpx.txt", "", ErrExtensionNotSupported, 0},

		{"right-extension-no-data.gpx", "", io.EOF, 0},
		{"right-extension-invalid-data.gpx", "invalid blob", io.EOF, 0},

		{"right-extension-valid-data.gpx", "<xml></xml>", nil, 0},

		{"right-file-1.gpx", string(sampleGpx), nil, 3},
		{"right-file-2-compressed.gpx.gz", string(sampleGpx), nil, 3},
		{"right-file-1.fit", string(sampleFit), nil, 260},
		{"right-file-2-compressed.fit.gz", string(sampleFit), nil, 260},
	} {
		data := []byte(test.data)
		var filePath string
		if path.Ext(test.fileName) == ".gz" {
			filePath = tmp.CreateGzip(test.fileName, data)
		} else {
			filePath = tmp.Create(test.fileName, data)
		}

		activity, err := File(filePath)
		require.Equal(t, test.err, err)

		// Assert data source is populated correctly
		if test.err == nil {
			fileDataSource := trainer.DataSource{
				Type: trainer.FileDataSource,
				Name: filePath,
			}
			assert.Equal(t, fileDataSource, activity.Metadata().DataSource)
			assert.Equal(t, test.dataPointsLen, len(activity.DataPoints()))
		}
	}
}

func TestOpenFileWithMissingFile(t *testing.T) {
	_, err := File("some-path.gpx")
	assert.NotNil(t, err)
}
