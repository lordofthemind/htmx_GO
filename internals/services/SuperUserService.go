package services

import (
	"context"
	"errors"

	"github.com/lordofthemind/htmx_GO/internals/repositories"
	"github.com/lordofthemind/htmx_GO/internals/types"
	"golang.org/x/crypto/bcrypt"
)

type SuperuserService interface {
	RegisterSuperuser(ctx context.Context, username, email, password string) error
	AuthenticateSuperuser(ctx context.Context, email, password string) (*types.Superuser, error)
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
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
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
