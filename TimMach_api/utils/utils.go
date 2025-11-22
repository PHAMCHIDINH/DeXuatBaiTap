package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RespondError trả JSON lỗi chuẩn và dừng xử lý request.
func RespondError(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{"error": msg})
}

// UserIDFromContext lấy userID do middleware auth đã set trước đó (string).
func UserIDFromContext(c *gin.Context) (string, bool) {
	keys := []string{"userID", "user_id"}
	for _, key := range keys {
		if raw, ok := c.Get(key); ok {
			if v, ok := raw.(string); ok && v != "" {
				return v, true
			}
		}
	}
	return "", false
}

func FormatUserID(seq int64, ts time.Time) string {
	date := ts.Format("20060102")
	return fmt.Sprintf("USER_%s_%03d", date, seq)
}
