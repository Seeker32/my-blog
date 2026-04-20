package initialization

import (
	"github.com/Seeker32/my-blog/internal/repository"
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository repository.UserRepository
}

func InitRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository: repository.NewUserRepository(db),
	}
}