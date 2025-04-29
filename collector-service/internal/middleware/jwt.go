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
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.secret), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if exp, ok := claims["exp"].(float64); !ok || int64(exp) < time.Now().Unix() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("userID", claims["user_id"].(string))
		ctx.Next()
	}
}
