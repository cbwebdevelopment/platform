package mongo

import (
	"time"

	"github.com/tidepool-org/platform/app"
	"github.com/tidepool-org/platform/errors"
)

type Config struct {
	Addresses  string         `json:"addresses"` // TODO: This should be an array, but configor does not support that. Bleech! Fix?
	Database   string         `json:"database"`
	Collection string         `json:"collection"`
	Username   *string        `json:"username"`
	Password   *string        `json:"password"`
	Timeout    *time.Duration `json:"timeout"`
	SSL        bool           `json:"ssl"`
}

func (c *Config) Validate() error {
	addresses := app.SplitStringAndRemoveWhitespace(c.Addresses, ",")
	if len(addresses) < 1 {
		return errors.New("mongo", "addresses is missing")
	}
	if c.Database == "" {
		return errors.New("mongo", "database is missing")
	}
	if c.Collection == "" {
		return errors.New("mongo", "collection is missing")
	}
	if c.Username != nil && *c.Username == "" {
		return errors.New("mongo", "username is empty")
	}
	if c.Password != nil && *c.Password == "" {
		return errors.New("mongo", "password is empty")
	}
	if c.Timeout != nil && *c.Timeout < 0 {
		return errors.New("mongo", "timeout is invalid")
	}
	return nil
}

func (c *Config) Clone() *Config {
	clone := &Config{
		Addresses:  c.Addresses,
		Database:   c.Database,
		Collection: c.Collection,
		SSL:        c.SSL,
	}
	if c.Username != nil {
		clone.Username = app.StringAsPointer(*c.Username)
	}
	if c.Password != nil {
		clone.Password = app.StringAsPointer(*c.Password)
	}
	if c.Timeout != nil {
		clone.Timeout = app.DurationAsPointer(*c.Timeout)
	}
	return clone
}
