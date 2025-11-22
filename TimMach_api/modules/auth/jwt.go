package auth

import (
	"context"
	"time"

	db "chidinh/db/sqlc"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService tạo access token cho user.
type TokenService interface {
	GenerateToken(ctx context.Context, user db.User) (string, error)
}

// JWTMaker là impl đơn giản của TokenService dùng HMAC.
type JWTMaker struct {
	Secret string
	TTL    time.Duration
}

func (m JWTMaker) GenerateToken(_ context.Context, user db.User) (string, error) {
	ttl := m.TTL
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}

	now := time.Now()
	uid := user.ID
	claims := jwt.MapClaims{
		"user_id": uid,
		"email":   user.Email,
		"exp":     now.Add(ttl).Unix(),
		"iat":     now.Unix(),
		"sub":     uid,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.Secret))
}
