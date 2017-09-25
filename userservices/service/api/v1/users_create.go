package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/tidepool-org/platform/profile"
	commonService "github.com/tidepool-org/platform/service"
	"github.com/tidepool-org/platform/user"
	"github.com/tidepool-org/platform/userservices/service"
)

func UsersCreate(serviceContext service.Context) {
	var usersCreateParameters *user.UserCreate
	var profileCreateParameters *profile.Profile

	var err error
	if err = serviceContext.Request().DecodeJsonPayload(&usersCreateParameters); err != nil {
		serviceContext.RespondWithError(commonService.ErrorJSONMalformed())
		return
	}

	newUser, err := serviceContext.UserStoreSession().CreateUser(usersCreateParameters)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to create user", err)
		return
	}

	fmt.Printf("New User: %+v", newUser)
	var value = fmt.Sprintf(`{
		"profile": {
			"fullName": "%s",
			"patient": {
				"birthday":"1986-01-01",
				"diagnosisDate":"1992-01-01"
			}
		}
	}`, usersCreateParameters.FullName)

	value = strings.Replace(value, "\n", "", -1)
	value = strings.Replace(value, "\t", "", -1)

	profileCreateParameters = &profile.Profile{
		UserID: newUser.ID,
		Value:  value,
	}

	newProfile, err := serviceContext.ProfileStoreSession().CreateProfile(profileCreateParameters)
	if err != nil {
		serviceContext.RespondWithInternalServerFailure("Unable to create user", err)
		return
	}

	fmt.Printf("New Profile%+v", newProfile)

	serviceContext.RespondWithStatusAndData(http.StatusOK, struct{}{})
}
