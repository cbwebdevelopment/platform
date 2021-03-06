package deduplicator

import (
	"github.com/tidepool-org/platform/data"
	"github.com/tidepool-org/platform/data/store"
	"github.com/tidepool-org/platform/data/types/upload"
	"github.com/tidepool-org/platform/log"
)

type Factory interface {
	CanDeduplicateDataset(dataset *upload.Upload) (bool, error)
	NewDeduplicatorForDataset(logger log.Logger, dataStoreSession store.Session, dataset *upload.Upload) (data.Deduplicator, error)

	IsRegisteredWithDataset(dataset *upload.Upload) (bool, error)
	NewRegisteredDeduplicatorForDataset(logger log.Logger, dataStoreSession store.Session, dataset *upload.Upload) (data.Deduplicator, error)
}
