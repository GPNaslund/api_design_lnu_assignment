package userwebhookdto

import dogdto "1dv027/aad/internal/dto/dog"

type NewDogDispatchDTO struct {
	NewDog dogdto.DogDTO `json:"new_dog"`
	Secret string        `json:"secret"`
}
