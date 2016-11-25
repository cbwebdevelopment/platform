package v1

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/tidepool-org/platform/auth/service/context"
)

var AcmeURL = os.Getenv("TIDEPOOL_OAUTH_ACME_URL")
var AcmeClientID = os.Getenv("TIDEPOOL_OAUTH_ACME_CLIENT_ID")
var AcmeClientSecret = os.Getenv("TIDEPOOL_OAUTH_ACME_CLIENT_SECRET")

var LoginURL = fmt.Sprintf("%s/v1/oauth2/login", AcmeURL)
var TokenURL = fmt.Sprintf("%s/v1/oauth2/token", AcmeURL)
var APIURL = fmt.Sprintf("%s/v1/users/self", AcmeURL)

const RedirectURL = "http://localhost:8080/Nile"

var _ConfigOnce sync.Once
var _Config *oauth2.Config

func Config() *oauth2.Config {
	if _Config == nil {
		_ConfigOnce.Do(func() {
			_Config = &oauth2.Config{
				ClientID:     AcmeClientID,
				ClientSecret: AcmeClientSecret,
				Endpoint: oauth2.Endpoint{
					AuthURL:  LoginURL,
					TokenURL: TokenURL,
				},
				RedirectURL: RedirectURL,
				Scopes:      []string{"offline_access"},
			}
		})
	}
	return _Config
}

// TODO: Fix Authentication Middleware returning JSON. We need to return something nice.

func (r *Router) GetProviderAuthorize(response rest.ResponseWriter, request *rest.Request) {
	ctx := context.MustNew(r, response, request)

	authCodeURL := Config().AuthCodeURL("no_its_you")

	ctx.Logger().WithField("headers", request.Header).Warn("Request Headers")

	http.Redirect(response.(http.ResponseWriter), request.Request, authCodeURL, http.StatusTemporaryRedirect)
}
