package userAuth

import (
	"context"
	"errors"
	"log"

	"github.com/Parachurami/ecommerce-app-api/internal/store"
	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
)

type Service interface {
	Login(*types.LoginUserParams) (*types.User, error)
	Register(*types.RegisterUserParams) (*types.User, *types.Profile, error)
}

type svc struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return &svc{
		store: store,
	}
}

func (s *svc) Login(params *types.LoginUserParams) (*types.User, error) {
	log.Print("Params: ", params.Password)
	ctx := context.Background()
	user, err := s.store.GetUserByEmail(params.Email, ctx)
	if err != nil {
		log.Print("Error logging in")
		return nil, errors.New("Invalid Credentials")
	}
	if params.Password != params.ConfirmPassword {
		return nil, errors.New("Passwords must match")
	}
	log.Printf("Hash value: %v, Hash length: %v", user.Password, len(user.Password))
	if err := utils.CompareHash(params.Password, user.Password); err != nil {
		log.Print("Error when comparing hashes")
		return nil, errors.New("Invalid Credentials")
	}
	return user, nil
}

func (s *svc) Register(params *types.RegisterUserParams) (*types.User, *types.Profile, error) {
	ctx := context.Background()
	userExists, err := s.store.GetUserByEmail(params.Email, ctx)
	log.Printf("Existing user: %v", userExists)
	if err == nil {
		log.Printf("User exists")
		return nil, nil, errors.New("User already exists")
	}
	user, profile, creationErr := s.store.CreateUser(params, ctx)
	if creationErr != nil {
		log.Printf("Error creating user: %v", creationErr.Error())
		return nil, nil, errors.New("Error creating user")
	}
	return user, profile, nil
}
