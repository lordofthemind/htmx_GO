package tokens

import (
	"time"

	"github.com/lordofthemind/htmx_GO/internals/configs"
	"github.com/o1egl/paseto"
)

// PasetoManager implements the TokenManager interface for PASETO tokens.
type PasetoManager struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoManager returns a new instance of PasetoManager.
func NewPasetoManager() *PasetoManager {
	return &PasetoManager{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(configs.TokenSymmetricKey),
	}
}

// GenerateToken generates a PASETO token for the given user ID.
func (p *PasetoManager) GenerateToken(userID string) (string, error) {
	token, err := p.paseto.Encrypt(p.symmetricKey, map[string]interface{}{
		"user_id":   userID,
		"expire_at": time.Now().Add(configs.TokenAccessDuration).Unix(),
	}, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken validates a PASETO token and returns the claims if valid.
func (p *PasetoManager) ValidateToken(tokenString string) (map[string]interface{}, error) {
	var payload map[string]interface{}
	err := p.paseto.Decrypt(tokenString, p.symmetricKey, &payload, nil)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
