package auth

import (
	"errors"
	"oqu/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	IncorrectToken = errors.New("Incorrect token")
	ExpiredErr     = errors.New("Token expired")
)

type claims struct {
	jwt.RegisteredClaims
	UserId    int
	Role      string
	TokenType string
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
		UserId:    userId,
		Role:      role,
		TokenType: "access",
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
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTtl)),
		},
		UserId:    userId,
		Role:      role,
		TokenType: "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JwtAuth) GenerateTokens(userId int, role string) (*models.Tokens, error) {
	access, err := j.newAccessToken(userId, role)
	if err != nil {
		return nil, err
	}

	refresh, err := j.newRefreshToken(userId, role)
	if err != nil {
		return nil, err
	}

	return &models.Tokens{Access: access, Refresh: refresh}, nil
}

func (j *JwtAuth) ParseToken(tokenString string) (*claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, IncorrectToken
		}
		return j.secret, nil
	}, jwt.WithExpirationRequired())

	if errors.Is(err, jwt.ErrTokenExpired) {
		return nil, ExpiredErr
	} else if err != nil {
		return nil, IncorrectToken
	}

	c, ok := token.Claims.(*claims)
	if !ok || !token.Valid {
		return nil, IncorrectToken
	}

	return c, nil
}

func (j *JwtAuth) RefreshAccessToken(refreshToken string) (string, error) {
	refreshClaims, err := j.ParseToken(refreshToken)
	if err != nil {
		return "", IncorrectToken
	}

	if refreshClaims.TokenType != "refresh" {
		return "", IncorrectToken
	}

	accessToken, err := j.newAccessToken(refreshClaims.UserId, refreshClaims.Role)
	if err != nil {
		return "", IncorrectToken
	}

	return accessToken, nil
}
