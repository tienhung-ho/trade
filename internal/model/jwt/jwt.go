package jwtmodel

import "github.com/golang-jwt/jwt/v5"

type AccountJwtClaims struct {
	ID       uint64 `json:"user_id"`
	WalletID uint64 `json:"wallet_id"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}
