package repository

import (
	"event-management-system/models"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindById(id int) (models.User, error)
	Save(user models.User) (models.User, error)
	Update(user models.User) (models.User, error)
	Delete(id int) (models.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{db: db}
}

func (u *userRepositoryImpl) FindAll() ([]models.User, error) {
	var users []models.User

	res := u.db.Find(&users) // Mengisi slice users dengan data dari database
	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return nil, res.Error // Kembalikan error jika ada
	}

	if res.RowsAffected == 0 {
		fmt.Println("no users found")
		return users, nil // Jika tidak ada pengguna ditemukan, kembalikan slice kosong dan nil error
	}

	return users, nil // Mengembalikan slice users yang ditemukan dan nil error
}

func (u *userRepositoryImpl) FindById(id int) (models.User, error) {
	var user models.User

	res := u.db.First(&user, id)
	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	if res.RowsAffected == 0 {
		fmt.Println("no users found")
		return user, nil
	}

	return user, nil
}

func (u *userRepositoryImpl) Save(user models.User) (models.User, error) {
	res := u.db.Create(&user)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return user, nil
}

func (u *userRepositoryImpl) Update(user models.User) (models.User, error) {
	res := u.db.Updates(&user)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return user, nil
}

func (u *userRepositoryImpl) Delete(id int) (models.User, error) {
	checkId, err := u.FindById(id)

	if err != nil {
		return checkId, err
	}

	var user models.User

	res := u.db.Delete(&user, id)

	if res.Error != nil {
		fmt.Println("err:", res.Error)
		return user, res.Error
	}

	return checkId, nil
}
