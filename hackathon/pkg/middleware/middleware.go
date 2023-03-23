package middleware

import (
	"encoding/json"
	"hackathon/pkg/auth"
	"hackathon/pkg/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	tokenTypeBearer       = "Bearer"
	errorCodeUnauthorized = 401
	headerAuthorization   = "Authorization"
)

func AuthenizationMiddleware(authHandler *auth.AuthHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwt := extractToken(ctx)
		accountInfor, err := authHandler.ValidateAccessToken(jwt)
		if err != nil {
			logrus.Error("token verifying error", err)
			ctx.JSON(http.StatusUnauthorized, nil)
			return
		}
		account := models.Account{}
		err = json.Unmarshal([]byte(accountInfor), &account)

		if err != nil {
			logrus.Error("Unmarshal account error", err)
			ctx.JSON(http.StatusUnauthorized, nil)
			return
		}
		ctx.Set(auth.AccountKey, account)
		ctx.Next()
	}
}

func extractToken(ctx *gin.Context) string {
	authHeader := ctx.GetHeader(headerAuthorization)
	s := strings.SplitN(authHeader, " ", 2)
	if len(s) != 2 || !strings.EqualFold(s[0], tokenTypeBearer) {
		return ""
	}
	return s[1]
}
