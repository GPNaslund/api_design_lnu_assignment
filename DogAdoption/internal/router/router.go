package router

import (
	"1dv027/aad/internal/config"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

type IoCContainer interface {
	Resolve(name string, lifecycle config.Lifecycle) any
}

type Handler interface {
	Handle(c *fiber.Ctx) error
}

type AuthMiddleware interface {
	AuthenticateRequest(c *fiber.Ctx) error
}

type QueryParamsMiddleware interface {
	ValidateQueryParams(c *fiber.Ctx) error
}

type Router struct {
	container IoCContainer
}

func NewRouter(container IoCContainer) Router {
	return Router{
		container: container,
	}
}

func (r Router) StartRouter() {
	app := fiber.New()
	basePath := os.Getenv("ROUTER_BASE_PATH")
	cfg := swagger.Config{
		BasePath: basePath,
		FilePath: "cmd/docs/swagger.json",
		Path:     "api/v1/swagger",
		Title:    "Swagger API Docs",
		CacheAge: 3600,
	}
	app.Use(swagger.New(cfg))

	base := app.Group(basePath)
	api := base.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		apiHandler := r.container.Resolve("ApiHandler", config.Transient).(Handler)
		return apiHandler.Handle(c)
	})

	auth := v1.Group("/auth")
	auth.Post("/login", func(c *fiber.Ctx) error {
		loginHandler := r.container.Resolve("AuthLoginHandler", config.Transient).(Handler)
		return loginHandler.Handle(c)
	})

	dogs := v1.Group("/dogs")
	dogs.Delete("/:id", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		deleteDogHandler := r.container.Resolve("DogDeleteHandler", config.Transient).(Handler)
		return deleteDogHandler.Handle(c)
	})
	dogs.Put("/:id", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		putDogHandler := r.container.Resolve("DogPutHandler", config.Transient).(Handler)
		return putDogHandler.Handle(c)
	})
	dogs.Get("/:id", func(c *fiber.Ctx) error {
		getDogByIdHandler := r.container.Resolve("DogGetByIdHandler", config.Transient).(Handler)
		return getDogByIdHandler.Handle(c)
	})
	dogs.Post("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		postDogsHandler := r.container.Resolve("DogPostHandler", config.Transient).(Handler)
		return postDogsHandler.Handle(c)
	})
	dogs.Get("/", func(c *fiber.Ctx) error {
		queryParamsMiddleware := r.container.Resolve("QueryParamsMiddleware", config.Transient).(QueryParamsMiddleware)
		return queryParamsMiddleware.ValidateQueryParams(c)
	}, func(c *fiber.Ctx) error {
		getDogsHandler := r.container.Resolve("DogGetHandler", config.Transient).(Handler)
		return getDogsHandler.Handle(c)
	})

	dogshelters := v1.Group("/dogshelters")
	dogshelters.Delete("/:id", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		deleteDogSheltersHandler := r.container.Resolve("DogShelterDeleteHandler", config.Transient).(Handler)
		return deleteDogSheltersHandler.Handle(c)
	})
	dogshelters.Put("/:id", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		putDogsheltersHandler := r.container.Resolve("DogShelterPutHandler", config.Transient).(Handler)
		return putDogsheltersHandler.Handle(c)
	})
	dogshelters.Get("/:id", func(c *fiber.Ctx) error {
		getDogSheltersByIdHandler := r.container.Resolve("DogShelterGetByIdHandler", config.Transient).(Handler)
		return getDogSheltersByIdHandler.Handle(c)
	})
	dogshelters.Post("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		postDogSheltersHandler := r.container.Resolve("DogShelterPostHandler", config.Transient).(Handler)
		return postDogSheltersHandler.Handle(c)
	})
	dogshelters.Get("/", func(c *fiber.Ctx) error {
		queryParamsMiddleware := r.container.Resolve("QueryParamsMiddleware", config.Transient).(QueryParamsMiddleware)
		return queryParamsMiddleware.ValidateQueryParams(c)
	}, func(c *fiber.Ctx) error {
		getDogsheltersHandler := r.container.Resolve("DogShelterGetHandler", config.Transient).(Handler)
		return getDogsheltersHandler.Handle(c)
	})

	users := v1.Group("/users")
	users.Post("/", func(c *fiber.Ctx) error {
		postUserHandler := r.container.Resolve("UserPostHandler", config.Transient).(Handler)
		return postUserHandler.Handle(c)
	})
	users.Delete("/:id", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		deleteUserHandler := r.container.Resolve("UserDeleteHandler", config.Transient).(Handler)
		return deleteUserHandler.Handle(c)
	})
	users.Get("/me", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		userGetMeHandler := r.container.Resolve("UserGetMeHandler", config.Transient).(Handler)
		return userGetMeHandler.Handle(c)
	})

	userwebhook := users.Group("/:id/webhook")
	userwebhook.Delete("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		deleteUserWebhookHandler := r.container.Resolve("UserWebhookDeleteHandler", config.Transient).(Handler)
		return deleteUserWebhookHandler.Handle(c)
	})
	userwebhook.Get("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		getUserWebhookHandler := r.container.Resolve("UserWebhookGetHandler", config.Transient).(Handler)
		return getUserWebhookHandler.Handle(c)
	})
	userwebhook.Post("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		postUserWebhookHandler := r.container.Resolve("UserWebhookPostHandler", config.Transient).(Handler)
		return postUserWebhookHandler.Handle(c)
	})
	userwebhook.Put("/", func(c *fiber.Ctx) error {
		authMiddleware := r.container.Resolve("AuthMiddleware", config.Transient).(AuthMiddleware)
		return authMiddleware.AuthenticateRequest(c)
	}, func(c *fiber.Ctx) error {
		putUserWebhookHandler := r.container.Resolve("UserWebhookPutHandler", config.Transient).(Handler)
		return putUserWebhookHandler.Handle(c)
	})

	app.Listen(os.Getenv("APPLICATION_PORT"))

}
