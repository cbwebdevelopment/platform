package auth

import (
	"time"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/log"
)

type Client interface {
	ServerToken() (string, error)

	ValidateToken(ctx Context, token string) (Details, error)

	GetStatus(ctx Context) (*Status, error)
}

type Context interface {
	Logger() log.Logger
	Request() *rest.Request

	AuthClient() Client
	AuthDetails() Details
}

type Details interface {
	Token() string

	IsServer() bool
	UserID() string
}

type Use struct {
	OAuth *struct {
		Provider string `json:"provider,omitempty"`
	} `json:"oauth,omitempty"`
}

type RestrictedTokenRequest struct {
	Use      Use  `json:"use,omitempty"`
	Duration *int `json:"duration,omitempty"`
}

type RestrictedToken struct {
	Use            Use       `json:"use,omitempty" bson:"use,omitempty"`
	ExpirationTime time.Time `json:"duration,omitempty" bson:"duration,omitempty"`

	CreatedTime    time.Time `json:"createdTime,omitempty" bson:"createdTime,omitempty"`
	CreatedUserID  string    `json:"createdUserID,omitempty" bson:"createdUserID,omitempty"`
	ModifiedTime   time.Time `json:"modifiedTime,omitempty" bson:"modifiedTime,omitempty"`
	ModifiedUserID string    `json:"modifiedUserID,omitempty" bson:"modifiedUserID,omitempty"`
}

type Status struct {
	Version   string
	Server    interface{}
	AuthStore interface{}
}

const TidepoolAuthTokenHeaderName = "X-Tidepool-Session-Token"
