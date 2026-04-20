package service

import (
	"errors"

	"github.com/Seeker32/my-blog/internal/model"
	"github.com/Seeker32/my-blog/internal/repository"
	"github.com/Seeker32/my-blog/internal/request"
	"github.com/Seeker32/my-blog/pkg"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(req request.UserRegisterRequest) (*model.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (service *userService) CreateUser(req request.UserRegisterRequest) (*model.User, error) {
	// 验证两次密码是否一致
	if req.Password != req.CheckPassword {
		return nil, errors.New("两次密码不一致")
	}

	// 检查用户名是否已存在
	exists, err := service.repo.ExistsByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		pkg.BLOG_Logger.Error("密码加密失败", zap.Error(err))
		return nil, err
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     2, // 默认为普通用户
	}

	return service.repo.CreateUser(&user)
}
