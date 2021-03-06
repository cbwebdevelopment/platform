package mongo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"fmt"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/message/store"
	"github.com/tidepool-org/platform/message/store/mongo"
	baseMongo "github.com/tidepool-org/platform/store/mongo"
	testMongo "github.com/tidepool-org/platform/test/mongo"
)

func NewMessage(groupID string, userID string) bson.M {
	return bson.M{
		"groupid":      groupID,
		"userid":       userID,
		"guid":         app.NewID(),
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
		"createdtime":  time.Now().UTC().Format(time.RFC3339),
		"modifiedtime": time.Now().UTC().Format(time.RFC3339),
		"messagetext":  "test",
	}
}

func NewMessages(groupID string, userID string) []interface{} {
	messages := []interface{}{}
	for count := 0; count < 2; count++ {
		messages = append(messages, NewMessage(groupID, userID))
	}
	parentMessage := NewMessage(groupID, userID)
	messages = append(messages, parentMessage)
	for count := 0; count < 2; count++ {
		message := NewMessage(groupID, userID)
		message["parentmessage"] = parentMessage["guid"]
		messages = append(messages, message)
	}
	return messages
}

func MarkMessagesDeleted(messages []interface{}) {
	for index, message := range messages {
		messages[index] = MarkMessageDeleted(message.(bson.M))
	}
}

func MarkMessageDeleted(message bson.M) bson.M {
	message["user"] = bson.M{
		"fullName": fmt.Sprintf("deleted user (%s)", message["userid"]),
	}
	delete(message, "userid")
	return message
}

func ValidateMessages(testMongoCollection *mgo.Collection, selector bson.M, expectedMessages []interface{}) {
	var actualMessages []interface{}
	Expect(testMongoCollection.Find(selector).Select(bson.M{"_id": 0}).All(&actualMessages)).To(Succeed())
	Expect(actualMessages).To(ConsistOf(expectedMessages...))
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
				var messages []interface{}

				BeforeEach(func() {
					testMongoSession = testMongo.Session().Copy()
					testMongoCollection = testMongoSession.DB(mongoConfig.Database).C(mongoConfig.Collection)
					messages = append(NewMessages(app.NewID(), app.NewID()), NewMessages(app.NewID(), app.NewID())...)
				})

				JustBeforeEach(func() {
					Expect(testMongoCollection.Insert(messages...)).To(Succeed())
				})

				AfterEach(func() {
					if testMongoSession != nil {
						testMongoSession.Close()
					}
				})

				Context("DeleteMessagesFromUser", func() {
					var deleteGroupID string
					var deleteUserID string
					var deleteUser *store.User
					var deleteMessages []interface{}

					BeforeEach(func() {
						deleteGroupID = app.NewID()
						deleteUserID = app.NewID()
						deleteUser = &store.User{
							ID:       deleteUserID,
							FullName: fmt.Sprintf("deleted user (%s)", deleteUserID),
						}
						deleteMessages = NewMessages(deleteGroupID, deleteUserID)
						messages = append(messages, NewMessages(deleteUserID, deleteGroupID)...)
					})

					JustBeforeEach(func() {
						Expect(testMongoCollection.Insert(deleteMessages...)).To(Succeed())
					})

					It("succeeds if it successfully removes messages", func() {
						Expect(mongoSession.DeleteMessagesFromUser(deleteUser)).To(Succeed())
					})

					It("returns an error if the user is missing", func() {
						Expect(mongoSession.DeleteMessagesFromUser(nil)).To(MatchError("mongo: user is missing"))
					})

					It("returns an error if the user id is missing", func() {
						deleteUser.ID = ""
						Expect(mongoSession.DeleteMessagesFromUser(deleteUser)).To(MatchError("mongo: user id is missing"))
					})

					It("returns an error if the session is closed", func() {
						mongoSession.Close()
						Expect(mongoSession.DeleteMessagesFromUser(deleteUser)).To(MatchError("mongo: session closed"))
					})

					It("has the correct stored messages", func() {
						ValidateMessages(testMongoCollection, bson.M{}, append(messages, deleteMessages...))
						Expect(mongoSession.DeleteMessagesFromUser(deleteUser)).To(Succeed())
						MarkMessagesDeleted(deleteMessages)
						ValidateMessages(testMongoCollection, bson.M{}, append(messages, deleteMessages...))
					})
				})

				Context("DestroyMessagesForUserByID", func() {
					var destroyGroupID string
					var destroyUserID string
					var destroyMessages []interface{}

					BeforeEach(func() {
						destroyGroupID = app.NewID()
						destroyUserID = app.NewID()
						destroyMessages = NewMessages(destroyGroupID, destroyUserID)
						messages = append(messages, NewMessages(destroyUserID, destroyGroupID)...)
					})

					JustBeforeEach(func() {
						Expect(testMongoCollection.Insert(destroyMessages...)).To(Succeed())
					})

					It("succeeds if it successfully removes messages", func() {
						Expect(mongoSession.DestroyMessagesForUserByID(destroyGroupID)).To(Succeed())
					})

					It("returns an error if the user id is missing", func() {
						Expect(mongoSession.DestroyMessagesForUserByID("")).To(MatchError("mongo: user id is missing"))
					})

					It("returns an error if the session is closed", func() {
						mongoSession.Close()
						Expect(mongoSession.DestroyMessagesForUserByID(destroyGroupID)).To(MatchError("mongo: session closed"))
					})

					It("has the correct stored messages", func() {
						ValidateMessages(testMongoCollection, bson.M{}, append(messages, destroyMessages...))
						Expect(mongoSession.DestroyMessagesForUserByID(destroyGroupID)).To(Succeed())
						ValidateMessages(testMongoCollection, bson.M{}, messages)
					})
				})
			})
		})
	})
})
