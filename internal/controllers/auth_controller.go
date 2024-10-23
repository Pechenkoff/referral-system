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

// Login обрабатывает вход пользователя
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
		ac.logger.Warn("failed to login", sl.Err(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}

// Register обрабатывает регистрацию нового пользователя
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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}
