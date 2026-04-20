package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/Seeker32/my-blog/initialization"
	"go.uber.org/zap"
)

//	@title			Swagger API
//	@version		1.0
//	@description	This is a blog API.
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	configPath, err := filepath.Abs("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to get config path: %v", err))
	}

	app, err := initialization.InitializeApp(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize app: %v", err))
	}
	defer app.Logger.Sync()

	app.Logger.Info("Starting blog service...")

	router := app.SetupRouter()

	addr := fmt.Sprintf(":%d", app.Config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		app.Logger.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.Logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	app.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		app.Logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	app.Logger.Info("Server exited")
}