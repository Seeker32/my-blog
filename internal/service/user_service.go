package service

import (
	"errors"

	"github.com/Seeker32/my-blog/internal/model"
	"github.com/Seeker32/my-blog/internal/request"
	"github.com/Seeker32/my-blog/pkg"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(req request.UserRegisterRequest) (model.User, error) {
	// 验证两次密码是否一致
	if req.Password != req.CheckPassword {
		return model.User{}, errors.New("两次密码不一致")
	}

	// 检查用户名是否已存在
	var count int64
	if err := pkg.BLOG_DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		pkg.BLOG_Logger.Error("查询用户失败", zap.Error(err))
		return model.User{}, err
	}
	if count > 0 {
		return model.User{}, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		pkg.BLOG_Logger.Error("密码加密失败", zap.Error(err))
		return model.User{}, err
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     2, // 默认为普通用户
	}

	if err := pkg.BLOG_DB.Create(&user).Error; err != nil {
		pkg.BLOG_Logger.Error("创建用户失败", zap.Error(err))
		return model.User{}, err
	}

	return user, nil
}
