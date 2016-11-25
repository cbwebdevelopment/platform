package store

import (
	"github.com/tidepool-org/platform/auth"
	"github.com/tidepool-org/platform/log"
	"github.com/tidepool-org/platform/store"
)

type Store interface {
	store.Store

	NewRestrictedTokensSession(log log.Logger) RestrictedTokensSession
}

type RestrictedTokensSession interface {
	store.Session

	Create(restrictedToken *auth.RestrictedTokenRequest) (*auth.RestrictedToken, error)
	Get(id string) (*auth.RestrictedToken, error)
	Delete(id string) error
}
