package users

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes gắn endpoint users vào router (có thể truyền auth middleware).
func RegisterUserRoutes(r *gin.RouterGroup, h *Controller) {
	group := r.Group("/users")
	group.GET("/me", h.GetMe)
}
