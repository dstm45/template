// Package controllers contient tous les controlleurs
package controllers

import "github.com/dstm45/template/pkg/services"

type UserController struct {
	UserService services.IUserService
}

func NewUserController(svc services.IUserService) *UserController {
	return &UserController{
		UserService: svc,
	}
}
