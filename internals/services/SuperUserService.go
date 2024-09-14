package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/htmx_GO/internals/repositories"
	"github.com/lordofthemind/htmx_GO/internals/types"
	"golang.org/x/crypto/bcrypt"
)

type SuperuserService interface {
	RegisterSuperuser(ctx context.Context, username, email, password string) error
	AuthenticateSuperuser(ctx context.Context, email, password string) (*types.SuperUserType, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, username, password string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
	Verify2FA(ctx context.Context, userID uuid.UUID, code string) error
	GetFilePath(fileID string) (string, error)
	GetRole(ctx context.Context, userID uuid.UUID) (string, error)
	UpdateRole(ctx context.Context, userID uuid.UUID, role string) error
	Enable2FA(ctx context.Context, userID uuid.UUID, isEnabled bool) error
	BulkUpdateSuperusers(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}) error
	SearchSuperusers(ctx context.Context, searchQuery string) ([]*types.SuperUserType, error)
}

type superuserService struct {
	repo repositories.SuperuserRepository
}

func NewSuperuserService(repo repositories.SuperuserRepository) SuperuserService {
	return &superuserService{repo: repo}
}

// RegisterSuperuser creates a new superuser with hashed password.
func (s *superuserService) RegisterSuperuser(ctx context.Context, username, email, password string) error {
	// Check if the email already exists
	_, err := s.repo.FindSuperuserByEmail(ctx, email)
	if err == nil {
		return errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Create superuser
	superuser := &types.SuperUserType{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return s.repo.CreateSuperuser(ctx, superuser)
}

// AuthenticateSuperuser verifies a superuser's credentials.
func (s *superuserService) AuthenticateSuperuser(ctx context.Context, email, password string) (*types.SuperUserType, error) {
	superuser, err := s.repo.FindSuperuserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("superuser not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(superuser.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return superuser, nil
}

// UpdateProfile updates the username and password of a superuser.
func (s *superuserService) UpdateProfile(ctx context.Context, userID uuid.UUID, username, password string) error {
	superuser, err := s.repo.FindSuperuserByID(ctx, userID)
	if err != nil {
		return err
	}

	if username != "" {
		superuser.Username = username
	}

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("failed to hash password")
		}
		superuser.Password = string(hashedPassword)
	}

	superuser.UpdatedAt = time.Now().Unix()
	return s.repo.UpdateSuperuser(ctx, superuser)
}

// SendPasswordResetEmail sends a reset token (placeholder functionality).
func (s *superuserService) SendPasswordResetEmail(ctx context.Context, email string) error {
	superuser, err := s.repo.FindSuperuserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Generate and send a reset token (placeholder logic)
	resetToken := "generated-reset-token"
	fmt.Printf("Sending password reset token to %s: %s\n", superuser.Email, resetToken)

	return s.repo.UpdateResetToken(ctx, superuser.ID, resetToken)
}

// ResetPassword resets the password of a superuser using a token.
func (s *superuserService) ResetPassword(ctx context.Context, token, password string) error {
	superuser, err := s.repo.FindSuperuserByResetToken(ctx, token)
	if err != nil {
		return errors.New("invalid reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	superuser.Password = string(hashedPassword)
	superuser.UpdatedAt = time.Now().Unix()

	return s.repo.UpdateSuperuser(ctx, superuser)
}

// Verify2FA verifies the 2FA code (placeholder logic).
func (s *superuserService) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	superuser, err := s.repo.FindSuperuserByID(ctx, userID)
	if err != nil {
		return err
	}

	// Verify 2FA code (placeholder)
	if code != "expected-2fa-code" {
		return errors.New("invalid 2FA code")
	}

	superuser.Is2FAEnabled = true
	return s.repo.UpdateSuperuser(ctx, superuser)
}

// Enable2FA enables or disables 2FA for a superuser.
func (s *superuserService) Enable2FA(ctx context.Context, userID uuid.UUID, isEnabled bool) error {
	return s.repo.Enable2FA(ctx, userID, isEnabled)
}

// GetFilePath returns the file path for a given file ID (placeholder functionality).
func (s *superuserService) GetFilePath(fileID string) (string, error) {
	// Placeholder logic for file path
	return fmt.Sprintf("./uploads/%s", fileID), nil
}

// GetRole retrieves the role of a superuser by their ID.
func (s *superuserService) GetRole(ctx context.Context, userID uuid.UUID) (string, error) {
	return s.repo.GetRoleByID(ctx, userID)
}

// UpdateRole updates the role of a superuser by their ID.
func (s *superuserService) UpdateRole(ctx context.Context, userID uuid.UUID, role string) error {
	return s.repo.UpdateSuperuserRole(ctx, userID, role)
}

// BulkUpdateSuperusers updates multiple superusers at once.
func (s *superuserService) BulkUpdateSuperusers(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}) error {
	return s.repo.BulkUpdateSuperusers(ctx, ids, updates)
}

// SearchSuperusers allows searching for superusers based on partial matches of full name, username, or email.
func (s *superuserService) SearchSuperusers(ctx context.Context, searchQuery string) ([]*types.SuperUserType, error) {
	return s.repo.SearchSuperusers(ctx, searchQuery)
}
