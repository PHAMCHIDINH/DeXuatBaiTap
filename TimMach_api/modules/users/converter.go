package users

import (
	db "chidinh/db/sqlc"
	"time"
)

// User is a domain representation detached from DB driver types.
type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
}

func toUserDomain(u db.User) User {
	return User{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt.Time,
	}
}

func toUserResponse(u User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
