package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONRecovery recovers panics and returns FastAPI-like {"detail":"..."}.
func JSONRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		slog.Error("panic recovered", "err", recovered, "path", c.Request.URL.Path)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": "Internal server error",
		})
	})
}
