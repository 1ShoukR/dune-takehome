package utils

import (
	"os"
	"time"

	"dune-takehome-server/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecret []byte

// getJWTSecret returns the JWT secret, loading it from env if not already loaded
func getJWTSecret() []byte {
	if len(jwtSecret) == 0 {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default-secret-key-for-development" // Fallback for now, could use panic if I wanted to
		}
		jwtSecret = []byte(secret)
	}
	return jwtSecret
}

// Claims represents JWT claims
type Claims struct {
	UserID               primitive.ObjectID `json:"user_id"`
	Email                string             `json:"email"`
	jwt.RegisteredClaims `json:"registered_claims"`
}

// GenerateJWT generates a JWT token for a user
func GenerateJWT(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dune-form-builder",
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ValidateJWT validates and parses a JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
