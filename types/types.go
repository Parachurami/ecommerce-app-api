package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	Id        uuid.NullUUID `json:"id"`
	Email     string        `json:"email"`
	UserRole  Role          `json:"userRole"`
	Password  string        `json:"-"`
	CreatedAt time.Time     `json:"createdAt"`
}

type RegisterUserParams struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required,min=6,max=130"`
	RoleType  Role   `json:"role" validate:"required"`
}

type LoginUserParams struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

type Profile struct {
	Id           uuid.NullUUID
	FirstName    string         `json:"firstName"`
	LastName     string         `json:"lastName"`
	Bio          sql.NullString `json:"userBio"`
	ProfileImage sql.NullString `json:"imageUrl"`
	Email        string         `json:"email"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type UpdateProfileParams struct {
	FirstName    string
	LastName     string
	Bio          string
	ProfileImage string
	Email        string
}
