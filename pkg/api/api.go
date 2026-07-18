package api

import (
	"github.com/dstm45/template/pkg/controllers"
	"github.com/dstm45/template/pkg/database"
	"github.com/dstm45/template/pkg/services"
)

type API struct {
	UserController *controllers.UserController
}

type Services struct {
	UserService services.UserService
}

func InitializeServices(queries *database.Queries) *Services {
	userService := services.NewUserService(queries)
	return &Services{
		UserService: userService,
	}
}

func NewAPI(svc *Services) *API {
	UserController := controllers.NewUserController(svc.UserService)
	return &API{
		UserController: UserController,
	}
}
