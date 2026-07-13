package api

import (
	"log"
	"net/http"

	userAuth "github.com/Parachurami/ecommerce-app-api/internal/auth"
	"github.com/Parachurami/ecommerce-app-api/internal/product"
	"github.com/Parachurami/ecommerce-app-api/internal/profile"
	"github.com/Parachurami/ecommerce-app-api/internal/store"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type application struct {
	db          *pgxpool.Pool
	config      config
	redisClient *redis.Client
}

func NewApp(db *pgxpool.Pool, redisClient *redis.Client, config config) *application {
	return &application{
		db:          db,
		config:      config,
		redisClient: redisClient,
	}
}

func (app *application) Mount() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.ClientIPFromRemoteAddr)
	dbStore := store.NewStore(app.db)
	router.Route("/api/v1", func(mainRouter chi.Router) {
		mainRouter.Get("/health", func(res http.ResponseWriter, req *http.Request) {
			res.Write([]byte("Everything is good"))
		})
		mainRouter.Route("/auth", func(r chi.Router) {
			authService := userAuth.NewService(dbStore)
			authHandler := userAuth.NewHandler(authService, app.redisClient)
			r.Post("/login", authHandler.LoginUser)
			r.Post("/register", authHandler.RegisterUser)
			r.Route("/logout", func(LogoutRouter chi.Router) {
				LogoutRouter.Use(utils.WithJWT)
				LogoutRouter.Post("/", authHandler.Logout)
			})
		})
		mainRouter.Route("/account", func(accountRouter chi.Router) {
			accountRouter.Use(utils.WithJWT)
			accountRouter.Route("/profile", func(profileRouter chi.Router) {
				profileService := profile.NewService(dbStore)
				profileHandler := profile.NewHandler(profileService)
				profileRouter.Patch("/", profileHandler.UpdateProfile)
				profileRouter.Get("/", profileHandler.GetProfile)
			})
		})
		mainRouter.Route("/admin", func(adminRouter chi.Router) {
			adminRouter.Use(utils.WithJWT)
			adminRouter.Route("/products", func(productRouter chi.Router) {
				productService := product.NewService(dbStore)
				productHandler := product.NewHandler(productService)
				productRouter.Post("/", productHandler.CreateProduct)
				productRouter.Get("/", productHandler.GetProducts)
			})
		})
	})
	return router
}

func (app *application) Run(handler http.Handler) error {
	srv := &http.Server{
		Addr:    app.config.address,
		Handler: handler,
	}
	log.Printf("Server running on %v", app.config.address)
	return srv.ListenAndServe()
}

type config struct {
	dsn     string
	address string
}

func NewConfig(dsn, address string) config {
	return config{
		dsn:     dsn,
		address: address,
	}
}
