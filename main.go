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
	"github.com/Seeker32/my-blog/internal/model"
	"github.com/Seeker32/my-blog/pkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 1. 初始化配置
	configPath, err := filepath.Abs("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("failed to get config path: %v", err))
	}
	config, err := initialization.InitViper(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to init config: %v", err))
	}
	pkg.BLOG_CONFIG = config

	// 2. 初始化日志
	logger := initialization.InitLogger(&config.Logger)
	pkg.BLOG_Logger = logger
	defer logger.Sync()

	logger.Info("Starting blog service...")

	// 3. 设置 Gin 模式
	gin.SetMode(config.Server.Mode)

	// 4. 初始化数据库
	dsn := config.Database.GetDatabaseUrl()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	pkg.BLOG_DB = db
	logger.Info("Database connected successfully")

	// 5. 自动迁移数据库表
	err = db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{})
	if err != nil {
		logger.Fatal("failed to migrate database", zap.Error(err))
	}
	logger.Info("Database migration completed")

	// 6. 初始化路由
	router := gin.Default()
	apiGroup := router.Group("/api/v1")
	apiGroup.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// 7. 启动服务器（支持优雅关闭）
	addr := fmt.Sprintf(":%d", config.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Info("Server starting", zap.String("addr", addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// 8. 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	// 9. 优雅关闭，最多等待 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
