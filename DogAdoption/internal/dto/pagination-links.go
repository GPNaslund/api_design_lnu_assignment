package dto

type PaginationLinksDTO struct {
	Self  string `json:"self"`
	Next  string `json:"next"`
	Prev  string `json:"previous"`
	First string `json:"first"`
	Last  string `json:"last"`
}
