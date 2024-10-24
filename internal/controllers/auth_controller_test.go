package controllers_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"referral-system/internal/controllers"
	"referral-system/internal/controllers/mocks"
	"referral-system/internal/entities"
	"referral-system/internal/infrastructure/logger/handlers/slogdiscard"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAuthController_Login_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockAuthService := mocks.NewAuthService(t)

	logger := slogdiscard.NewDiscardLogger()
	authController := controllers.NewAuthController(mockAuthService, logger)

	router.POST("/auth/login", authController.Login)
	var mockReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	mockReq.Email, mockReq.Password = "example@mail.com", "test_password"
	mockUser := &entities.User{
		Email:          mockReq.Email,
		HashedPassword: mockReq.Password,
	}

	mockAuthService.On("LoginUser", mockReq.Email, mockReq.Password).
		Return(mockUser, "valid_token", nil)

	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{"email": "example@mail.com", "password": "test_password"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "valid_token")
	assert.Contains(t, w.Body.String(), "example@mail.com")

	mockAuthService.AssertExpectations(t)
}

func TestAuthController_Login_InvalidCredentials(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockAuthService := mocks.NewAuthService(t)

	logger := slogdiscard.NewDiscardLogger()
	authController := controllers.NewAuthController(mockAuthService, logger)
	router.POST("/auth/login", authController.Login)

	// Настраиваем mock-ответ для метода LoginUser
	mockAuthService.On("LoginUser", "example@mail.com", "wrong_password").
		Return(nil, "", errors.New("Invalid credentials"))

	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{"email": "example@mail.com", "password": "wrong_password"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid credentials")

	mockAuthService.AssertExpectations(t)
}

func TestAuthController_Login_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockAuthService := mocks.NewAuthService(t)

	logger := slogdiscard.NewDiscardLogger()
	authController := controllers.NewAuthController(mockAuthService, logger)
	router.POST("/auth/login", authController.Login)

	// Создаем HTTP-запрос с некорректными данными
	req, _ := http.NewRequest("POST", "/auth/login", strings.NewReader(`{"email": "example@mail.com"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request")

	mockAuthService.AssertNotCalled(t, "LoginUser")
}
