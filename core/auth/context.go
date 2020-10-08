package auth

import (
	"context"
)

// This method should be consumed within the application to extract user data from the request
// This user data is attached to the user via the above AddUserToContext method, and usually injected
// at the edge of the application.

// Warning: if something other than a UserData struct has been ser on the user key
// this will cause a panic, which is expected. These issues should be caught in development.
func GetUserFromContext(ctx context.Context) *UserData {
	val := ctx.Value("user")

	if val == nil {
		return nil
	}

	return val.(*UserData)
}

func GetUserIDFromContext(ctx context.Context) string {
	userData := GetUserFromContext(ctx)

	if userData == nil {
		return ""
	}

	return userData.Id
}

func GetUserAuthTokenFromContext(ctx context.Context) string {
	userData := GetUserFromContext(ctx)

	if userData == nil {
		return ""
	}

	return userData.AuthToken
}
