package users

import "github.com/gin-gonic/gin"

// RegisterUserRoutes gắn endpoint users vào router (có thể truyền auth middleware).
func RegisterUserRoutes(r *gin.RouterGroup, h *Controller, auth gin.HandlerFunc) {
	group := r.Group("/users")
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)

	if auth != nil {
		group.GET("/me", auth, h.GetMe)
	} else {
		group.GET("/me", h.GetMe)
	}
}
