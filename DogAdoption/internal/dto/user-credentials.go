package dto

import "1dv027/aad/internal/model"

type UserCredentials struct {
	Username string
	UserRole model.UserRole
	Id       int
}
