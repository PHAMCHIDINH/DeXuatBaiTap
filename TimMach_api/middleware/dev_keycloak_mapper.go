package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	db "chidinh/db/sqlc"
	"chidinh/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

// DevKeycloakMapper parses a Bearer token (without verification), maps Keycloak users to DB rows,
// and sets the internal userID into the Gin context. Only intended for local/dev usage.
func DevKeycloakMapper(queries *db.Queries) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondError(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			utils.RespondError(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}
		fmt.Println("Parsed token:", token)

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(c, http.StatusUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}
		fmt.Println("Parsed claims:", claims)

		sub, _ := claims.GetSubject()
		if sub == "" {
			utils.RespondError(c, http.StatusUnauthorized, "Missing sub in token")
			c.Abort()
			return
		}
		fmt.Println("Subject:", sub)

		ctx := c.Request.Context()

		user, err := queries.GetUserByKeycloakID(ctx, &sub)
		switch {
		case err == nil:
			c.Set("userID", user.ID)
			c.Set("userEmail", user.Email)
			c.Next()
			return
		case errors.Is(err, pgx.ErrNoRows):
			// try by email
			email := ""
			if v, ok := claims["email"].(string); ok {
				email = v
			}
			if email == "" {
				if v, ok := claims["preferred_username"].(string); ok {
					email = v
				}
			}
			if email == "" {
				utils.RespondError(c, http.StatusUnauthorized, "Token missing email")
				c.Abort()
				return
			}

			userByEmail, err := queries.GetUserByEmail(ctx, email)
			switch {
			case err == nil:
				user, err := queries.AttachKeycloakID(ctx, db.AttachKeycloakIDParams{
					ID:         userByEmail.ID,
					KeycloakID: &sub,
				})
				if err != nil {
					utils.RespondError(c, http.StatusInternalServerError, "cannot attach keycloak id")
					c.Abort()
					return
				}
				c.Set("userID", user.ID)
				c.Set("userEmail", user.Email)
			case errors.Is(err, pgx.ErrNoRows):
				seq, seqErr := queries.NextUserSeq(ctx)
				if seqErr != nil {
					utils.RespondError(c, http.StatusInternalServerError, "cannot allocate user id")
					c.Abort()
					return
				}
				newID := utils.FormatUserID(seq, time.Now())
				user, err := queries.CreateKeycloakUser(ctx, db.CreateKeycloakUserParams{
					ID:         newID,
					Email:      email,
					KeycloakID: &sub,
				})
				if err != nil {
					utils.RespondError(c, http.StatusInternalServerError, "cannot create user")
					c.Abort()
					return
				}
				c.Set("userID", user.ID)
				c.Set("userEmail", user.Email)
			default:
				utils.RespondError(c, http.StatusInternalServerError, "cannot fetch user by email")
				c.Abort()
				return
			}
		default:
			utils.RespondError(c, http.StatusInternalServerError, "cannot fetch user")
			c.Abort()
			return
		}

		c.Next()
	}
}
