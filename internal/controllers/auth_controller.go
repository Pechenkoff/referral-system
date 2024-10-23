package controllers

import (
	"log/slog"
	"net/http"
	"referral-system/internal/infrastructure/logger/sl"
	"referral-system/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
	logger      *slog.Logger
}

// NewAuthController создает новый AuthController
func NewAuthController(authService services.AuthService, logger *slog.Logger) *AuthController {
	return &AuthController{authService: authService, logger: logger}
}

// Login godoc
// @Summary Вход пользователя
// @Description Вход пользователя с получением JWT токена
// @Tags auth
// @Accept json
// @Produce json
// @Param email body string true "Email пользователя"
// @Param password body string true "Пароль"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ac.logger.Warn("failed to bind request", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, token, err := ac.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		ac.logger.Error("failed to login", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Register godoc
// @Summary Регистрация нового пользователя
// @Description Регистрация нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param name body string true "Имя пользователя"
// @Param email body string true "Email пользователя"
// @Param password body string true "Пароль"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		ac.logger.Warn("failed to bind request", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := ac.authService.RegisterUser(req.Name, req.Email, req.Password)
	if err != nil {
		ac.logger.Error("failed to register user", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}
