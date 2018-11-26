package provider

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"strings"
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

	for _, tt := range []struct {
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

		{"upper-case-ext.GPX", string(sampleGpx), nil, 3},
		{"upper-case-ext.GPX.GZ", string(sampleGpx), nil, 3},
		{"upper-case-ext.FIT", string(sampleFit), nil, 260},
		{"upper-case-ext.FIT.GZ", string(sampleFit), nil, 260},
	} {
		t.Run(fmt.Sprintf("Provide from file %s", tt.fileName), func(t *testing.T) {
			data := []byte(tt.data)
			var filePath string
			if strings.ToLower(path.Ext(tt.fileName)) == ".gz" {
				filePath = tmp.CreateGzip(tt.fileName, data)
			} else {
				filePath = tmp.Create(tt.fileName, data)
			}

			activity, err := File(filePath)
			require.Equal(t, tt.err, err)

			// Assert data source is populated correctly
			if tt.err == nil {
				fileDataSource := trainer.DataSource{
					Type: trainer.FileDataSource,
					Name: filePath,
				}
				assert.Equal(t, fileDataSource, activity.Metadata().DataSource)
				assert.Equal(t, tt.dataPointsLen, len(activity.DataPoints()))
			}
		})
	}
}

func TestOpenFileWithMissingFile(t *testing.T) {
	_, err := File("some-path.gpx")
	assert.NotNil(t, err)
}
