package dogdto

import "1dv027/aad/internal/dto"

type DogsAndPaginationLinksDTO struct {
	Dogs            []DogDTO               `json:"dogs"`
	PaginationLinks dto.PaginationLinksDTO `json:"pagination_links"`
}
