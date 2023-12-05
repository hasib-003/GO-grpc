package entity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func (claims *JwtClaims) NewToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, _ := token.SignedString([]byte("secret"))
	return t
}
func RefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
