package entity

import "github.com/golang-jwt/jwt"

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}
