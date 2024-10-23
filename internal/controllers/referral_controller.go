package controllers

import (
	"log/slog"
	"net/http"
	"referral-system/internal/infrastructure/logger/sl"
	"referral-system/internal/services"
	"time"

	"github.com/gin-gonic/gin"
)

type ReferralController struct {
	referralService services.ReferralService
	logger          *slog.Logger
}

// NewReferralController создает новый ReferralController
func NewReferralController(referralService services.ReferralService, logger *slog.Logger) *ReferralController {
	return &ReferralController{referralService: referralService, logger: logger}
}

// CreateReferralCode обрабатывает создание нового реферального кода
func (rc *ReferralController) CreateReferralCode(c *gin.Context) {
	var req struct {
		Code      string `json:"code" binding:"required"`
		ExpiresIn int64  `json:"expires_in" binding:"required"` // Время жизни в секундах
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		rc.logger.Warn("failed to bind request", sl.Err(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Получаем user_id из контекста (передан JWT миддлварой)
	userID, exists := c.Get("user_id")
	if !exists {
		rc.logger.Warn("Unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	expiresIn := time.Duration(req.ExpiresIn) * time.Second

	// Создаем реферальный код
	referral, err := rc.referralService.CreateReferralCode(int(userID.(float64)), req.Code, expiresIn)
	if err != nil {
		rc.logger.Error("failed to create referral code", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"referral_code": referral,
	})
}

// DeleteReferralCode обрабатывает удаление реферального кода
func (rc *ReferralController) DeleteReferralCode(c *gin.Context) {
	// Получаем user_id из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		rc.logger.Warn("unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Удаляем реферальный код
	err := rc.referralService.DeleteReferralCode(int(userID.(float64)))
	if err != nil {
		rc.logger.Error("failed to delete referral code", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Referral code deleted successfully",
	})
}

// GetReferralsByUserID возвращает список рефералов по ID реферера
func (rc *ReferralController) GetReferralsByUserID(c *gin.Context) {
	// Получаем user_id из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		rc.logger.Warn("unauthorized user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Получаем список рефералов
	referrals, err := rc.referralService.GetReferralsByReferrerID(int(userID.(float64)))
	if err != nil {
		rc.logger.Error("failed to get referral code", sl.Err(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"referrals": referrals,
	})
}
