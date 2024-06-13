package model

import "fmt"

type UserRole string

const (
	ADMIN      UserRole = "admin"
	DOGSHELTER UserRole = "dog_shelter"
	USER       UserRole = "user"
)

func StringToUserRole(roleString string) (UserRole, error) {
	switch roleString {
	case string(ADMIN):
		return ADMIN, nil
	case string(DOGSHELTER):
		return DOGSHELTER, nil
	case string(USER):
		return USER, nil
	default:
		return "", fmt.Errorf("invalid UserRole: %s", roleString)
	}
}
