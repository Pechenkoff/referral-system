package routes

import (
	"net/http"
	"referral-system/internal/controllers"
	"referral-system/internal/middlewares"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, authController *controllers.AuthController, referralController *controllers.ReferralController, jwtSecret string) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Маршруты для аутентификации
	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	// Защищенные маршруты
	protected := router.Group("/referrals")
	protected.Use(middlewares.JWTMiddleware(jwtSecret))
	{
		protected.POST("/", referralController.CreateReferralCode)
		protected.DELETE("/", referralController.DeleteReferralCode)
		protected.GET("/list", referralController.GetReferralsByUserID)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
	})

	router.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "method not allowed"})
	})
}
