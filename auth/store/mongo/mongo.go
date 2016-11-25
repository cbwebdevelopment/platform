package mongo

import (
	"github.com/tidepool-org/platform/auth"
	"github.com/tidepool-org/platform/auth/store"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/store/mongo"
)

type Store struct {
	*mongo.Store
}

func New(lgr log.Logger, cfg *mongo.Config) (*Store, error) {
	str, err := mongo.New(lgr, cfg)
	if err != nil {
		return nil, err
	}

	return &Store{
		Store: str,
	}, nil
}

func (s *Store) NewRestrictedTokensSession(lgr log.Logger) store.RestrictedTokensSession {
	return &RestrictedTokensSession{
		Session: s.Store.NewSession(lgr, "restrictedTokens"),
	}
}

type RestrictedTokensSession struct {
	*mongo.Session
}

func (r *RestrictedTokensSession) Create(restrictedToken *auth.RestrictedTokenRequest) (*auth.RestrictedToken, error) {
	// if err := d.validateDataset(dataset); err != nil {
	// 	return err
	// }

	// if d.IsClosed() {
	// 	return errors.New("mongo", "session closed")
	// }

	// startTime := time.Now()

	// dataset.CreatedTime = d.Timestamp()
	// dataset.CreatedUserID = d.AgentUserID()

	// dataset.ByUser = dataset.CreatedUserID

	// // TODO: Consider upsert instead to prevent multiples being created?

	// query := bson.M{
	// 	"_userId":  dataset.UserID,
	// 	"uploadId": dataset.UploadID,
	// 	"type":     dataset.Type,
	// }
	// count, err := d.C().Find(query).Count()
	// if err == nil {
	// 	if count > 0 {
	// 		err = errors.New("mongo", "dataset already exists")
	// 	} else {
	// 		err = d.C().Insert(dataset)
	// 	}
	// }

	// loggerFields := log.Fields{"userId": dataset.UserID, "datasetId": dataset.UploadID, "duration": time.Since(startTime) / time.Microsecond}
	// d.Logger().WithFields(loggerFields).WithError(err).Debug("CreateDataset")

	// if err != nil {
	// 	return errors.Wrap(err, "mongo", "unable to create dataset")
	// }

	return nil, nil
}

func (r *RestrictedTokensSession) Get(id string) (*auth.RestrictedToken, error) {
	return nil, nil
}

func (r *RestrictedTokensSession) Delete(id string) error {
	return nil
}
