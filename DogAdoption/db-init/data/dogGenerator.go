package data

import (
	"math/rand/v2"
	"time"
)

type GeneratedDog struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	BirthDate    time.Time `json:"birthDate"`
	Breed        string    `json:"breed"`
	IsNeutered   bool      `json:"isNeutered"`
	ShelterID    int       `json:"shelterId"`
	ImageURL     string    `json:"imageUrl,omitempty"`
	AdoptionFee  int       `json:"adoptionFee"`
	IsAdopted    bool      `json:"isAdopted"`
	FriendlyWith string    `json:"friendlyWith,omitempty"`
	Gender       string    `json:"gender"`
}

func GenerateDog(amountOfShelters int) GeneratedDog {
	var dog GeneratedDog
	dog.Name = names[rand.IntN(len(names))]
	dog.Description = description[rand.IntN(len(description))]
	dog.BirthDate = generateRandomDate(2010)
	dog.Breed = breeds[rand.IntN(len(breeds))]
	dog.IsNeutered = isNeutered[rand.IntN(2)]
	dog.ShelterID = rand.IntN(amountOfShelters) + 1
	dog.ImageURL = imageUrls[rand.IntN(len(imageUrls))]
	dog.AdoptionFee = 5000 + rand.IntN(5000)
	dog.IsAdopted = isAdopted[rand.IntN(2)]
	dog.FriendlyWith = friendlyWith[rand.IntN(len(friendlyWith))]
	dog.Gender = gender[rand.IntN(2)]

	return dog
}

var (
	names = []string{
		"Bailey",
		"Rufus",
		"Happy",
		"Peppy",
		"Luna",
		"Lucky",
		"Rocky",
		"Sadie",
		"Milo",
		"Ruby",
	}
	description = []string{
		"Enjoys long walks in the park and loves to play fetch. Not a fan of rainy days.",
		"Adores cuddling and quiet evenings. Slightly wary of strangers but warms up quickly.",
		"High-energy and loves to explore. Needs plenty of exercise and enjoys agility training.",
		"Independent spirit, but loyal. Prefers not to be left alone for too long.",
		"Great with children and other pets. Loves being the center of attention.",
		"Shy and reserved, especially around loud noises. Prefers calm environments.",
		"Gentle giant who loves nothing more than lounging by your side. Not keen on excessive exercise.",
		"Curious and playful, always getting into mischief. Loves toys and puzzles.",
		"Protective and alert, making a great watchdog. Takes a while to make new friends.",
		"Adventurous and loves outdoor activities. Not a fan of staying indoors for too long.",
	}
	breeds     = []string{"Labrador", "Beagle", "Bulldog", "Poodle", "Shepherd", "Finnish Lapphund"}
	isNeutered = []bool{true, false}
	imageUrls  = []string{
		"https://shelterimages.example.com/dogs/dog_1.jpg",
		"https://shelterimages.example.com/dogs/dog_2.jpg",
		"https://shelterimages.example.com/dogs/dog_3.jpg",
		"https://shelterimages.example.com/dogs/dog_4.jpg",
		"https://shelterimages.example.com/dogs/dog_5.jpg",
		"https://shelterimages.example.com/dogs/dog_6.jpg",
		"https://shelterimages.example.com/dogs/dog_7.jpg",
		"https://shelterimages.example.com/dogs/dog_8.jpg",
		"https://shelterimages.example.com/dogs/dog_9.jpg",
		"https://shelterimages.example.com/dogs/dog_10.jpg",
	}
	friendlyWith = []string{"Children and dogs", "Cats", "City and children", "Dogs", "None", "Cats and dogs"}
	isAdopted    = []bool{true, false}
	gender       = []string{"male", "female"}
)

func generateRandomDate(startYear int) time.Time {
	start := time.Date(startYear, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Now()

	days := end.Sub(start).Hours() / 24

	randomDays := rand.IntN(int(days))

	randomDate := start.AddDate(0, 0, randomDays)

	return randomDate
}
