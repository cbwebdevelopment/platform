package data

import (
	"strconv"

	"github.com/tidepool-org/platform/errors"
)

type Deduplicator interface {
	Name() string
	Version() string

	RegisterDataset() error

	AddDatasetData(datasetData []Datum) error
	DeduplicateDataset() error

	DeleteDataset() error
}

type DeduplicatorDescriptor struct {
	Name    string `bson:"name,omitempty"`
	Version string `bson:"version,omitempty"`
	Hash    string `bson:"hash,omitempty"`
}

func NewDeduplicatorDescriptor() *DeduplicatorDescriptor {
	return &DeduplicatorDescriptor{}
}

func (d *DeduplicatorDescriptor) IsRegisteredWithAnyDeduplicator() bool {
	return d.Name != ""
}

func (d *DeduplicatorDescriptor) IsRegisteredWithNamedDeduplicator(name string) bool {
	return d.Name == name
}

func (d *DeduplicatorDescriptor) RegisterWithDeduplicator(deduplicator Deduplicator) error {
	if d.Name != "" {
		return errors.Newf("data", "deduplicator descriptor already registered with %s", strconv.Quote(d.Name))
	}
	if d.Version != "" {
		return errors.New("data", "deduplicator descriptor already registered with unknown deduplicator")
	}

	d.Name = deduplicator.Name()
	d.Version = deduplicator.Version()
	return nil
}
