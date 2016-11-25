package acme

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/tidepool-org/platform/user/service"
)

func Devices(serviceContext service.Context) {

	fmt.Fprintf(os.Stderr, "OAuthToken=%#v\n", OAuthToken)

	tokenSource := Config().TokenSource(serviceContext.Request().Context(), OAuthToken)

	originalToken, err := tokenSource.Token()

	fmt.Fprintf(os.Stderr, "NEW originalToken=%#v\n", originalToken)
	fmt.Fprintf(os.Stderr, "NEW err=%#v\n", err)

	httpClient := Config().Client(serviceContext.Request().Context(), originalToken)

	devicesURL, err := url.Parse(fmt.Sprintf("%s/devices", APIURL))
	query := devicesURL.Query()
	query.Set("startDate", "2017-06-01T00:00:00")
	query.Set("endDate", "2017-08-01T00:00:00")
	devicesURL.RawQuery = query.Encode()

	requestURL := devicesURL.String()

	fmt.Fprintf(os.Stderr, "requestURL=%#v\n", requestURL)

	request, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to create request", err)
		return
	}

	response, err := httpClient.Do(request)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to perform request", err)
		return
	}
	defer response.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(response.Body)

	newToken, err := tokenSource.Token()
	if newToken.AccessToken != originalToken.AccessToken ||
		newToken.TokenType != originalToken.TokenType ||
		newToken.RefreshToken != originalToken.RefreshToken ||
		newToken.Expiry != originalToken.Expiry {
		fmt.Fprintf(os.Stderr, "Saving newToken!!!")
		OAuthToken = newToken
	}

	fmt.Fprintf(os.Stderr, "AFTER newToken=%#v\n", newToken)
	fmt.Fprintf(os.Stderr, "AFTER err=%#v\n", err)

	responseObject := struct {
		Body string
	}{
		Body: string(bodyBytes),
	}

	devicesTemplate, err := template.New("devices").Parse(devicesHTML)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to parse template", err)
		return
	}
	if devicesTemplate == nil {
		serviceContext.RespondWithInternalServerFailure("Parsed template is missing", err)
		return
	}

	responseWriter := serviceContext.Response().(http.ResponseWriter)
	responseWriter.Header().Add("Content-Type", "text/html")

	devicesTemplate.Execute(responseWriter, responseObject)
}

const devicesHTML string = `
<html>
<head>
<title>Acme OAuth Devices</title>
</head>
<body>
<h2>Acme OAuth Devices</h2>
</p>
{{.Body}}
</p>
<a href="/oauth/acme/home">Home</a>
</body>
</html>
`
