package v1

import (
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"

	"github.com/tidepool-org/platform/auth/service/context"
)

func (r *Router) GetProviderRedirect(response rest.ResponseWriter, request *rest.Request) {
	ctx := context.MustNew(r, response, request)

	ctx.RespondWithStatusAndData(http.StatusOK, nil)
}
