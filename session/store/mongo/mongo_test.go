package mongo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/session/store"
	"github.com/tidepool-org/platform/session/store/mongo"
	baseMongo "github.com/tidepool-org/platform/store/mongo"
	testMongo "github.com/tidepool-org/platform/test/mongo"
)

func NewBaseSession() bson.M {
	now := time.Now()
	return bson.M{
		"_id":       app.NewID(),
		"duration":  86400,
		"expiresAt": now.Add(86400 * time.Second).Unix(),
		"createdAt": now.Unix(),
		"time":      now.Unix(),
	}
}

func NewServerSession() bson.M {
	session := NewBaseSession()
	session["isServer"] = true
	session["serverId"] = app.NewID()
	return session
}

func NewUserSession(userID string) bson.M {
	session := NewBaseSession()
	session["isServer"] = false
	session["userId"] = userID
	return session
}

func NewServerSessions() []interface{} {
	sessions := []interface{}{}
	sessions = append(sessions, NewServerSession(), NewServerSession(), NewServerSession())
	return sessions
}

func NewUserSessions(userID string) []interface{} {
	sessions := []interface{}{}
	sessions = append(sessions, NewUserSession(userID), NewUserSession(userID), NewUserSession(userID))
	return sessions
}

func ValidateSessions(testMongoCollection *mgo.Collection, selector bson.M, expectedSessions []interface{}) {
	var actualSessions []interface{}
	Expect(testMongoCollection.Find(selector).All(&actualSessions)).To(Succeed())
	Expect(actualSessions).To(ConsistOf(expectedSessions...))
}

var _ = Describe("Mongo", func() {
	var mongoConfig *baseMongo.Config
	var mongoStore *mongo.Store
	var mongoSession store.Session

	BeforeEach(func() {
		mongoConfig = &baseMongo.Config{
			Addresses:  testMongo.Address(),
			Database:   testMongo.Database(),
			Collection: testMongo.NewCollectionName(),
			Timeout:    app.DurationAsPointer(5 * time.Second),
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
		It("returns an error if unsuccessful", func() {
			var err error
			mongoStore, err = mongo.New(nil, nil)
			Expect(err).To(HaveOccurred())
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
				var sessions []interface{}

				BeforeEach(func() {
					testMongoSession = testMongo.Session().Copy()
					testMongoCollection = testMongoSession.DB(mongoConfig.Database).C(mongoConfig.Collection)
					sessions = append(NewServerSessions(), NewUserSessions(app.NewID())...)
				})

				JustBeforeEach(func() {
					Expect(testMongoCollection.Insert(sessions...)).To(Succeed())
				})

				AfterEach(func() {
					if testMongoSession != nil {
						testMongoSession.Close()
					}
				})

				Context("DestroySessionsForUserByID", func() {
					var destroyUserID string
					var destroySessions []interface{}

					BeforeEach(func() {
						destroyUserID = app.NewID()
						destroySessions = NewUserSessions(destroyUserID)
					})

					JustBeforeEach(func() {
						Expect(testMongoCollection.Insert(destroySessions...)).To(Succeed())
					})

					It("succeeds if it successfully removes sessions", func() {
						Expect(mongoSession.DestroySessionsForUserByID(destroyUserID)).To(Succeed())
					})

					It("returns an error if the user id is missing", func() {
						Expect(mongoSession.DestroySessionsForUserByID("")).To(MatchError("mongo: user id is missing"))
					})

					It("returns an error if the session is closed", func() {
						mongoSession.Close()
						Expect(mongoSession.DestroySessionsForUserByID(destroyUserID)).To(MatchError("mongo: session closed"))
					})

					It("has the correct stored sessions", func() {
						ValidateSessions(testMongoCollection, bson.M{}, append(sessions, destroySessions...))
						Expect(mongoSession.DestroySessionsForUserByID(destroyUserID)).To(Succeed())
						ValidateSessions(testMongoCollection, bson.M{}, sessions)
					})
				})
			})
		})
	})
})
