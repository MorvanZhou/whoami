package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"whoami/backend-go/internal/config"
	"whoami/backend-go/internal/service"
)

type StoreAPIHandler struct {
	cfg    *config.Config
	tokens *service.StoreTokenService
}

func NewStoreAPIHandler(cfg *config.Config, tokens *service.StoreTokenService) *StoreAPIHandler {
	return &StoreAPIHandler{cfg: cfg, tokens: tokens}
}

func (h *StoreAPIHandler) Get(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "token is required"})
		return
	}

	plain, ok, err := h.tokens.Consume(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to consume token"})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"detail": "Token not found, expired, or already used"})
		return
	}

	body := "WHOAMI_API_KEY=" + plain + "\nWHOAMI_ENDPOINT=" + h.cfg.FrontendURL + "\n"
	c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(body))
}
