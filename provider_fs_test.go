package trainer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileExt(t *testing.T) {
	for _, test := range []struct {
		input    string
		expected fileExt
	}{
		{"some.gpx", extGpx},
		{"file.tmp.gpx", extGpx},
		{"file.gpx.tmp", extUnknown},
		{"activity.txt", extUnknown},
		{"", extUnknown},
	} {
		actual := getFileExt(test.input)
		assert.Equal(t, test.expected, actual)
	}
}
