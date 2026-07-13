package product

import (
	"context"
	"log"

	"github.com/Parachurami/ecommerce-app-api/internal/store"
	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/google/uuid"
)

type Service interface {
	CreateProduct(context.Context, uuid.UUID, *types.CreateProductParams) (*types.Product, error)
	GetProducts(context.Context, uuid.UUID) ([]types.Product, error)
}

type svc struct {
	store *store.Store
}

func NewService(store *store.Store) Service {
	return &svc{
		store: store,
	}
}

func (s *svc) CreateProduct(ctx context.Context, userId uuid.UUID, params *types.CreateProductParams) (*types.Product, error) {
	product, err := s.store.CreateProduct(userId, params, ctx)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *svc) GetProducts(ctx context.Context, userId uuid.UUID) ([]types.Product, error) {
	user, getUserErr := s.store.GetUserById(userId, ctx)
	if getUserErr != nil {
		log.Print("Error getting user: ", getUserErr)
		return nil, utils.UserNotFound
	}
	if user == nil {
		return nil, utils.UserNotFound
	}
	products, getProductsErr := s.store.GetProducts(userId, ctx)
	if getProductsErr != nil {
		log.Print("Error fetching products from products service: ", getProductsErr)
		return nil, getProductsErr
	}
	return products, nil
}
