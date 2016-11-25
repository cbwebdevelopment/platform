package v1

import (
	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/auth/service"
	"github.com/tidepool-org/platform/errors"
)

type Router struct {
	service.Service
}

func NewRouter(svc service.Service) (*Router, error) {
	if svc == nil {
		return nil, errors.New("v1", "service is missing")
	}

	return &Router{
		Service: svc,
	}, nil
}

func (r *Router) Routes() []*rest.Route {
	return []*rest.Route{
	// rest.Post("/v1/restrictedtokens", r.CreateRestrictedToken),
	// rest.Get("/v1/restrictedtokens/:token", r.GetRestrictedToken),
	// rest.Delete("/v1/restrictedtokens/:token", r.DeleteRestrictedToken),
	// rest.Get("/v1/oauth/:id/authorize", r.GetProviderAuthorize),
	// rest.Get("/v1/oauth/:id/redirect", r.GetProviderRedirect),
	}
}

/*
NEW:

OAuthProvider
	id: 1234567890
	slug: 'acme'
	name: 'Acme'
	clientId:
	clientSecret:

OAuthProviderSession
	providerId: <above>
	restrictedToken: ...
*/

/*
POST /v1/ratokens create one with X-Tidepool-Session-Token
	request includes what the token is for
	oauth:
		provider: TBD
	returns json with token in body

GET /v1/ratokens/token
	body is empty
	returns json with
		user info
		request info (oauth.provider)

DELETE /v1/ratokens/token
	deletes the token
	returns nothing
*/
