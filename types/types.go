package types

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Role string

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    time.Duration
	RtExpires    time.Duration
}

var (
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
    CONSTRAINT fk_product FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
*/

type Product struct {
	Id           uuid.NullUUID  `json:"id"`
	UserId       uuid.NullUUID  `json:"userId"`
	Name         string         `json:"title"`
	Description  string         `json:"desciption"`
	Budget       float64        `json:"budget"`
	Skills       []string       `json:"skills"`
	Duration     time.Duration  `json:"duration"`
	Expiration   time.Duration  `json:"expiration"`
	ImageUrl     sql.NullString `json:"imageUrl"`
	Deliverables []string       `json:"deliverables"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type CreateProductParams struct {
	Name         string         `json:"name" validate:"required" msg:"name is required"`
	Desciption   string         `json:"description" validate:"required"  msg:"name is required"`
	Budget       float64        `json:"budget" validate:"required"  msg:"description is required"`
	Skills       []string       `json:"skills" validate:"required"  msg:"skills is required"`
	Duration     time.Duration  `json:"duration" validate:"required"  msg:"duration is required"`
	Expiration   time.Duration  `json:"expiration" validate:"required"  msg:"expiration is required"`
	ImageUrl     sql.NullString `json:"image" validate:"required"  msg:"image is required"`
	Deliverables []string       `json:"deliverables" validate:"required"  msg:"deliverables is required"`
}

type UpdateProductParams struct {
	Name         *string
	Desciption   *string
	Budget       *float64
	Skills       *[]string
	Duration     *time.Duration
	Expiration   *time.Duration
	ImageUrl     *string
	Deliverables *[]string
}
