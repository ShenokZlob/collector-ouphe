// go test -v github.com/ShenokZlob/collector-ouphe/collector-service/internal/middleware
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	secret := "test-secret"
	mw := NewJWTMiddleware(secret, nil)

	t.Run("missing Bearer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		// ctx.Request = req

		r.Use(mw.Authorization())
		r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.here")
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		// ctx.Request = req

		r.Use(mw.Authorization())
		r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("valid token", func(t *testing.T) {
		token := generateTestJWT(secret, "user123", time.Now().Add(time.Hour))

		req := httptest.NewRequest("GET", "/", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		// ctx.Request = req

		var receivedUserID string
		r.Use(mw.Authorization())
		r.GET("/", func(c *gin.Context) {
			receivedUserID = c.GetString("userID")
			c.Status(http.StatusOK)
		})

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "user123", receivedUserID)
	})
}

func generateTestJWT(secret, userID string, exp time.Time) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
