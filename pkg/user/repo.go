package user

import (
	"context"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	CreateUser(context.Context, *User) (*User, error)
	UpdateUser(context.Context,*User,string ) error
	GetUserById(context.Context, string) (*User, error)
	GetUserByEmail(context.Context, string) (*User, error)
}
