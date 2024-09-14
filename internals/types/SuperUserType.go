package types

import (
	"github.com/google/uuid"
)

// Superuser represents a user with administrative privileges.
type SuperUserType struct {
	ID               uuid.UUID `bson:"_id,omitempty" json:"id"`
	FullName         string    `bson:"full_name" json:"full_name" validate:"required,min=3,max=32"`
	Username         string    `bson:"username" json:"username" validate:"required,min=3,max=32"`
	Email            string    `bson:"email" json:"email" validate:"required,email"`
	Password         string    `bson:"password" json:"password" validate:"required,min=6"`
	Role             string    `bson:"role" json:"role" validate:"required"`
	CreatedAt        int64     `bson:"created_at" json:"created_at"`
	UpdatedAt        int64     `bson:"updated_at" json:"updated_at"`
	Is2FAEnabled     bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled"`
	ResetToken       string    `bson:"reset_token" json:"reset_token"`
	PermissionGroups []string  `bson:"permission_groups" json:"permission_groups"`
}

// // Superuser represents a user with administrative privileges.
// type Superuser struct {
// 	ID                        uuid.UUID `bson:"_id,omitempty" json:"id"`
// 	Username                  string    `bson:"username" json:"username" validate:"required,min=3,max=32"`
// 	Email                     string    `bson:"email" json:"email" validate:"required,email"`
// 	Password                  string    `bson:"password" json:"password" validate:"required,min=6"`
// 	Role                      string    `bson:"role" json:"role" validate:"required"`
// 	CreatedAt                 int64     `bson:"created_at" json:"created_at"`
// 	UpdatedAt                 int64     `bson:"updated_at" json:"updated_at"`
// 	Is2FAEnabled              bool      `bson:"is_2fa_enabled" json:"is_2fa_enabled"`
// 	ResetToken                string    `bson:"reset_token" json:"reset_token"`
// 	LastLoginAt               int64     `bson:"last_login_at" json:"last_login_at"`
// 	AccountLocked             bool      `bson:"account_locked" json:"account_locked"`
// 	FailedAttempts            int       `bson:"failed_attempts" json:"failed_attempts"`
// 	LastLoginIP               string    `bson:"last_login_ip" json:"last_login_ip"`
// 	PermissionGroups          []string  `bson:"permission_groups" json:"permission_groups"`
// 	PasswordLastChangedAt     int64     `bson:"password_last_changed_at" json:"password_last_changed_at"`
// 	PasswordChangeRequired    bool      `bson:"password_change_required" json:"password_change_required"`
// 	EmailNotificationsEnabled bool      `bson:"email_notifications_enabled" json:"email_notifications_enabled"`
// 	Bio                       string    `bson:"bio" json:"bio"`
// 	TwoFactorAuthMethod       string    `bson:"two_factor_auth_method" json:"two_factor_auth_method"`
// 	RecoveryEmail             string    `bson:"recovery_email" json:"recovery_email"`
// 	Timezone                  string    `bson:"timezone" json:"timezone"`
// 	LastActionAt              int64     `bson:"last_action_at" json:"last_action_at"`
// }
