package usecase

import (
	"errors"
	"event-management-system/models"
	"event-management-system/utils"
	"event-management-system/utils/service"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthenticationUseCase interface {
	LoginAdmin(identifier, password string) (string, models.User, error)
	LoginUser(identifier, password string) (string, models.User, error)
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

// func (a *authenticationUseCase) LoginAdmin(input, password string) (string, models.User, error) {
// 	var user models.User
// 	var err error

// 	// Pengecekan apakah input berupa email atau username
// 	if utils.IsEmail(input) {
// 		user, err = a.userUseCase.FindUserByEmail(input)
// 	} else {
// 		user, err = a.userUseCase.FindUserByUsername(input)
// 	}

// 	// if strings.ToLower(user.Role) != "organization" && strings.ToLower(user.Role) != "admin" {
// 	// 	return "", user, errors.New("invalid credentialsss")
// 	// }

// 	hasRole := false
// 	for _, role := range user.Role {
// 		if strings.ToLower(role.Name) != "admin" && strings.ToLower(role.Name) != "super admin" && strings.ToLower(role.Name) != "event organizer" {
// 			hasRole = true
// 			fmt.Println("Role status", hasRole)
// 			break
// 		}
// 	}

// 	if !hasRole {
// 		return "", user, errors.New("you are not authorized to login as admin")
// 	}

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return "", user, errors.New("user not found")
// 		}
// 		return "", user, err
// 	}

// 	// Validasi password
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
// 		return "", user, errors.New("invalid credentials") // Hindari pesan yang detail
// 	}

// 	// Buat token
// 	token, err := a.jwtService.CreateToken(user)
// 	if err != nil {
// 		return "", user, err
// 	}

// 	return token, user, nil
// }

func (a *authenticationUseCase) LoginAdmin(input, password string) (string, models.User, error) {
	var user models.User
	var err error

	// Pengecekan apakah input berupa email atau username
	if utils.IsEmail(input) {
		user, err = a.userUseCase.FindUserByEmail(input)
	} else {
		user, err = a.userUseCase.FindUserByUsername(input)
	}

	// Validasi jika user tidak ditemukan
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", user, errors.New("user not found")
		}
		return "", user, err
	}

	// Validasi role user
	hasValidRole := false
	for _, role := range user.Role {
		fmt.Println("Role:", role.Name)
		if strings.ToLower(role.Name) == "admin" || strings.ToLower(role.Name) == "super admin" || strings.ToLower(role.Name) == "event organizer" {
			hasValidRole = true
			break
		}
	}

	if !hasValidRole {
		return "", user, errors.New("you are not authorized to login as admin")
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", user, errors.New("invalid credentials") // Hindari pesan yang terlalu detail
	}

	// Buat token
	token, err := a.jwtService.CreateToken(user)
	if err != nil {
		return "", user, err
	}

	return token, user, nil
}

func (a *authenticationUseCase) LoginUser(input, password string) (string, models.User, error) {
	var user models.User
	var err error

	// Pengecekan apakah input berupa email atau username
	if utils.IsEmail(input) {
		user, err = a.userUseCase.FindUserByEmail(input)
	} else {
		user, err = a.userUseCase.FindUserByUsername(input)
	}

	// if strings.ToLower(user.Role) != "user" {
	// 	return "", user, errors.New("invalid credentials")
	// }

	hasRole := false
	for _, role := range user.Role {
		if strings.ToLower(role.Name) == "admin" || strings.ToLower(role.Name) == "super admin" || strings.ToLower(role.Name) == "event organizer" {
			hasRole = true
			break
		}
	}

	if hasRole {
		return "", user, errors.New("you are not authorized to login as regular user")

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
