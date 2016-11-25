package acme

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"golang.org/x/oauth2"

	"github.com/tidepool-org/platform/user/service"
)

var OAuthToken *oauth2.Token

// TODO: What does any access token ultimately give DIRECT access to
// Is it just authenticatedUserID?

// TODO: SUAT should contain
// authenticatedUserID string
// createdTime time.Time
// createdUserID string
// modifiedTime time.Time
// modifiedUserID string
// expirationTime time.Time
// source string ("oauth", etc...)
// payload interface{}
//   for oauth usage, payload should include:
//     providerID string
//     state string - id.New()

func Redirect(serviceContext service.Context) {

	// TODO: What does any failure below actually mean?
	//   certainly log.Error
	//   redirect to nice error page?
	//   record error in dataSource status (if possible, some failures occur without authenticatedUserID - namely without valid SUAT)
	//   NOT error for missing/expired SUAT (could happen in normal course with poor network)

	// ---

	// TODO: Get Single Use Access Token (SUAT) cookie; if missing, then error
	// TODO: Get SUAT from store
	//			does store check if SUAT is expired - I think so, it should also automatically delete

	// TODO: Delete SUAT - not sure this should go here, maybe later down

	fmt.Fprintf(os.Stderr, "redirectURL=%s\n", serviceContext.Request().URL)

	query := serviceContext.Request().URL.Query()

	// TODO: Should we ensure/force that ONLY "state", and "code" or "error" are in query?

	queryState := query.Get("state")
	fmt.Fprintf(os.Stderr, "queryState=%s\n", queryState)

	// TODO: Compare queryState against SUAT state, if not exact match, then fail

	queryError := query.Get("error")
	if queryError != "" {
		fmt.Fprintf(os.Stderr, "queryError=%s\n", queryError)
		// TODO: What to do? We got an error from Provider
		return
	}

	queryCode := query.Get("code")
	fmt.Fprintf(os.Stderr, "queryCode=%s\n", queryCode)
	if queryCode == "" {
		// TODO: Error - 301 to error page?
		return
	}

	// TODO: get oauth config via providerID, if no match, then fail

	token, err := Config().Exchange(serviceContext.Request().Context(), queryCode)
	if err != nil {
		// TODO: Error - 301 to error page?
		return
	}

	fmt.Fprintf(os.Stderr, "token=%#v\n", token)

	bytes, err := json.Marshal(token)
	fmt.Fprintf(os.Stderr, "tokenJSON=%#v\n", string(bytes))

	// Remember token
	OAuthToken = token

	// serviceContext.Response().Header().Set("Location", "/oauth/acme/home")
	// serviceContext.Response().WriteHeader(301)

	responseObject := struct {
	}{}

	redirectTemplate, err := template.New("redirect").Parse(redirectHTML)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to parse template", err)
		return
	}
	if redirectTemplate == nil {
		serviceContext.RespondWithInternalServerFailure("Parsed template is missing", err)
		return
	}

	responseWriter := serviceContext.Response().(http.ResponseWriter)
	responseWriter.Header().Add("Content-Type", "text/html")

	redirectTemplate.Execute(responseWriter, responseObject)
}

const redirectHTML string = `
<html>
<head>
<title>Acme OAuth Redirect</title>
</head>
<body onLoad='loaded();'>
</body>
<script>
	function loaded() {
		window.opener.parent.location.href = '/oauth/acme/home';
		window.close();
	}
</script>
</html>
`
