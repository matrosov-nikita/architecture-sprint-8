package get_reports

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"reports/usecases/reports"
	"reports/usecases/token_verifier"
)

type tokenVerifier interface {
	VerifyToken(ctx context.Context, token string) error
}

type Handler struct {
	tokenVerifier tokenVerifier
}

func NewHandler(tokenVerifier tokenVerifier) *Handler {
	return &Handler{tokenVerifier: tokenVerifier}
}

func (h *Handler) GetReports(c *gin.Context) {
	token, exists := c.Get("access_token")
	if !exists {
		slog.Error("access token is missing in request headers")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token found"})
		return
	}

	tokenString, ok := token.(string)
	if !ok {
		slog.Error("access token is not string")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Can not parse token as string"})
		return
	}

	if err := h.tokenVerifier.VerifyToken(c, tokenString); err != nil {
		l := slog.With("err", err)

		if errors.Is(err, token_verifier.ErrRoleAccessDenied) {
			l.Error("verify token failed")
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		l.Error("invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	c.JSON(http.StatusOK, reports.GetReports())
}
