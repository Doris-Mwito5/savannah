package main

import (
	"context"
	config "github/Doris-Mwito5/savannah-pos/config"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/loggers"
	"github/Doris-Mwito5/savannah-pos/web/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var (
	PORT = "8080"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, relying on system env")
	}

	if err := config.LoadEnvConfig(); err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}


	loggers.InitLogger("savannah-pos")
	// Init db
	dB := db.InitDB()
	defer dB.Close()

	loggers.Info("Starting rh-backend")

	//domain store
	domainStore := domain.NewStore()

	appRouter := routes.BuildRouter(
		dB,
		domainStore,
	)

	server := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: appRouter,
	}

	done := make(chan struct{})

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		loggers.Info("Process terminated...shutting down")

		if err := server.Shutdown(context.Background()); err != nil {
			loggers.Fatalf("Server shut down error: %v", err)
		}

		close(done)
	}()

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			loggers.Info("Server shut down")
		} else {
			loggers.Fatal("Server shut down unexpectedly!")
		}
	}

	timeout := 30 * time.Second
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	code := 0
	select {
	case <-sigint:
		code = 1
		loggers.Info("Process forcibly terminated")
	case <-time.After(timeout):
		code = 1
		loggers.Info("Shutdown timeout. Forcibly shutting down...")
	case <-done:
		loggers.Info("Shutdown completed...")
	}

	loggers.Info("Server exiting...")

	os.Exit(code)

}
