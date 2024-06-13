package dogdto

import "time"

type UpdateDogDTO struct {
	Name         *string    `json:"name"`
	Description  *string    `json:"description"`
	BirthDate    *time.Time `json:"birth_date"`
	Breed        *string    `json:"breed"`
	IsNeutered   *bool      `json:"is_neutered"`
	ImageUrl     *string    `json:"image_url"`
	AdoptionFee  *int       `json:"adoption_fee"`
	IsAdopted    *bool      `json:"is_adopted"`
	FriendlyWith *string    `json:"friendly_with"`
	Gender       *string    `json:"gender"`
}
