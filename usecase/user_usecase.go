package usecase

import (
	"event-management-system/models"
	"event-management-system/repository"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	FindAllUser() ([]models.User, error)
	FindAllEventUser() ([]models.User, error)
	FindUserById(id int) (models.User, error)
	CreateUser(input models.User) (models.User, error)
	UpdateUser(id models.GetCustomerDetailInput, input models.User) (models.User, error)
	DeleteUserById(id int) (models.User, error)
	FindUserByUsername(username string) (models.User, error)
	FindUserByEmail(email string) (models.User, error)
	FinByParams(params string, value bool) ([]models.User, error)
}

type userUseCaseImpl struct {
	userRepository repository.UserRepository
}

func NewUserUseCase(userRepository repository.UserRepository) UserUseCase {
	return &userUseCaseImpl{userRepository: userRepository}
}

func (uc *userUseCaseImpl) FindAllUser() ([]models.User, error) {
	users, err := uc.userRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *userUseCaseImpl) FindUserById(id int) (models.User, error) {
	user, err := uc.userRepository.FindById(id)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (uc *userUseCaseImpl) FindUserByUsername(username string) (models.User, error) {
	user, err := uc.userRepository.FindBySingle("username", username)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (uc *userUseCaseImpl) FindUserByEmail(email string) (models.User, error) {
	user, err := uc.userRepository.FindBySingle("email", email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (uc *userUseCaseImpl) CreateUser(input models.User) (models.User, error) {
	user := models.User{}

	var missingFields []string

	if strings.TrimSpace(input.Username) == "" {
		missingFields = append(missingFields, "Username")
	}

	if strings.TrimSpace(input.Email) == "" {
		missingFields = append(missingFields, "Email")
	}

	if strings.TrimSpace(input.Password) == "" {
		missingFields = append(missingFields, "Password")
	}

	// if strings.TrimSpace(input.Role) == "" {
	// 	missingFields = append(missingFields, "Role")
	// }

	if len(missingFields) > 0 {
		return user, fmt.Errorf("inputan %s tidak boleh string kosong", strings.Join(missingFields, ", "))
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user, err
	}

	user.Username = strings.ToLower(input.Username)
	user.Email = strings.ToLower(input.Email)
	user.Password = string(hashPassword)
	user.Role = input.Role
	user.IsActive = input.IsActive

	saveUser, err := uc.userRepository.Save(user)

	if err != nil {
		return saveUser, err
	}

	return saveUser, nil
}

func (uc *userUseCaseImpl) UpdateUser(inputId models.GetCustomerDetailInput, user models.User) (models.User, error) {
	newId, _ := strconv.Atoi(inputId.Id)
	checkUser, err := uc.FindUserById(newId)
	if err != nil {
		return checkUser, err
	}

	if strings.TrimSpace(user.Username) != "" {
		checkUser.Username = user.Username
	}

	if strings.TrimSpace(user.Email) != "" {
		checkUser.Email = user.Email
	}

	if strings.TrimSpace(user.Password) != "" {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return user, err
		}

		checkUser.Password = string(hashPassword)
	}

	// if strings.TrimSpace(user.Role) != "" {
	// 	checkUser.Role = user.Role
	// }

	if user.IsActive != nil {
		checkUser.IsActive = user.IsActive
	}

	checkUser.UpdatedAt = time.Now()

	updateUser, err := uc.userRepository.Update(checkUser)
	if err != nil {
		return updateUser, err
	}

	return updateUser, nil
}

func (uc *userUseCaseImpl) DeleteUserById(id int) (models.User, error) {
	user, err := uc.userRepository.Delete(id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (uc *userUseCaseImpl) FinByParams(params string, value bool) ([]models.User, error) {

	users, err := uc.userRepository.FindByArray(params, value)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (uc *userUseCaseImpl) FindAllEventUser() ([]models.User, error) {
	users, err := uc.userRepository.FindAllUserEvent()
	if err != nil {
		return nil, err
	}

	return users, nil
}
