package initialization

import "github.com/Seeker32/my-blog/internal/api/handler"

type Handlers struct {
	UserHandler handler.UserHandler
}

func InitHandlers(svc *Service) *Handlers {
	return &Handlers{
		UserHandler: handler.NewUserHandler(svc.UserService),
	}
}