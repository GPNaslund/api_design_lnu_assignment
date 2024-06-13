package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func CreateDogSheltersSchema(conn *pgx.Conn) error {
	ctx := context.Background()
	query := `
	CREATE TABLE IF NOT EXISTS DogShelters (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		website TEXT,
		country TEXT NOT NULL,
		city TEXT NOT NULL,
		address TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);
	`

	_, err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating DogShelters schema: %v", err)
	}
	return nil
}

func CreateDogsSchema(conn *pgx.Conn) error {
	ctx := context.Background()
	query := `
	CREATE TABLE IF NOT EXISTS Dogs (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		birth_date DATE NOT NULL,
		breed TEXT NOT NULL,
		is_neutered BOOLEAN NOT NULL,
		shelter_id INTEGER NOT NULL,
		image_url TEXT,
		adoption_fee INTEGER NOT NULL,
		is_adopted BOOLEAN NOT NULL,
		friendly_with TEXT,
		gender TEXT NOT NULL CHECK (gender IN ('male', 'female')),
		FOREIGN KEY (shelter_id) REFERENCES DogShelters(id) ON DELETE CASCADE
	);
	`

	_, err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating Dog schema: %v", err)
	}
	return nil
}

func CreateAdminsSchema(conn *pgx.Conn) error {
	ctx := context.Background()
	query := `
	CREATE TABLE IF NOT EXISTS Admins (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	)
	`
	_, err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating Admin schema: %v", err)
	}
	return nil
}

func CreateUsersSchema(conn *pgx.Conn) error {
	ctx := context.Background()
	query := `
	CREATE TABLE IF NOT EXISTS Users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`
	_, err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating User schema: %v", err)
	}
	return nil
}

func CreateUserWebhooksSchema(conn *pgx.Conn) error {
	ctx := context.Background()
	query := `
	CREATE TABLE IF NOT EXISTS UserWebhooks (
		id SERIAL PRIMARY KEY,
		webhook_endpoint TEXT NOT NULL,
		client_secret TEXT NOT NULL,
		webhook_actions TEXT[] NOT NULL,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES Users(id) ON DELETE CASCADE
	);
	`

	_, err := conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("error creating Webhook schema: %v", err)
	}
	return nil
}
