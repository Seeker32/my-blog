package repository

import (
	"github.com/Seeker32/my-blog/internal/model"
	"github.com/Seeker32/my-blog/pkg"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	ExistsByUsername(username string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo userRepository) CreateUser(user *model.User) (*model.User, error) {
	if err := repo.db.Create(user).Error; err != nil {
		pkg.BLOG_Logger.Error("创建用户失败", zap.Error(err))
		return nil, err
	}
	return user, nil
}

func (repo userRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	if err := repo.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		pkg.BLOG_Logger.Error("查询用户失败", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}
