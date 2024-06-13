package dogshelterdto

import "1dv027/aad/internal/dto"

type DogSheltersAndPaginationLinksDTO struct {
	DogShelterData  []DogShelterDTO        `json:"dog_shelter_data"`
	PaginationLinks dto.PaginationLinksDTO `json:"pagination_links"`
}
