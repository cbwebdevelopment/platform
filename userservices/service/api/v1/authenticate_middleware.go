package v1

import (
	commonService "github.com/tidepool-org/platform/service"
	"github.com/tidepool-org/platform/userservices/client"
	"github.com/tidepool-org/platform/userservices/service"
)

func Authenticate(handler service.HandlerFunc) service.HandlerFunc {
	return func(context service.Context) {
		authenticationToken := context.Request().Header.Get(client.TidepoolAuthenticationTokenHeaderName)
		if authenticationToken == "" {
			context.RespondWithError(commonService.ErrorAuthenticationTokenMissing())
			return
		}

		authenticationDetails, err := context.UserServicesClient().ValidateAuthenticationToken(context, authenticationToken)
		if err != nil {
			if client.IsUnauthorizedError(err) {
				context.RespondWithError(commonService.ErrorUnauthenticated())
			} else {
				context.RespondWithInternalServerFailure("Unable to validate authentication token", err, authenticationToken)
			}
			return
		}

		context.SetAuthenticationDetails(authenticationDetails)

		handler(context)
	}
}
