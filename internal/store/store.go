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

func (s *Store) GetUserById(id int64, ctx context.Context) (*types.User, error) {
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
		if scanError := rows.Scan(id, email, password, createdAt); scanError != nil {
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
