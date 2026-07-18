package services

import (
	"context"

	"github.com/dstm45/template/pkg/database"
)

type UserService interface {
	CreateUser(ctx context.Context) error
	GetUser(ctx context.Context) error
	DeleteUser(ctx context.Context) error
	UpdateUser(ctx context.Context) error
}

type userService struct {
	DB *database.Queries
}

func NewUserService(queries *database.Queries) *userService {
	return &userService{
		DB: queries,
	}
}

func (svc *userService) CreateUser(ctx context.Context) error { return nil }
func (svc *userService) GetUser(ctx context.Context) error    { return nil }
func (svc *userService) DeleteUser(ctx context.Context) error { return nil }
func (svc *userService) UpdateUser(ctx context.Context) error { return nil }
