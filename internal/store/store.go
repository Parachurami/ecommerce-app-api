package store

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Parachurami/ecommerce-app-api/types"
	"github.com/Parachurami/ecommerce-app-api/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func NewStore(db *pgxpool.Pool) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUserById(id uuid.UUID, ctx context.Context) (*types.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	rows, queryError := s.db.Query(ctx, query, id)
	if queryError != nil {
		return nil, queryError
	}

	user := new(types.User)
	defer rows.Close()
	for rows.Next() {
		id := &user.Id
		email := &user.Email
		password := &user.Password
		createdAt := &user.CreatedAt
		role := &user.UserRole
		if scanError := rows.Scan(id, email, password, role, createdAt); scanError != nil {
			return nil, scanError
		}
	}
	if !user.Id.Valid {
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (s *Store) GetUserByEmail(email string, ctx context.Context) (*types.User, error) {
	query := "SELECT * FROM users WHERE email = $1"
	rows, queryError := s.db.Query(ctx, query, email)
	if queryError != nil {
		return nil, queryError
	}
	user := new(types.User)
	defer rows.Close()
	for rows.Next() {
		id := &user.Id
		email := &user.Email
		password := &user.Password
		role := &user.UserRole
		createdAt := &user.CreatedAt
		if scanError := rows.Scan(id, email, password, role, createdAt); scanError != nil {
			log.Print("Error scanning", scanError)
			return nil, scanError
		}

	}

	if !user.Id.Valid {
		log.Printf("User ID: %v", user)
		return nil, errors.New("User not found")
	}
	return user, nil
}

func (s *Store) CreateUser(params *types.RegisterUserParams, ctx context.Context) (*types.User, *types.Profile, error) {
	queryString := "INSERT INTO users(email, password, role) VALUES($1, $2, $3) RETURNING id"
	var id uuid.NullUUID
	s.db.QueryRow(ctx,
		queryString,
		params.Email,
		params.Password,
		params.RoleType).Scan(&id)

	profile, createProfileErr := s.CreateProfile(id.UUID, params.Email, params.FirstName, params.LastName, "", ctx)

	if createProfileErr != nil {

		log.Print("Error creating profile: ", createProfileErr.Error())
		return nil, nil, createProfileErr
	}

	return &types.User{
		Id:        id,
		Email:     params.Email,
		UserRole:  params.RoleType,
		Password:  params.Password,
		CreatedAt: time.Now(),
	}, profile, nil
}

func (s *Store) GetProfileById(id uuid.UUID, ctx context.Context) (*types.Profile, error) {
	queryString := "SELECT * FROM profile WHERE userId = $1"
	row := s.db.QueryRow(ctx, queryString, id)
	if row == nil {
		log.Printf("Profile with id: %v, does not exist", id)
		return nil, errors.New("Profile does not exist")
	}
	profile := new(types.Profile)
	if err := row.Scan(&profile.Id,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.ProfileImage,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt); err != nil {
		log.Print("Error scanning row: ", err.Error())
		return nil, err
	}
	return profile, nil
}

func (s *Store) CreateProfile(id uuid.UUID, email, firstName, lastName, image string, ctx context.Context) (*types.Profile, error) {
	query, args, err := utils.Psql.
		Insert("profile").
		Columns("userId", "email", "firstName", "lastName", "profileImage").
		Values(id, email, firstName, lastName, image).Suffix("RETURNING *").ToSql()
	if err != nil {
		log.Print("Error converting query: ", err)
		return nil, err
	}
	row := s.db.QueryRow(ctx, query, args...)
	if row == nil {
		log.Print("Error creating profile")
		return nil, errors.New("Error Creating Profile")
	}
	profile := new(types.Profile)
	if err := row.Scan(&profile.Id, &profile.FirstName, &profile.LastName, &profile.Email, &profile.ProfileImage, &profile.Bio, &profile.CreatedAt, &profile.UpdatedAt); err != nil {
		log.Print("Error returning new profile: ", err.Error())
		return nil, err
	}
	return profile, nil
}

func (s *Store) UpdateProfile(id uuid.UUID, params *types.UpdateProfileParams, ctx context.Context) (*types.Profile, error) {
	builder := utils.Psql.Update("profile").Where(squirrel.Eq{"userId": id})
	if params.FirstName != "" {
		builder = builder.Set("firstName", params.FirstName)
	}
	if params.LastName != "" {
		builder = builder.Set("lastName", params.LastName)
	}
	if params.Bio != "" {
		builder = builder.Set("bio", params.Bio)
	}
	if params.ProfileImage != "" {
		builder = builder.Set("profileImage", params.ProfileImage)
	}
	builder = builder.Set("updatedAt", time.Now()).Suffix("RETURNING *")
	sql, args, sqlErr := builder.ToSql()
	if sqlErr != nil {
		log.Print("Error converting to sql:", sqlErr.Error())
		return nil, sqlErr
	}
	profile := new(types.Profile)
	if err := s.db.QueryRow(ctx, sql, args...).Scan(&profile.Id,
		&profile.FirstName,
		&profile.LastName,
		&profile.Email,
		&profile.ProfileImage,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt); err != nil {
		return nil, err
	}
	return profile, nil
}

func (store *Store) CreateProduct(userId uuid.UUID, params *types.CreateProductParams, ctx context.Context) (*types.Product, error) {
	/*
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		    userId UUID NOT NULL,
		    name TEXT NOT NULL,
		    description TEXT NOT NULL,
		    budget FLOAT NOT NULL,
		    skills TEXT[],
		    duration INTERVAL NOT NULL,
		    expiration INTERVAL NOT NULL,
		    image_url TEXT,
		    deliverables TEXT[],
		    createdAt TIMESTAMPTZ DEFAULT NOW(),
		    updatedAt TIMESTAMPTZ DEFAULT NOW(),
	*/
	builder := utils.Psql.
		Insert("products").
		Columns("userId", "name", "description", "budget", "skills", "duration", "expiration", "image_url", "deliverables", "createdAt", "updatedAt").
		Values(userId, params.Name, params.Desciption, params.Budget, params.Skills, params.Duration, params.Expiration, params.ImageUrl, params.Deliverables, time.Now(), time.Now()).
		Suffix("RETURNING *")
	sql, args, sqlErr := builder.ToSql()
	if sqlErr != nil {
		log.Print("Could not execute query: ", sqlErr.Error())
		return nil, errors.New("Could not execute query")
	}
	product := new(types.Product)
	if err := store.db.QueryRow(ctx, sql, args...).Scan(&product.Id, &product.UserId, &product.Name, &product.Description, &product.Budget, &product.Skills, &product.Duration, &product.Expiration, &product.ImageUrl, &product.Deliverables, &product.CreatedAt, &product.UpdatedAt); err != nil {
		log.Print("Could not create product: ", err.Error())
		return nil, errors.New("Could not create product")
	}
	return product, nil
}

func (store *Store) UpdateProduct(userId uuid.UUID, params *types.UpdateProductParams, ctx context.Context) (*types.Product, error) {
	builder := utils.Psql.Update("products").Where(squirrel.Eq{"userId": userId})
	if params.Budget != nil {
		builder = builder.Set("budget", *params.Budget)
	}
	if params.Deliverables != nil {
		builder = builder.Set("deliverables", *params.Deliverables)
	}
	if params.Desciption != nil {
		builder = builder.Set("description", *params.Desciption)
	}
	if params.ImageUrl != nil {
		builder = builder.Set("image_url", *params.ImageUrl)
	}
	if params.Name != nil {
		builder = builder.Set("name", *params.Name)
	}
	if params.Duration != nil {
		builder = builder.Set("duration", *params.Duration)
	}
	if params.Expiration != nil {
		builder = builder.Set("expiration", *params.Expiration)
	}
	if params.Skills != nil {
		builder = builder.Set("skills", *params.Skills)
	}
	builder = builder.Set("updatedAt", time.Now()).Suffix("RETURNING *")
	sql, args, sqlErr := builder.ToSql()
	if sqlErr != nil {
		log.Print("Error converting to sql: ", sqlErr)
		return nil, utils.InternalServerError
	}
	rows, rowErr := store.db.Query(ctx, sql, args...)
	if rowErr != nil {
		log.Print("There was a row error: ", rowErr.Error())
		return nil, utils.InternalServerError
	}
	product := new(types.Product)
	for rows.Next() {
		scanErr := rows.Scan(&product.Id, &product.UserId, &product.Name, &product.Description,
			&product.Budget, &product.Budget, &product.Skills, &product.Duration, &product.Expiration, &product.ImageUrl,
			&product.Deliverables, &product.CreatedAt, &product.UpdatedAt)
		if scanErr != nil {
			log.Print("Error scanning product after update: ", scanErr)
			return nil, utils.InternalServerError
		}
	}
	return product, nil
}

func (store *Store) GetProducts(userId uuid.UUID, ctx context.Context) ([]types.Product, error) {
	builder := utils.Psql.Select("*").From("products").Where(squirrel.Eq{"userId": userId})
	sql, args, sqlErr := builder.ToSql()
	if sqlErr != nil {
		log.Print("Error converting to sql: ", sqlErr)
		return nil, utils.InternalServerError
	}
	rows, queryErr := store.db.Query(ctx, sql, args...)
	if queryErr != nil {
		log.Print("Error querying products: ", queryErr)
		return nil, utils.InternalServerError
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		/*
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			    userId UUID NOT NULL,
			    name TEXT NOT NULL,
			    description TEXT NOT NULL,
			    budget FLOAT NOT NULL,
			    skills TEXT[],
			    duration INTERVAL NOT NULL,
			    expiration INTERVAL NOT NULL,
			    image_url TEXT,
			    deliverables TEXT[],
			    createdAt TIMESTAMPTZ DEFAULT NOW(),
			    updatedAt TIMESTAMPTZ DEFAULT NOW(),
		*/
		product := new(types.Product)
		if err := utils.ScanRow(rows, product); err != nil {
			log.Print("Error scanning row: ", err)
			return nil, utils.InternalServerError
		}
		products = append(products, *product)
	}
	return products, nil
}

func (store *Store) DeleteProductById(userId, productId uuid.UUID, ctx context.Context) error {
	builder := utils.Psql.Delete("products").Where(squirrel.Eq{"userId": userId, "id": productId})
	sql, args, sqlErr := builder.ToSql()
	if sqlErr != nil {
		log.Print("Error converting builder to sql: ", sqlErr)
		return utils.InternalServerError
	}
	if _, executionErr := store.db.Exec(ctx, sql, args...); executionErr != nil {
		log.Print("Error executing query: ", executionErr)
		return utils.InternalServerError
	}
	return nil
}

func (store *Store) DeleteProductsByIds(ctx context.Context, userId uuid.UUID, productIds uuid.UUIDs) error {
	var err error = nil
	builder := utils.Psql.Delete("products").Where(squirrel.Eq{"userId": userId, "id": productIds})
	sql, args, err := builder.ToSql()
	_, err = store.db.Exec(ctx, sql, args...)
	return err
}
