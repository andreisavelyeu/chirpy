package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pass), nil
}

func CheckPasswordHash(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return err
	}
	return nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	subject, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(subject)

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	splittedAuthorizaton := strings.Split(authorization, " ")

	if len(splittedAuthorizaton) != 2 {
		return "", errors.New("missing bearer")
	}

	bearer := splittedAuthorizaton[1]

	return bearer, nil
}

func MakeRefreshToken() (string, error) {
	// 32 bytes required
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	splittedAuthorizaton := strings.Split(authorization, " ")
	if len(splittedAuthorizaton) != 2 {
		return "", errors.New("missing api key")
	}
	if splittedAuthorizaton[0] != "ApiKey" {
		return "", errors.New("missing api key")
	}

	apiKey := splittedAuthorizaton[1]

	return apiKey, nil
}
