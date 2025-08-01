package models

import (
	"database/sql/driver"
	"time"
)

type UserType string

const (
	UserTypeGuest UserType = "guest"
	UserTypeCast  UserType = "cast"
	UserTypeAdmin UserType = "admin"
)

func (ut UserType) Value() (driver.Value, error) {
	return string(ut), nil
}

func (ut *UserType) Scan(value interface{}) error {
	*ut = UserType(value.(string))
	return nil
}

type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	UserType     UserType  `json:"user_type"`
	Name         string    `json:"name"`
	Phone        *string   `json:"phone,omitempty"`
	BirthDate    *time.Time `json:"birth_date,omitempty"`
	ProfileImage *string   `json:"profile_image,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRegistration struct {
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required,min=8"`
	UserType  UserType  `json:"user_type" binding:"required,oneof=guest cast"`
	Name      string    `json:"name" binding:"required"`
	Phone     *string   `json:"phone"`
	BirthDate *time.Time `json:"birth_date" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserProfile struct {
	ID           int              `json:"id"`
	Email        string           `json:"email"`
	UserType     UserType         `json:"user_type"`
	Name         string           `json:"name"`
	Phone        *string          `json:"phone,omitempty"`
	BirthDate    *time.Time       `json:"birth_date,omitempty"`
	ProfileImage *string          `json:"profile_image,omitempty"`
	CastProfile  *CastProfile     `json:"cast_profile,omitempty"`
	Rating       *float64         `json:"rating,omitempty"`
	ReviewCount  int              `json:"review_count"`
}