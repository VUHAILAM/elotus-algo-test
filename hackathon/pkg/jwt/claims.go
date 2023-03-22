package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	jwt.StandardClaims
	Email string `json:"email"`
	Type  string `json:"type"`
	Scope string `json:"scope"`
}
