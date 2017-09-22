package mongo

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/tidepool-org/platform/errors"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/store/mongo"
	"github.com/tidepool-org/platform/user"
	"github.com/tidepool-org/platform/user/store"
)

func New(logger log.Logger, config *Config) (*Store, error) {
	if config == nil {
		return nil, errors.New("mongo", "config is missing")
	}

	baseStore, err := mongo.New(logger, config.Config)
	if err != nil {
		return nil, err
	}

	config = config.Clone()
	if err = config.Validate(); err != nil {
		return nil, errors.Wrap(err, "mongo", "config is invalid")
	}

	return &Store{
		Store:  baseStore,
		config: config,
	}, nil
}

type Store struct {
	*mongo.Store
	config *Config
}

func (s *Store) NewSession(logger log.Logger) store.Session {
	return &Session{
		Session: s.Store.NewSession(logger),
		config:  s.config,
	}
}

type Session struct {
	*mongo.Session
	config *Config
}

func (s *Session) GetUserByID(userID string) (*user.User, error) {
	if userID == "" {
		return nil, errors.New("mongo", "user id is missing")
	}

	if s.IsClosed() {
		return nil, errors.New("mongo", "session closed")
	}

	startTime := time.Now()

	users := []*user.User{}
	selector := bson.M{
		"userid": userID,
	}
	err := s.C().Find(selector).Limit(2).All(&users)

	loggerFields := log.Fields{"userId": userID, "duration": time.Since(startTime) / time.Microsecond}
	s.Logger().WithFields(loggerFields).WithError(err).Debug("GetUserByID")

	if err != nil {
		return nil, errors.Wrap(err, "mongo", "unable to get user by id")
	}

	if usersCount := len(users); usersCount == 0 {
		return nil, nil
	} else if usersCount > 1 {
		s.Logger().WithField("userId", userID).Warn("Multiple users found for user id")
	}

	user := users[0]

	if meta, ok := user.Private["meta"]; ok && meta.ID != "" {
		user.ProfileID = &meta.ID
	}

	return user, nil
}

func (s *Session) DeleteUser(user *user.User) error {
	if user == nil {
		return errors.New("mongo", "user is missing")
	}
	if user.ID == "" {
		return errors.New("mongo", "user id is missing")
	}

	if s.IsClosed() {
		return errors.New("mongo", "session closed")
	}

	startTime := time.Now()

	user.DeletedTime = s.Timestamp()
	user.DeletedUserID = s.AgentUserID()

	selector := bson.M{
		"userid": user.ID,
	}
	err := s.C().Update(selector, user)

	loggerFields := log.Fields{"userId": user.ID, "duration": time.Since(startTime) / time.Microsecond}
	s.Logger().WithFields(loggerFields).WithError(err).Debug("DeleteUser")

	if err != nil {
		return errors.Wrap(err, "mongo", "unable to delete user")
	}
	return nil
}

func (s *Session) CreateUser(details *user.UserCreate) (*user.User, error) {
	if details == nil {
		return nil, errors.New("mongo", "user is missing")
	}
	if details.Email == "" {
		return nil, errors.New("mongo", "user email is missing")
	}
	if details.Password == "" {
		return nil, errors.New("mongo", "user password is missing")
	}
	if details.Emails == nil {
		details.Emails = append(details.Emails, details.Email)
	}

	var newUser = &user.User{Email: details.Email, Emails: details.Emails, Roles: details.Roles}

	var err error

	if newUser.ID, err = s.generateUniqueHash([]string{newUser.Email, details.Password}, 10); err != nil {
		return nil, errors.New("hash", "error generating id")
	}
	if newUser.Hash, err = s.generateUniqueHash([]string{newUser.Email, details.Password, newUser.ID}, 24); err != nil {
		return nil, errors.New("hash", "error generating hash")
	}

	newUser.PasswordHash = s.HashPassword(newUser.ID, details.Password)
	newUser.EmailVerified = true

	if s.IsClosed() {
		return nil, errors.New("mongo", "session closed")
	}

	startTime := time.Now()

	err = s.C().Insert(newUser)

	createdUser, err := s.GetUserByID(newUser.ID)
	loggerFields := log.Fields{"userId": createdUser.ID, "duration": time.Since(startTime) / time.Microsecond}
	s.Logger().WithFields(loggerFields).WithError(err).Debug("CreateUser")

	if err != nil {
		return nil, errors.Wrap(err, "mongo", "unable to create user")
	}
	return createdUser, nil
}

func (s *Session) DestroyUserByID(userID string) error {
	if userID == "" {
		return errors.New("mongo", "user id is missing")
	}

	if s.IsClosed() {
		return errors.New("mongo", "session closed")
	}

	startTime := time.Now()

	selector := bson.M{
		"userid": userID,
	}
	err := s.C().Remove(selector)

	loggerFields := log.Fields{"userId": userID, "duration": time.Since(startTime) / time.Microsecond}
	s.Logger().WithFields(loggerFields).WithError(err).Debug("DestroyUserByID")

	if err != nil {
		return errors.Wrap(err, "mongo", "unable to destroy user by id")
	}
	return nil
}

// TODO: This really isn't the right place for this, but we shouldn't be using a
// password hash algorithm with an external salt, but instead something like bcrypt

// TODO: We should use a constant-time password matching algorithm

func (s *Session) PasswordMatches(user *user.User, password string) bool {
	return user.PasswordHash == s.HashPassword(user.ID, password)
}

// TODO: Do away with external salt and use hash algorithm with internal salt (eg. bcrypt/scrypt)

func (s *Session) HashPassword(userID string, password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	hash.Write([]byte(s.config.PasswordSalt))
	hash.Write([]byte(userID))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *Session) generateUniqueHash(strings []string, length int) (string, error) {
	if len(strings) > 0 && length > 0 {
		hash := sha256.New()

		for i := range strings {
			hash.Write([]byte(strings[i]))
		}

		max := big.NewInt(9999999999)
		//add some randomness
		n, err := rand.Int(rand.Reader, max)

		if err != nil {
			return "", err
		}
		hash.Write([]byte(n.String()))
		//and use unix nano
		hash.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10)))
		hashString := hex.EncodeToString(hash.Sum(nil))
		return string([]rune(hashString)[0:length]), nil
	}

	return "", errors.New("hash", "both strings and length are required")
}
