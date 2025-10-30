package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func secret() []byte { return []byte(os.Getenv("JWT_SECRET")) }

func expiresAt() time.Time {
	h := 72
	if v := os.Getenv("JWT_EXPIRES_HOURS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			h = n
		}
	}
	return time.Now().Add(time.Duration(h) * time.Hour)
}

func GenerateToken(uid uint, email, role string) (string, error) {
	claims := Claims{
		UserID: uid, Email: email, Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret())
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return secret(), nil
	})
	if err != nil {
		return nil, err
	}
	if c, ok := token.Claims.(*Claims); ok && token.Valid {
		return c, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
