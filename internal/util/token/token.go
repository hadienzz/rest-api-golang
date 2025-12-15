package token

import (
	"errors"
	"go-fiber-api/internal/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(c *fiber.Ctx, userID uuid.UUID) (string, error) {
	cfg := config.Get()

	now := time.Now().UTC()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
			Issuer:    "go-fiber-api",
			Subject:   "auth-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtKey))
	if err != nil {
		return "", err
	}

	SetAuthToken(c, tokenString)

	return tokenString, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	cfg := config.Get()

	claims := &CustomClaims{}

	parsedToken, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(cfg.JwtKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	customClaims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("could not parse claims")
	}

	return customClaims, nil
}
