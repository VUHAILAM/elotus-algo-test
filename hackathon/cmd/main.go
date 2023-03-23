package main

import (
	"context"
	"encoding/json"
	"hackathon/configs"
	"hackathon/pkg/account"
	"hackathon/pkg/auth"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	router := gin.Default()
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		logrus.Error("Init db error", err)
		return
	}

	authConfig, err := configs.LoadAuthConfig()
	if err != nil {
		logrus.Error("Load auth config err", err)
		return
	}

	authHandler := auth.NewAuthHandler(authConfig)
	accountService := account.NewAccountService(db, authHandler)

	router.POST("/login", func(ctx *gin.Context) {
		req := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		err := json.NewDecoder(ctx.Request.Body).Decode(&req)
		if err != nil {
			logrus.Error("Parse Login request error", err)
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		accessToken, refreshToken, err := accountService.Login(context.Background(), req.Username, req.Password)
		if err != nil {
			logrus.Error("Login error", err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	})

	router.POST("/register", func(ctx *gin.Context) {
		req := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}

		err := json.NewDecoder(ctx.Request.Body).Decode(&req)
		if err != nil {
			logrus.Error("Parse Register request error", err)
			ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		err = accountService.Register(context.Background(), req.Username, req.Password)
		if err != nil {
			logrus.Error("Register error", err)
			ctx.JSON(http.StatusInternalServerError, nil)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		logrus.Info("timeout of 5 seconds.")
	}
	logrus.Info("Server exiting")
}
