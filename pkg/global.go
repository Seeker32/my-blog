package pkg

import (
	"github.com/Seeker32/my-blog/initialization"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	BLOG_Logger *zap.Logger
	BLOG_CONFIG *initialization.Config
	BLOG_DB     *gorm.DB
)
