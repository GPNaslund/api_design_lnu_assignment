package dto

type EntryPointLinksDTO struct {
	OpenApi           string `json:"open_api"`
	DogsUrl           string `json:"dogs_url"`
	DogSheltersUrl    string `json:"dog_shelters_url"`
	AuthenticationUrl string `json:"authentication_url"`
	UsersUrl          string `json:"users_url"`
}
