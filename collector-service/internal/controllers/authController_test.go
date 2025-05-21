package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/mocks"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// go test github.com/ShenokZlob/collector-ouphe/collector-service/internal/controllers -run UnitTestSuite
type UnitTestSuite struct {
	suite.Suite
	router          *gin.Engine
	authController  *AuthController
	authServiceMock *mocks.MockAuthServicer
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (its *UnitTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	its.router = gin.New()

	its.router.POST("/register", func(ctx *gin.Context) { its.authController.Register(ctx) })
	its.router.GET("/user/telegram/:telegram_id", func(ctx *gin.Context) { its.authController.Who(ctx) })
	its.router.POST("/login", func(ctx *gin.Context) { its.authController.Login(ctx) })
}

func (its *UnitTestSuite) SetupTest() {
	its.authServiceMock = &mocks.MockAuthServicer{}
	its.authController = NewAuthController(its.authServiceMock, logger.SilentLogger{})
}

func (its *UnitTestSuite) TestRegister() {
	// Mocking the HTTP context
	reqModel := &models.User{TelegramID: 123, FirstName: "John", Username: "@john"}
	resModel := &models.User{ID: "abc", TelegramID: 123, FirstName: "John", Username: "@john"}
	its.authServiceMock.On("Register", reqModel).Return(resModel, nil)

	// Making a request to the Register endpoint
	body := `{"telegram_id":123,"first_name":"John","username":"@john"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	its.router.ServeHTTP(w, req)

	// Asserting the response
	its.Equal(http.StatusCreated, w.Code)
	var resp struct {
		User UserResponse `json:"user"`
	}
	its.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	its.Equal("abc", resp.User.ID)
	its.Equal(int64(123), resp.User.TelegramID)
	its.Equal("John", resp.User.FirstName)
	its.Equal("@john", resp.User.Username)

	its.authServiceMock.AssertExpectations(its.T())
}

func (its *UnitTestSuite) TestWho() {
	// Mocking the HTTP context
	telegramIdString := "123"
	respModel := &models.User{ID: "abc", TelegramID: 123, FirstName: "John", Username: "@john"}
	its.authServiceMock.On("Who", telegramIdString).Return(respModel, nil)

	// Making a request to the Who endpoint
	req := httptest.NewRequest(http.MethodGet, "/user/telegram/123", nil)
	w := httptest.NewRecorder()

	its.router.ServeHTTP(w, req)

	// Asserting the response
	its.Equal(http.StatusOK, w.Code)
	var resp UserResponse
	its.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	its.Equal("abc", resp.ID)
	its.Equal(int64(123), resp.TelegramID)
	its.Equal("John", resp.FirstName)
	its.Equal("@john", resp.Username)

	its.authServiceMock.AssertExpectations(its.T())
}

func (its *UnitTestSuite) TestLogin() {
	// Mocking the HTTP context
	reqModel := &models.User{TelegramID: 123, FirstName: "John", Username: "@john"}
	its.authServiceMock.On("Login", reqModel).Return(nil)

	// Making a request to the Login endpoint
	body := `{"telegram_id":123,"first_name":"John","username":"@john"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	its.router.ServeHTTP(w, req)

	// Asserting the response
	its.Equal(http.StatusOK, w.Code)
	var resp struct {
		Token string `json:"token"`
	}
	its.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	its.NotEmpty(resp.Token)

	its.authServiceMock.AssertExpectations(its.T())
}
