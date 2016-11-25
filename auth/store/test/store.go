package test

import (
	"github.com/tidepool-org/platform/auth/store"
	"github.com/tidepool-org/platform/log"
	testStore "github.com/tidepool-org/platform/store/test"
)

type Store struct {
	*testStore.Store
	NewRestrictedTokensSessionInvocations int
	NewRestrictedTokensSessionInputs      []log.Logger
	NewRestrictedTokensSessionOutputs     []store.RestrictedTokensSession
}

func NewStore() *Store {
	return &Store{
		Store: testStore.NewStore(),
	}
}

func (s *Store) NewRestrictedTokensSession(lgr log.Logger) store.RestrictedTokensSession {
	s.NewRestrictedTokensSessionInvocations++

	s.NewRestrictedTokensSessionInputs = append(s.NewRestrictedTokensSessionInputs, lgr)

	if len(s.NewRestrictedTokensSessionOutputs) == 0 {
		panic("Unexpected invocation of NewRestrictedTokensSession on Store")
	}

	output := s.NewRestrictedTokensSessionOutputs[0]
	s.NewRestrictedTokensSessionOutputs = s.NewRestrictedTokensSessionOutputs[1:]
	return output
}

func (s *Store) UnusedOutputsCount() int {
	return len(s.NewRestrictedTokensSessionOutputs)
}
