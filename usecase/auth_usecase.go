package usecase

import (
	"errors"
	"event-management-system/models"
	"event-management-system/utils"
	"event-management-system/utils/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationUseCase interface {
	Login(identifier, password string) (string, models.User, error)
}

type authenticationUseCase struct {
	userUseCase UserUseCase
	jwtService  service.JwtService
}

func NewAuthenticationUseCase(uc UserUseCase, jwtService service.JwtService) AuthenticationUseCase {
	return &authenticationUseCase{userUseCase: uc, jwtService: jwtService}
}

// func (a *authenticationUseCase) Login(input, password string) (string, models.User, error) {
// 	var user models.User
// 	var err error

// 	if utils.IsEmail(input) {
// 		user, err = a.userUseCase.FindUserByEmail(input)
// 	} else {
// 		user, err = a.userUseCase.FindUserByUsername(input)
// 	}

// 	if err != nil {
// 		return "", user, err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		return "", user, errors.New("invalid password")
// 	}

// 	token, err := a.jwtService.CreateToken(user)
// 	if err != nil {
// 		return "", user, err
// 	}

// 	return token, user, nil
// }

func (a *authenticationUseCase) Login(input, password string) (string, models.User, error) {
	var user models.User
	var err error

	// Pengecekan apakah input berupa email atau username
	if utils.IsEmail(input) {
		user, err = a.userUseCase.FindUserByEmail(input)
	} else {
		user, err = a.userUseCase.FindUserByUsername(input)
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", user, errors.New("user not found")
		}
		return "", user, err
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", user, errors.New("invalid credentials") // Hindari pesan yang detail
	}

	// Buat token
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", user, err
	}

	return token, user, nil
}
