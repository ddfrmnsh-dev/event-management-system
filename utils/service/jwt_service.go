package service

import (
	"event-management-system/config"
	"event-management-system/models"
	modelUtil "event-management-system/utils/model_util"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	CreateToken(user models.User) (string, error)
	VerifyToken(tokenString string) (modelUtil.JwtPayloadClaim, error)
}

type jwtService struct {
	cfg config.TokenConfig
}

func NewJwtService(cfg config.TokenConfig) JwtService {
	return &jwtService{cfg: cfg}
}

// func (j *jwtService) CreateToken(user models.User) (string, error) {
// 	tokenKey := j.cfg.JwtSignatureKey
// 	newId := strconv.Itoa(user.Id)
// 	claims := modelUtil.JwtPayloadClaim{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Issuer:    j.cfg.ApplicationName,
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
// 		},
// 		UserId: newId,
// 		Role:   user.Role,
// 	}

// 	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
// 	token, err := jwtNewClaim.SignedString(tokenKey)
// 	if err != nil {
// 		return "", nil
// 	}

// 	return token, nil
// }

func (j *jwtService) CreateToken(user models.User) (string, error) {
	// Konversi token key ke byte array
	tokenKey := []byte(j.cfg.JwtSignatureKey)

	// Ambil ID user sebagai string
	newId := strconv.Itoa(user.Id)

	// Konversi slice Role menjadi slice string
	var roleNames []string
	for _, role := range user.Role {
		roleNames = append(roleNames, role.Name)
	}

	// Buat claims
	claims := modelUtil.JwtPayloadClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.cfg.ApplicationName,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.cfg.AccessTokenLifeTime)),
		},
		UserId: newId,
		Role:   roleNames, // Gunakan array string untuk roles
	}

	// Buat token baru
	jwtNewClaim := jwt.NewWithClaims(j.cfg.JwtSigningMethod, claims)
	token, err := jwtNewClaim.SignedString(tokenKey)
	if err != nil {
		return "", err // Kembalikan error jika terjadi
	}

	return token, nil
}

func (j *jwtService) VerifyToken(tokenString string) (modelUtil.JwtPayloadClaim, error) {
	tokenParse, err := jwt.ParseWithClaims(tokenString, &modelUtil.JwtPayloadClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.cfg.JwtSignatureKey, nil
	})

	if err != nil {
		return modelUtil.JwtPayloadClaim{}, err
	}

	claim, ok := tokenParse.Claims.(*modelUtil.JwtPayloadClaim)

	if !ok {
		return modelUtil.JwtPayloadClaim{}, fmt.Errorf("error claim")
	}

	return *claim, nil
}
