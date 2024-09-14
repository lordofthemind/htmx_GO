package tokens

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lordofthemind/htmx_GO/internals/configs"
)

// JWTManager implements the TokenManager interface for JWT tokens.
type JWTManager struct{}

// NewJWTManager returns a new instance of JWTManager.
func NewJWTManager() *JWTManager {
	return &JWTManager{}
}

// GenerateToken generates a JWT token for the given user ID.
func (j *JWTManager) GenerateToken(userID string) (string, error) {
	// Define the token claims
	claims := jwt.MapClaims{
		"user_id":   userID,
		"expire_at": time.Now().Add(configs.TokenAccessDuration).Unix(),
	}

	// Create the token using the signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the symmetric key from the configuration
	tokenString, err := token.SignedString([]byte(configs.TokenSymmetricKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid.
func (j *JWTManager) ValidateToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(configs.TokenSymmetricKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Check if the token is valid and return the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("invalid token")
	}
}
