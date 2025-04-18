package auth

import (
	"context"
)

// AuthRepository defines the interface for authentication related database operations.
type AuthRepository interface {
	// CreateUser inserts a new user into the database.
	CreateUser(ctx context.Context, user *User) error
	// FindByUsername retrieves a user by their username.
	// Returns ErrUserNotFound if the user does not exist.
	FindByUsername(ctx context.Context, username string) (*User, error)
	// FindByEmailOrUsername retrieves a user by their email or username.
	// Returns ErrUserNotFound if the user does not exist.
	FindByEmailOrUsername(ctx context.Context, email, username string) (*User, error)
	// FindByID retrieves a user by their ID.
	// Returns ErrUserNotFound if the user does not exist.
	FindByID(ctx context.Context, id string) (*User, error)
}
