package dogshelterdto

type NewDogShelterDTO struct {
	Name     *string `json:"name"`
	Website  *string `json:"website"`
	Country  *string `json:"country"`
	City     *string `json:"city"`
	Address  *string `json:"address"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}
