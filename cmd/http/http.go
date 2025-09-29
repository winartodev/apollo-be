package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	config2 "github.com/winartodev/apollo-be/config"
	_ "github.com/winartodev/apollo-be/docs"
	middleware2 "github.com/winartodev/apollo-be/infrastructure/middleware"
	"github.com/winartodev/apollo-be/infrastructure/routes"
	"github.com/winartodev/apollo-be/modules/auth"
	"github.com/winartodev/apollo-be/modules/country"
	"github.com/winartodev/apollo-be/modules/user"
)

// @title			Apollo API
// @version		1.0
// @description	This is the Apollo server.
// @host			localhost:8081
// @BasePath		/api
// @security		Definitions.apikey BearerAuth
// @in				header
// @name			Authorization
// @schemes		http https
func main() {
	cfg, err := config2.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := cfg.Database.SetupConnection()
	if err != nil {
		panic(err)
	}

	redis, err := cfg.Redis.SetupConnection()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	e.HideBanner = true

	e.Validator = &config2.CustomValidator{
		Validator: validator.New(),
	}

	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader("${time_rfc3339} | ${level} | ${short_file}:${line} |")

	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${method} | ${uri} | ${status} | ${latency_human}\n",
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5, // Compression level
	}))
	e.Use(middleware2.GetAppPlatform())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authHandler, err := auth.InitializeAuthAPI(db, redis, &cfg.SMTP, &cfg.OTP)
	if err != nil {
		panic(err)
	}

	otpHandler, err := auth.InitializeOtpAPI(db, redis, &cfg.SMTP, &cfg.OTP)
	if err != nil {
		panic(err)
	}

	userHandler, err := user.InitializeUserAPI(db)
	if err != nil {
		panic(err)
	}

	countryHandler, err := country.InitializeCountryAPI()

	if err := routes.RegisterHandler(e, authHandler, userHandler, otpHandler, countryHandler); err != nil {
		panic(err)
	}

	shutdownChan := make(chan struct{})

	go func() {
		if err := e.Start(fmt.Sprintf("0.0.0.0:%v", cfg.Http.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			e.Logger.Fatalf("shutting down the server: %v", err)
			close(shutdownChan)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		e.Logger.Info("Received shutdown signal")
	case <-shutdownChan:
		e.Logger.Error("Server crashed, initiating shutdown")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	e.Logger.Info("shutting down server")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Errorf("HTTP server shutdown error: %v", err)
	}

	e.Logger.Info("Closing database connections...")
	if db != nil {
		if err := db.Close(); err != nil {
			e.Logger.Errorf("Database close error: %v", err)
		}
	}

	e.Logger.Info("Closing Redis connection...")
	if redis != nil {
		if err := redis.Close(); err != nil {
			e.Logger.Errorf("Redis close error: %v", err)
		}
	}

	e.Logger.Info("Server exited properly")
}
