package services

import (
	"testing"

	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/mocks"
	"github.com/ShenokZlob/collector-ouphe/collector-service/internal/models"
	"github.com/ShenokZlob/collector-ouphe/pkg/logger"
	"github.com/stretchr/testify/suite"
)

// go test github.com/ShenokZlob/collector-ouphe/collector-service/internal/services -run UnitTestSuite
type UnitTestSuite struct {
	suite.Suite
	authService *AuthService
	authRepMock *mocks.MockAuthRepositorer
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, &UnitTestSuite{})
}

func (uts *UnitTestSuite) SetupTest() {
	authRepMock := mocks.MockAuthRepositorer{}
	authService := NewAuthService(&authRepMock, logger.SilentLogger{})

	uts.authService = authService
	uts.authRepMock = &authRepMock
}

func (uts *UnitTestSuite) TestRegister() {
	var userInfoReq, userInfoResp *models.User
	userInfoReq = &models.User{
		TelegramID: 12345,
		FirstName:  "Danya",
	}
	userInfoResp = &models.User{
		ID:         "skjfaoijah3",
		TelegramID: 12345,
		FirstName:  "Danya",
	}

	uts.authRepMock.On("CreateUser", userInfoReq).Return(userInfoResp, nil)

	actual, err := uts.authService.Register(userInfoReq)
	uts.Equal(int64(12345), actual.TelegramID)
	uts.Nil(err)
}

func (uts *UnitTestSuite) TestWho() {
	telegramIdString := "12345"
	telegramIdInt64 := int64(12345)
	userInfoResp := &models.User{
		ID:         "skjfaoijah3",
		TelegramID: 12345,
		FirstName:  "Danya",
	}

	uts.authRepMock.On("FindUserByTelegramID", telegramIdInt64).Return(userInfoResp, nil)

	user, err := uts.authService.Who(telegramIdString)
	uts.Equal(telegramIdInt64, user.TelegramID)
	uts.Nil(err)
}

func (uts *UnitTestSuite) TestLogin() {
	telegramIdInt64 := int64(12345)
	userInfoReq := &models.User{
		TelegramID: 12345,
		FirstName:  "Oleg",
	}
	userInfoResp := &models.User{
		ID:         "skjfaoijah3",
		TelegramID: 12345,
		FirstName:  "Oleg",
	}

	uts.authRepMock.On("FindUserByTelegramID", telegramIdInt64).
		Return(userInfoResp, nil)

	err := uts.authService.Login(userInfoReq)
	uts.Nil(err)
}
