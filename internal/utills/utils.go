package utills

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dorik33/DeNet/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(storedPassword, providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(providedPassword))
}
func GenerateToken(id int, email string, secretKey []byte, duration time.Duration) (string, error) {
	claims := models.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			ID:        strconv.Itoa(id),
		},
		Email: email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func ValidateToken(tokenStr string, secretKey []byte) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&models.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		},
	)

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*models.UserClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return 0, fmt.Errorf("token expired")
	}

	userID, err := strconv.Atoi(claims.ID)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID in token")
	}

	return userID, nil
}
