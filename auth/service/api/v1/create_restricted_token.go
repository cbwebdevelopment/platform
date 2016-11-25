package v1

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/auth/service/context"
)

func (r *Router) CreateRestrictedToken(response rest.ResponseWriter, request *rest.Request) {
	ctx := context.MustNew(r, response, request)

	ctx.RespondWithStatusAndData(http.StatusCreated, map[string]string{"restrictedtoken": "XYZ"})
}
