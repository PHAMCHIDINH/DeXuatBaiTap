package users

import "time"

// UserResponse đại diện dữ liệu user trả về cho client.
type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
