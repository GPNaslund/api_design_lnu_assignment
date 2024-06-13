package main

import (
	"1dv027/aad/db-init/data"
	db "1dv027/aad/db-init/schema"
	"1dv027/aad/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {

	envPath, err := filepath.Abs("../.env")
	if err != nil {
		log.Fatalf("Error getting absolute path to env")
	}
	err = godotenv.Load(envPath)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to a database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	err = db.CreateDogSheltersSchema(conn)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = db.CreateDogsSchema(conn)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	shelterQuery := `
	INSERT INTO DogShelters (
		name, website, country, city, address, username, password
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	cryptoService, err := service.NewCryptographyService(os.Getenv("CRYPTO_KEY"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create cryptography service")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	happyDogShelterPassword, err := cryptoService.HashPassword(os.Getenv("DOGSHELTER1_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash happy dogshelter password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	_, err = conn.Exec(ctx, shelterQuery, "Happy dogs shelter", "https://www.happydogsshelter.com", "Sweden", "Stockholm", "Happy dog street 1", "testdogshelter", happyDogShelterPassword)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add shelter to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	wowDogShelterPassword, err := cryptoService.HashPassword(os.Getenv("DOGSHELTER2_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash wow dogshelter password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, shelterQuery, "Wow dog shelter", "https://www.wowdogsshelter.com", "Norway", "Oslo", "Wow street 2", "wowdogshelter", wowDogShelterPassword)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add shelter to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	outdoorspalacePassword, err := cryptoService.HashPassword(os.Getenv("DOGSHELTER3_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash outdoors palace password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, shelterQuery, "Outdoors palace", "https://www.outdoorspalace.com", "Sweden", "Stockholm", "Outdoors creek 2", "outdoorspalace", outdoorspalacePassword)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add shelter to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	for i := 0; i < 20; i++ {
		dog := data.GenerateDog(1)
		dogQuery := `
		INSERT INTO Dogs (
			name, description, birth_date, breed, is_neutered, shelter_id, image_url, adoption_fee, is_adopted, friendly_with, gender
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		_, err = conn.Exec(ctx, dogQuery,
			dog.Name,
			dog.Description,
			dog.BirthDate,
			dog.Breed,
			dog.IsNeutered,
			dog.ShelterID,
			dog.ImageURL,
			dog.AdoptionFee,
			dog.IsAdopted,
			dog.FriendlyWith,
			dog.Gender,
		)
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed to add dog to table")
			fmt.Fprint(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	err = db.CreateAdminsSchema(conn)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create admin table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	adminPassword, err := cryptoService.HashPassword(os.Getenv("ADMIN_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash admin password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, `INSERT INTO Admins (username, password) VALUES ($1, $2)`, "testadmin", adminPassword)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add admin to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = db.CreateUsersSchema(conn)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create users table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	testuserPassword, err := cryptoService.HashPassword(os.Getenv("USER1_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash test user password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, `INSERT INTO Users (username, password) VALUES ($1, $2)`, "testuser", testuserPassword)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add testuser to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	testuser2Password, err := cryptoService.HashPassword(os.Getenv("USER2_PASSWORD"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to hash test user password")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	_, err = conn.Exec(ctx, `INSERT INTO Users (username, password) VALUES ($1, $2)`, "testuser2", testuser2Password)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add testuser2 to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	err = db.CreateUserWebhooksSchema(conn)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to create webhooks table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}

	webhookUrl := os.Getenv("WEBHOOK_URL")
	clientSecret, err := cryptoService.EncryptPlainText(os.Getenv("WEBHOOK1_SECRET"))
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to encrypt webhooks client secret")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
	webhookActions := []string{"new_dog_added"}

	_, err = conn.Exec(ctx, `INSERT INTO UserWebhooks (webhook_endpoint, client_secret, webhook_actions, user_id) VALUES ($1, $2, $3, $4)`, webhookUrl, clientSecret, webhookActions, 1)
	if err != nil {
		fmt.Fprint(os.Stderr, "Failed to add testuser2 to table")
		fmt.Fprint(os.Stderr, err.Error())
		os.Exit(1)
	}
}
