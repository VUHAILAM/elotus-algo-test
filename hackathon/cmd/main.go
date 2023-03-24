package main

import (
	"context"
	"encoding/json"
	"fmt"
	"hackathon/configs"
	"hackathon/pkg/account"
	"hackathon/pkg/auth"
	"hackathon/pkg/files"
	"hackathon/pkg/middleware"
	"hackathon/pkg/models"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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

	imageService := files.NewImageService(db)
	router.Use(middleware.AuthenizationMiddleware(authHandler)).POST("/upload", func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}

		if file.Size > 8000000 {
			ctx.String(http.StatusBadRequest, "file larger than 8 MB")
			return
		}

		contentType := file.Header.Get("Content-type")
		fmt.Printf("Content-type: %s", contentType)

		s := strings.Split(contentType, "/")

		if s[0] != "image" {
			ctx.String(http.StatusBadRequest, "not an image")
			return
		}

		acc, ok := ctx.Get(auth.AccountKey)
		if !ok {
			logrus.Error("can not get account from context")
			ctx.String(http.StatusInternalServerError, "can not get account from context")
		}
		username := acc.(models.Account).Username

		imageMetadata := models.Image{
			Username:    username,
			Name:        file.Filename,
			ContentType: contentType,
			Size:        file.Size,
		}

		if err := ctx.SaveUploadedFile(file, "/tmp"); err != nil {
			ctx.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		if err := imageService.Save(context.Background(), imageMetadata); err != nil {
			ctx.String(http.StatusInternalServerError, "save image information to database error: %s", err.Error())
			return
		}
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
