package config

import (
	dataaccess "1dv027/aad/internal/data-access"
	"1dv027/aad/internal/handlers"
	apihandler "1dv027/aad/internal/handlers/api"
	authhandler "1dv027/aad/internal/handlers/auth"
	doghandler "1dv027/aad/internal/handlers/dog"
	dogshelterhandler "1dv027/aad/internal/handlers/dog-shelter"
	"1dv027/aad/internal/handlers/middleware"
	userhandler "1dv027/aad/internal/handlers/user"
	usermehandler "1dv027/aad/internal/handlers/user/me"
	userwebhookhandler "1dv027/aad/internal/handlers/user/webhook"
	"1dv027/aad/internal/repository"
	"1dv027/aad/internal/service"
	apiservice "1dv027/aad/internal/service/api"
	authservice "1dv027/aad/internal/service/auth"
	dogsheltersservice "1dv027/aad/internal/service/dog-shelter"
	dogsservice "1dv027/aad/internal/service/dogs"
	usersservice "1dv027/aad/internal/service/users"
	userwebhookservice "1dv027/aad/internal/service/users/webhook"
	"1dv027/aad/internal/webhook"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ContainerConfig struct {
	DatabaseConnector     *pgxpool.Pool
	CryptographySecretKey string
	BasePath              string
	JwtSigningKey         string
}

// Setup for the IoC container
func SetupContainer(config ContainerConfig) *Container {
	c := NewContainer()

	// Data access layer
	c.ProvideSingleton("AdminsDataAccess", func() any {
		return dataaccess.NewAdminsDataAccess(config.DatabaseConnector)
	})
	c.ProvideSingleton("DogsDataAccess", func() any {
		return dataaccess.NewDogsDataAccess(config.DatabaseConnector)
	})
	c.ProvideSingleton("DogSheltersDataAccess", func() any {
		return dataaccess.NewDogSheltersDataAccess(config.DatabaseConnector)
	})
	c.ProvideSingleton("UserWebhooksDataAccess", func() any {
		return dataaccess.NewUsersWebhookDataAccess(config.DatabaseConnector)
	})
	c.ProvideSingleton("UsersDataAccess", func() any {
		return dataaccess.NewUserDataAccess(config.DatabaseConnector)
	})

	// Repository layer
	c.ProvideSingleton("DogSheltersRepository", func() any {
		dogSheltersDataAccess := c.Resolve("DogSheltersDataAccess", Singleton).(repository.DogSheltersDataAccess)
		return repository.NewDogSheltersRepository(dogSheltersDataAccess)
	})
	c.ProvideSingleton("DogsRepository", func() any {
		dogsDataAccess := c.Resolve("DogsDataAccess", Singleton).(repository.DogsDataAccess)
		return repository.NewDogsRepository(dogsDataAccess)
	})
	c.ProvideSingleton("LoginRepository", func() any {
		adminsDataAccess := c.Resolve("AdminsDataAccess", Singleton).(repository.GetAdminsDataAccess)
		dogSheltersDataAccess := c.Resolve("DogSheltersDataAccess", Singleton).(repository.GetDogSheltersDataAccess)
		usersDataAccess := c.Resolve("UsersDataAccess", Singleton).(repository.GetUsersDataAccess)
		return repository.NewLoginRepository(adminsDataAccess, dogSheltersDataAccess, usersDataAccess)
	})
	c.ProvideSingleton("UserWebhooksRepository", func() any {
		userWebhooksDataAccess := c.Resolve("UserWebhooksDataAccess", Singleton).(repository.UserWebhooksDataAccess)
		return repository.NewUserWebhooksRepository(userWebhooksDataAccess)
	})
	c.ProvideSingleton("UsersRepository", func() any {
		usersDataAcess := c.Resolve("UsersDataAccess", Singleton).(repository.UsersDataAccess)
		return repository.NewUsersRepository(usersDataAcess)
	})

	// Service layer
	/// Util
	c.ProvideSingleton("JwtGenerator", func() any {
		return service.NewJwtService(config.JwtSigningKey)
	})
	c.ProvideSingleton("CryptographyService", func() any {
		cryptoService, err := service.NewCryptographyService(config.CryptographySecretKey)
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed to create cryptography service")
			os.Exit(1)
		}
		return cryptoService
	})
	c.ProvideSingleton("HateoasLinkGenerator", func() any {
		return service.NewHateoasLinkGenerator(config.BasePath)
	})
	c.ProvideSingleton("WebhookDispatcher", func() any {
		userWebhooksRepo := c.Resolve("UserWebhooksRepository", Singleton).(webhook.UserWebhooksRepository)
		cryptoService := c.Resolve("CryptographyService", Singleton).(webhook.CryptographyService)
		return webhook.NewWebhookDispatcher(userWebhooksRepo, cryptoService)
	})

	// Api
	c.ProvideSingleton("ApiService", func() any {
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(apiservice.ApiLinksGenerator)
		return apiservice.NewApiService(linkGenerator)
	})

	/// Auth
	c.ProvideSingleton("AuthLoginService", func() any {
		loginRepository := c.Resolve("LoginRepository", Singleton).(authservice.LoginRepository)
		jwtGenerator := c.Resolve("JwtGenerator", Singleton).(authservice.JwtGenerator)
		cryptoService := c.Resolve("CryptographyService", Singleton).(authservice.CryptographyService)
		return authservice.NewLoginService(loginRepository, jwtGenerator, cryptoService)
	})
	/// Dogs
	c.ProvideSingleton("DogsDeleteService", func() any {
		dogsRepo := c.Resolve("DogsRepository", Singleton).(dogsservice.DeleteDogRepository)
		return dogsservice.NewDeleteDogService(dogsRepo)
	})
	c.ProvideSingleton("DogsGetByIdService", func() any {
		dogsRepo := c.Resolve("DogsRepository", Singleton).(dogsservice.GetDogByIdRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsservice.GetDogByIdLinkGenerator)
		return dogsservice.NewGetDogByIdService(dogsRepo, linkGenerator)
	})
	c.ProvideSingleton("DogsGetService", func() any {
		dogsRepo := c.Resolve("DogsRepository", Singleton).(dogsservice.GetDogsRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsservice.GetDogsLinkGenerator)
		return dogsservice.NewGetDogsService(dogsRepo, linkGenerator)
	})
	c.ProvideSingleton("DogsPostService", func() any {
		dogsRepo := c.Resolve("DogsRepository", Singleton).(dogsservice.PostDogsRepository)
		webhookDispatcher := c.Resolve("WebhookDispatcher", Singleton).(dogsservice.NewDogWebhookDispatcher)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsservice.PostDogsLinkGenerator)
		return dogsservice.NewPostDogService(dogsRepo, linkGenerator, webhookDispatcher)
	})
	c.ProvideSingleton("DogsPutService", func() any {
		dogsRepo := c.Resolve("DogsRepository", Singleton).(dogsservice.PutDogsRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsservice.PutDogsLinkGenerator)
		return dogsservice.NewPutDogService(dogsRepo, linkGenerator)
	})
	/// DogShelters
	c.ProvideSingleton("DogSheltersDeleteService", func() any {
		dogSheltersRepo := c.Resolve("DogSheltersRepository", Singleton).(dogsheltersservice.DeleteDogSheltersRepository)
		return dogsheltersservice.NewDeleteDogSheltersService(dogSheltersRepo)
	})
	c.ProvideSingleton("DogSheltersGetByIdService", func() any {
		dogSheltersRepo := c.Resolve("DogSheltersRepository", Singleton).(dogsheltersservice.GetDogSheltersByIdRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsheltersservice.GetDogSheltersByIdLinkGenerator)
		return dogsheltersservice.NewGetDogSheltersByIdService(dogSheltersRepo, linkGenerator)
	})
	c.ProvideSingleton("DogSheltersGetService", func() any {
		dogSheltersRepo := c.Resolve("DogSheltersRepository", Singleton).(dogsheltersservice.GetDogSheltersRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsheltersservice.GetDogSheltersLinkGenerator)
		return dogsheltersservice.NewGetDogSheltersService(dogSheltersRepo, linkGenerator)
	})
	c.ProvideSingleton("DogSheltersPostService", func() any {
		dogSheltersRepo := c.Resolve("DogSheltersRepository", Singleton).(dogsheltersservice.PostDogSheltersRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsheltersservice.PostDogSheltersLinkGenerator)
		cryptoService := c.Resolve("CryptographyService", Singleton).(dogsheltersservice.PostDogSheltersCryptographyService)
		return dogsheltersservice.NewPostDogSheltersService(dogSheltersRepo, linkGenerator, cryptoService)
	})
	c.ProvideSingleton("DogSheltersPutService", func() any {
		dogSheltersRepo := c.Resolve("DogSheltersRepository", Singleton).(dogsheltersservice.PutDogSheltersRepository)
		linkGenerator := c.Resolve("HateoasLinkGenerator", Singleton).(dogsheltersservice.PutDogSheltersLinkGenerator)
		return dogsheltersservice.NewPutDogSheltersService(dogSheltersRepo, linkGenerator)
	})
	/// Users
	c.ProvideSingleton("UsersDeleteService", func() any {
		userRepo := c.Resolve("UsersRepository", Singleton).(usersservice.DeleteUsersRepository)
		return usersservice.NewDeleteUsersService(userRepo)
	})
	c.ProvideSingleton("UsersGetMeService", func() any {
		userRepo := c.Resolve("UsersRepository", Singleton).(usersservice.GetUsersMeRepository)
		return usersservice.NewGetUsersMeService(userRepo)
	})
	c.ProvideSingleton("UsersPostService", func() any {
		userRepo := c.Resolve("UsersRepository", Singleton).(usersservice.PostUsersRepository)
		cryptoService := c.Resolve("CryptographyService", Singleton).(usersservice.PostUsersCryptographyService)
		return usersservice.NewPostUsersService(userRepo, cryptoService)
	})
	//// Webhooks
	c.ProvideSingleton("UserWebhooksDataValidator", func() any {
		return userwebhookservice.NewWebhookDataValidatorService()
	})
	c.ProvideSingleton("UserWebhooksDeleteService", func() any {
		userWebhookRepo := c.Resolve("UserWebhooksRepository", Singleton).(userwebhookservice.DeleteUserWebhookRepository)
		return userwebhookservice.NewDeleteWebhookService(userWebhookRepo)
	})
	c.ProvideSingleton("UserWebhooksGetService", func() any {
		userWebhookRepo := c.Resolve("UserWebhooksRepository", Singleton).(userwebhookservice.GetUserWebhookRepository)
		return userwebhookservice.NewGetUserWebhookService(userWebhookRepo)
	})
	c.ProvideSingleton("UserWebhooksPostService", func() any {
		userWebhookRepo := c.Resolve("UserWebhooksRepository", Singleton).(userwebhookservice.PostUserWebhooksRepository)
		dataValidator := c.Resolve("UserWebhooksDataValidator", Singleton).(userwebhookservice.PostUserWebhooksDataValidator)
		cryptoService := c.Resolve("CryptographyService", Singleton).(userwebhookservice.PostUserWebhooksCryptographyService)
		return userwebhookservice.NewPostUserWebhookService(userWebhookRepo, dataValidator, cryptoService)
	})
	c.ProvideSingleton("UserWebhookPutService", func() any {
		userWebhookRepo := c.Resolve("UserWebhooksRepository", Singleton).(userwebhookservice.PutUserWebhooksRepository)
		dataValidator := c.Resolve("UserWebhooksDataValidator", Singleton).(userwebhookservice.PutUserWebhooksDataValidator)
		cryptoService := c.Resolve("CryptographyService", Singleton).(userwebhookservice.PutUserWebhooksCryptographyService)
		return userwebhookservice.NewPutUserWebhookService(userWebhookRepo, dataValidator, cryptoService)
	})
	// Handlers
	/// Util
	c.ProvideSingleton("RequestBodyValidator", func() any {
		return handlers.NewRequestBodyValidator()
	})

	/// Api
	c.ProvideTransient("ApiHandler", func() any {
		apiService := c.Resolve("ApiService", Singleton).(apihandler.ApiService)
		return apihandler.NewApiHandler(apiService)
	})

	/// Auth
	c.ProvideTransient("AuthLoginHandler", func() any {
		authService := c.Resolve("AuthLoginService", Singleton).(authhandler.AuthService)
		reqBodyValidator := c.Resolve("RequestBodyValidator", Singleton).(authhandler.RequestBodyValidator)
		return authhandler.NewLoginHandler(authService, reqBodyValidator)
	})
	/// Dog
	c.ProvideTransient("DogDeleteHandler", func() any {
		service := c.Resolve("DogsDeleteService", Singleton).(doghandler.DeleteDogService)
		return doghandler.NewDeleteDogHandler(service)
	})
	c.ProvideTransient("DogGetByIdHandler", func() any {
		service := c.Resolve("DogsGetByIdService", Singleton).(doghandler.GetDogByIdService)
		return doghandler.NewGetDogByIdHandler(service)
	})
	c.ProvideTransient("DogGetHandler", func() any {
		service := c.Resolve("DogsGetService", Singleton).(doghandler.GetDogsService)
		return doghandler.NewGetDogsHandler(service)
	})
	c.ProvideTransient("DogPostHandler", func() any {
		service := c.Resolve("DogsPostService", Singleton).(doghandler.PostDogService)
		return doghandler.NewPostDogHandler(service)
	})
	c.ProvideTransient("DogPutHandler", func() any {
		service := c.Resolve("DogsPutService", Singleton).(doghandler.PutDogService)
		return doghandler.NewPutDogHandler(service)
	})
	/// DogShelter
	c.ProvideTransient("DogShelterDeleteHandler", func() any {
		service := c.Resolve("DogSheltersDeleteService", Singleton).(dogshelterhandler.DeleteDogShelterService)
		return dogshelterhandler.NewDeleteShelterHandler(service)
	})
	c.ProvideTransient("DogShelterGetByIdHandler", func() any {
		service := c.Resolve("DogSheltersGetByIdService", Singleton).(dogshelterhandler.GetDogShelterByIdService)
		return dogshelterhandler.NewGetDogShelterByIdHandler(service)
	})
	c.ProvideTransient("DogShelterGetHandler", func() any {
		service := c.Resolve("DogSheltersGetService", Singleton).(dogshelterhandler.GetDogShelterService)
		return dogshelterhandler.NewGetDogShelterHandler(service)
	})
	c.ProvideTransient("DogShelterPostHandler", func() any {
		service := c.Resolve("DogSheltersPostService", Singleton).(dogshelterhandler.PostDogShelterService)
		return dogshelterhandler.NewPostShelterHandler(service)
	})
	c.ProvideTransient("DogShelterPutHandler", func() any {
		service := c.Resolve("DogSheltersPutService", Singleton).(dogshelterhandler.PutDogShelterService)
		return dogshelterhandler.NewPutDogShelterHandler(service)
	})
	/// Middleware
	c.ProvideTransient("AuthMiddleware", func() any {
		jwtGenerator := c.Resolve("JwtGenerator", Singleton).(middleware.JwtService)
		return middleware.NewAuthMiddleware(jwtGenerator)
	})
	c.ProvideTransient("QueryParamsMiddleware", func() any {
		return middleware.NewQueryParamsValidator()
	})
	/// User
	//// Me
	c.ProvideTransient("UserGetMeHandler", func() any {
		userMeService := c.Resolve("UsersGetMeService", Singleton).(usermehandler.GetUserMeService)
		return usermehandler.NewGetUserMeHandler(userMeService)
	})
	//// Webhook
	c.ProvideTransient("UserWebhookDeleteHandler", func() any {
		service := c.Resolve("UserWebhooksDeleteService", Singleton).(userwebhookhandler.DeleteUserWebhookService)
		return userwebhookhandler.NewDeleteUserWebhookHandler(service)
	})
	c.ProvideTransient("UserWebhookGetHandler", func() any {
		service := c.Resolve("UserWebhooksGetService", Singleton).(userwebhookhandler.GetUserWebhookService)
		return userwebhookhandler.NewGetUserWebhookHandler(service)
	})
	c.ProvideTransient("UserWebhookPostHandler", func() any {
		service := c.Resolve("UserWebhooksPostService", Singleton).(userwebhookhandler.PostUserWebhookService)
		return userwebhookhandler.NewPostUserWebhookHandler(service)
	})
	c.ProvideTransient("UserWebhookPutHandler", func() any {
		service := c.Resolve("UserWebhookPutService", Singleton).(userwebhookhandler.PutUserWebhookService)
		return userwebhookhandler.NewPutUserWebhookHandler(service)
	})
	/// User
	c.ProvideTransient("UserDeleteHandler", func() any {
		service := c.Resolve("UsersDeleteService", Singleton).(userhandler.DeleteUserService)
		return userhandler.NewDeleteUserHandler(service)
	})
	c.ProvideTransient("UserPostHandler", func() any {
		service := c.Resolve("UsersPostService", Singleton).(userhandler.PostUserService)
		return userhandler.NewPostUserHandler(service)
	})

	return c
}
