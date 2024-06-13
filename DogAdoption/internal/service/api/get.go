package apiservice

import "1dv027/aad/internal/dto"

type ApiLinksGenerator interface {
	GenerateEntryPointLinks() dto.EntryPointLinksDTO
}

type ApiService struct {
	linksGenerator ApiLinksGenerator
}

func NewApiService(linksGenerator ApiLinksGenerator) ApiService {
	return ApiService{
		linksGenerator: linksGenerator,
	}
}

func (a ApiService) GetEntryPointLinks() dto.EntryPointLinksDTO {
	return a.linksGenerator.GenerateEntryPointLinks()
}
