package userdto

type NewUserDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}
