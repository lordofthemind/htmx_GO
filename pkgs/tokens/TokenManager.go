package tokens

import "github.com/lordofthemind/htmx_GO/internals/configs"

type TokenManager interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

func NewTokenManager() TokenManager {
	if configs.UseJWT {
		return NewJWTManager()
	}
	return NewPasetoManager()
}
