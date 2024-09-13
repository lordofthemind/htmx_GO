package types

import (
	"github.com/google/uuid"
)

// Superuser represents a user with administrative privileges.
type Superuser struct {
	ID           uuid.UUID `bson:"_id,omitempty" json:"id"`
	Username     string    `bson:"username" json:"username" validate:"required,min=3,max=32"`
	Email        string    `bson:"email" json:"email" validate:"required,email"`
	Password     string    `bson:"password" json:"password" validate:"required,min=6"`
	Role         string    `bson:"role" json:"role" validate:"required"`
	CreatedAt    int64     `bson:"created_at" json:"created_at"`
	UpdatedAt    int64     `bson:"updated_at" json:"updated_at"`
	Is2FAEnabled bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled"`
	ResetToken   string    `bson:"reset_token" json:"reset_token"`
}
