package modelutil

import "github.com/golang-jwt/jwt/v5"

type JwtPayloadClaim struct {
	jwt.RegisteredClaims
	UserId string `json:"UserId"`
	Role   string `json:"role"`
}
