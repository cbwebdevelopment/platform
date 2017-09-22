package v1

import (
	"fmt"
	"net/http"

	commonService "github.com/tidepool-org/platform/service"
	"github.com/tidepool-org/platform/user"
	"github.com/tidepool-org/platform/userservices/service"
)

// type UsersCreateParameters struct {
// 	Email    *string   `json:"username,omitempty"`
// 	Emails   *[]string `json:"emails,omitempty"`
// 	Password *string   `json:"password,omitempty"`
// 	Roles    *[]string `json:"roles,omitempty"`
// }

// type ProfileCreateParameters struct {
// 	FullName *string `json:"fullName"`
// }

func UsersCreate(serviceContext service.Context) {
	var usersCreateParameters *user.UserCreate
	// var profileCreateParameters *profile.Profile

	var err error
	if err = serviceContext.Request().DecodeJsonPayload(&usersCreateParameters); err != nil {
		serviceContext.RespondWithError(commonService.ErrorJSONMalformed())
		return
	}

	// if err = serviceContext.Request().DecodeJsonPayload(&profileCreateParameters); err != nil {
	// 	serviceContext.RespondWithError(commonService.ErrorJSONMalformed())
	// 	return
	// }

	fmt.Print(usersCreateParameters)

	newUser, err := serviceContext.UserStoreSession().CreateUser(usersCreateParameters)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to create user", err)
		return
	}

	fmt.Printf("%+v", newUser)

	serviceContext.RespondWithStatusAndData(http.StatusOK, struct{}{})
}
