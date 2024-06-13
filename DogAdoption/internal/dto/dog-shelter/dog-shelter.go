package dogshelterdto

type DogShelterDTO struct {
	Id      int                `json:"id"`
	Name    string             `json:"name"`
	Website string             `json:"website"`
	Country string             `json:"country"`
	City    string             `json:"city"`
	Address string             `json:"address"`
	Links   DogShelterDtoLinks `json:"links"`
}

type DogShelterDtoLinks struct {
	SelfLink string `json:"self_link"`
	DogsLink string `json:"dogs_link"`
}
