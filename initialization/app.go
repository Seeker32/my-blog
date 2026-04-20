package initialization

import (
	"fmt"

	"github.com/Seeker32/my-blog/internal/api/router"
	"github.com/Seeker32/my-blog/pkg"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Config     *pkg.Config
	Logger     *zap.Logger
	Repository *Repository
	Service    *Service
	Handlers   *Handlers
}

func InitializeApp(configPath string) (*App, error) {
	config, err := InitViper(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to init config: %w", err)
	}
	pkg.BLOG_CONFIG = config

	logger := InitLogger(&config.Logger)
	pkg.BLOG_Logger = logger

	db, err := InitDB(config, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %w", err)
	}

	repo := InitRepository(db)

	svc := InitService(repo)

	hdl := InitHandlers(svc)

	return &App{
		Config:     config,
		Logger:     logger,
		Repository: repo,
		Service:    svc,
		Handlers:   hdl,
	}, nil
}

func (app *App) SetupRouter() *gin.Engine {
	gin.SetMode(app.Config.Server.Mode)
	engine := gin.Default()

	router.SetupRoutes(engine, app.Handlers.UserHandler)

	return engine
}