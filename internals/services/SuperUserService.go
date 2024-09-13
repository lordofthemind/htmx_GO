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

// SuperuserService defines the methods for the superuser-related functionalities.
type SuperuserService interface {
	RegisterSuperuser(ctx context.Context, username, email, password string) error
	AuthenticateSuperuser(ctx context.Context, email, password string) (*types.Superuser, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, username, password string) error
	SendPasswordResetEmail(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, password string) error
	Verify2FA(ctx context.Context, userID uuid.UUID, code string) error
	ListRoles(ctx context.Context) ([]string, error)
	GetUserActivityLogs(ctx context.Context, userID uuid.UUID) ([]types.UserActivityLog, error)
	GetFilePath(fileID string) (string, error)
}

type superuserService struct {
	repo repositories.SuperuserRepository
}

func NewSuperuserService(repo repositories.SuperuserRepository) SuperuserService {
	return &superuserService{repo: repo}
}

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
	superuser := &types.Superuser{
		ID:        uuid.New(), // Generate a new UUID
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	return s.repo.CreateSuperuser(ctx, superuser)
}

func (s *superuserService) AuthenticateSuperuser(ctx context.Context, email, password string) (*types.Superuser, error) {
	// Find superuser by email
	superuser, err := s.repo.FindSuperuserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("superuser not found")
	}

	// Check the password
	err = bcrypt.CompareHashAndPassword([]byte(superuser.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return superuser, nil
}

// Update user profile.
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

// Send password reset email (just a placeholder for now).
func (s *superuserService) SendPasswordResetEmail(ctx context.Context, email string) error {
	superuser, err := s.repo.FindSuperuserByEmail(ctx, email)
	if err != nil {
		return err
	}

	// Here you would typically generate a reset token and send an email
	resetToken := "reset-token-placeholder"
	fmt.Printf("Sending password reset token to %s: %s\n", superuser.Email, resetToken)
	return nil
}

// Reset user password.
func (s *superuserService) ResetPassword(ctx context.Context, token, password string) error {
	// Token validation logic should go here (this is a placeholder)
	if token != "valid-token" {
		return errors.New("invalid reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	superuser, err := s.repo.FindSuperuserByResetToken(ctx, token)
	if err != nil {
		return err
	}

	superuser.Password = string(hashedPassword)
	superuser.UpdatedAt = time.Now().Unix()

	return s.repo.UpdateSuperuser(ctx, superuser)
}

// Verify 2FA code.
func (s *superuserService) Verify2FA(ctx context.Context, userID uuid.UUID, code string) error {
	superuser, err := s.repo.FindSuperuserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 2FA verification logic (placeholder)
	if code != "expected-2fa-code" {
		return errors.New("invalid 2FA code")
	}

	// Mark the 2FA as verified
	superuser.Is2FAEnabled = true
	return s.repo.UpdateSuperuser(ctx, superuser)
}

// List user roles (for role-based access control).
func (s *superuserService) ListRoles(ctx context.Context) ([]string, error) {
	roles, err := s.repo.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// Get user activity logs.
func (s *superuserService) GetUserActivityLogs(ctx context.Context, userID uuid.UUID) ([]types.UserActivityLog, error) {
	logs, err := s.repo.FindActivityLogsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// Get file path by file ID.
func (s *superuserService) GetFilePath(fileID string) (string, error) {
	// This is a placeholder. Normally, you'd look up the file in the database.
	filePath := fmt.Sprintf("./uploads/%s", fileID)
	return filePath, nil
}
