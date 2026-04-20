package pkg

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BLOG_Logger *zap.Logger
	BLOG_CONFIG *Config
	BLOG_DB     *gorm.DB
)