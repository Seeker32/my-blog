package initialization

import (
	"github.com/Seeker32/my-blog/internal/model"
	"github.com/Seeker32/my-blog/pkg"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *pkg.Config, logger *zap.Logger) (*gorm.DB, error) {
	dsn := cfg.Database.GetDatabaseUrl()

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logger.Error("failed to connect database", zap.Error(err))
		return nil, err
	}

	pkg.BLOG_DB = db
	logger.Info("Database connected successfully")

	err = db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Article{})
	if err != nil {
		logger.Error("failed to migrate database", zap.Error(err))
		return nil, err
	}

	logger.Info("Database migration completed")
	return db, nil
}