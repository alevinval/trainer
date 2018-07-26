package trainer

import "time"

type (
	// Metadata is used to describe an Activity
	Metadata struct {
		Name       string
		Time       time.Time
		DataSource DataSource
	}

	// DataSource is used to describe the origin of a resource
	DataSource struct {
		Type DataSourceType
		Name string
	}

	// DataSourceType is used to describe the type of a DataSource.
	// Resources may come from files, databases, APIs, etc...
	DataSourceType string
)

const (
	// FileDataSource indicates a resource originates from a file.
	FileDataSource DataSourceType = "file"
)

func newDataSource(dataSourceType DataSourceType, name string) DataSource {
	return DataSource{Type: dataSourceType, Name: name}
}
