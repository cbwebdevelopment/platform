package mongo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/base64"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/crypto"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/permission/store"
	"github.com/tidepool-org/platform/permission/store/mongo"
	baseMongo "github.com/tidepool-org/platform/store/mongo"
	testMongo "github.com/tidepool-org/platform/test/mongo"
)

func NewPermission(groupID string, userID string) bson.M {
	encryptedGroupID, err := crypto.EncryptWithAES256UsingPassphrase([]byte(groupID), []byte("secret"))
	Expect(err).ToNot(HaveOccurred())

	return bson.M{
		"groupId": base64.StdEncoding.EncodeToString(encryptedGroupID),
		"userId":  userID,
		"permissions": bson.M{
			"upload": bson.M{},
			"view":   bson.M{},
		},
	}
}

func NewPermissions(userID string) []interface{} {
	permissions := []interface{}{}
	permissions = append(permissions, NewPermission(app.NewID(), userID), NewPermission(userID, app.NewID()))
	return permissions
}

func ValidatePermissions(testMongoCollection *mgo.Collection, selector bson.M, expectedPermissions []interface{}) {
	var actualPermissions []interface{}
	Expect(testMongoCollection.Find(selector).Select(bson.M{"_id": 0}).All(&actualPermissions)).To(Succeed())
	Expect(actualPermissions).To(ConsistOf(expectedPermissions...))
}

var _ = Describe("Mongo", func() {
	var mongoConfig *mongo.Config
	var mongoStore *mongo.Store
	var mongoSession store.Session

	BeforeEach(func() {
		mongoConfig = &mongo.Config{
			Config: &baseMongo.Config{
				Addresses:  testMongo.Address(),
				Database:   testMongo.Database(),
				Collection: testMongo.NewCollectionName(),
				Timeout:    app.DurationAsPointer(5 * time.Second),
			},
			Secret: "secret",
		}
	})

	AfterEach(func() {
		if mongoSession != nil {
			mongoSession.Close()
		}
		if mongoStore != nil {
			mongoStore.Close()
		}
	})

	Context("New", func() {
		It("returns an error if logger is missing", func() {
			var err error
			mongoStore, err = mongo.New(nil, mongoConfig)
			Expect(err).To(MatchError("mongo: logger is missing"))
			Expect(mongoStore).To(BeNil())
		})

		It("returns an error if base config is missing", func() {
			var err error
			mongoConfig.Config = nil
			mongoStore, err = mongo.New(log.NewNull(), mongoConfig)
			Expect(err).To(MatchError("mongo: config is missing"))
			Expect(mongoStore).To(BeNil())
		})

		It("returns an error if base config is invalid", func() {
			var err error
			mongoConfig.Config.Addresses = ""
			mongoStore, err = mongo.New(log.NewNull(), mongoConfig)
			Expect(err).To(MatchError("mongo: config is invalid; mongo: addresses is missing"))
			Expect(mongoStore).To(BeNil())
		})

		It("returns an error if config is missing", func() {
			var err error
			mongoStore, err = mongo.New(log.NewNull(), nil)
			Expect(err).To(MatchError("mongo: config is missing"))
			Expect(mongoStore).To(BeNil())
		})

		It("returns an error if config is invalid", func() {
			var err error
			mongoConfig.Secret = ""
			mongoStore, err = mongo.New(log.NewNull(), mongoConfig)
			Expect(err).To(MatchError("mongo: config is invalid; mongo: secret is missing"))
			Expect(mongoStore).To(BeNil())
		})

		It("returns a new store and no error if successful", func() {
			var err error
			mongoStore, err = mongo.New(log.NewNull(), mongoConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(mongoStore).ToNot(BeNil())
		})
	})

	Context("with a new store", func() {
		BeforeEach(func() {
			var err error
			mongoStore, err = mongo.New(log.NewNull(), mongoConfig)
			Expect(err).ToNot(HaveOccurred())
			Expect(mongoStore).ToNot(BeNil())
		})

		Context("NewSession", func() {
			It("returns a new session if no logger specified", func() {
				mongoSession = mongoStore.NewSession(nil)
				Expect(mongoSession).ToNot(BeNil())
				Expect(mongoSession.Logger()).ToNot(BeNil())
			})

			It("returns a new session if logger specified", func() {
				logger := log.NewNull()
				mongoSession = mongoStore.NewSession(logger)
				Expect(mongoSession).ToNot(BeNil())
				Expect(mongoSession.Logger()).To(Equal(logger))
			})
		})

		Context("with a new session", func() {
			BeforeEach(func() {
				mongoSession = mongoStore.NewSession(log.NewNull())
				Expect(mongoSession).ToNot(BeNil())
			})

			Context("with persisted data", func() {
				var testMongoSession *mgo.Session
				var testMongoCollection *mgo.Collection
				var permissions []interface{}

				BeforeEach(func() {
					testMongoSession = testMongo.Session().Copy()
					testMongoCollection = testMongoSession.DB(mongoConfig.Database).C(mongoConfig.Collection)
					permissions = NewPermissions(app.NewID())
				})

				JustBeforeEach(func() {
					Expect(testMongoCollection.Insert(permissions...)).To(Succeed())
				})

				AfterEach(func() {
					if testMongoSession != nil {
						testMongoSession.Close()
					}
				})

				Context("DestroyPermissionsForUserByID", func() {
					var destroyUserID string
					var destroyPermissions []interface{}

					BeforeEach(func() {
						destroyUserID = app.NewID()
						destroyPermissions = NewPermissions(destroyUserID)
					})

					JustBeforeEach(func() {
						Expect(testMongoCollection.Insert(destroyPermissions...)).To(Succeed())
					})

					It("succeeds if it successfully removes permissions", func() {
						Expect(mongoSession.DestroyPermissionsForUserByID(destroyUserID)).To(Succeed())
					})

					It("returns an error if the user id is missing", func() {
						Expect(mongoSession.DestroyPermissionsForUserByID("")).To(MatchError("mongo: user id is missing"))
					})

					It("returns an error if the session is closed", func() {
						mongoSession.Close()
						Expect(mongoSession.DestroyPermissionsForUserByID(destroyUserID)).To(MatchError("mongo: session closed"))
					})

					It("has the correct stored permissions", func() {
						ValidatePermissions(testMongoCollection, bson.M{}, append(permissions, destroyPermissions...))
						Expect(mongoSession.DestroyPermissionsForUserByID(destroyUserID)).To(Succeed())
						ValidatePermissions(testMongoCollection, bson.M{}, permissions)
					})
				})
			})
		})
	})
})
