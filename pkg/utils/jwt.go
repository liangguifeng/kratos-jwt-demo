package utils

import jwt2 "github.com/golang-jwt/jwt/v4"

type MyClaims struct {
	UserName string `json:"username"`
	jwt2.RegisteredClaims
}
