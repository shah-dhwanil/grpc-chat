package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `db:"id"`
	Name string `db:"name"`
	PrimaryEmail string `db:"primary_email"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type CreateUserRequest struct {
	Name string `db:"name"`
	PrimaryEmail string `db:"primary_email"`
}

type UpdateUserRequest struct {
	Name *string `db:"name"`
	PrimaryEmail *string `db:"primary_email"`
}
