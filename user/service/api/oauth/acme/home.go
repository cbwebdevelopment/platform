package acme

import (
	"html/template"
	"net/http"

	"github.com/tidepool-org/platform/user/service"
)

func Home(serviceContext service.Context) {
	homeTemplate, err := template.New("redirect").Parse(homeHTML)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to parse template", err)
		return
	}
	if homeTemplate == nil {
		serviceContext.RespondWithInternalServerFailure("Parsed template is missing", err)
		return
	}

	responseWriter := serviceContext.Response().(http.ResponseWriter)
	responseWriter.Header().Add("Content-Type", "text/html")

	homeTemplate.Execute(responseWriter, OAuthToken)
}

const homeHTML string = `
<html>
<head>
<title>Acme OAuth Home</title>
</head>
<body>
<h2>Acme OAuth Home</h2>
<table>
<tr><td>AccessToken</td><td>{{.AccessToken}}</td></tr>
<tr><td>TokenType</td><td>{{.TokenType}}</td></tr>
<tr><td>RefreshToken</td><td>{{.RefreshToken}}</td></tr>
<tr><td>Expiry</td><td>{{.Expiry}}</td></tr>
</table>
</p>
<a href="/oauth/acme/devices">Devices</a>
</body>
</html>
`
