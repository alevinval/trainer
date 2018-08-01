package testutil

import "github.com/alevinval/trainer/internal/trainer"

type MockActivity struct {
	MockMetadata   *trainer.Metadata
	MockDataPoints trainer.DataPointList
}

func (ma *MockActivity) Metadata() *trainer.Metadata {
	if ma.MockMetadata != nil {
		return ma.MockMetadata
	}
	return &trainer.Metadata{}
}

func (ma *MockActivity) DataPoints() trainer.DataPointList {
	if ma.MockDataPoints != nil {
		return ma.MockDataPoints
	}
	return trainer.DataPointList{}
}
