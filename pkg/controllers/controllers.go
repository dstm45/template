// Package controllers contient tous les controlleurs
package controllers

import "github.com/dstm45/template/pkg/services"

type UserController struct {
	UserService services.UserService
}

func NewUserController(svc services.UserService) *UserController {
	return &UserController{
		UserService: svc,
	}
}
