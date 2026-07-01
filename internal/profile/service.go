package profile

import (
	"context"
	"errors"
	"log"

	"github.com/Parachurami/ecommerce-app-api/internal/store"
	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/google/uuid"
)

type Service interface {
	UpdateProfile(uuid.UUID, *types.UpdateProfileParams, context.Context) (*types.Profile, error)
	GetProfile(uuid.UUID, context.Context) (*types.Profile, error)
}

type svc struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return &svc{
		store: store,
	}
}

func (s *svc) UpdateProfile(id uuid.UUID, params *types.UpdateProfileParams, ctx context.Context) (*types.Profile, error) {
	profile, err := s.store.UpdateProfile(id, params, ctx)
	if err != nil {
		log.Print("Error updating profile: ", err.Error())
		return nil, err
	}
	return profile, nil
}

func (s *svc) GetProfile(id uuid.UUID, ctx context.Context) (*types.Profile, error) {
	profile, err := s.store.GetProfileById(id, ctx)
	if err != nil {
		log.Print("Error fetching profile: ", err.Error())
		return nil, errors.New("Could not fetch profile")
	}
	return profile, nil
}
