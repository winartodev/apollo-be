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
	"github.com/winartodev/apollo-be/core/config"
	"github.com/winartodev/apollo-be/core/helper"
	apolloMiddleware "github.com/winartodev/apollo-be/core/middleware"
	"github.com/winartodev/apollo-be/core/routes"
	_ "github.com/winartodev/apollo-be/docs"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := cfg.Database.SetupConnection()
	if err != nil {
		panic(err)
	}

	jwt, err := helper.NewJWT()
	if err != nil {
		panic(err)
	}

	redisClient, err := cfg.Redis.SetupConnection()
	if err != nil {
		panic(err)
	}

	redisUtil, err := helper.NewRedisUtil(redisClient)
	if err != nil {
		panic(err)
	}

	databaseUtil, err := helper.NewDatabaseUtil(db)
	if err != nil {
		panic(err)
	}

	apolloMiddleware.NewMiddleware(jwt)

	e := echo.New()

	e.HideBanner = true

	e.Validator = &helper.CustomValidator{
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
	e.Use(apolloMiddleware.GetAppPlatform())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	authHandler, err := auth.InitializeAuthAPI(databaseUtil, redisUtil, jwt, &cfg.Smtp, &cfg.OTP)
	if err != nil {
		panic(err)
	}

	otpHandler, err := auth.InitializeOtpAPI(databaseUtil, redisUtil, jwt, &cfg.Smtp, &cfg.OTP)
	if err != nil {
		panic(err)
	}

	userHandler, err := user.InitializeUserAPI(databaseUtil, redisUtil, jwt)
	if err != nil {
		panic(err)
	}
	
	countryHandler, err := country.InitializeCountryAPI()

	if err := routes.RegisterHandler(e, authHandler, userHandler, otpHandler, countryHandler); err != nil {
		panic(err)
	}

	shutdownChan := make(chan struct{})

	go func() {
		if err := e.Start(fmt.Sprintf(":%v", cfg.Http.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	if redisClient != nil {
		if err := redisClient.Close(); err != nil {
			e.Logger.Errorf("Redis close error: %v", err)
		}
	}

	e.Logger.Info("Server exited properly")
}
