package models

import (
	"context"
	"homework/pkg/errors"
)

const (
	errContextIsNil = "context is nil"
	errNoUser       = "no user"
)

var ctxUser ctxAuthKey = struct{}{}

type ctxAuthKey struct{}

type SignUpInput struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

type SignInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Admin    bool   `json:"admin"`
}

func (u *User) ToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxUser, u)
}

func UserFromContext(ctx context.Context) (*User, error) {
	if ctx == nil {
		return nil, errors.New(errContextIsNil)
	}

	rawUser := ctx.Value(ctxUser)
	if rawUser == nil {
		return nil, errors.New(errNoUser)
	}

	user := rawUser.(*User)
	if user == nil {
		return nil, errors.New(errNoUser)
	}

	return user, nil
}
