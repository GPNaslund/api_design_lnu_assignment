package model

import "time"

type Dog struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	BirthDate    time.Time `json:"birth_date"`
	Breed        string    `json:"breed"`
	IsNeutered   bool      `json:"is_neutered"`
	ShelterId    int       `json:"shelter_id"`
	ImageUrl     string    `json:"image_url"`
	AdoptionFee  int       `json:"adoption_fee"`
	IsAdopted    bool      `json:"is_adopted"`
	FriendlyWith string    `json:"friendly_with"`
	Gender       string    `json:"gender"`
}

func (d *Dog) ToJson() map[string]any {
	return map[string]any{
		"id":            d.Id,
		"name":          d.Name,
		"description":   d.Description,
		"birth_date":    d.BirthDate,
		"breed":         d.Breed,
		"is_neutered":   d.IsNeutered,
		"shelter_id":    d.ShelterId,
		"image_url":     d.ImageUrl,
		"adoption_fee":  d.AdoptionFee,
		"is_adopted":    d.IsAdopted,
		"friendly_with": d.FriendlyWith,
		"gender":        d.Gender,
	}
}
