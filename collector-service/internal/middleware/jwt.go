package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Needs initialization logger
type JWTMiddleware struct {
	secret string
	logger logger.Logger
}

func NewJWTMiddleware(secret string, logger logger.Logger) *JWTMiddleware {
	return &JWTMiddleware{
		secret: secret,
		logger: logger,
	}
}

func (m *JWTMiddleware) Authorization() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		m.logger.Info("JWT middleware triggered")
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			m.logger.Info("Authorization header does not start with Bearer", logger.String("header", authHeader))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})
		if err != nil {
			m.logger.Error("Failed to parse JWT token", logger.Error(err), logger.String("token", tokenStr))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			m.logger.Error("Invalid JWT claims", logger.String("claims", tokenStr))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			m.logger.Error("Invalid or missing user_id in JWT claims")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: Check if the token is valid
		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			m.logger.Warn("JWT token has expired", logger.Int("expiration", int(exp)), logger.String("current_time", time.Now().String()))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}
