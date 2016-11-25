package acme

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sync"

	"golang.org/x/oauth2"

	"github.com/tidepool-org/platform/user/service"
)

var AcmeURL = os.Getenv("TIDEPOOL_OAUTH_ACME_URL")
var AcmeClientID = os.Getenv("TIDEPOOL_OAUTH_ACME_CLIENT_ID")
var AcmeClientSecret = os.Getenv("TIDEPOOL_OAUTH_ACME_CLIENT_SECRET")

var LoginURL = fmt.Sprintf("%s/v1/oauth2/login", AcmeURL)
var TokenURL = fmt.Sprintf("%s/v1/oauth2/token", AcmeURL)
var APIURL = fmt.Sprintf("%s/v1/users/self", AcmeURL)

const RedirectURL = "http://localhost:8080/Nile"

// ctx := context.Background()

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

// // Redirect user to consent page to ask for permission
// // for the scopes specified above.
// url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
// fmt.Printf("Visit the URL for the auth dialog: %v", url)

// // Use the authorization code that is pushed to the redirect
// // URL. Exchange will do the handshake to retrieve the
// // initial access token. The HTTP Client returned by
// // conf.Client will refresh the token as necessary.
// var code string
// if _, err := fmt.Scan(&code); err != nil {
//     log.Fatal(err)
// }
// tok, err := conf.Exchange(ctx, code)
// if err != nil {
//     log.Fatal(err)
// }

// client := conf.Client(ctx, tok)
// client.Get("...")

func Start(serviceContext service.Context) {
	startTemplate, err := template.New("start").Parse(startHTML)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to parse template", err)
		return
	}
	if startTemplate == nil {
		serviceContext.RespondWithInternalServerFailure("Parsed template is missing", err)
		return
	}

	// redirectURL, err := url.Parse(RedirectURL)
	// if err != nil {
	// 	serviceContext.RespondWithInternalServerFailure("Unable to parse redirectURL", err)
	// 	return
	// }
	// redirectQuery := redirectURL.Query()
	// // redirectQuery.Set("user_id", "my-test-user-id")
	// redirectURL.RawQuery = redirectQuery.Encode()

	// loginURL, err := url.Parse(LoginURL)
	// if err != nil {
	// 	serviceContext.RespondWithInternalServerFailure("Unable to parse login URL", err)
	// 	return
	// }
	// loginQuery := loginURL.Query()
	// loginQuery.Set("client_id", ClientID)
	// loginQuery.Set("redirect_uri", RedirectURL)
	// loginQuery.Set("response_type", "code")
	// loginQuery.Set("scope", "offline_access")
	// loginQuery.Set("state", "yo_its_me")
	// loginURL.RawQuery = loginQuery.Encode()

	// TODO: state should include reference to:
	// 	session token
	//  oauth provider ("acme")

	authCodeURL := Config().AuthCodeURL("no_its_you")

	fmt.Fprintf(os.Stderr, "authCodeURL=%s\n", authCodeURL)

	startData := struct {
		URL string
	}{
		authCodeURL,
	}

	response := serviceContext.Response().(http.ResponseWriter)
	response.Header().Add("Content-Type", "text/html")

	startTemplate.Execute(response, startData)
}

const startHTML string = `
<html>
<head>
<title>Acme OAuth Start</title>
</head>
<body>
<h2>Acme OAuth Start</h2>
</p>
Click the link to begin the Acme OAuth workflow:
<a href='#' onClick='authenticate();'>Authenticate</a>
</body>
<script>
	function authenticate() {
		var width = 1080;
		var height = 840;
		var left = window.screenX + (window.innerWidth / 2) - (width / 2);
		var top = window.screenY + (window.innerHeight / 2) - (height / 2);

		window.open('{{.URL}}', 'oauth', 'left=' + left + ', top=' + top + ', width=' + width + ', height=' + height);
		// window.open('{{.URL}}', 'oauth', 'toolbar=no, location=no, directories=no, status=no, menubar=no, scrollbars=yes, resizable=yes, copyhistory=no, width=' + width + ', height=' + height + ', top=' + top + ', left=' + left);
	}
</script>
</html>
`
