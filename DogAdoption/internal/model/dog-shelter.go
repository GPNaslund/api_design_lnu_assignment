package model

type DogShelter struct {
	Id       int
	Name     string
	Website  string
	Country  string
	City     string
	Address  string
	Username string
	Password string
}

func (d *DogShelter) ToJson() map[string]any {
	return map[string]any{
		"id":       d.Id,
		"name":     d.Name,
		"website":  d.Website,
		"country":  d.Country,
		"city":     d.City,
		"address":  d.Address,
		"username": d.Username,
		"password": d.Password,
	}
}
