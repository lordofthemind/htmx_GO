package repositories

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lordofthemind/htmx_GO/internals/types"
)

type inMemorySuperuserRepo struct {
	data map[uuid.UUID]*types.SuperUserType
	mu   sync.RWMutex
}

// NewInMemorySuperuserRepository initializes an in-memory Superuser repository.
func NewInMemorySuperuserRepository() SuperuserRepository {
	return &inMemorySuperuserRepo{
		data: make(map[uuid.UUID]*types.SuperUserType),
	}
}

// CreateSuperuser creates a new superuser in memory.
func (r *inMemorySuperuserRepo) CreateSuperuser(ctx context.Context, superuser *types.SuperUserType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	superuser.CreatedAt = time.Now().Unix()
	superuser.UpdatedAt = time.Now().Unix()
	if superuser.ID == uuid.Nil {
		superuser.ID = uuid.New()
	}
	r.data[superuser.ID] = superuser
	return nil
}

// FindSuperuserByEmail finds a superuser by email in memory.
func (r *inMemorySuperuserRepo) FindSuperuserByEmail(ctx context.Context, email string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.data {
		if su.Email == email {
			return su, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// FindSuperuserByID finds a superuser by ID in memory.
func (r *inMemorySuperuserRepo) FindSuperuserByID(ctx context.Context, id uuid.UUID) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if su, ok := r.data[id]; ok {
		return su, nil
	}
	return nil, errors.New("superuser not found")
}

// UpdateSuperuser updates a superuser's details in memory.
func (r *inMemorySuperuserRepo) UpdateSuperuser(ctx context.Context, superuser *types.SuperUserType) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[superuser.ID]; !ok {
		return errors.New("superuser not found")
	}
	superuser.UpdatedAt = time.Now().Unix()
	r.data[superuser.ID] = superuser
	return nil
}

// FindSuperuserByUsername finds a superuser by username in memory.
func (r *inMemorySuperuserRepo) FindSuperuserByUsername(ctx context.Context, username string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.data {
		if su.Username == username {
			return su, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// FindSuperuserByResetToken finds a superuser by reset token in memory.
func (r *inMemorySuperuserRepo) FindSuperuserByResetToken(ctx context.Context, token string) (*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, su := range r.data {
		if su.ResetToken == token {
			return su, nil
		}
	}
	return nil, errors.New("superuser not found")
}

// GetRoleByID returns the role of a superuser by ID in memory.
func (r *inMemorySuperuserRepo) GetRoleByID(ctx context.Context, id uuid.UUID) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if su, ok := r.data[id]; ok {
		return su.Role, nil
	}
	return "", errors.New("superuser not found")
}

// DeleteSuperuserByID deletes a superuser by their ID in memory.
func (r *inMemorySuperuserRepo) DeleteSuperuserByID(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[id]; ok {
		delete(r.data, id)
		return nil
	}
	return errors.New("superuser not found")
}

// ListSuperusers lists superusers with pagination in memory.
func (r *inMemorySuperuserRepo) ListSuperusers(ctx context.Context, limit, skip int64) ([]*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var superusers []*types.SuperUserType
	count := int64(0)
	for _, su := range r.data {
		if count >= skip {
			superusers = append(superusers, su)
		}
		if int64(len(superusers)) == limit {
			break
		}
		count++
	}
	return superusers, nil
}

// UpdateResetToken updates the reset token for a superuser in memory.
func (r *inMemorySuperuserRepo) UpdateResetToken(ctx context.Context, id uuid.UUID, token string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if su, ok := r.data[id]; ok {
		su.ResetToken = token
		su.UpdatedAt = time.Now().Unix()
		r.data[id] = su
		return nil
	}
	return errors.New("superuser not found")
}

// Enable2FA enables or disables 2FA for a superuser in memory.
func (r *inMemorySuperuserRepo) Enable2FA(ctx context.Context, id uuid.UUID, isEnabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if su, ok := r.data[id]; ok {
		su.Is2FAEnabled = isEnabled
		su.UpdatedAt = time.Now().Unix()
		r.data[id] = su
		return nil
	}
	return errors.New("superuser not found")
}

// SoftDeleteSuperuser marks a superuser as archived instead of permanently deleting in memory.
func (r *inMemorySuperuserRepo) SoftDeleteSuperuser(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if su, ok := r.data[id]; ok {
		su.UpdatedAt = time.Now().Unix()
		r.data[id] = su
		return nil
	}
	return errors.New("superuser not found")
}

// SearchSuperusers allows partial search by full_name, username, or email in memory.
func (r *inMemorySuperuserRepo) SearchSuperusers(ctx context.Context, searchQuery string) ([]*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*types.SuperUserType
	for _, su := range r.data {
		if matchString(su.FullName, searchQuery) || matchString(su.Username, searchQuery) || matchString(su.Email, searchQuery) {
			results = append(results, su)
		}
	}
	return results, nil
}

// Helper function to match strings with a basic case-insensitive contains check.
func matchString(field, query string) bool {
	return field != "" && len(query) > 0 && (len(query) <= len(field)) && containsIgnoreCase(field, query)
}

// containsIgnoreCase checks if a string contains another string, ignoring case.
func containsIgnoreCase(a, b string) bool {
	return len(b) > 0 && len(a) >= len(b) && (a == b || a[len(a)-len(b):] == b)
}

// FindAll2FAEnabledSuperusers finds all superusers with 2FA enabled in memory.
func (r *inMemorySuperuserRepo) FindAll2FAEnabledSuperusers(ctx context.Context) ([]*types.SuperUserType, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var superusers []*types.SuperUserType
	for _, su := range r.data {
		if su.Is2FAEnabled {
			superusers = append(superusers, su)
		}
	}
	return superusers, nil
}

// UpdateSuperuserRole updates the role of a superuser in memory.
func (r *inMemorySuperuserRepo) UpdateSuperuserRole(ctx context.Context, id uuid.UUID, role string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if su, ok := r.data[id]; ok {
		su.Role = role
		su.UpdatedAt = time.Now().Unix()
		r.data[id] = su
		return nil
	}
	return errors.New("superuser not found")
}

// BulkUpdateSuperusers updates multiple superusers at once in memory.
func (r *inMemorySuperuserRepo) BulkUpdateSuperusers(ctx context.Context, ids []uuid.UUID, updates map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		if su, ok := r.data[id]; ok {
			for key, value := range updates {
				switch key {
				case "username":
					su.Username = value.(string)
				case "password":
					su.Password = value.(string)
				case "role":
					su.Role = value.(string)
				case "is_2fa_enabled":
					su.Is2FAEnabled = value.(bool)
				}
			}
			su.UpdatedAt = time.Now().Unix()
			r.data[id] = su
		}
	}
	return nil
}
