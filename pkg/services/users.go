package services

import (
	"context"

	"github.com/dstm45/template/pkg/database"
)

type IUserService interface {
	CreateUser(ctx context.Context) error
	GetUser(ctx context.Context) error
	DeleteUser(ctx context.Context) error
	UpdateUser(ctx context.Context) error
}

type UserService struct {
	DB *database.Queries
}

func NewUserService(queries *database.Queries) *UserService {
	return &UserService{
		DB: queries,
	}
}

func (svc UserService) CreateUser(ctx context.Context) error { return nil }
func (svc UserService) GetUser(ctx context.Context) error    { return nil }
func (svc UserService) DeleteUser(ctx context.Context) error { return nil }
func (svc UserService) UpdateUser(ctx context.Context) error { return nil }
