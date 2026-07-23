package auth

import (
	"errors"
	"time"

	"github.com/HenryNg101/auth-service/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"userID"`
	Role   string `json:"role"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

var jwtSecret = []byte(config.LoadJWTSecret()) // Load secret from .env

func GenerateAccessToken(userID uint, role string) (string, error) {
	return generateToken(userID, role, "access", 15*time.Minute)
}

func GenerateRefreshToken(userID uint, role string) (string, error) {
	return generateToken(userID, role, "refresh", 7*24*time.Hour)
}

func generateToken(userID uint, role, tokenType string, duration time.Duration) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("JWT secret is required. You have to set it in .env file in root folder using JWT_SECRET variable")
	}
	claims := Claims{
		UserID: userID,
		Role:   role,
		Type:   tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret is required")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
