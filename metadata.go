package trainer

import "time"

const (
	FileSource Source = "file"
)

type (
	Metadata struct {
		Name       string
		Time       time.Time
		Source     Source
		SourceName string
	}

	Source string
)
