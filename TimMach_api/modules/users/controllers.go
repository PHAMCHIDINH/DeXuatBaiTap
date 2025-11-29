package users

import (
	db "chidinh/db/sqlc"
	"chidinh/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Controller gom các dependency cho module users.
type Controller struct {
	Queries *db.Queries
}

func NewController(queries *db.Queries) *Controller {
	return &Controller{
		Queries: queries,
	}
}

// GET /users/me
func (h *Controller) GetMe(c *gin.Context) {
	userID, ok := utils.UserIDFromContext(c)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "missing user in context")
		return
	}

	user, err := h.Queries.GetUserByID(c, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		utils.RespondError(c, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "cannot fetch user")
		return
	}

	c.JSON(http.StatusOK, toUserResponse(toUserDomain(user)))
}
