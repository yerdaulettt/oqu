package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	incorrectToken = errors.New("Incorrect token")
	expiredErr     = errors.New("Token expired")
)

type claims struct {
	jwt.RegisteredClaims
	UserId int
	Role   string
}

type pair struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type JwtAuth struct {
	secret     []byte
	accessTtl  time.Duration
	refreshTtl time.Duration
}

func NewJwtAuth(s []byte, aTtl, rTtl time.Duration) *JwtAuth {
	return &JwtAuth{
		secret:     s,
		accessTtl:  aTtl,
		refreshTtl: rTtl,
	}
}

func (j *JwtAuth) newAccessToken(userId int, role string) (string, error) {
	now := time.Now()

	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTtl)),
		},
		UserId: userId,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtAuth) newRefreshToken(userId int, role string) (string, error) {
	now := time.Now()

	c := claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTtl)),
		},
		UserId: userId,
		Role:   role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtAuth) GenerateTokens(userId int, role string) (*pair, error) {
	access, err := j.newAccessToken(userId, role)
	if err != nil {
		return nil, err
	}

	refresh, err := j.newRefreshToken(userId, role)
	if err != nil {
		return nil, err
	}

	return &pair{Access: access, Refresh: refresh}, nil
}

func (j *JwtAuth) ParseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, incorrectToken
		}
		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, expiredErr
	} else if err != nil {
		return nil, incorrectToken
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, incorrectToken
	}

	return c, nil
}
