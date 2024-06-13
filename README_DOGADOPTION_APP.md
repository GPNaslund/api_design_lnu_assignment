# DogAdoptionAPI
## Info
The api is thought of being a national/global dog adoption api, where dogshelters can register (through admins only), and register dogs that are up for adoption or already adopted. There is a possibility for a user to register a generic user account, and to register a webhook to be notified when a new dog is added. The user registration and handling is rudimentary and implemented with the purpose of being able to register a webhook. In a real world scenario, more information would be collected and handled better. The webhook functionality is rudimentary and supports only dispatching when a new dog is added, but the structure is built to be able to extend this feature.

## Database seeding
- go run db-init/main.go for creation of database schemas and database seeding

## Run the application
- go run cmd/main.go to start the application/server

## Run postman tests
- The postman test suite is meant to be ran from the whole collection, the first folder of tests from the /auth folder sets necessary api keys.

## .env
There is "sensitive" information in the dogman collection (admin password etc). This data is set by the .env variable, and to allow examinators to test the api without any other configuration, this was needed. In a real world scenario any sensitive information would be passed securely and then inserted in the collection for testing.
