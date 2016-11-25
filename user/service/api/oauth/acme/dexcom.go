package acme

import (
	"github.com/tidepool-org/platform/user/service"
)

func Routes() []service.Route {
	return []service.Route{
		service.MakeRoute("GET", "/oauth/acme/start", Start), // TODO: Temporary starting page
		service.MakeRoute("GET", "/Nile", Redirect),          // "/oauth/acme/redirect"
		service.MakeRoute("GET", "/oauth/acme/devices", Devices),
		service.MakeRoute("GET", "/oauth/acme/home", Home),
	}
}
