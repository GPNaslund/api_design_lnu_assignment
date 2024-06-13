package dogshelterdto

import "1dv027/aad/internal/model"

type GetDogSheltersQueryResponseDTO struct {
	DogShelters          []model.DogShelter
	TotalAmountAvailable int
}
