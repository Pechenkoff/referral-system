package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"referral-system/internal/config"
	"referral-system/internal/controllers"
	"referral-system/internal/infrastructure/logger/handlers/slogpretty"
	"referral-system/internal/infrastructure/logger/sl"
	"referral-system/internal/repositories/postgres"
	"referral-system/internal/routes"
	"referral-system/internal/services"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for a referral system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// считываем флаги
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	migrationPath := flag.String("migration", "file://db/migrations", "Path to the migration directory")
	flag.Parse()

	// читаем конфигурационный файл
	cfg := config.MustLoadConfig(*configPath)

	// создаем логгер
	logger := setupPrettyLog()

	// Подключаемся к базе данных
	dbConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	poolConfig, err := pgxpool.ParseConfig(dbConnStr)
	if err != nil {
		panic(fmt.Errorf("unable to parse connection string: %v\n", err))
	}

	poolConfig.MaxConns = 10
	poolConfig.MaxConnLifetime = time.Hour

	dbConn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		panic(fmt.Errorf("Unable to connect to database: %v\n", err))
	}
	defer dbConn.Close()

	// Проводим миграции
	postgres.MustRunMigration(dbConnStr, *migrationPath)

	// создаем копии репозиториев
	userRepo := postgres.NewPostgresUserRepository(dbConn)
	referralCodeRepo := postgres.NewPostgresReferralCodeRepository(dbConn)
	referralRepo := postgres.NewPostgresReferralRepository(dbConn)

	// создаем копии сервисов
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	referralService := services.NewReferralService(referralCodeRepo, userRepo, referralRepo)

	// создаем контроллеры
	authController := controllers.NewAuthController(authService, logger)
	referralController := controllers.NewReferralController(referralService, logger)

	// создаем копию роутера
	router := gin.Default()
	routes.RegisterRoutes(router, authController, referralController, cfg.JWTSecret)

	// подключаем Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		WriteTimeout: time.Duration(cfg.Timeouts.WriteTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.Timeouts.ReadTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Timeouts.IdleTimeout) * time.Second,
	}

	// релизуем gracefull отключение сервера
	errChan := make(chan error, 1)

	go func() {
		logger.Info("Starting server on", "port", cfg.Port)
		errChan <- server.ListenAndServe()
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case sig := <-sigint:
		logger.Debug("Caught signal", "signal", sig)
	case err := <-errChan:
		logger.Error("error listen and serve", sl.Err(err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", sl.Err(err))
	}

	logger.Info("Server stopped gracefully")

}

func setupPrettyLog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
