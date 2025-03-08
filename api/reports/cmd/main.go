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

	"github.com/gin-gonic/gin"

	"reports/handlers/get_reports"
	"reports/middlewares"
	"reports/usecases/token_verifier"
)

const (
	readHeaderTimeout           = 5 * time.Second
	shutdownTimeout             = 2 * time.Second
	signalChannelBufferCapacity = 1
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	handler := get_reports.NewHandler(token_verifier.MustCreateTokenVerifier())

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent) // 204 No Content
			return
		}
		c.Next()
	})

	router.GET("/reports", middlewares.AuthMiddleware(), handler.GetReports)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to Reports API!",
		})
	})
	srv := &http.Server{
		Addr:              ":8000",
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	sigchan := make(chan os.Signal, signalChannelBufferCapacity)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидание сигнала от системы (например, при завершении работы)
	<-sigchan
	fmt.Println("Shutdown signal received, initiating graceful shutdown...")

	// Отмена контекста при получении сигнала
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownCancel()

	// Ожидаем завершения работы сервера с таймаутом
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully shut down")
	}
}
