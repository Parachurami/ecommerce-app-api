package api

import (
	"log"
	"net/http"

	userAuth "github.com/Parachurami/ecommerce-app-api/internal/auth"
	"github.com/Parachurami/ecommerce-app-api/internal/profile"
	"github.com/Parachurami/ecommerce-app-api/internal/store"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type application struct {
	db     *pgxpool.Pool
	config config
}

func NewApp(db *pgxpool.Pool, config config) *application {
	return &application{
		db:     db,
		config: config,
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
			authHandler := userAuth.NewHandler(authService)
			r.Post("/login", authHandler.LoginUser)
			r.Post("/register", authHandler.RegisterUser)
		})
		mainRouter.Route("/account", func(accountRouter chi.Router) {
			accountRouter.Route("/profile", func(profileRouter chi.Router) {
				profileService := profile.NewService(dbStore)
				profileHandler := profile.NewHandler(profileService)
				profileRouter.Patch("/", utils.WithJWT(profileHandler.UpdateProfile))
				profileRouter.Get("/", utils.WithJWT(profileHandler.GetProfile))
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
