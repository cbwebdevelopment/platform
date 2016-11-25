package test

import (
	"github.com/tidepool-org/platform/auth"
	testStore "github.com/tidepool-org/platform/store/test"
)

type CreateOutput struct {
	RestrictedToken *auth.RestrictedToken
	Error           error
}

type GetOutput struct {
	RestrictedToken *auth.RestrictedToken
	Error           error
}

type RestrictedTokensSession struct {
	*testStore.Session
	CreateInvocations int
	CreateInputs      []*auth.RestrictedTokenRequest
	CreateOutputs     []CreateOutput
	GetInvocations    int
	GetInputs         []string
	GetOutputs        []GetOutput
	DeleteInvocations int
	DeleteInputs      []string
	DeleteOutputs     []error
}

func NewRestrictedTokensSession() *RestrictedTokensSession {
	return &RestrictedTokensSession{
		Session: testStore.NewSession(),
	}
}

func (r *RestrictedTokensSession) Create(restrictedToken *auth.RestrictedTokenRequest) (*auth.RestrictedToken, error) {
	r.CreateInvocations++

	r.CreateInputs = append(r.CreateInputs, restrictedToken)

	if len(r.CreateOutputs) == 0 {
		panic("Unexpected invocation of Create on RestrictedTokensSession")
	}

	output := r.CreateOutputs[0]
	r.CreateOutputs = r.CreateOutputs[1:]
	return output.RestrictedToken, output.Error
}

func (r *RestrictedTokensSession) Get(id string) (*auth.RestrictedToken, error) {
	r.GetInvocations++

	r.GetInputs = append(r.GetInputs, id)

	if len(r.GetOutputs) == 0 {
		panic("Unexpected invocation of Get on RestrictedTokensSession")
	}

	output := r.GetOutputs[0]
	r.GetOutputs = r.GetOutputs[1:]
	return output.RestrictedToken, output.Error
}

func (r *RestrictedTokensSession) Delete(id string) error {
	r.DeleteInvocations++

	r.DeleteInputs = append(r.DeleteInputs, id)

	if len(r.DeleteOutputs) == 0 {
		panic("Unexpected invocation of Delete on RestrictedTokensSession")
	}

	output := r.DeleteOutputs[0]
	r.DeleteOutputs = r.DeleteOutputs[1:]
	return output
}

func (r *RestrictedTokensSession) UnusedOutputsCount() int {
	return r.Session.UnusedOutputsCount() +
		len(r.CreateOutputs) +
		len(r.GetOutputs) +
		len(r.DeleteOutputs)
}
