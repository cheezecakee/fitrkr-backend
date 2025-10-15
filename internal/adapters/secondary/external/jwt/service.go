// Package jwt
package jwt

import (
	"time"

	"github.com/cheezecakee/logr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWT interface {
	MakeJWT(userID uuid.UUID, roles []string) (string, error)
	ValidateJWT(tokenString string) (uuid.UUID, error)
}

type UserClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	SecretKey []byte
	ExpiresIn time.Duration
	Claims    UserClaims
}

func NewJWTManager(secretKey []byte, expiresIn time.Duration) JWT {
	return &JWTManager{SecretKey: secretKey, ExpiresIn: expiresIn}
}

func (j *JWTManager) MakeJWT(userID uuid.UUID, roles []string) (string, error) {
	claims := &UserClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "fitrkr",
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		logr.Get().Errorf("err signing token: %v", err)
		return "", err
	}
	logr.Get().Info("Token generated")
	return ss, nil
}

func (j *JWTManager) ValidateJWT(tokenString string) (uuid.UUID, error) {
	var userClaims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, nil
	}

	userID, err := userClaims.GetSubject()
	if err != nil {
		logr.Get().Errorf("error getting user id: %v", err)
		return uuid.Nil, err
	}

	logr.Get().Info("Token validated")
	return uuid.MustParse(userID), err
}
