// Package jwt
package jwt

import (
	"errors"
	"time"

	"github.com/cheezecakee/logr"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/cheezecakee/fitrkr-athena/internal/core/domain/user"
)

var ErrInvalidToken = errors.New("invalid token")

type JWT interface {
	MakeJWT(userID uuid.UUID, roles []string) (string, error)
	ValidateJWT(tokenString string) (*AuthenticatedUser, error)
}

type AuthenticatedUser struct {
	UserID uuid.UUID
	Roles  user.Roles
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

func NewJWTManager(secretKey string, expiresIn time.Duration) JWT {
	if expiresIn == 0 {
		expiresIn = 15 * time.Minute
	}
	return &JWTManager{SecretKey: []byte(secretKey), ExpiresIn: expiresIn}
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

func (j *JWTManager) ValidateJWT(tokenString string) (*AuthenticatedUser, error) {
	var userClaims UserClaims

	token, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(userClaims.Subject)
	if err != nil {
		logr.Get().Errorf("error getting user id: %v", err)
		return nil, err
	}

	logr.Get().Info("Token validated")
	return &AuthenticatedUser{UserID: userID, Roles: user.StringsToRoles(userClaims.Roles)}, err
}
