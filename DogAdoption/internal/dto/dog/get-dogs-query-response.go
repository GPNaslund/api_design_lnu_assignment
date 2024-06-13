package dogdto

import "1dv027/aad/internal/model"

type GetDogsQueryResponseDTO struct {
	Dogs                 []model.Dog
	TotalAmountAvailable int
}
