package auth

import "context"

type UserData struct {
	AuthToken string
	Id        string
}

// Authorizer is meant to take an incoming request context and add user request data
// to the incoming context
type Authorizer interface {
	AddUserToContext(ctx context.Context) (context.Context, error)
}
