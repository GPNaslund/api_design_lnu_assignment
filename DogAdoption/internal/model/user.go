package model

type User struct {
	Id       int
	Username string
	Password string
}

func (u *User) ToJson() map[string]any {
	return map[string]any{
		"id":       u.Id,
		"username": u.Username,
		"password": u.Password,
	}
}
