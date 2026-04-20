package initialization

import "github.com/Seeker32/my-blog/internal/service"

type Service struct {
	UserService service.UserService
}

func InitService(repo *Repository) *Service {
	return &Service{
		UserService: service.NewUserService(repo.UserRepository),
	}
}