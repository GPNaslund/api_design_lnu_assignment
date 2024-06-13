package dto

type PaginationParams struct {
	Page  *int
	Limit *int
}

type DogsFilterParams struct {
	Breed      *string
	Gender     *string
	IsNeutered *string
	IsAdopted  *string
	ShelterId  *int
}

type DogShelterFilterParams struct {
	Country *string
	City    *string
	Name    *string
}

type QueryParams struct {
	Pagination       *PaginationParams
	DogsFilter       *DogsFilterParams
	DogShelterFilter *DogShelterFilterParams
}
